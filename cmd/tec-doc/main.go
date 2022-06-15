package main

import (
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"tec-doc/internal/config"
	"tec-doc/internal/logger"
	"tec-doc/internal/service"
)

func main() {
	var (
		err    error
		conf   *config.Config
		egroup *errgroup.Group
		log    zerolog.Logger
		srvc   *service.Service
	)

	if conf, err = config.NewConfigEnv(); err != nil {
		panic(err)
	}
	log = logger.NewLogger(conf)
	srvc = service.NewService(conf, log)

	egroup = new(errgroup.Group)
	egroup.Go(srvc.Start)
	if err = egroup.Wait(); err != nil {
		log.Error().Msg(err.Error())
		srvc.Stop()
	}
}
