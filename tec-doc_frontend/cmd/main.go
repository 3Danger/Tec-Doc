package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"tec-doc/internal/client"
	"tec-doc/internal/config"
)

func main() {
	var (
		err          error
		clientTecDoc *client.Client
	)
	cnf := config.Get()
	clientTecDoc = client.New(cnf)
	fmt.Printf("%+v\n", *cnf)
	err = clientTecDoc.Run(context.Background())
	if err != nil {
		log.Error().Err(err).Send()
	}
}
