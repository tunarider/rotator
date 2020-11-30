package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Action int

const (
	ActionArchive Action = 0
)

type Target struct {
	Name       string
	Path       string
	Regexp     string
	DateFormat string `yaml:"date_format"`
	Retention  int
	Group      bool
	GroupDate  bool `yaml:"group_date"`
	Action     Action
}

type Config struct {
	Targets []Target
}

func ParseConfig(configPath string) (*Config, error) {
	c := &Config{}
	if configPath == "" {
		return c, errors.New("config not found")
	}
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return c, errors.New("failed to read config file")
	}
	err = yaml.Unmarshal(configFile, c)
	return c, err
}
