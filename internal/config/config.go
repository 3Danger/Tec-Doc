package config

type Config struct {
	Server struct {
		Port	string	`yaml:"port" envconfig:"SERVER_PORT"`
		Host	string	`yaml:"host" envconfig:"SERVER_HOST"`
	} `yaml:"server"`
	Database struct {
		Username	string	`yaml:"user" envconfig:"DB_USERNAME"`
		Password	string	`yaml:"pass" envconfig:"DB_PASSWORD"`
	} `yaml:"database"`
}
