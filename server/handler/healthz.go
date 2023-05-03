package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// HealthStatus is a struct that represents the health status of the server.
type HealthStatus struct {
	Status string `json:"status"`
}

// Healthz is a handler that returns the health status of the server.
func Healthz(w http.ResponseWriter, r *http.Request) {
	hs := HealthStatus{
		Status: http.StatusText(http.StatusOK),
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(hs); err != nil {
		log.Printf("%s: %v", logMsgWriteResponse, err)
	}
}
