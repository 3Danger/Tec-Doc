package main

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"strings"
	"tec-doc/internal/config"
	l "tec-doc/internal/logger"
	"tec-doc/internal/service"
	"tec-doc/internal/web/externalserver"
	"tec-doc/internal/web/internalserver"
	"tec-doc/pkg/sig"
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

	egroup, ctx := errgroup.WithContext(context.Background())
	egroup.Go(func() error {
		return sig.Listen(ctx)
	})

	srvc := service.New(conf, logger)

	internalServ := internalserver.New(conf.InternalServPort)
	externalServ := externalserver.New(conf.ExternalServPort, srvc, logger)

	srvc.SetInternalServer(internalServ)
	srvc.SetExternalServer(externalServ)

	egroup.Go(func() error {
		return srvc.Start(ctx)
	})

	if err = egroup.Wait(); err != nil {
		logger.Error().Err(err).Send()
	}
	srvc.Stop()
}
