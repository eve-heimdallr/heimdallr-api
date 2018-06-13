package common

import (
	"net/url"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

// Config stores the environment configuration for the current execution
type Config struct {
	ServePort           int
	DatabaseURL         *url.URL
	ESIBaseURL          *url.URL
	SSOAuthorizationURL *url.URL
	SSOTokenURL         *url.URL
	SSORedirectURL      *url.URL
	SSOReturnBaseURL    *url.URL
	SSOScopes           string
	SSOClientID         string
	SSOSecretKey        string
}

// GetConfigFromEnvironment builds a config container based on the environment
func GetConfigFromEnvironment() (*Config, error) {
	var err error
	config := &Config{}

	if config.ServePort, err = strconv.Atoi(os.Getenv("PORT")); err != nil {
		return nil, errors.Wrap(err, "failed fetching PORT from environment")
	}

	if config.DatabaseURL, err = url.Parse(os.Getenv("DATABASE_URL")); err != nil {
		return nil, errors.Wrap(err, "failed fetching DATABASE_URL from environment")
	}

	if config.ESIBaseURL, err = url.Parse(os.Getenv("ESI_BASE_URL")); err != nil {
		return nil, errors.Wrap(err, "failed fetching ESI_BASE_URL from environment")
	}

	if config.SSOAuthorizationURL, err = url.Parse(os.Getenv("SSO_AUTH_URL")); err != nil {
		return nil, errors.Wrap(err, "failed fetching SSO_AUTH_URL from environment")
	}

	if config.SSOTokenURL, err = url.Parse(os.Getenv("SSO_TOKEN_URL")); err != nil {
		return nil, errors.Wrap(err, "failed fetching SSO_TOKEN_URL from environment")
	}

	if config.SSORedirectURL, err = url.Parse(os.Getenv("SSO_REDIRECT_URL")); err != nil {
		return nil, errors.Wrap(err, "failed fetching SSO_REDIRECT_URL from environment")
	}

	if config.SSOReturnBaseURL, err = url.Parse(os.Getenv("SSO_RETURN_BASE_URL")); err != nil {
		return nil, errors.Wrap(err, "failed fetching SSO_RETURN_BASE_URL from environment")
	}

	config.SSOScopes = os.Getenv("SSO_SCOPES")
	config.SSOClientID = os.Getenv("SSO_CLIENT_ID")
	config.SSOSecretKey = os.Getenv("SSO_SECRET_KEY")

	return config, nil
}
