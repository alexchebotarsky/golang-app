package handler

import (
	"encoding/json"
	"net/http"
)

type HealthStatus struct {
	Status string `json:"status"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	hs := HealthStatus{
		Status: http.StatusText(http.StatusOK),
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(hs)
	handleWritingErr(err)
}
