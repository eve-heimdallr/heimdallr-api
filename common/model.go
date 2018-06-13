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

type User struct {
	ID            int
	CharacterID   int
	CharacterName string
	CreatedOn     time.Time
	Admin         bool
}

type UISession struct {
	ID           string
	CharacterID  int
	OAuthTokenID int
	Expiration   time.Time
}

func NewUISession(characterID int, oauthTokenID int) UISession {
	sessionIDBytes := make([]byte, 64)
	rand.Read(sessionIDBytes)
	sessionID := base64.StdEncoding.EncodeToString(sessionIDBytes)

	expiration := time.Now().Add(uiSessionLifetime)

	return UISession{
		ID:           sessionID,
		CharacterID:  characterID,
		OAuthTokenID: oauthTokenID,
		Expiration:   expiration,
	}
}

type BackgroundSession struct {
	CharacterID  int
	OAuthTokenID int
}

type OAuthAccessToken struct {
	ID           int
	AccessToken  string
	RefreshToken string
	Expiration   time.Time
}
