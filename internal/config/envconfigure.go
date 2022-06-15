package config

import "github.com/kelseyhightower/envconfig"

/*
	export TEC_DOC_PORT=4242
	export TEC_DOC_DEBUG_LEVEL=DEBUG
	export TEC_DOC_DEBUG_TM_FORMAT=MS
*/
func NewConfigEnv() (*Config, error) {
	conf := new(Config)
	err1 := envconfig.Process("TEC_DOC", &conf.Server)
	err2 := envconfig.Process("TEC_DOC_DEBUG", &conf.Debug)
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	return conf, nil
}
