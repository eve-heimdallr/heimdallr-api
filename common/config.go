package common

import (
	"net/url"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

// Config stores the environment configuration for the current execution
type Config struct {
	ServePort   int
	DatabaseURL *url.URL
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

	return config, nil
}
