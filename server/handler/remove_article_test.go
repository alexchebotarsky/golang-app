package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/goodleby/pure-go-server/client/database"
)

type fakeArticleRemover struct {
	articles []database.Article
}

func (m *fakeArticleRemover) RemoveArticle(id string) error {
	for i, article := range m.articles {
		if article.ID == id {
			m.articles = append(m.articles[:i], m.articles[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Article with id '%s' not found", id)
}

func TestRemoveArticle(t *testing.T) {
	type args struct {
		articleRemover *fakeArticleRemover
		req            *http.Request
	}
	tests := []struct {
		name         string
		args         args
		wantStatus   int
		wantErr      bool
		wantArticles []database.Article
	}{
		{
			name: "should remove the article from the database and return nil",
			args: args{
				articleRemover: &fakeArticleRemover{
					articles: []database.Article{
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
			wantArticles: []database.Article{
				{
					ID:          "other_id",
					Title:       "Other test title",
					Description: "Other test description",
					Body:        "Other test body",
				},
			},
		},
		{
			name: "should return an internal error if it fails no article with the id found in the database",
			args: args{
				articleRemover: &fakeArticleRemover{
					articles: []database.Article{
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
			wantArticles: []database.Article{
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
			handler := RemoveArticle(tt.args.articleRemover)
			handler(w, tt.args.req)

			if w.Code != tt.wantStatus {
				t.Fatalf("RemoveArticle() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if !reflect.DeepEqual(tt.args.articleRemover.articles, tt.wantArticles) {
				t.Fatalf("RemoveArticle() articles = %v, want %v", tt.args.articleRemover.articles, tt.wantArticles)
			}

			// If we expect an error, we just need to check the response body is not empty.
			if tt.wantErr && w.Body.Len() == 0 {
				t.Fatalf("RemoveArticle() response body is empty, want error")
			}
		})
	}
}
