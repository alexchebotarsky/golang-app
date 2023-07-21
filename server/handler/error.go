package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type errorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

func HandleError(ctx context.Context, w http.ResponseWriter, err error, statusCode int, shouldLog bool) {
	span := trace.SpanFromContext(ctx)

	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

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
