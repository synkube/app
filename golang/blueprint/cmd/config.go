package cmd

import (
	"github.com/synkube/app/core/data"
)

// Config represents the structure of the configuration file
type Config struct {
	AppName      string              `yaml:"appName"`
	Version      string              `yaml:"version"`
	ServerConfig []data.ServerConfig `yaml:"serverConfig"`
	DbConfig     data.DbConfig       `yaml:"dbConfig"`
}
