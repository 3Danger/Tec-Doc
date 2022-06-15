package config

import "github.com/kelseyhightower/envconfig"

/*
	export TEC_DOC_LOG_LEVEL=DEBUG
	export TEC_DOC_LOG_TM_FORMAT=MS
*/

func NewConfigEnv() (conf *Config, err error) {
	conf = new(Config)
	err = envconfig.Process("TEC_DOC", conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
