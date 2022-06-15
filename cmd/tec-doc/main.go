package main

import (
	"context"
	"golang.org/x/sync/errgroup"
	"tec-doc/internal/config"
	"tec-doc/internal/logger"
	"tec-doc/internal/service"
	"time"
)

func main() {
	conf, err := config.NewConfigEnv()
	if err != nil {
		panic(err)
	}
	log := logger.NewLogger(conf)
	srvc := service.NewService(conf, log)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	eg := new(errgroup.Group)
	eg.Go(func() error {
		return srvc.Start(ctx)
	})
	if err = eg.Wait(); err != nil {
		log.Error().Msg(err.Error())
	}

	if err = srvc.Stop(); err != nil {
		log.Error().Msg(err.Error())
	}
	//eGroupServer.Go()

}
