package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goodleby/golang-app/client/example"
)

type ExampleDataFetcher interface {
	FetchExampleData(ctx context.Context) (*example.ExampleData, error)
}

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

		err = json.NewEncoder(w).Encode(exampleData)
		handleWritingErr(err)
	}
}
