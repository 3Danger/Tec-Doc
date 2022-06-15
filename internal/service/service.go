package service

import (
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

func (s *Service) Start() error {
	s.log.Info().Msg("start")
	return nil
}

func (s *Service) Stop() {
	s.log.Info().Msg("stop")
}
