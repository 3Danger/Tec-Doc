package main

import (
	"golang.org/x/sync/errgroup"
	"tec-doc/internal/config"
	"tec-doc/internal/logger"
	"tec-doc/internal/service"
)

func main() {
	conf, err := config.NewConfigEnv()
	if err != nil {
		panic(err)
	}
	log := logger.NewLogger(conf)
	srvc := service.NewService(conf, log)

	eg := new(errgroup.Group)
	eg.Go(srvc.Start)
	if err = eg.Wait(); err != nil {
		log.Error().Msg(err.Error())
	}

	if err = srvc.Stop(); err != nil {
		log.Error().Msg(err.Error())
	}
	//eGroupServer.Go()

}
