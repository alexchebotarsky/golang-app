package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/goodleby/golang-app/tracing"
)

type errorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

func HandleError(ctx context.Context, w http.ResponseWriter, err error, statusCode int, shouldLog bool) {
	span := tracing.SpanFromContext(ctx)

	span.RecordError(err)

	if shouldLog {
		slog.Error(fmt.Sprintf("Handler error: %v", err), "status", statusCode)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err = json.NewEncoder(w).Encode(errorResponse{
		Error:      fmt.Sprintf("%v", err),
		StatusCode: statusCode,
	})
	handleWritingErr(err)
}
