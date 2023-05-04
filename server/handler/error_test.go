package handler

import (
	"errors"
	"net/http/httptest"
	"testing"
)

func Test_handleError(t *testing.T) {
	type args struct {
		err        error
		statusCode int
		shouldLog  bool
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantBody   string
	}{
		{
			name: "should set passed status and write passed error",
			args: args{
				err:        errors.New("test error"),
				statusCode: 500,
				shouldLog:  false,
			},
			wantStatus: 500,
			wantBody:   "test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			handleError(w, tt.args.err, tt.args.statusCode, tt.args.shouldLog)

			if w.Code != tt.wantStatus {
				t.Fatalf("handleError() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if w.Body.String() != tt.wantBody {
				t.Fatalf("handleError() response body = %v, want %v", w.Body.String(), tt.wantBody)
			}
		})
	}
}
