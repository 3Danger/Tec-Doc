package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
	"tec-doc/internal/config"
	l "tec-doc/internal/logger"
	"tec-doc/internal/service"
	"tec-doc/internal/web/internalserver"
)

func initConfig() (*config.Config, *zerolog.Logger, error) {
	var (
		conf   *config.Config
		logger zerolog.Logger
		err    error
	)
	// Init Config
	conf = new(config.Config)
	if err = envconfig.Process("TEC_DOC", conf); err != nil {
		return nil, nil, err
	}

	// Init Logger
	logger, err = l.InitLogger(strings.ToLower(conf.LogLevel))
	if err != nil {
		return nil, nil, err
	}
	return conf, &logger, nil
}

func main() {
	conf, logger, err := initConfig()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	srvc := service.NewService(conf, logger)
	srvr := internalserver.New(conf.InternalServAddress, srvc)

	if err = srvr.Start(); err != nil {
		logger.Error().Err(err).Send()
	}
	if err = srvr.Stop(); err != nil {
		logger.Error().Err(err).Send()
	}

}
