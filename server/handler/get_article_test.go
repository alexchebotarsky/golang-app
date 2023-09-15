package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/goodleby/golang-server/client/database"
)

type fakeArticleSelector struct {
	articles   []database.Article
	shouldFail bool
}

func (f *fakeArticleSelector) SelectArticle(ctx context.Context, id string) (*database.Article, error) {
	if f.shouldFail {
		return nil, errors.New("test error")
	}

	for _, a := range f.articles {
		if a.ID == id {
			return &a, nil
		}
	}

	return nil, &database.ErrNotFound{Err: fmt.Errorf("article with id %q not found", id)}
}

func TestGetArticle(t *testing.T) {
	type args struct {
		articleSelector ArticleSelector
		req             *http.Request
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
		wantBody   *database.Article
	}{
		{
			name: "should return an article from the database",
			args: args{
				articleSelector: &fakeArticleSelector{
					articles: []database.Article{
						{
							ID:          "other_id",
							Title:       "Other test title",
							Description: "Other test description",
							Body:        "Other test body",
						},
						{
							ID:          "test_id",
							Title:       "Test title",
							Description: "Test description",
							Body:        "Test body",
						},
					},
					shouldFail: false,
				},
				req: addChiURLParams(httptest.NewRequest(http.MethodGet, "/articles/test_id", nil), map[string]string{
					"id": "test_id",
				}),
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
			wantBody: &database.Article{
				ID:          "test_id",
				Title:       "Test title",
				Description: "Test description",
				Body:        "Test body",
			},
		},
		{
			name: "should return an err not found if no article with the id found in the database",
			args: args{
				articleSelector: &fakeArticleSelector{
					articles: []database.Article{
						{
							ID:          "test_id",
							Title:       "Test title",
							Description: "Test description",
							Body:        "Test body",
						},
					},
					shouldFail: false,
				},
				req: addChiURLParams(httptest.NewRequest(http.MethodGet, "/articles/some_id", nil), map[string]string{
					"id": "some_id",
				}),
			},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
			wantBody:   nil,
		},
		{
			name: "should return an internal error if failed to select article from the database",
			args: args{
				articleSelector: &fakeArticleSelector{
					articles:   []database.Article{},
					shouldFail: true,
				},
				req: addChiURLParams(httptest.NewRequest(http.MethodGet, "/articles/some_id", nil), map[string]string{
					"id": "some_id",
				}),
			},
			wantErr:    true,
			wantStatus: http.StatusInternalServerError,
			wantBody:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler := GetArticle(tt.args.articleSelector)
			handler(w, tt.args.req)

			if w.Code != tt.wantStatus {
				t.Fatalf("GetArticle() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// If we expect an error, we just need to check the response body is not empty.
			if tt.wantErr {
				if w.Body.Len() == 0 {
					t.Fatalf("GetArticle() response body is empty, want error")
				}
				return
			}

			// Decode the response body into a database.Article struct for comparison.
			var resBody database.Article
			if err := json.NewDecoder(w.Body).Decode(&resBody); err != nil {
				t.Fatalf("GetArticle() error json decoding response body: %v", err)
			}

			if !reflect.DeepEqual(&resBody, tt.wantBody) {
				t.Fatalf("GetArticle() response body = %v, want %v", resBody, tt.wantBody)
			}
		})
	}
}
