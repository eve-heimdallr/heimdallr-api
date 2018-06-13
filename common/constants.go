package common

import (
	"os"
	"time"
)

// HeimdallrVersion is the version string for Heimdallr
const HeimdallrVersion = "0.0.1"

// HeimdallrUserAgent is the HTTP User Agent that Heimdallr should use in requests
const HeimdallrUserAgent = "eve-heimdallr/" + HeimdallrVersion

const uiSessionLifetime = 24 * time.Hour

// JWTSecret is the secret used to sign JWTs
var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

func init() {
	if string(JWTSecret) == "" {
		JWTSecret = []byte("notsosecret")
	}
}
