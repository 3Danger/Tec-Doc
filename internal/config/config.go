package config

import _ "github.com/kelseyhightower/envconfig"

type Config struct {
	LogLevel            string `envconfig:"LOG_LEVEL"`
	InternalServAddress string `envconfig:"INTERNAL_SERV_ADDRESS" default:":8000"`
	ExternalServAddress string `envconfig:"EXTERNAL_SERV_ADDRESS" default:":8050"`
}
