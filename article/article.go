package article

import "errors"

type Article struct {
	Payload
	ID int `json:"id" db:"id"`
}

type Payload struct {
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Body        string `json:"body" db:"body"`
}

func (p *Payload) Validate() error {
	if p.Title == "" {
		return errors.New("article payload Title is empty")
	}

	if p.Description == "" {
		return errors.New("article payload Description is empty")
	}

	if p.Body == "" {
		return errors.New("article payload Body is empty")
	}

	return nil
}
