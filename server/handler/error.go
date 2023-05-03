package handler

import (
	"io"
	"log"
	"net/http"
)

// handleError is a helper function that writes an error to the response.
func handleError(w http.ResponseWriter, err error, statusCode int, shouldLog bool) {
	if shouldLog {
		log.Print(err)
	}

	w.WriteHeader(statusCode)
	if _, err := io.WriteString(w, err.Error()); err != nil {
		log.Printf("%s: %v", logMsgWriteResponse, err)
	}
}
