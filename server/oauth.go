package server

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/http"
	"net/url"
	"time"

	"github.com/eve-heimdallr/heimdallr-api/common"
	"github.com/eve-heimdallr/heimdallr-api/db"
	"github.com/eve-heimdallr/heimdallr-api/esi"
	"github.com/pkg/errors"
)

type oAuthStartHandler struct {
	db             *sql.DB
	ssoURL         *url.URL
	ssoCallbackURL *url.URL
	ssoScopes      string
	ssoClientID    string
	ssoSecretKey   string
}

func newOAuthStartHandler(config *common.Config, db *sql.DB) (*oAuthStartHandler, error) {
	handler := &oAuthStartHandler{
		db:             db,
		ssoURL:         config.SSOAuthorizationURL,
		ssoCallbackURL: config.SSORedirectURL,
		ssoScopes:      config.SSOScopes,
		ssoClientID:    config.SSOClientID,
		ssoSecretKey:   config.SSOSecretKey,
	}

	if handler.ssoURL == nil || handler.ssoURL.String() == "" {
		return nil, errors.New("SSO authorization URL not specified in configuration")
	}

	if handler.ssoClientID == "" {
		return nil, errors.New("SSO client ID not specified in configuration")
	}

	if handler.ssoSecretKey == "" {
		return nil, errors.New("SSO secret key not specified in configuration")
	}

	return handler, nil
}

func (h oAuthStartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tx, err := h.db.Begin()
	defer tx.Rollback()
	if err != nil {
		common.WriteHTTPError(w, http.StatusInternalServerError, errors.Wrap(err, "could not open db transaction"))
		return
	}

	stateBytes := make([]byte, 16)
	rand.Read(stateBytes)
	state := base64.RawURLEncoding.EncodeToString(stateBytes)
	fingerprint := common.OAuthFingerprint(r)

	session := common.OAuthSession{State: state, Fingerprint: fingerprint, Timestamp: time.Now()}
	common.LogDebug().Printf("Saving auth session %v", session)
	err = db.SaveOAuthSession(tx, session)
	if err != nil {
		common.WriteHTTPError(w, http.StatusInternalServerError, errors.Wrap(err, "could not save oauth session"))
		return
	}

	err = tx.Commit()
	if err != nil {
		common.WriteHTTPError(w, http.StatusInternalServerError, errors.Wrap(err, "could not commit db changes"))
		return
	}

	ssoRedirectURL, _ := url.Parse(h.ssoURL.String())
	query := ssoRedirectURL.Query()
	query.Set("response_type", "code")
	query.Set("redirect_uri", h.ssoCallbackURL.String())
	query.Set("client_id", h.ssoClientID)
	query.Set("scope", h.ssoScopes)
	query.Set("state", state)
	ssoRedirectURL.RawQuery = query.Encode()

	common.LogDebug().Printf("Redirecting auth to %s", ssoRedirectURL.String())
	w.Header().Set("Location", ssoRedirectURL.String())
	w.WriteHeader(http.StatusFound)
	w.Write([]byte("Redirecting to " + ssoRedirectURL.String()))
}

type oAuthCallbackHandler struct {
	db               *sql.DB
	esi              *esi.Context
	ssoTokenURL      *url.URL
	ssoReturnBaseURL *url.URL
	ssoClientID      string
	ssoSecretKey     string
}

func newOAuthCallbackHandler(config *common.Config, db *sql.DB) (*oAuthCallbackHandler, error) {
	handler := &oAuthCallbackHandler{
		db:               db,
		ssoTokenURL:      config.SSOTokenURL,
		ssoReturnBaseURL: config.SSOReturnBaseURL,
		ssoClientID:      config.SSOClientID,
		ssoSecretKey:     config.SSOSecretKey,
	}

	if handler.ssoTokenURL == nil || handler.ssoTokenURL.String() == "" {
		return nil, errors.New("SSO token URL not specified in configuration")
	}

	if handler.ssoReturnBaseURL == nil || handler.ssoReturnBaseURL.String() == "" {
		return nil, errors.New("SSO return base URL not specified in configuration")
	}

	if handler.ssoClientID == "" {
		return nil, errors.New("SSO client ID not specified in configuration")
	}

	if handler.ssoSecretKey == "" {
		return nil, errors.New("SSO secret key not specified in configuration")
	}

	esiContext, err := esi.NewContext(config, "")
	if err != nil {
		return nil, err
	}
	handler.esi = esiContext

	return handler, nil
}

func (h oAuthCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tx, err := h.db.Begin()
	defer tx.Rollback()
	if err != nil {
		common.WriteHTTPError(w, http.StatusInternalServerError, errors.Wrap(err, "could not open db transaction"))
		return
	}

	authCode := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	fingerprint := common.OAuthFingerprint(r)

	if ok, validationErr := db.ValidateOAuthSession(tx, state, fingerprint); validationErr != nil {
		common.WriteHTTPError(w, http.StatusInternalServerError, errors.Wrap(validationErr, "error validating oauth session"))
		return
	} else if !ok {
		common.WriteHTTPError(w, http.StatusUnauthorized, errors.New("authentication session expired or invalid, please try again"))
		return
	}

	token, err := common.OAuthVerifyAuthCode(authCode, h.ssoClientID, h.ssoSecretKey, h.ssoTokenURL)
	if err != nil {
		common.WriteHTTPError(w, http.StatusUnauthorized, errors.Wrap(err, "failed token verification"))
		return
	}

	tokenID, err := db.InsertOAuthAccessToken(tx, token.AccessToken, token.RefreshToken, token.Expiration)
	if err != nil {
		common.WriteHTTPError(w, http.StatusInternalServerError, errors.Wrap(err, "error saving verified oauth token"))
		return
	}

	esiContext := h.esi.WithAccessToken(token.AccessToken)
	userDetails, err := esiContext.GetCurrentUserDetails()
	if err != nil {
		common.WriteHTTPError(w, http.StatusInternalServerError, errors.Wrap(err, "failed retrieving user details from access token"))
		return
	}

	if err = db.PopulateUser(tx, userDetails.CharacterID, userDetails.CharacterName); err != nil {
		common.WriteHTTPError(w, http.StatusInternalServerError, errors.Wrap(err, "failed populating user table with logged in user"))
		return
	}

	uiSession := common.NewUISession(userDetails.CharacterID, userDetails.CharacterName, int(tokenID))
	db.InsertUISession(tx, uiSession)

	jwt := uiSession.JWT()
	jwtStr, err := jwt.SignedString(common.JWTSecret)
	if err != nil {
		common.WriteHTTPError(w, http.StatusInternalServerError, errors.Wrap(err, "failed composing JWT"))
		return
	}

	returnURL := h.ssoReturnBaseURL // TODO: add custom return path / queryargs / fragment
	returnURLQueryValues := url.Values{"token": []string{jwtStr}}
	returnURL.RawQuery = returnURLQueryValues.Encode()
	common.LogDebug().Printf("Auth successful, redirecting to: %s", returnURL.String())

	w.Header().Set("Location", returnURL.String())
	w.WriteHeader(http.StatusFound)
	w.Write([]byte("Auth successful, redirecting..."))
	tx.Commit()
}
