package main

import (
	"net/http"
	"encoding/json"
)


// Error structure to send on failure
type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}
// sends a JSON message with the given error code
func sendJsonError(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: status, Text: http.StatusText(status)}); err != nil {
		panic(err)
	}
}

// sends a JSON message containing the given object
func sendJsonData(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
