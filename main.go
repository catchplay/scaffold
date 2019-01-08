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
				err := scaffold.New().Generate()
				//fmt.Printf("error:%+v\n", err)
				if err == nil {
					fmt.Println("Success Created. Please excute `make up` to start service.")
				}

				return err
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
