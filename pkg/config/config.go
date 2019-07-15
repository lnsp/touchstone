package config

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/lnsp/touchstone/pkg/benchmark"
	"github.com/lnsp/touchstone/pkg/benchmark/suites"
	"github.com/lnsp/touchstone/pkg/util"
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

func (c *Config) Matrix() (*benchmark.Matrix, error) {
	b := suites.All
	for _, f := range c.Filter {
		b = benchmark.Filter(b, f)
	}
	m := &benchmark.Matrix{
		OCIs:  c.OCIs,
		CRIs:  c.CRIs,
		Items: b,
		Runs:  c.Runs,
	}
	return m, nil
}

func (c *Config) MapOutput(dir string) (io.WriteCloser, error) {
	if c.Output == "" {
		return os.Stdout, nil
	}
	return util.GetOutputTarget(dir + c.Output), nil
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
	if config.Runs < 1 {
		return nil, errors.New("runs must be larger than 0")
	}
	return config, nil
}
