package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ServerAddress   string `envconfig:"EXTERNAL_SERV_PORT" default:":8050"`
	FrontendAddress string `envconfig:"FRONTEND_PORT" default:":8001"`
}

func Get() (config *Config) {
	config = new(Config)
	if err := envconfig.Process("TEC_DOC", config); err != nil {
		log.Error().Err(err).Send()
		return nil
	}
	return config
}
