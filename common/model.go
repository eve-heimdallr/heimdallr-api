package common

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

// OAuthSession holds data representing an in-progress authentiation
type OAuthSession struct {
	State       string
	Fingerprint string
	Timestamp   time.Time
}

// User holds data pertaining to a remembered user of the application
type User struct {
	ID            int
	CharacterID   int
	CharacterName string
	CreatedOn     time.Time
	Admin         bool
}

// UISession is a fast-expiring, non-renewing session used for UI interactions
type UISession struct {
	ID            string
	CharacterID   int
	CharacterName string
	OAuthTokenID  int
	Expiration    time.Time
}

// NewUISession creates a new UI Session for a given character with a token ID
func NewUISession(characterID int, characterName string, oauthTokenID int) UISession {
	sessionIDBytes := make([]byte, 64)
	rand.Read(sessionIDBytes)
	sessionID := base64.StdEncoding.EncodeToString(sessionIDBytes)

	expiration := time.Now().Add(uiSessionLifetime)

	return UISession{
		ID:            sessionID,
		CharacterID:   characterID,
		CharacterName: characterName,
		OAuthTokenID:  oauthTokenID,
		Expiration:    expiration,
	}
}

// BackgroundSession is a long-living, renewing session used for polling for updates
type BackgroundSession struct {
	CharacterID  int
	OAuthTokenID int
}

// OAuthAccessToken is a saved Oauth token from a successful authentication
type OAuthAccessToken struct {
	ID           int
	AccessToken  string
	RefreshToken string
	Expiration   time.Time
}
