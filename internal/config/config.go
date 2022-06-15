package config

import _ "github.com/kelseyhightower/envconfig"

type Config struct {
	LogLevel      string `envconfig:"LOG_LEVEL"`
	LogTimeFormat string `envconfig:"LOG_TM_FORMAT"`

	ServerPort string `envconfig:"SERVER_PORT" required:"true"`

	DbUsername string `envconfig:"DATABASE_USERNAME"`
	DbPassword string `envconfig:"DATABASE_PASSWORD"`
	DbName     string `envconfig:"DATABASE_NAME"`
	DbPort     string `envconfig:"DATABASE_PORT"`
}
