package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "server",
		Usage: "cms and chatbot server",
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "start cms and chatbot service",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "port",
						Aliases:  []string{"p"},
						Usage:    "Listen to port",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "config",
						Aliases:  []string{"c"},
						Usage:    "Load configuration file",
						Required: true,
					},
					&cli.BoolFlag{
						Name:    "debug",
						Aliases: []string{"d"},
						Usage:   "Debug mode",
					},
				},
				Action: func(ctx *cli.Context) error {
					server := ctx.String("port")
					config := ctx.String("config")
					debug := ctx.Bool("debug")

					return run_server(server, config, debug)
				},
			},
			{
				Name:  "webhook",
				Usage: "start webhook",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "port",
						Aliases:  []string{"p"},
						Usage:    "Listen to port",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "config",
						Aliases:  []string{"c"},
						Usage:    "Load configuration file",
						Required: true,
					},
					&cli.BoolFlag{
						Name:    "debug",
						Aliases: []string{"d"},
						Usage:   "Debug mode",
					},
				},
				Action: func(ctx *cli.Context) error {
					server := ctx.String("port")
					config := ctx.String("config")
					debug := ctx.Bool("debug")

					return run_webhook(server, config, debug)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
