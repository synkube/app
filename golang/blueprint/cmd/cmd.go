package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"github.com/synkube/app/core/common"
	"github.com/synkube/app/core/data"
	"github.com/urfave/cli/v2"
)

var cfg Config

func Start(args []string, buildInfo string) error {
	app := &cli.App{
		Name:  "app",
		Usage: "A blueprint Golang application",
		Action: func(c *cli.Context) error {
			// Default action if no subcommand is provided
			return runApplication(c)
		},
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run the application",
				Action: func(c *cli.Context) error {
					return runApplication(c)
				},
			},
			{
				Name:  "info",
				Usage: "Information about how to use this application",
				Action: func(c *cli.Context) error {
					fmt.Println("Usage info for this application:")
					fmt.Println("- Use the '--config' flag to specify a configuration file.")
					fmt.Println("- Use the 'info' command to get information about usage.")
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Usage: "Load configuration from `FILE`",
			},
		},
	}

	return app.Run(os.Args)
}

func runApplication(c *cli.Context) error {
	if err := initConfig(c.String("config")); err != nil {
		return err
	}
	fmt.Println("Running the application with arguments:", c.Args().Slice())
	common.PrettyPrintYAML(cfg)
	data.InitializeDB(cfg.DbConfig)

	StartServers(cfg.ServerConfig)
	return nil
}

func initConfig(cfgFile string) error {
	if cfgFile != "" {
		fmt.Printf("Loading config file from: %s\n", cfgFile)
		viper.SetConfigFile(cfgFile)
	} else {
		fmt.Println("Loading default config file")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %s", err)
		return err
	}

	fmt.Println("Using config file:", viper.ConfigFileUsed())

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("Unable to decode into struct: %s", err)
		return err
	}

	return nil
}
