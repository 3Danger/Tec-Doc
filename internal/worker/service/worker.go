package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
	"tec-doc/internal/tec-doc/config"
	"tec-doc/internal/tec-doc/store/postgres"
	"tec-doc/pkg/clients/content"
	"tec-doc/pkg/model"
	"time"
)

type service struct {
	log           *zerolog.Logger
	conf          *config.Config
	store         postgres.Store
	contentClient content.ClientSource
}

func New(ctx context.Context, conf *config.Config, log *zerolog.Logger) *service {
	store, err := postgres.NewStore(ctx, &conf.Postgres)
	if err != nil {
		log.Error().Err(err).Send()
		return nil
	}

	return &service{
		log:   log,
		conf:  conf,
		store: store,
		contentClient: content.ClientSource{
			ClientJsonRPC: content.New("contentClient",
				*log,
				"http://source.content-card.svc.k8s.stage-dp/source/migration"),
		},
	}
}

func (s *service) TaskWorkerRun(ctx context.Context, conf config.WorkerConfig) error {
	s.log.Info().Msg("starting product card worker")
	tick := time.NewTicker(conf.Timer)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
			err := s.RunProductCreation(ctx, conf)
			if err != nil {
				s.log.Err(err).Send()
				time.Sleep(5 * time.Second)
			}
		}
	}
}

func (s *service) RunProductCreation(ctx context.Context, conf config.WorkerConfig) error {
	uploadID, err := s.store.GetOldestTask(ctx, nil)
	if err != nil {
		return fmt.Errorf("can't get oldest task: %w", err)
	}

	err = s.store.UpdateTaskStatus(ctx, nil, uploadID, postgres.StatusProcess)
	if err != nil {
		return fmt.Errorf("can't update task status: %w", err)
	}

	//TODO offset
	products, err := s.store.GetProductsBufferWithStatus(ctx, nil, uploadID, postgres.StatusNew, conf.Offset, 0)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_ = s.store.UpdateTaskStatus(ctx, nil, uploadID, postgres.StatusCompleted)
			return nil
		}
		return fmt.Errorf("can't get products from buffer: %w", err)
	}

	for _, pr := range products {
		if err := s.CreateProductCard(ctx, uploadID, pr); err != nil {
			return fmt.Errorf("can't create product card: %w", err)
		}
	}

	return nil
}

func (s *service) CreateProductCard(ctx context.Context, uploadID int64, product model.Product) error {
	tx, err := s.store.Transaction(ctx)
	if err != nil {
		return fmt.Errorf("can't init transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	productStatus := postgres.StatusCompleted

	//TODO content Client, serviceUUID, X-Int-Supplier-Id
	err = s.contentClient.Migration(ctx, nil)
	if err != nil {
		productStatus = postgres.StatusError
	}

	if err = s.store.UpdateProductStatus(ctx, tx, product.ID, productStatus); err != nil {
		return fmt.Errorf("can't update product status: %w", err)
	}

	var (
		failed    int64 = 0
		processed int64 = 1
	)

	if productStatus == postgres.StatusError {
		failed = 1
		processed = 0
	}

	if err := s.store.UpdateTaskProductsNumber(ctx, tx, uploadID, failed, processed); err != nil {
		return fmt.Errorf("can't update task product number: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("can't commit tx: %w", err)
	}

	return nil
}
