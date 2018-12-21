package main

import (
	"fmt"
	"log"
	"os"

	"github.com/catchplay/scaffold/scaffold"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0-rc"
	app.Usage = "Generate scaffold project layout for Go."
	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   " Generate scaffold project layout",
			Action: func(c *cli.Context) error {
				scfd := scaffold.New()

				err := scfd.Generate()
				if err != nil {
					fmt.Printf("%+v\n", err)
				}

				return scfd.Generate()
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
