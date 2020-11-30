package main

import (
	"github.com/tunarider/rotator/internal/cmd"
	"github.com/tunarider/rotator/internal/config"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:    "Rotator",
		Usage:   "rotate files",
		Version: "v0.1.2",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "config",
				Aliases: []string{"c"},
			},
			&cli.BoolFlag{
				Name:    "dry",
				Aliases: []string{"d"},
			},
		},
		Action: func(c *cli.Context) error {
			conf, err := config.ParseConfig(c.Path("config"))
			if err != nil {
				return cli.Exit(err, 1)
			}
			return cmd.Rotate(conf, c.Bool("dry"))
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
