package config

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

func InitConfig(cfgFile string, cfg *Config) error {
	return data.LoadConfig(cfgFile, cfg)
}
