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
	"tec-doc/pkg/metrics"
	"tec-doc/pkg/sig"
)

// @title Tec-Doc API
// @version 1.0
// @descriptionAPI Tec-Doc server

// @host localhost:8002
// @schemes http
// @BasePath /
func initConfig() (*config.Config, *zerolog.Logger, *metrics.Metrics, error) {
	var conf config.Config
	if err := envconfig.Process("", &conf); err != nil {
		return nil, nil, nil, err
	}
	logger, err := l.InitLogger(conf.LogLevel)
	if err != nil {
		return nil, nil, nil, err
	}

	mts := metrics.NewMetrics("external", "HttpServer")

	return &conf, &logger, mts, nil
}

func main() {
	conf, logger, mts, err := initConfig()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	erg, ctx := errgroup.WithContext(context.Background())
	erg.Go(func() error {
		return sig.Listen(ctx)
	})

	svc := service.New(ctx, conf, logger, mts)
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
