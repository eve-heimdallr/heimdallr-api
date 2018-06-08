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
	db            *sql.DB
	config        *common.Config
	uptimeHandler http.Handler
}

// NewServer initializes a new Server
func NewServer(config *common.Config) (*Server, error) {
	var err error

	if config.ServePort == 0 {
		return nil, errors.New("Serve port not specified")
	}

	server := &Server{config: config}
	if server.db, err = db.GetDatabase(config); err != nil {
		return nil, errors.Wrap(err, "Could not initialize DB")
	}

	server.uptimeHandler = newUptimeHandler()

	return server, nil
}

// ServeAsync starts the Heimdallr API server asynchronously and returns a channel feeding errors back
func (s Server) ServeAsync() <-chan error {
	addrString := fmt.Sprintf(":%d", s.config.ServePort)
	router := mux.NewRouter()

	router.Handle("/", s.uptimeHandler)

	errChan := make(chan error)
	go func() {
		common.LogInfo().Printf("Starting Heimdallr server on %s", addrString)
		errChan <- http.ListenAndServe(addrString, router)
	}()
	return errChan
}
