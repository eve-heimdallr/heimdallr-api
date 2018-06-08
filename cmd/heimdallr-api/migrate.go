package main

import (
	"github.com/eve-heimdallr/heimdallr-api/common"
	"github.com/urfave/cli"
)

var migrateCommand = cli.Command{
	Name:    "migrate",
	Aliases: []string{"m"},
	Action: func(ctx *cli.Context) error {
		common.SetDebugEnabled(ctx.GlobalBool("debug"))
		common.LogDebug().Print("Enabled debug logging")
		common.LogInfo().Print("Running database migrations")
		return nil
	},
}
