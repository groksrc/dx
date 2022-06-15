package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var app = cli.NewApp()

func info() {
	app.Name = "dx"
	app.Usage = "Generate a CLI for your dev team"
	app.Author = "groksrc"
	app.Version = "0.0.1"
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name: "init",
			Aliases: []string{"i"},
			Usage: "Initialize your project",
			Action: func(c *cli.Context) error {
				fmt.Printf("Inited %q", c.Args().Get(0))
				return nil
			},
		},
		{
			Name: "cli",
			Aliases: []string{"c"},
			Usage: "Create your company CLI",
			Subcommands: []cli.Command{
				{
					Name: "create",
					Category: "cli",
					Action: func(c *cli.Context) error {
						fmt.Printf("not implemented. Open a PR %q", c.Args().Get(0))
						return nil
					},
				},
				{
					Name: "update",
					Category: "cli",
					Action: func(c *cli.Context) error {
						fmt.Printf("not implemented. Open a PR %q", c.Args().Get(0))
						return nil
					},
				},
			},
		},
	}
}

func main() {
	info()
	commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}