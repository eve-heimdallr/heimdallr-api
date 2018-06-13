package db

import (
	"database/sql"
	"time"

	"github.com/eve-heimdallr/heimdallr-api/common"
)

// SaveOAuthSession remembers one in-prograss authentication
func SaveOAuthSession(tx *sql.Tx, session common.OAuthSession) error {
	_, err := tx.Exec(`
		INSERT INTO oauth_session(state, fingerprint, timestamp) VALUES ($1, $2, $3)
	`, session.State, session.Fingerprint, session.Timestamp)
	return err
}

// ValidateOAuthSession checks whether the database contains an auth session with the given state and fingerprint
func ValidateOAuthSession(tx *sql.Tx, state string, fingerprint string) (bool, error) {
	row := tx.QueryRow(`
		SELECT 1
		  FROM oauth_session
		 WHERE state=$1 AND fingerprint=$2 AND timestamp >= (now() - interval '10 minute')`,
		state, fingerprint)

	var found int
	if err := row.Scan(&found); err == nil {
		return true, nil
	} else if err == sql.ErrNoRows {
		return false, nil
	} else {
		return false, err
	}
}

// InsertOAuthAccessToken inserts a new bearer token row and returns the row's ID
func InsertOAuthAccessToken(tx *sql.Tx, accessToken string, refreshToken string, expiration time.Time) (int, error) {
	row := tx.QueryRow(`
		INSERT INTO oauth_bearer_token(access_token, refresh_token, expiration) VALUES ($1, $2, $3) RETURNING id
	`, accessToken, refreshToken, expiration)

	rowID := 0
	if err := row.Scan(&rowID); err != nil {
		return 0, err
	}

	return rowID, nil
}

// InsertUISession creates a new database entry for the given UI session
func InsertUISession(tx *sql.Tx, session common.UISession) error {
	_, err := tx.Exec(`
		INSERT INTO user_ui_session(id, character_id, token_id, expiration) VALUES ($1, $2, $3, $4)
	`, session.ID, session.CharacterID, session.OAuthTokenID, session.Expiration)
	return err
}
