package main

import (
	"os"
	"os/signal"

	"github.com/eve-heimdallr/heimdallr-api/common"
	"github.com/eve-heimdallr/heimdallr-api/server"
	"github.com/urfave/cli"
)

var serveCommand = cli.Command{
	Name:    "serve",
	Aliases: []string{"s"},
	Action: func(ctx *cli.Context) error {
		var err error
		common.SetDebugEnabled(ctx.GlobalBool("debug"))

		config, err := common.GetConfigFromEnvironment()
		if err != nil {
			return err
		}
		server, err := server.NewServer(config)
		if err != nil {
			return err
		}

		interruptChan := make(chan os.Signal, 1)
		signal.Notify(interruptChan, os.Interrupt)
		errChan := server.ServeAsync()

		select {
		case <-interruptChan:
			common.LogInfo().Print("received keyboard interrupt, ending execution")
			// TODO: do cleanup here?
			return nil
		case err = <-errChan:
			common.LogError().Print("fatal error running HTTP server: " + err.Error())
			return err
		}
	},
}
