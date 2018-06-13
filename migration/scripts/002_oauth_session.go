package migrationscripts

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up002, Down002)
}

// Up002 applies this migration
func Up002(tx *sql.Tx) error {
	_, err := tx.Exec(`
    CREATE TABLE oauth_session(
      state VARCHAR(128) NOT NULL,
      fingerprint TEXT NOT NULL,
      timestamp TIMESTAMP NOT NULL,
      PRIMARY KEY (state)
    );
  `)
	return err
}

// Down002 downgrades this migration
func Down002(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE oauth_session;")
	return err
}
