package db

import (
	"database/sql"

	"github.com/eve-heimdallr/heimdallr-api/common"
	_ "github.com/lib/pq" // Import postgresql driver
)

// GetDatabase builds a *sql.DB out of the given configuration
func GetDatabase(config *common.Config) (*sql.DB, error) {
	return sql.Open("postgres", config.DatabaseURL.String())
}
