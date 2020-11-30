package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Action string

const (
	ActionArchive Action = "archive"
)

type GroupBy string

const (
	GroupByDate = "date"
	GroupByName = "name"
	GroupByFile = "file"
)

type Target struct {
	Name       string
	Path       string
	Regexp     string
	DateFormat string `yaml:"date_format"`
	Retention  int
	GroupBy    GroupBy `yaml:"group_by"`
	Action     Action
	Remove     bool
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
