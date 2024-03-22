package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHealth(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantBody   *healthResponse
	}{
		{
			name: "should return status OK",
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/_healthz", nil),
			},
			wantStatus: http.StatusOK,
			wantBody: &healthResponse{
				Status: http.StatusText(http.StatusOK),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			Health(w, tt.args.req)

			if w.Code != tt.wantStatus {
				t.Fatalf("Health() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var resBody healthResponse
			err := json.NewDecoder(w.Body).Decode(&resBody)
			if err != nil {
				t.Fatalf("Health() error json decoding response body: %v", err)
			}

			if !reflect.DeepEqual(&resBody, tt.wantBody) {
				t.Fatalf("Health() response body = %v, want %v", resBody, tt.wantBody)
			}
		})
	}
}
