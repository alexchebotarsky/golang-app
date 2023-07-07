package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

func HandleError(w http.ResponseWriter, err error, statusCode int, shouldLog bool) {
	if shouldLog {
		log.Printf("Handler error: %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(errorResponse{
		Error:      err.Error(),
		StatusCode: statusCode,
	}); err != nil {
		log.Printf("%s: %v", logMsgWriteResponse, err)
	}
}
