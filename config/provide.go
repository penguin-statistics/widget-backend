package config

import (
	"github.com/jinzhu/configor"
)

// StaticService describes a service that serves static files
type StaticService struct {
	Endpoint string `yaml:"endpoint"`
	Root string `yaml:"root"`
}

// Config describes all configuration options for this app
type ConfigT struct {
	Server struct {
		Listen string `yaml:"listen"`
	} `yaml:"server"`

	Static struct {
		Widget  StaticService `yaml:"widget"`
		Docs    StaticService `yaml:"docs"`
	} `yaml:"static"`

	Upstream struct {
		Meta struct{
			Servers []string `yaml:"servers"`
		} `yaml:"meta"`
	} `yaml:"upstream"`
}

var C ConfigT

// New initializes config files
func init() {
	var config ConfigT
	c := configor.New(&configor.Config{
		// set env prefix. now for e.g. Config.Server.Listen, use WIDGET_BACKEND_SERVER_LISTEN to specify its value
		ENVPrefix:            "WIDGET_BACKEND",
		ErrorOnUnmatchedKeys: false,
	})
	err := c.Load(&config, "config.yml")
	if err != nil {
		panic(err)
	}
	C = config
}