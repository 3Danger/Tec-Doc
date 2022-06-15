package service

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"tec-doc/internal/config"
)

type Service struct {
	conf *config.Config
	log  zerolog.Logger
}

func NewService(conf *config.Config, log zerolog.Logger) *Service {
	return &Service{
		conf: conf,
		log:  log,
	}
}

func (s *Service) Start(ctx context.Context) error {
	if ctx == nil {
		return errors.New("ctx is empty")
	}
	select {
	case <-ctx.Done():
		return nil
	default:
		s.log.Info().Str("", "").Msg("start on port: " + s.conf.ServerPort)
	}
	return nil
}

func (s *Service) Stop() error {
	s.log.Info().Msg("stop")
	return nil
}
