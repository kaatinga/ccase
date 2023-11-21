package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"

	"github.com/kaatinga/ccase/internal/lsc"
	"github.com/kaatinga/ccase/internal/settings"

	"github.com/urfave/cli/v2"
)

var version = "unknown"

func init() {
	file, err := os.Open("VERSION")
	if err != nil {
		return
	}

	var data []byte
	data, err = io.ReadAll(file)
	if err != nil {
		return
	}

	if len(data) > 0 {
		version = string(bytes.TrimSpace(data))
	}
}

func main() {
	app := &cli.App{
		Name:           "A change case CLI tool",
		Description:    "Ccase helps to change the name of files.",
		Version:        version,
		Compiled:       time.Now(),
		DefaultCommand: "lsc",
		Authors: []*cli.Author{
			{
				Name: "Michael Gunkoff",
			},
		},
		HelpName: "Ccase",
		Usage:    "choose the case you need to change to",
		Commands: []*cli.Command{
			{
				Name:   "lower_snake_case",
				Action: lsc.UpdateFiles,
				Aliases: []string{
					"lsc",
				},
			},
			// {
			// 	Name:   "camelCase",
			// 	Action: ,
			// },
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Usage: "provide a valid path to work with git repository",
				Action: func(context *cli.Context, s string) error {
					if s != "" {
						settings.Path = s
					}

					return settings.DefinePaths()
				},
			},
			// &cli.StringFlag{
			// 	Name:  "mask",
			// 	Usage: "provide a valid mask to work with files",
			// 	Action: func(context *cli.Context, s string) error {
			// 		if s != "" {
			// 			settings.Mask = s
			// 		}
			//
			// 		return nil
			// 	},
			// },
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
