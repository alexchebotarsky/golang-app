package article

import "errors"

type Article struct {
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Body        string `json:"body" db:"body"`
	ID          string `json:"id" db:"id"`
}

func (a *Article) Validate() error {
	if a.Title == "" {
		return errors.New("article Title is empty")
	}

	if a.Description == "" {
		return errors.New("article Description is empty")
	}

	if a.Body == "" {
		return errors.New("article Body is empty")
	}

	if a.ID == "" {
		return errors.New("article ID is empty")
	}

	return nil
}
