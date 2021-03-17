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

type Global struct {
	DateFormat *string `yaml:"date_format"`
	Retention  *int
	GroupBy    *GroupBy `yaml:"group_by"`
	Remove     *bool
}

type Target struct {
	Name       string
	Path       string
	Regexp     string
	Action     Action
	DateFormat *string `yaml:"date_format"`
	Retention  *int
	GroupBy    *GroupBy `yaml:"group_by"`
	Remove     *bool
}

type Config struct {
	Global  Global
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
	for i, _ := range c.Targets {
		if c.Targets[i].DateFormat == nil {
			c.Targets[i].DateFormat = c.Global.DateFormat
		}
		if c.Targets[i].Retention == nil {
			c.Targets[i].Retention = c.Global.Retention
		}
		if c.Targets[i].GroupBy == nil {
			c.Targets[i].GroupBy = c.Global.GroupBy
		}
		if c.Targets[i].Remove == nil {
			c.Targets[i].Remove = c.Global.Remove
		}
	}
	return c, err
}
