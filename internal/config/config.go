package config

import _ "github.com/kelseyhightower/envconfig"

type Config struct {
	Debug struct {
		Level      string `envconfig:"LEVEL"`
		TimeFormat string `envconfig:"TM_FORMAT"`
	}
	Server struct {
		Port string `envconfig:"PORT" required:"true"`
	}
}
