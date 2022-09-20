package service

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
	"tec-doc/internal/tec-doc/config"
	"tec-doc/internal/tec-doc/store/postgres"
	"tec-doc/pkg/clients/contentCard"
	"tec-doc/pkg/clients/tecdoc"
	"tec-doc/pkg/model"
	"time"
)

type Enricher interface {
	Enrichment(products []model.Product) (productsEnriched []model.ProductEnriched, err error)
	ConvertToCharacteristics(pe *model.ProductEnriched) *model.ProductCharacteristics
}

type Service interface {
	TaskWorkerRun(ctx context.Context) (err error)
}

type service struct {
	log           *zerolog.Logger
	conf          *config.Config
	store         postgres.Store
	contentClient contentCard.ContentCardClient
	enricher      Enricher // for Enrichment products
}

func New(ctx context.Context, conf *config.Config, log *zerolog.Logger) Service {
	var (
		err           error
		store         postgres.Store
		contentClient contentCard.ContentCardClient
	)
	if store, err = postgres.NewStore(ctx, &conf.Postgres); err != nil {
		log.Error().Err(err).Str("worker", "can't create store").Send()
		return nil
	}

	if contentClient, err = contentCard.New(&conf.Content); err != nil {
		log.Error().Err(err).Str("worker", "can't create client of contentCard").Send()
		return nil
	}

	return &service{
		log:           log,
		conf:          conf,
		store:         store,
		contentClient: contentClient,
		enricher:      tecdoc.NewClient(conf.TecDoc.URL, conf.TecDoc, log),
	}
}

func (s *service) TaskWorkerRun(ctx context.Context) (err error) {
	s.log.Info().Msg("starting product card worker")
	var (
		logger           = s.log.With().Str("worker", "card").Logger()
		tick             = time.NewTicker(time.Second * 1)
		productsEnriched []model.ProductEnriched
		supplierStr      string
	)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
			if productsEnriched, supplierStr, err = s.getProductsEnriched(ctx); err != nil {
				logger.Err(err).Msg("can't enrich the products")
				continue
			}
			if len(productsEnriched) > 0 {
				if err = s.runProductCreation(productsEnriched, supplierStr); err != nil {
					logger.Err(err).Msg("can't create card")
				}
				if err = s.UpdateStatus(ctx, productsEnriched); err != nil {
					logger.Err(err).Msg("can't update status on store")
				}
			}
		}
	}
}

func (s *service) getProductsEnriched(ctx context.Context) (productsEnriched []model.ProductEnriched, supplier string, err error) {
	var (
		uploadID int64
		products []model.Product
		tx       postgres.Transaction
	)
	if uploadID, supplier, err = s.store.GetOldestTask(ctx, nil); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", nil
		}
		s.log.Error().Err(err).Str("store", "can't get oldest task")
		return nil, "", err
	}

	if products, err = s.store.GetProductsBufferWithStatus(ctx, nil, uploadID, 1000, 0, postgres.StatusNew); err != nil {
		s.log.Error().Err(err).Str("store", "can't get products from buffer")
		return nil, "", err
	}

	if len(products) == 0 {
		if tx, err = s.store.Transaction(ctx); err != nil {
			s.log.Error().Err(err).Str("store", "can't create a transaction")
			return nil, "", err
		}
		if err = s.store.MoveProductsToHistoryByUploadId(ctx, tx, uploadID); err != nil {
			s.log.Error().Err(err).Str("store", "can't move products from buffer to history")
			_ = tx.Rollback(ctx)
			return nil, "", err
		}
		if err = s.store.UpdateTaskStatus(ctx, tx, uploadID, postgres.StatusCompleted); err != nil {
			s.log.Error().Err(err).Str("store", "can't update status task")
			_ = tx.Rollback(ctx)
			return nil, "", err
		}
		return nil, "", tx.Commit(ctx)
	}
	if productsEnriched, err = s.enricher.Enrichment(products); err != nil {
		return nil, "", err
	}
	return productsEnriched, supplier, nil
}

func (s *service) runProductCreation(productsEnriched []model.ProductEnriched, supplierIdStr string) error {
	uploader, err := s.contentClient.Upload(supplierIdStr)
	if err != nil {
		return err
	}

	for i := range productsEnriched {
		if productsEnriched[i].Status != postgres.StatusError {
			body, err := s.makeUploadBody(&productsEnriched[i])
			if err != nil {
				return err
			}
			if err = uploader(body); err != nil {
				s.log.Error().Err(err).Str("CreateProductCard", "can't upload")
				productsEnriched[i].Status = postgres.StatusError
				productsEnriched[i].ErrorResponse = "не удалось сформировать карточку"
			} else {
				productsEnriched[i].Status = postgres.StatusCompleted
			}
		}
	}

	return nil
}

func (s *service) UpdateStatus(ctx context.Context, productsEnriched []model.ProductEnriched) (err error) {
	tx, err := s.store.Transaction(ctx)
	if err != nil {
		s.log.Error().Err(err).Str("store", "can't create a transaction")
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()
	var processed, failed int64
	for i := range productsEnriched {
		if err = s.store.UpdateProductBuffer(ctx, tx, &productsEnriched[i].Product); err != nil {
			s.log.Error().Err(err).Str("store", "can't update product status")
			return err
		}
		if productsEnriched[i].Status == postgres.StatusCompleted {
			processed++
			continue
		}
		failed++
	}

	if err = s.store.UpdateTaskProductsNumber(ctx, tx, productsEnriched[0].UploadID, failed, processed); err != nil {
		s.log.Error().Err(err).Str("store", "can't update task product number")
		return err
	}
	if err = s.store.UpdateTaskStatus(ctx, tx, productsEnriched[0].UploadID, postgres.StatusProcess); err != nil {
		s.log.Error().Err(err).Str("store", "can't update task status")
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		s.log.Error().Err(err).Str("store", "can't call Commit on transaction")
		return err
	}
	return nil
}
