package main

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"strings"
	"tec-doc/internal/tec-doc/config"
	l "tec-doc/internal/tec-doc/logger"
	"tec-doc/internal/worker/service"
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
	errGr, ctx := errgroup.WithContext(context.Background())
	srvc := service.New(ctx, conf, logger)

	//TODO timer
	errGr.Go(func() error {
		return srvc.TaskWorkerRun(ctx, conf.Worker)
	})

	if err := errGr.Wait(); err != nil {
		logger.Error().Err(err).Send()
	}
}
