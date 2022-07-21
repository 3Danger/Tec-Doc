package service

import (
	"context"
	"github.com/rs/zerolog"
	"tec-doc/internal/tec-doc/config"
	"tec-doc/internal/tec-doc/model"
	"tec-doc/internal/tec-doc/store/postgres"
	"tec-doc/pkg/clients/tecdoc"
	"time"
)

//go:generate -source=service.go -destination=mock/service_mock.go

type Store interface {
	CreateTask(ctx context.Context, tx postgres.Transaction, supplierID int64, userID int64, ip string, uploadDate time.Time) (int64, error)
	SaveIntoBuffer(ctx context.Context, tx postgres.Transaction, products []model.Product) error
	GetSupplierTaskHistory(ctx context.Context, tx postgres.Transaction, supplierID int64, limit int, offset int) ([]model.Task, error)
	GetProductsFromBuffer(ctx context.Context, tx postgres.Transaction, uploadID int64) ([]model.Product, error)
	SaveProductsToHistory(ctx context.Context, tx postgres.Transaction, products []model.Product) error
	DeleteFromBuffer(ctx context.Context, tx postgres.Transaction, uploadID int64) error
	GetProductsHistory(ctx context.Context, tx postgres.Transaction, uploadID int64, limit int, offset int) ([]model.Product, error)
	Transaction(ctx context.Context) (postgres.Transaction, error)
}

type TecDocClient interface {
	GetArticles(ctx context.Context, tecDocCfg config.TecDocConfig, dataSupplierID int, article string) ([]model.Article, error)
	GetBrand(ctx context.Context, tecDocCfg config.TecDocConfig, brandName string) (*model.Brand, error)
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
	tecDocClient   TecDocClient
}

func New(conf *config.Config, log *zerolog.Logger) *Service {
	store, err := postgres.NewStore(&conf.Postgres)
	if err != nil {
		log.Error().Err(err).Send()
		return nil
	}
	log.Info().Msg("create service")
	return &Service{
		conf:         conf,
		log:          log,
		database:     store,
		tecDocClient: tecdoc.NewClient(conf.TecDoc.URL, conf.TecDoc.Timeout),
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
