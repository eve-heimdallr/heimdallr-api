package esi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// CurrentUserDetails contains the current user's character ID and name
type CurrentUserDetails struct {
	CharacterID   int
	CharacterName string
}

type esiSSOVerifyResponse struct {
	// http://eveonline-third-party-documentation.readthedocs.io/en/latest/sso/obtaincharacterid.html
	CharacterID        int    `json:"CharacterID"`
	CharacterName      string `json:"CharacterName"`
	ExpiresOnStr       string `json:"ExpiresOn"` // 2014-05-23T15:01:15.182864Z
	Scopes             string `json:"scopes"`
	TokenType          string `json:"TokenType"`
	CharacterOwnerHash string `json:"CharacterOwnerHash"`
}

// GetCurrentUserDetails recovers the user details for the authenticated user
func (ctx Context) GetCurrentUserDetails() (*CurrentUserDetails, error) {
	queryPath, _ := url.Parse("/verify/")
	request, _ := http.NewRequest("GET", ctx.baseURL.ResolveReference(queryPath).String(), http.NoBody)
	ctx.applyHeaders(request)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "error executing current user details request")
	}

	bodyData, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK status code response to current user details request; code=%d body=%s", response.StatusCode, string(bodyData))
	}

	var body esiSSOVerifyResponse
	json.Unmarshal(bodyData, &body)

	return &CurrentUserDetails{CharacterID: body.CharacterID, CharacterName: body.CharacterName}, nil
}
