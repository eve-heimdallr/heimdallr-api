package main

import (
	"errors"

	"github.com/eve-heimdallr/heimdallr-api/common"
	"github.com/eve-heimdallr/heimdallr-api/db"
	"github.com/eve-heimdallr/heimdallr-api/migration"
	"github.com/urfave/cli"
)

var migrateCommand = cli.Command{
	Name:    "migrate",
	Aliases: []string{"m"},
	Flags: []cli.Flag{
		cli.Int64Flag{Name: "up-to", Usage: "migrate up to this version"},
		cli.Int64Flag{Name: "down-to", Usage: "migrate down to this version"},
	},
	Usage: "migrate to the most recent version, or to the version specified by --up-to or --down-to",
	Action: func(ctx *cli.Context) error {
		var err error
		common.SetDebugEnabled(ctx.GlobalBool("debug"))

		upTo := ctx.Int64("up-to")
		downTo := ctx.Int64("down-to")

		if upTo > 0 && downTo > 0 {
			return errors.New("may not specify both --up-to and --down-to")
		}

		config, err := common.GetConfigFromEnvironment()
		if err != nil {
			return err
		}
		db, err := db.GetDatabase(config)
		if err != nil {
			return err
		}

		if upTo > 0 {
			common.LogInfo().Printf("migrating up to version %d", upTo)
			err = migration.UpTo(db, upTo)
		} else if downTo > 0 {
			common.LogInfo().Printf("migrating down to version %d", downTo)
			err = migration.DownTo(db, downTo)
		} else {
			common.LogInfo().Print("migrating to latest version")
			err = migration.Up(db)
		}

		if err == nil {
			common.LogInfo().Print("migration successful")
		} else {
			common.LogError().Printf("error during migration: %v", err)
		}

		return nil
	},
}
