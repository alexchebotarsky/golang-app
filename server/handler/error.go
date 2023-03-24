package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func handleError(w http.ResponseWriter, err error, statusCode int, shouldLog bool) {
	errMsg := fmt.Sprintf("Encountered handler error, status %d %s: %v", statusCode, http.StatusText(statusCode), err)

	if shouldLog {
		log.Print(errMsg)
	}

	w.WriteHeader(statusCode)
	if _, err := io.WriteString(w, errMsg); err != nil {
		log.Printf("%s: %v", logMsgWriteResponse, err)
	}
}
