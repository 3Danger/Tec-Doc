package service

import (
	"context"
	"github.com/rs/zerolog"
	"tec-doc/internal/tec-doc/config"
	"tec-doc/internal/tec-doc/model"
	"tec-doc/internal/tec-doc/store/postgres"
	"tec-doc/internal/tec-doc/web/externalserver"
	"tec-doc/internal/tec-doc/web/internalserver"
	"tec-doc/pkg/clients/services"
	"tec-doc/pkg/clients/tecdoc"
	"tec-doc/pkg/metrics"
	"time"
)

type Store interface {
	CreateTask(ctx context.Context, tx postgres.Transaction, supplierID int64, userID int64, ip string, uploadDate time.Time) (int64, error)
	SaveIntoBuffer(ctx context.Context, tx postgres.Transaction, products []model.Product) error
	GetSupplierTaskHistory(ctx context.Context, tx postgres.Transaction, supplierID int64, limit int, offset int) ([]model.Task, error)
	GetProductsBuffer(ctx context.Context, tx postgres.Transaction, uploadID int64, limit int, offset int) ([]model.Product, error)
	SaveProductsToHistory(ctx context.Context, tx postgres.Transaction, products []model.Product) error
	DeleteFromBuffer(ctx context.Context, tx postgres.Transaction, uploadID int64) error
	GetProductsHistory(ctx context.Context, tx postgres.Transaction, uploadID int64, limit int, offset int) ([]model.Product, error)
	Transaction(ctx context.Context) (postgres.Transaction, error)
	Stop()
}

type TecDocClient interface {
	GetArticles(dataSupplierID int, article string) ([]model.Article, error)
	GetBrand(brandName string) (*model.Brand, error)
}

type Server interface {
	Start() error
	Stop() error
}

type Service struct {
	conf            *config.Config
	log             *zerolog.Logger
	abacClient      services.ABAC
	suppliersClient services.Suppliers
	database        Store
	tecDocClient    TecDocClient

	externalServer Server
	internalServer Server
}

func New(ctx context.Context, conf *config.Config, log *zerolog.Logger, mts *metrics.Metrics) *Service {
	store, err := postgres.NewStore(ctx, &conf.Postgres)
	if err != nil {
		log.Error().Err(err).Send()
		return nil
	}

	svc := Service{
		conf:         conf,
		log:          log,
		database:     store,
		tecDocClient: tecdoc.NewClient(conf.TecDoc.URL, conf.TecDoc),
	}
	svc.internalServer = internalserver.New(conf.InternalServPort)
	svc.externalServer = externalserver.New(conf.ExternalServPort, &svc, log, mts)
	return &svc
}

func (s *Service) Start(ctx context.Context) error {
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

	if err = s.internalServer.Stop(); err != nil {
		s.log.Error().Err(err).Msg("error stopping internal_server")
	}
	s.log.Info().Msg("stopping internal server")

	if err = s.externalServer.Stop(); err != nil {
		s.log.Error().Err(err).Msg("error stopping external_server")
	}
	s.log.Info().Msg("stopping external server")

	s.database.Stop()
	s.log.Info().Msg("stopping database")
}

func (s *Service) Scope() *config.Scope {
	return &s.conf.Scope
}

func (s *Service) Abac() services.ABAC {
	return s.abacClient
}

func (s *Service) Suppliers() services.Suppliers {
	return s.suppliersClient
}
