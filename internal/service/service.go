package service

import (
	"context"
	"github.com/rs/zerolog"
	"tec-doc/internal/config"
)

type Service struct {
	conf *config.Config
	log  *zerolog.Logger
}

func NewService(conf *config.Config, log *zerolog.Logger) *Service {
	log.Info().Msg("create service")
	return &Service{
		conf: conf,
		log:  log,
	}
}

func (s *Service) Start(ctx context.Context) error {
	s.log.Info().Msg("starting service")

	//for {} ...
	select {
	case <-ctx.Done():
		s.log.Info().Msg("service done")
		return nil
	default:
		// Do something ....
	}
	return nil
}

func (s *Service) Stop() {
	s.log.Info().Msg("stopping service")
}
