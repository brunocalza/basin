package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DBS struct {
		Path    string `yaml:"path"`
		Replica struct {
			URL       string `yaml:"url"`
			AccessKey string `yaml:"access_key"`
			Frequency string `yaml:"frequency"`
		} `yaml:"replica"`
	} `yaml:"dbs"`
}

func setupConfig(path string) (*Config, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return &Config{}, err
	}

	conf := &Config{}
	if err := yaml.Unmarshal(buf, conf); err != nil {
		return &Config{}, err
	}

	return conf, nil
}
