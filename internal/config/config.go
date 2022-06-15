package config

import _ "github.com/kelseyhightower/envconfig"

type Config struct {
	LogLevel      string `envconfig:"LOG_LEVEL"`
	LogTimeFormat string `envconfig:"LOG_TM_FORMAT"`
}
