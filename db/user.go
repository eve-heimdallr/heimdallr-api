package db

import (
	"database/sql"
	"time"

	"github.com/eve-heimdallr/heimdallr-api/common"
)

// GetUser retrieves a user based on their character ID
func GetUser(tx *sql.Tx, characterID int) (*common.User, error) {
	user := common.User{}
	row := tx.QueryRow(`
		SELECT character_id, character_name, created_on, admin
		  FROM h_user
		 WHERE character_id=$1;
	`, characterID)
	if err := row.Scan(&user.CharacterID, &user.CharacterName, &user.CreatedOn, &user.Admin); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

// PopulateUser creates a non-admin user with the given name, if a user with the ID does not exist
func PopulateUser(tx *sql.Tx, characterID int, characterName string) error {
	user, err := GetUser(tx, characterID)
	if err != nil {
		return err
	}
	if user != nil {
		return nil
	}

	_, err = tx.Exec(`
		INSERT INTO h_user(character_id, character_name, created_on, admin) VALUES ($1, $2, $3, false);
	`, characterID, characterName, time.Now())
	return err
}
