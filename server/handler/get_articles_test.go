package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/goodleby/golang-app/article"
)

type fakeAllArticlesSelector struct {
	articles   []article.Article
	shouldFail bool
}

func (f *fakeAllArticlesSelector) SelectAllArticles(ctx context.Context) ([]article.Article, error) {
	if f.shouldFail {
		return nil, errors.New("test error")
	}

	return f.articles, nil
}

func TestGetAllArticles(t *testing.T) {
	type args struct {
		allArticlesSelector AllArticlesSelector
		req                 *http.Request
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
		wantBody   []article.Article
	}{
		{
			name: "should return all articles from the database",
			args: args{
				allArticlesSelector: &fakeAllArticlesSelector{
					articles: []article.Article{
						{
							ID:          "test_id",
							Title:       "Test title",
							Description: "Test description",
							Body:        "Test body",
						},
						{
							ID:          "other_id",
							Title:       "Other test title",
							Description: "Other test description",
							Body:        "Other test body",
						},
					},
					shouldFail: false,
				},
				req: httptest.NewRequest(http.MethodGet, "/articles", nil),
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
			wantBody: []article.Article{
				{
					ID:          "test_id",
					Title:       "Test title",
					Description: "Test description",
					Body:        "Test body",
				},
				{
					ID:          "other_id",
					Title:       "Other test title",
					Description: "Other test description",
					Body:        "Other test body",
				},
			},
		},
		{
			name: "should return empty array if no articles in the database",
			args: args{
				allArticlesSelector: &fakeAllArticlesSelector{
					articles:   []article.Article{},
					shouldFail: false,
				},
				req: httptest.NewRequest(http.MethodGet, "/articles", nil),
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
			wantBody:   []article.Article{},
		},
		{
			name: "should return an internal error if it fails to select articles from the database",
			args: args{
				allArticlesSelector: &fakeAllArticlesSelector{
					articles:   []article.Article{},
					shouldFail: true,
				},
				req: httptest.NewRequest(http.MethodGet, "/articles", nil),
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
			wantBody:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler := GetAllArticles(tt.args.allArticlesSelector)
			handler(w, tt.args.req)

			if w.Code != tt.wantStatus {
				t.Fatalf("GetAllArticles() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// If we expect an error, we just need to check the response body is not empty.
			if tt.wantErr {
				if w.Body.Len() == 0 {
					t.Fatalf("GetAllArticles() response body is empty, want error")
				}
				return
			}

			// Decode the response body into a article.Article struct for comparison.
			var resBody []article.Article
			if err := json.NewDecoder(w.Body).Decode(&resBody); err != nil {
				t.Fatalf("GetAllArticles() error json decoding response body: %v", err)
			}

			if !reflect.DeepEqual(resBody, tt.wantBody) {
				t.Fatalf("GetAllArticles() response body = %v, want %v", resBody, tt.wantBody)
			}
		})
	}
}
