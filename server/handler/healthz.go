package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthStatus struct {
	Status string `json:"status"`
}

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
