package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"tec-doc/frontend/client"
	"tec-doc/frontend/config"
)

func main() {
	var (
		err          error
		clientTecDoc *client.Client
	)
	clientTecDoc = client.New(config.Get())
	err = clientTecDoc.Run(context.Background())
	if err != nil {
		log.Error().Err(err).Send()
	}
}
