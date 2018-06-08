package migrationscripts

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up001, Down001)
}

// Up001 applies this migration
func Up001(tx *sql.Tx) error {
	_, err := tx.Exec(`
    CREATE TABLE users(
      character_id INT NOT NULL,
      character_name VARCHAR(128) NOT NULL,
      created_on TIMESTAMP,
      admin BOOLEAN,
      PRIMARY KEY (character_id)
    );
  `)
	return err
}

// Down001 downgrades this migration
func Down001(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE users;")
	return err
}
