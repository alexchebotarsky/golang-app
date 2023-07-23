package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/goodleby/golang-server/client/example"
)

// ExampleDataFetcher is an interface that fetches example data.
type ExampleDataFetcher interface {
	FetchExampleData(ctx context.Context) (*example.ExampleData, error)
}

// GetExampleData is a handler that fetches example data from example client.
func GetExampleData(exampleFetcher ExampleDataFetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		exampleData, err := exampleFetcher.FetchExampleData(ctx)
		if err != nil {
			HandleError(ctx, w, fmt.Errorf("error fetching articles: %v", err), http.StatusInternalServerError, true)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(exampleData); err != nil {
			log.Printf("%s: %v", logMsgWriteResponse, err)
		}
	}
}
