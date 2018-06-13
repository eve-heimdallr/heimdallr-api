package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type oauthVerifyResponseBody struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// OAuthVerifyAuthCode accomplishes the "verify" step in OAuth2 "code grant" flow
func OAuthVerifyAuthCode(authCode string, clientID string, clientSecret string, ssoTokenURL *url.URL) (*OAuthAccessToken, error) {
	body := url.Values{
		"grant_type": []string{"authorization_code"},
		"code":       []string{authCode},
	}
	request, _ := http.NewRequest("POST", ssoTokenURL.String(), strings.NewReader(body.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(clientID, clientSecret)

	LogDebug().Printf("Verifying auth code... url=%s body=%s", ssoTokenURL.String(), body.Encode())

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		LogError().Panic(err)
		return nil, err
	}

	responseBytes, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK response code verifying auth code; status=%d body=`%s`", response.StatusCode, string(responseBytes))
	}

	var responseBody oauthVerifyResponseBody
	if err = json.Unmarshal(responseBytes, &responseBody); err != nil {
		return nil, errors.Wrap(err, "error parsing verify response")
	}

	token := OAuthAccessToken{
		AccessToken:  responseBody.AccessToken,
		RefreshToken: responseBody.RefreshToken,
		Expiration:   time.Now().Add(time.Duration(responseBody.ExpiresIn) * time.Second),
	}

	return &token, nil
}

// OAuthFingerprint creates a client fingerprint based on an HTTP request, for redundant security
func OAuthFingerprint(r *http.Request) string {
	return fmt.Sprintf("%s++%s", r.RemoteAddr, r.UserAgent())
}
