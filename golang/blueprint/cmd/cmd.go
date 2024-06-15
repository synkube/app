package cmd

import (
	"log"
	"os"

	"github.com/synkube/app/blueprint/config"
	"github.com/synkube/app/blueprint/data"
	"github.com/synkube/app/core/common"
	coreData "github.com/synkube/app/core/data"
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
	if err := config.InitConfig(c.String("config"), &cfg); err != nil {
		return err
	}
	log.Println("Running the application with arguments:", c.Args().Slice())

	ds = data.Initialize(&cfg)
	dm := data.NewDataModel(ds)
	StartServers(cfg.ServerConfig, dm)
	return nil
}
