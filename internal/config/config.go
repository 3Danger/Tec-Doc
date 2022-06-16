package config

import _ "github.com/kelseyhightower/envconfig"

type Config struct {
	LogLevel string `envconfig:"LOG_LEVEL"`
}
