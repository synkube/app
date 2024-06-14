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
	Indexer      Indexer             `yaml:"indexer"`
	Chain        data.Chain          `yaml:"chain"`
}

type Indexer struct {
	StartBlock    int `yaml:"startBlock"`
	EndBlock      int `yaml:"endBlock"`
	BatchSize     int `yaml:"batchSize"`
	MaxWorkers    int `yaml:"maxWorkers"`
	MaxRetries    int `yaml:"maxRetries"`
	RetryInterval int `yaml:"retryInterval"`
	RetryBackoff  int `yaml:"retryBackoff"`
}

func InitConfig(cfgFile string, cfg *Config) error {
	return data.LoadConfig(cfgFile, &cfg)
}
