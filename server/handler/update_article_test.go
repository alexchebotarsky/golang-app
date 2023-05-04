package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/goodleby/pure-go-server/client/database"
)

type fakeArticleUpdater struct {
	articles []database.Article
}

func (m *fakeArticleUpdater) UpdateArticle(id string, article database.Article) error {
	for i, a := range m.articles {
		if a.ID == id {
			m.articles[i] = article
			return nil
		}
	}

	return fmt.Errorf("Article with id '%s' not found", id)
}

func TestUpdateArticle(t *testing.T) {
	type args struct {
		articleUpdater *fakeArticleUpdater
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
			name: "should return the updated article when successful",
			args: args{
				articleUpdater: &fakeArticleUpdater{
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
				req: addChiURLParams(httptest.NewRequest(http.MethodPatch, "/articles/test_id", makeJSONBody(t, &database.Article{
					ID:          "test_id",
					Title:       "Updated title",
					Description: "Updated description",
					Body:        "Updated body",
				})), map[string]string{
					"id": "test_id",
				}),
			},
			wantStatus: http.StatusNoContent,
			wantErr:    false,
			wantArticles: []database.Article{
				{
					ID:          "test_id",
					Title:       "Updated title",
					Description: "Updated description",
					Body:        "Updated body",
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
			name: "should return an internal error if no article with the id found in the database",
			args: args{
				articleUpdater: &fakeArticleUpdater{
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
				req: addChiURLParams(httptest.NewRequest(http.MethodPatch, "/articles/some_id", makeJSONBody(t, &database.Article{
					ID:          "some_id",
					Title:       "Updated title",
					Description: "Updated description",
					Body:        "Updated body",
				})), map[string]string{
					"id": "some_id",
				}),
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
		{
			name: "should return a bad request error if wrong body is provided",
			args: args{
				articleUpdater: &fakeArticleUpdater{
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
				req: addChiURLParams(httptest.NewRequest(http.MethodPatch, "/articles/some_id", nil), map[string]string{
					"id": "some_id",
				}),
			},
			wantStatus: http.StatusBadRequest,
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
			handler := UpdateArticle(tt.args.articleUpdater)
			handler(w, tt.args.req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetArticle() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if !reflect.DeepEqual(tt.args.articleUpdater.articles, tt.wantArticles) {
				t.Errorf("GetArticle() articles = %v, want %v", tt.args.articleUpdater.articles, tt.wantArticles)
			}

			// If we expect an error, we just need to check the response body is not empty.
			if tt.wantErr && w.Body.Len() == 0 {
				t.Errorf("GetArticle() response body is empty, want error")
			}
		})
	}
}
