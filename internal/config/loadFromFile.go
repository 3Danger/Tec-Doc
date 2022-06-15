package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func LoadFromFile(filepath string) (*Config, error) {
	var (
		err  error
		data []byte
		file *os.File
		conf = new(Config)
	)

	if file, err = os.Open(filepath); err != nil {
		return nil, err
	}
	defer file.Close()
	if data, err = ioutil.ReadAll(file); err != nil {
		return nil, err
	} else if err = json.Unmarshal(data, conf); err != nil {
		return nil, err
	}
	return conf, nil
}
