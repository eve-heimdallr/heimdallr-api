package esi

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/eve-heimdallr/heimdallr-api/common"
)

// Context is a wrapper for configuration related to an ESI "session"
type Context struct {
	baseURL     *url.URL
	accessToken string
}

// NewContext creates a new ESI Context based on configuration and a given access token
func NewContext(config *common.Config, accessToken string) (*Context, error) {
	context := &Context{
		baseURL:     config.ESIBaseURL,
		accessToken: accessToken,
	}

	if context.baseURL == nil || context.baseURL.String() == "" {
		return nil, errors.New("ESI base URL not specified in configuration")
	}

	return context, nil
}

// WithAccessToken creates a copy of this context, but with a different access token
func (ctx Context) WithAccessToken(accessToken string) Context {
	return Context{
		baseURL:     ctx.baseURL,
		accessToken: accessToken,
	}
}

func (ctx Context) applyHeaders(request *http.Request) {
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctx.accessToken))
	request.Header.Set("User-Agent", common.HeimdallrUserAgent)
}
