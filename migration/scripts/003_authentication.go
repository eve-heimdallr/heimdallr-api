package migrationscripts

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up003, Down003)
}

// Up003 applies this migration
func Up003(tx *sql.Tx) error {
	_, err := tx.Exec(`
    CREATE TABLE oauth_bearer_token(
			id SERIAL NOT NULL,
			access_token TEXT NOT NULL,
			refresh_token TEXT NOT NULL,
			expiration TIMESTAMP NOT NULL,
			PRIMARY KEY (id)
		);

		CREATE TABLE user_ui_session (
			id VARCHAR(128) NOT NULL,
			character_id INT NOT NULL REFERENCES h_user(character_id),
			token_id INT NOT NULL REFERENCES oauth_bearer_token(id),
			expiration TIMESTAMP NOT NULL,
			PRIMARY KEY (id)
		);

		CREATE TABLE user_background_session (
			character_id INT NOT NULL REFERENCES h_user(character_id),
			token_id INT NOT NULL REFERENCES oauth_bearer_token(id)
		);
  `)
	return err
}

// Down003 downgrades this migration
func Down003(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE oauth_bearer_token, user_ui_session, user_background_session;")
	return err
}
