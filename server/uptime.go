package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eve-heimdallr/heimdallr-api/common"
)

type uptimeHandler struct {
	start time.Time
}

func newUptimeHandler() *uptimeHandler {
	common.LogInfo().Print("initializing uptime handler " + time.Now().String())
	return &uptimeHandler{
		start: time.Now(),
	}
}

type uptimeResponse struct {
	UptimeSeconds float64 `json:"uptime_seconds"`
}

func (h uptimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := uptimeResponse{time.Now().Sub(h.start).Seconds()}
	data, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
