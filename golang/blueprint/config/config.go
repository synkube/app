package config

import (
	"log"

	"github.com/spf13/viper"
	"github.com/synkube/app/core/common"
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
	if cfgFile != "" {
		log.Printf("Loading config file from: %s\n", cfgFile)
		viper.SetConfigFile(cfgFile)
	} else {
		log.Println("Loading default config file")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %s", err)
		return err
	}

	log.Println("Using config file:", viper.ConfigFileUsed())

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("Unable to decode into struct: %s", err)
		return err
	}
	common.PrettyPrintYAML(cfg)

	return nil
}
