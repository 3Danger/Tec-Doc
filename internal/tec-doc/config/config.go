package config

import (
	_ "github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	InternalServPort string `envconfig:"INTERNAL_SERV_PORT" required:"true"`
	ExternalServPort string `envconfig:"EXTERNAL_SERV_PORT" required:"true"`
	LogLevel         string `envconfig:"LOG_LEVEL" default:"debug"`
	Postgres         PostgresConfig
	Content          ContentClientConfig
	TecDoc           TecDocClientConfig
	Scope            Scope
}

type PostgresConfig struct {
	Username string        `envconfig:"POSTGRES_USERNAME" required:"true"`
	Password string        `envconfig:"POSTGRES_PASSWORD" required:"true"`
	Host     string        `envconfig:"POSTGRES_HOST" required:"true"`
	Port     string        `envconfig:"POSTGRES_PORT" required:"true"`
	DbName   string        `envconfig:"POSTGRES_DB" required:"true"`
	Timeout  time.Duration `envconfig:"POSTGRES_TIMEOUT" default:"30s"`
	MaxConns int32         `envconfig:"POSTGRES_MAX_CONNECTIONS" default:"100"`
	MinConns int32         `envconfig:"POSTGRES_MIN_CONNECTIONS" default:"10"`
}

type ContentClientConfig struct {
	URL     string        `envconfig:"CONTENT_CLIENT_URL"`
	Timeout time.Duration `envconfig:"CONTENT_CLIENT_TIMEOUT" default:"30s"`
}

type TecDocClientConfig struct {
	URL        string        `envconfig:"TEC_DOC_CLIENT_URL"`
	Timeout    time.Duration `envconfig:"TEC_DOC_CLIENT_TIMEOUT" default:"30s"`
	XApiKey    string        `envconfig:"TEC_DOC_CLIENT_API_KEY"`
	ProviderId int           `envconfig:"TEC_DOC_CLIENT_PROVIDER_ID"`
}

type Scope struct {
	Scope                   string `envconfig:"SCOPE" default:"suppliers"`
	ContentFeatureKey       string `envconfig:"CONTENT_FEATURE_KEY" default:"content"`
	UpdateContentFeatureKey string `envconfig:"UPDATE_CONTENT_FEATURE_KEY" default:"content/update"`
}
