package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "heimdallr-api"
	app.Usage = "main executable for running server-side Heimdallr processes"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{migrateCommand, serveCommand}
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "debug", Usage: "enable debug-level logging"},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
