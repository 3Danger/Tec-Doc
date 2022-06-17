package config

import (
	"time"

	_ "github.com/kelseyhightower/envconfig"
)

type Config struct {
	LogLevel            string `envconfig:"LOG_LEVEL" default:"debug"`
	ListenInternal      string `envconfig:"LISTEN_INTERNAL" default:":8000"`
	PostgresConfig      PostgresConfig
	ContentClientConfig ContentClientConfig
	TecDocConfig        TecDocConfig
}

type PostgresConfig struct {
	Username string        `envconfig:"POSTGRES_USERNAME"`
	Password string        `envconfig:"POSTGRES_PASSWORD"`
	Host     string        `envconfig:"POSTGRES_HOST"`
	Port     string        `envconfig:"POSTGRES_PORT"`
	DbName   string        `envconfig:"POSTGRES_DB"`
	Timeout  time.Duration `envconfig:"POSTGRES_TIMEOUT" default:"30s"`
	MaxConns int32         `envconfig:"MAX_CONNECTIONS" default:"100"`
	MinConns int32         `envconfig:"MIN_CONNECTIONS" default:"10"`
}

type ContentClientConfig struct {
	URL     string        `envconfig:"CONTENT_URL"`
	Timeout time.Duration `envconfig:"CONTENT_TIMEOUT" default:"30s"`
}

type TecDocConfig struct {
	URL     string        `envconfig:"TECDOC_URL"`
	Timeout time.Duration `envconfig:"TECDOC_TIMEOUT" default:"30s"`
}
