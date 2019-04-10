package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config is the top level app configuration
type Config struct {
	Service *PostgresConfig `yaml:"postgres"`
}

// NewConfig gets the app config from config file
func NewConfig(file string) *Config {
	cfg := &Config{}
	if err := load(cfg, file); err != nil {
		panic(fmt.Sprintf("failed to load config file %s", err.Error()))
	}
	return cfg
}

// PostgresConfig contains info for connecting to a postgress database
type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	DB       string `yaml:"db"`
}

func load(config interface{}, fname string) error {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return err
	}
	return nil
}
