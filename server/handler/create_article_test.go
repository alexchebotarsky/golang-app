package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/goodleby/pure-go-server/client/database"
)

type fakeArticleCreator struct {
	articles   []database.Article
	shouldFail bool
}

func (m *fakeArticleCreator) CreateArticle(article database.Article) error {
	if m.shouldFail {
		return errors.New("test error")
	}

	m.articles = append(m.articles, article)

	return nil
}

func TestCreateArticle(t *testing.T) {
	type args struct {
		articleCreator *fakeArticleCreator
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
			name: "should create the passed article in the database",
			args: args{
				articleCreator: &fakeArticleCreator{
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
				req: addChiURLParams(httptest.NewRequest(http.MethodPatch, "/articles/test_id", makeJSONBody(t, &database.Article{
					ID:          "new_id",
					Title:       "New title",
					Description: "New description",
					Body:        "New body",
				})), map[string]string{
					"id": "test_id",
				}),
			},
			wantStatus: http.StatusNoContent,
			wantErr:    false,
			wantArticles: []database.Article{
				{
					ID:          "test_id",
					Title:       "Test title",
					Description: "Test description",
					Body:        "Test body",
				},
				{
					ID:          "new_id",
					Title:       "New title",
					Description: "New description",
					Body:        "New body",
				},
			},
		},
		{
			name: "should return an internal error if it fails to create the article in the database",
			args: args{
				articleCreator: &fakeArticleCreator{
					articles: []database.Article{
						{
							ID:          "test_id",
							Title:       "Test title",
							Description: "Test description",
							Body:        "Test body",
						},
					},
					shouldFail: true,
				},
				req: addChiURLParams(httptest.NewRequest(http.MethodPatch, "/articles/test_id", makeJSONBody(t, &database.Article{
					ID:          "new_id",
					Title:       "New title",
					Description: "New description",
					Body:        "New body",
				})), map[string]string{
					"id": "test_id",
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
			},
		},
		{
			name: "should return a bad request error if wrong body is provided",
			args: args{
				articleCreator: &fakeArticleCreator{
					articles: []database.Article{
						{
							ID:          "test_id",
							Title:       "Test title",
							Description: "Test description",
							Body:        "Test body",
						},
					},
					shouldFail: true,
				},
				req: addChiURLParams(httptest.NewRequest(http.MethodPatch, "/articles/test_id", nil), map[string]string{
					"id": "test_id",
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler := CreateArticle(tt.args.articleCreator)
			handler(w, tt.args.req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetArticle() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if !reflect.DeepEqual(tt.args.articleCreator.articles, tt.wantArticles) {
				t.Errorf("GetArticle() articles = %v, want %v", tt.args.articleCreator.articles, tt.wantArticles)
			}

			// If we expect an error, we just need to check the response body is not empty.
			if tt.wantErr && w.Body.Len() == 0 {
				t.Errorf("GetArticle() response body is empty, want error")
			}
		})
	}
}
