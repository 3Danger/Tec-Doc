package service

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
	"io"
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
		tick     = time.NewTicker(time.Second * 10)
		uplaodID int64
		pe       []model.ProductEnriched
	)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
			if pe, uplaodID, err = s.getProductsEnriched(ctx); err != nil {
				time.Sleep(5 * time.Second)
				continue
			}
			if len(pe) > 0 {
				if err = s.runProductCreation(ctx, pe, uplaodID); err != nil {
					time.Sleep(5 * time.Second)
				}
			}
		}
	}
}

func (s *service) getProductsEnriched(ctx context.Context) (pe []model.ProductEnriched, uploadID int64, err error) {
	if uploadID, err = s.store.GetOldestTask(ctx, nil); errors.Is(err, pgx.ErrNoRows) {
		return nil, 0, nil
	}
	if err != nil {
		s.log.Error().Err(err).Str("store", "can't get oldest task")
		return nil, 0, err
	}

	var products []model.Product
	products, err = s.store.GetProductsBufferWithStatus(ctx, nil, uploadID, 1000, 0, postgres.StatusNew)
	if err != nil {
		s.log.Error().Err(err).Str("store", "can't get products from buffer")
		return
	}

	if len(products) == 0 {
		var tx postgres.Transaction
		if tx, err = s.store.Transaction(ctx); err != nil {
			s.log.Error().Err(err).Str("store", "can't create a transaction")
			return
		}
		defer func() { _ = tx.Rollback(ctx) }()
		if err = s.store.MoveProductsToHistoryByUploadId(ctx, tx, uploadID); err != nil {
			s.log.Error().Err(err).Str("store", "can't move products from buffer to history")
			return
		}
		if err = s.store.UpdateTaskStatus(ctx, tx, uploadID, postgres.StatusCompleted); err != nil {
			s.log.Error().Err(err).Str("store", "can't update status task")
			return
		}
		err = tx.Commit(ctx)
		return nil, 0, err
	}
	pe, err = s.enricher.Enrichment(products)
	return
}

func (s *service) runProductCreation(ctx context.Context, pe []model.ProductEnriched, uploadID int64) (err error) {
	var tx postgres.Transaction
	if tx, err = s.store.Transaction(ctx); err != nil {
		s.log.Error().Err(err).Str("store", "can't create a transaction")
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var processed, failed int64
	for i := range pe {
		if pe[i].Status != postgres.StatusError {
			if err = s.CreateProductCard(&pe[i]); err != nil {
				s.log.Error().Err(err).Str("CreateProductCard", "can't upload")
				pe[i].ErrorResponse = "не удалось сформировать карточку"
				pe[i].Status = postgres.StatusError
			}
		}
		if err = s.store.UpdateProductStatus(ctx, tx, pe[i].ID, pe[i].Status); err != nil {
			s.log.Error().Err(err).Str("store", "can't update product status")
			return err
		}
		if pe[i].Status == postgres.StatusError {
			if err = s.store.UpdateProductBufferErrorResponse(ctx, tx, pe[i].Product); err != nil {
				s.log.Error().Err(err).Str("store", "can't update product_buffer field errorResponse")
			}
			failed++
			continue
		}
		processed++
	}

	if err = s.store.UpdateTaskProductsNumber(ctx, tx, uploadID, failed, processed); err != nil {
		s.log.Error().Err(err).Str("store", "can't update task product number")
		return err
	}
	if err = s.store.UpdateTaskStatus(ctx, tx, uploadID, postgres.StatusProcess); err != nil {
		s.log.Error().Err(err).Str("store", "can't update task status")
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		s.log.Error().Err(err).Str("store", "can't call Commit on transaction")
		return err
	}
	return nil
}

func (s *service) CreateProductCard(product *model.ProductEnriched) error {
	//TODO что это, для чего это ?
	//TODO content Client, serviceUUID, X-Int-Supplier-Id
	//if err = s.contentClient.Migration(ctx, nil); err != nil {
	//	productStatus = postgres.StatusError
	//}
	var (
		body io.Reader
		err  error
	)
	characteristics := s.enricher.ConvertToCharacteristics(product)
	if body, err = s.makeUploadBody(characteristics); err != nil {
		return err
	}
	return s.contentClient.Upload(body)
}
