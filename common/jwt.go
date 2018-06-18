package common

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// HeimdallrClaims contains the standard JWT properties for Heimdallr
type HeimdallrClaims struct {
	CharacterID   int    `json:"cid"`
	CharacterName string `json:"cn"`
	SessionID     string `json:"sid"`
	jwt.StandardClaims
}

// JWT gives a frontend-useable version of the session
func (s UISession) JWT() *jwt.Token {
	claims := HeimdallrClaims{
		CharacterID:   s.CharacterID,
		CharacterName: s.CharacterName,
		SessionID:     s.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: s.Expiration.Unix(),
			Issuer:    "heimdallr@" + HeimdallrVersion,
			IssuedAt:  time.Now().Unix(),
		},
	}

	logDebug.Print("Creating claims", claims)
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}
