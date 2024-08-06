package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandleError(t *testing.T) {
	type args struct {
		err        error
		statusCode int
		shouldLog  bool
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantBody   *errorResponse
	}{
		{
			name: "should set passed status and write status text error",
			args: args{
				err:        errors.New("Test error"),
				statusCode: 500,
				shouldLog:  false,
			},
			wantStatus: 500,
			wantBody: &errorResponse{
				Error:      "Internal Server Error: 500",
				StatusCode: 500,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			HandleError(context.TODO(), w, tt.args.err, tt.args.statusCode, tt.args.shouldLog)

			if w.Code != tt.wantStatus {
				t.Fatalf("HandleError() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var resBody errorResponse
			err := json.NewDecoder(w.Body).Decode(&resBody)
			if err != nil {
				t.Fatalf("HandleError() error json decoding response body: %v", err)
			}

			if !reflect.DeepEqual(&resBody, tt.wantBody) {
				t.Fatalf("HandleError() response body = %v, want %v", resBody, tt.wantBody)
			}
		})
	}
}
