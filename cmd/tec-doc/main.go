package main

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"tec-doc/internal/tec-doc/config"
	l "tec-doc/internal/tec-doc/logger"
	"tec-doc/internal/tec-doc/service"
	"tec-doc/pkg/sig"
)

// todo: генерация метрик на этом уровне и прокидываем дальше в сервис и server
func initConfig() (*config.Config, *zerolog.Logger, error) {
	var conf config.Config
	if err := envconfig.Process("", &conf); err != nil {
		return nil, nil, err
	}
	logger, err := l.InitLogger(conf.LogLevel)
	if err != nil {
		return nil, nil, err
	}
	return &conf, &logger, nil
}

func main() {
	conf, logger, err := initConfig()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	erg, ctx := errgroup.WithContext(context.Background())
	erg.Go(func() error {
		return sig.Listen(ctx)
	})

	svc := service.New(ctx, conf, logger)
	logger.Info().Msg("service starting..")
	erg.Go(func() error {
		return svc.Start(ctx)
	})

	if err = erg.Wait(); err != nil {
		logger.Error().Err(err).Send()
	}
	svc.Stop()
	logger.Info().Msg("graceful shutdown done")
}
