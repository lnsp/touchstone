package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	OCIs   []string `yaml:"oci"`
	CRIs   []string `yaml:"cri"`
	Output string   `yaml:"output"`
	Filter []string `yaml:"filter"`
	Runs   int      `yaml:"runs"`
	Scale  int      `yaml:"scale"`
}

func Parse(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}
