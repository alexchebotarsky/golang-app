package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/goodleby/golang-server/tracing"
)

type errorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

func HandleError(ctx context.Context, w http.ResponseWriter, err error, statusCode int, shouldLog bool) {
	span := tracing.SpanFromContext(ctx)

	span.RecordError(err)

	if shouldLog {
		log.Printf("Handler error: %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err = json.NewEncoder(w).Encode(errorResponse{
		Error:      err.Error(),
		StatusCode: statusCode,
	})
	handleWritingErr(err)
}
