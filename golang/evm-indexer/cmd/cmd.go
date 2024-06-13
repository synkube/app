package cmd

import (
	"log"
	"os"

	"github.com/synkube/app/core/common"
	coreData "github.com/synkube/app/core/data"
	"github.com/synkube/app/evm-indexer/config"
	"github.com/synkube/app/evm-indexer/data"
	"github.com/synkube/app/evm-indexer/indexer"
	"github.com/urfave/cli/v2"
)

var cfg config.Config
var ds *coreData.DataStore

func Start(args []string, buildInfo common.BuildInfo) error {
	app := &cli.App{
		Name:    buildInfo.Name(),
		Version: buildInfo.Version(),
		Usage:   buildInfo.Description(),
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
					log.Println("Usage info for this application:")
					log.Println("- Use the '--config' flag to specify a configuration file.")
					log.Println("- Use the 'info' command to get information about usage.")
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
	// ctx, cancel := context.WithCancel(context.Background())
	if err := config.InitConfig(c.String("config"), &cfg); err != nil {
		return err
	}
	log.Println("Running the application with arguments:", c.Args().Slice())

	ds = data.Initialize(&cfg)
	if ds != nil {
		ds.CheckConnection()
	}
	go StartServers(cfg.ServerConfig)
	log.Println("Starting the indexer")
	indexer.StartIndexing(cfg.Chain, ds, cfg.Indexer)
	return nil
}
