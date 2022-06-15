package config

import "github.com/kelseyhightower/envconfig"

/*
	export TEC_DOC_LOG_LEVEL=DEBUG
	export TEC_DOC_LOG_TM_FORMAT=MS
	export TEC_DOC_SERVER_PORT=4242
	export TEC_DOC_DATABASE_USERNAME=NAN
	export TEC_DOC_DATABASE_PASSWORD=NAN
	export TEC_DOC_DATABASE_NAME=NAN
	export TEC_DOC_DATABASE_PORT=5432
*/

func NewConfigEnv() (conf *Config, err error) {
	conf = new(Config)
	err = envconfig.Process("TEC_DOC", conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
