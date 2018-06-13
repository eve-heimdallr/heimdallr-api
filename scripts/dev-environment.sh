#!/bin/bash

# Assumed UI running at port 8888, API at port 8080

export PORT=8080
export DATABASE_URL=postgres://heimdallr:heimdallr-secret@localhost:5432/heimdallr?sslmode=disable
export ESI_BASE_URL=https://esi.evetech.net/v1/
export SSO_AUTH_URL=https://login.eveonline.com/oauth/authorize
export SSO_TOKEN_URL=https://login.eveonline.com/oauth/token
export SSO_RETURN_BASE_URL=http://localhost:8888
export SSO_REDIRECT_URL=http://localhost:8080/oauth/callback
export SSO_SCOPES=esi-mail.read_mail.v1
export SSO_CLIENT_ID=XXX-PLACEHOLDER-XXX
export SSO_SECRET_KEY=XXX-PLACEHOLDER-XXX
export JWT_SECRET=still-not-very-secret
