package data

import (
	"log"

	"github.com/spf13/viper"
	"github.com/synkube/app/core/common"
)

type ServerConfig struct {
	Type string `yaml:"type"`
	Port int    `yaml:"port"`
}

type DbConfig struct {
	Type       string           `yaml:"type"`
	Clean      bool             `yaml:"clean"`
	Postgres   PostgresConfig   `yaml:"postgres,omitempty"`
	SQLite     SQLiteConfig     `yaml:"sqlite,omitempty"`
	MySQL      MySQLConfig      `yaml:"mysql,omitempty"`
	ClickHouse ClickhouseConfig `yaml:"clickhouse,omitempty"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type SQLiteConfig struct {
	File string `yaml:"file"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type ClickhouseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type RPC struct {
	URL  string `yaml:"url"`
	Type string `yaml:"type"`
}

type Chain struct {
	ID      int    `yaml:"id"`
	Name    string `yaml:"name"`
	Network string `yaml:"network"`
	RPCs    []RPC  `yaml:"rpcs"`
}

func LoadConfig(cfgFile string, cfg interface{}) error {
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
