package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/goodleby/golang-server/article"
)

type fakeArticleDeleter struct {
	articles []article.Article
}

func (f *fakeArticleDeleter) DeleteArticle(ctx context.Context, id string) error {
	for i, article := range f.articles {
		if article.ID == id {
			f.articles = append(f.articles[:i], f.articles[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("article with id %q not found", id)
}

func TestDeleteArticle(t *testing.T) {
	type args struct {
		articleDeleter *fakeArticleDeleter
		req            *http.Request
	}
	tests := []struct {
		name         string
		args         args
		wantStatus   int
		wantErr      bool
		wantArticles []article.Article
	}{
		{
			name: "should remove the article from the database and return nil",
			args: args{
				articleDeleter: &fakeArticleDeleter{
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
				},
				req: addChiURLParams(httptest.NewRequest(http.MethodDelete, "/articles/test_id", nil), map[string]string{
					"id": "test_id",
				}),
			},
			wantStatus: http.StatusNoContent,
			wantErr:    false,
			wantArticles: []article.Article{
				{
					ID:          "other_id",
					Title:       "Other test title",
					Description: "Other test description",
					Body:        "Other test body",
				},
			},
		},
		{
			name: "should return an internal error if no article with the id found in the database",
			args: args{
				articleDeleter: &fakeArticleDeleter{
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
				},
				req: addChiURLParams(httptest.NewRequest(http.MethodDelete, "/articles/some_id", nil), map[string]string{"id": "some_id"}),
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
			wantArticles: []article.Article{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler := DeleteArticle(tt.args.articleDeleter)
			handler(w, tt.args.req)

			if w.Code != tt.wantStatus {
				t.Fatalf("DeleteArticle() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if !reflect.DeepEqual(tt.args.articleDeleter.articles, tt.wantArticles) {
				t.Fatalf("DeleteArticle() articles = %v, want %v", tt.args.articleDeleter.articles, tt.wantArticles)
			}

			// If we expect an error, we just need to check the response body is not empty.
			if tt.wantErr && w.Body.Len() == 0 {
				t.Fatalf("DeleteArticle() response body is empty, want error")
			}
		})
	}
}
