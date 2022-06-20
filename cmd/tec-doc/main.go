package main

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"os"
	"strings"
	"tec-doc/internal/config"
	l "tec-doc/internal/logger"
	s "tec-doc/internal/service"
	"tec-doc/internal/web/internalserver"
	"time"
)

func TestExcelling() {
	array := []s.DummyXLSX{
		{"1", "John", 792853.37833},
		{"2", "Mitchel", 10202.13123},
		{"3", "Superman", +777.7777},
		{"4", "Dummy", 0101010.10110},
	}
	err := s.ToExcel("someExcelFile", "", array)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	os.Exit(0)
}

func main() {
	TestExcelling()
	TestConnection()
	var (
		err    error
		conf   *config.Config
		ctxeg  context.Context
		egroup *errgroup.Group
		logger *zerolog.Logger
		svc    *s.Service
	)

	// init config & logger
	conf, logger, err = initConfig()
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	svc = s.NewService(conf, logger)
	egroup, ctxeg = errgroup.WithContext(context.Background())

	egroup.Go(func() error {
		return svc.Start(ctxeg)
	})

	if err = egroup.Wait(); err != nil {
		log.Error().Msg(err.Error())
	}
	svc.Stop()
}

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

func TestConnection() {
	//Init
	conf := new(config.Config)
	if err := envconfig.Process("TEC_DOC", conf); err != nil {
		log.Err(err).Send()
	}

	//Start server
	serv := internalserver.New(conf.InternalServAddress)
	go func() {
		err := serv.Start()
		if err != nil {
			log.Error().Err(err).Send()
		}
	}()

	// Stop server
	ctx, closer := context.WithTimeout(context.Background(), time.Second*2500)
	defer closer()
	go func(ctx context.Context) {
		<-ctx.Done()
		err := serv.Stop()
		if err != nil {
			log.Error().Err(err).Send()
		}
	}(ctx)
	<-ctx.Done()

	// When Done
	time.Sleep(time.Second)
	fmt.Println("Done!")
	os.Exit(0)
}
