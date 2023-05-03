package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantBody   HealthStatus
	}{
		{
			name: "should return status OK",
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/_healthz", nil),
			},
			wantStatus: http.StatusOK,
			wantBody: HealthStatus{
				Status: http.StatusText(http.StatusOK),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			Healthz(w, tt.args.req)

			if w.Code != tt.wantStatus {
				t.Errorf("Healthz() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var wantBody bytes.Buffer
			if err := json.NewEncoder(&wantBody).Encode(tt.wantBody); err != nil {
				t.Errorf("Healthz() error json encoding wantBody: %v", err)
			}

			if w.Body.String() != wantBody.String() {
				t.Errorf("Healthz() response body = %v, want %v", w.Body.String(), wantBody.Bytes())
			}
		})
	}
}
