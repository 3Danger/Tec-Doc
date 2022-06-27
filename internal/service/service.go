package service

import (
	"context"
	"github.com/rs/zerolog"
	"tec-doc/internal/config"
	"tec-doc/internal/model"
	"tec-doc/internal/store/postgres"
	"time"
)

type Store interface {
	CreateTask(ctx context.Context, supplierID int, userID int, ip string, uploadDate time.Time) (int64, error)
	SaveIntoBuffer(ctx context.Context, products []model.Product) error
	GetSupplierTaskHistory(ctx context.Context, supplierID int, limit int, offset int) ([]model.Task, error)
	GetProductsFromBuffer(ctx context.Context, uploadID int) ([]model.Product, error)
	SaveProductsToHistory(ctx context.Context, products []model.Product) error
	DeleteFromBuffer(ctx context.Context, uploadID int) error
	GetProductsHistory(ctx context.Context, uploadID int, limit int, offset int) ([]model.Product, error)
}

type Server interface {
	Start() error
	Stop() error
}

type Service struct {
	conf           *config.Config
	log            *zerolog.Logger
	externalServer Server
	internalServer Server
	database       Store
	//tec_doc_client Client
}

func New(conf *config.Config, log *zerolog.Logger) *Service {
	store, err := postgres.NewStore(&conf.PostgresConfig)
	if err != nil {
		log.Error().Err(err).Send()
		return nil
	}
	log.Info().Msg("create service")
	return &Service{
		conf:     conf,
		log:      log,
		database: store,
		//tec_doc_client Client,
	}
}

func (s *Service) SetInternalServer(internalServer Server) {
	s.internalServer = internalServer
}

func (s *Service) SetExternalServer(externalServer Server) {
	s.externalServer = externalServer
}

func (s *Service) Start(ctx context.Context) error {
	s.log.Info().Msg("starting service")
	errChan := make(chan error, 2)
	go func() {
		errChan <- s.internalServer.Start()
	}()
	go func() {
		errChan <- s.externalServer.Start()
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		return err
	}
}

func (s *Service) Stop() {
	var err error

	s.log.Info().Msg("stopping service")
	if err = s.internalServer.Stop(); err != nil {
		s.log.Error().Err(err).Msg("error stopping internal_server")
	}
	if err = s.externalServer.Stop(); err != nil {
		s.log.Error().Err(err).Msg("error stopping external_server")
	}
}
