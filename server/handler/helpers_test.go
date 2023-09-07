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

func addChiURLParams(req *http.Request, params map[string]string) *http.Request {
	chiCtx := chi.NewRouteContext()
	for k, v := range params {
		chiCtx.URLParams.Add(k, v)
	}

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
}

func makeJSONBody(t *testing.T, bodyStruct any) io.Reader {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(bodyStruct)
	if err != nil {
		t.Fatalf("error encoding json body: %v", err)
	}

	return &buf
}
