package config

import (
	"time"

	_ "github.com/kelseyhightower/envconfig"
)

type Config struct {
	InternalServPort string `envconfig:"INTERNAL_SERV_PORT" required:"true"`
	ExternalServPort string `envconfig:"EXTERNAL_SERV_PORT" required:"true"`
	LogLevel         string `envconfig:"LOG_LEVEL" default:"debug"`
	ListenInternal   string `envconfig:"LISTEN_INTERNAL" default:":8000"`
	Postgres         PostgresConfig
	Content          ContentClientConfig
	TecDoc           TecDocConfig
}

type PostgresConfig struct {
	Username string        `envconfig:"USERNAME" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Host     string        `envconfig:"HOST" required:"true"`
	Port     string        `envconfig:"PORT" required:"true"`
	DbName   string        `envconfig:"DB" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" default:"30s"`
	MaxConns int32         `envconfig:"MAX_CONNECTIONS" default:"100"`
	MinConns int32         `envconfig:"MIN_CONNECTIONS" default:"10"`
}

type ContentClientConfig struct {
	URL     string        `envconfig:"URL"`
	Timeout time.Duration `envconfig:"TIMEOUT" default:"30s"`
}

type TecDocConfig struct {
	URL        string        `envconfig:"URL"`
	Timeout    time.Duration `envconfig:"TIMEOUT" default:"30s"`
	XApiKey    string        `envconfig:"API_KEY"`
	ProviderId int           `envconfig:"PROVIDER_ID"`
}
