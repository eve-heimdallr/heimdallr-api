package migration

import (
	"database/sql"

	"github.com/eve-heimdallr/heimdallr-api/common"
	_ "github.com/eve-heimdallr/heimdallr-api/migration/scripts" // Apply Go migrations
	"github.com/pressly/goose"
)

// Up is a wrapper around goose.Up
func Up(db *sql.DB) error {
	goose.SetLogger(common.LogInfo())
	return goose.Up(db, ".")
}

// UpTo is a wrapper around goose.UpTo
func UpTo(db *sql.DB, version int64) error {
	goose.SetLogger(common.LogInfo())
	return goose.UpTo(db, ".", version)
}

// DownTo is a wrapper around goose.DownTo
func DownTo(db *sql.DB, version int64) error {
	goose.SetLogger(common.LogInfo())
	return goose.DownTo(db, ".", version)
}
