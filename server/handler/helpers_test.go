package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

// addChiURLParams adds the given url parameters to the request context.
func addChiURLParams(req *http.Request, params map[string]string) *http.Request {
	chiCtx := chi.NewRouteContext()
	for k, v := range params {
		chiCtx.URLParams.Add(k, v)
	}

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
}

// makeJSONBody encodes the given struct into a JSON body reader.
func makeJSONBody(t *testing.T, bodyStruct any) io.Reader {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(bodyStruct)
	if err != nil {
		t.Errorf("error encoding json body: %v", err)
	}

	return &buf
}
