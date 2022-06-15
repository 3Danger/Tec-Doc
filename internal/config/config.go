package config

type Config struct {
	LogLevel			string	`envconfig:"LOG_LEVEL" default:"debug"`
	ListenInternal		string	`envconfig:"LISTEN_INTERNAL" default:":8000"`
	PostgresConfig		PostgresConfig
	ContentClientConfig	ContentClientConfig
	TecDocConfig		TecDocConfig
}

type PostgresConfig struct {
	Username	string	`envconfig:"POSTGRES_USERNAME"`
	Password	string	`envconfig:"POSTGRES_PASSWORD"`
	Host		string	`envconfig:"POSTGRES_HOST"`
	Port		string	`envconfig:"POSTGRES_PORT"`
	DbName		string	`envconfig:"POSTGRES_DB"`
}

type ContentClientConfig struct {

}

type TecDocConfig struct {

} 
