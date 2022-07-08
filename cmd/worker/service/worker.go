package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
	s "tec-doc/cmd/worker/postgres"
	"tec-doc/internal/config"
	"tec-doc/internal/model"
	"tec-doc/internal/store/postgres"
	"time"
)

type ContentClient interface {
	CreateCard(ctx context.Context, product model.Product) error
}

type Store interface {
	GetOldestTask(ctx context.Context, tx postgres.Transaction) (int64, int64, int64, error)
	GetProductFromBuffer(ctx context.Context, tx postgres.Transaction, uploadID, offset int64) (model.Product, error)
	UpdateProductStatus(ctx context.Context, tx postgres.Transaction, productID int64, status int) error
	UpdateTaskProductsNumber(ctx context.Context, tx postgres.Transaction, uploadID, productsFailed, productsProcessed int64) (bool, error)
	UpdateTaskStatus(ctx context.Context, tx postgres.Transaction, uploadID int64, status int) error
	Transaction(ctx context.Context) (postgres.Transaction, error)
}

type service struct {
	log           *zerolog.Logger
	conf          *config.Config
	store         Store
	contentClient ContentClient
}

func New(conf *config.Config, log *zerolog.Logger) *service {
	store, err := s.NewStore(&conf.PostgresConfig)
	if err != nil {
		log.Error().Err(err).Send()
		return nil
	}

	return &service{
		log:   log,
		conf:  conf,
		store: store,
	}
}

func (s *service) TaskWorkerRun(ctx context.Context, timer time.Duration) error {
	s.log.Info().Msg("starting product card worker")
	tick := time.NewTicker(timer)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
			err := s.CreateProductCard(ctx)
			//TODO check whether any undone tasks left
			if errors.Is(err, errors.New("no tasks left")) {
				s.log.Info().Msg("product card worker sleep")
				time.Sleep(5 * time.Minute)
			} else if err != nil {
				return fmt.Errorf("can't create product card: %w", err)
			}
		}
	}
}

func (s *service) CreateProductCard(ctx context.Context) error {
	tx, err := s.store.Transaction(ctx)
	if err != nil {
		return fmt.Errorf("can't init transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	uploadID, offset, total, err := s.store.GetOldestTask(ctx, tx)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("can't get oldest task: %w", errors.New("no tasks left"))
		}

	}

	err = s.store.UpdateTaskStatus(ctx, tx, uploadID, postgres.StatusProcess)
	if err != nil {
		return fmt.Errorf("can't update task status: %w", err)
	}

	var failed, processed int64

	for offset < total {
		ok := true
		product, err := s.store.GetProductFromBuffer(ctx, tx, uploadID, offset)
		if err != nil {
			return fmt.Errorf("can't get products from buffer: %w", err)
		}

		//TODO content Client
		err = s.contentClient.CreateCard(ctx, product)
		if err != nil {
			ok = false
		}

		var productStatus int
		if ok {
			processed += 1
			status = postgres.StatusCompleted
		} else {
			failed += 1
			status = postgres.StatusError
		}

		err = s.store.UpdateProductStatus(ctx, tx, product.ID, productStatus)
		if err != nil {
			return fmt.Errorf("can't update product status: %w", err)
		}

		offset += 1
	}

	correct, err := s.store.UpdateTaskProductsNumber(ctx, tx, uploadID, failed, processed)

	if correct == false || err != nil {
		if err != nil {
			return fmt.Errorf("can't update task product number: %w", err)
		}
	}
	err = s.store.UpdateTaskStatus(ctx, tx, uploadID, postgres.StatusCompleted)
	if err != nil {
		return fmt.Errorf("can't update task status: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("can't commit tx: %w", err)
	}

	return nil
}
