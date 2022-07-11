package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ServerHost   string `envconfig:"EXTERNAL_SERV_HOST" required:"true"`
	ServerPort   string `envconfig:"EXTERNAL_SERV_PORT" required:"true"`
	FrontendPort string `envconfig:"FRONTEND_ADDRESS" default:"8002"`
}

func Get() (config *Config) {
	config = new(Config)
	if err := envconfig.Process("TEC_DOC", config); err != nil {
		log.Error().Err(err).Send()
		return nil
	}
	return config
}
