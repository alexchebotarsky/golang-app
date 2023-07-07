package handler

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"reflect"
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
		wantBody   *errorResponse
	}{
		{
			name: "should set passed status and write passed error",
			args: args{
				err:        errors.New("test error"),
				statusCode: 500,
				shouldLog:  false,
			},
			wantStatus: 500,
			wantBody: &errorResponse{
				Error:      "test error",
				StatusCode: 500,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			HandleError(w, tt.args.err, tt.args.statusCode, tt.args.shouldLog)

			if w.Code != tt.wantStatus {
				t.Fatalf("handleError() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// Decode the response body into a database.Article struct for comparison.
			var resBody errorResponse
			if err := json.NewDecoder(w.Body).Decode(&resBody); err != nil {
				t.Fatalf("GetArticle() error json decoding response body: %v", err)
			}

			if !reflect.DeepEqual(&resBody, tt.wantBody) {
				t.Fatalf("GetArticle() response body = %v, want %v", resBody, tt.wantBody)
			}
		})
	}
}
