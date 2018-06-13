package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/eve-heimdallr/heimdallr-api/common"
	"github.com/eve-heimdallr/heimdallr-api/db"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Server is an encapsulation of everything it takes to run a Heimdallr API server
type Server struct {
	db                   *sql.DB
	config               *common.Config
	uptimeHandler        http.Handler
	oauthStartHandler    http.Handler
	oauthCallbackHandler http.Handler
}

// NewServer initializes a new Server
func NewServer(config *common.Config) (*Server, error) {
	var err error

	if config.ServePort == 0 {
		return nil, errors.New("serve port not specified")
	}

	server := &Server{config: config}
	if server.db, err = db.GetDatabase(config); err != nil {
		return nil, errors.Wrap(err, "could not initialize DB")
	}

	server.uptimeHandler = newUptimeHandler()
	server.oauthStartHandler, err = newOAuthStartHandler(config, server.db)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize OAuth start handler")
	}
	server.oauthCallbackHandler, err = newOAuthCallbackHandler(config, server.db)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize OAuth callback handler")
	}
	return server, nil
}

// ServeAsync starts the Heimdallr API server asynchronously and returns a channel feeding errors back
func (s Server) ServeAsync() <-chan error {
	errChan := make(chan error)
	go func() {
		addrString := fmt.Sprintf(":%d", s.config.ServePort)
		router := mux.NewRouter()

		router.Handle("/", s.uptimeHandler)
		router.Handle("/oauth/start", s.oauthStartHandler)
		router.Handle("/oauth/callback", s.oauthCallbackHandler)

		common.LogInfo().Printf("Starting Heimdallr server on %s", addrString)
		errChan <- http.ListenAndServe(addrString, router)
	}()
	return errChan
}
