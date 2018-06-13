package common

import (
	"encoding/json"
	"net/http"
)

type httpErrorBody struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

// WriteHTTPError is a utility function for writing a standard error
func WriteHTTPError(w http.ResponseWriter, status int, err error) {
	LogError().Printf("error response code=%d error=%v", status, err)
	data, _ := json.Marshal(httpErrorBody{Code: status, Error: err.Error()})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
