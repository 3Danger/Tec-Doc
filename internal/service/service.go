package service

import (
	"context"
	"database/sql"
	"github.com/rs/zerolog"
	"tec-doc/internal/config"
)

type Server interface {
	Start() error
	Stop() error
}

type Service struct {
	conf           *config.Config
	log            *zerolog.Logger
	externalServer Server
	internalServer Server
	products       map[int]*Product

	//TODO initialise it
	database *sql.DB
}

func New(conf *config.Config, log *zerolog.Logger) *Service {
	log.Info().Msg("create service")
	return &Service{
		conf:     conf,
		log:      log,
		products: make(map[int]*Product),
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
