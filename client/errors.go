package client

import (
	"errors"
	"strings"
)

type ErrNotFound struct {
	Err error
}

func (e ErrNotFound) Error() string {
	return e.Err.Error()
}

func (e ErrNotFound) Unwrap() error {
	return e.Err
}

type ErrUnauthorized struct {
	Err error
}

func (e ErrUnauthorized) Error() string {
	return e.Err.Error()
}

func (e ErrUnauthorized) Unwrap() error {
	return e.Err
}

type ErrMultiple struct {
	Errs []error
}

func (e ErrMultiple) Error() string {
	errStrings := make([]string, 0, len(e.Errs))

	for _, err := range e.Errs {
		errStrings = append(errStrings, err.Error())
	}

	return strings.Join(errStrings, " | ")
}

func (e ErrMultiple) Unwrap() error {
	return errors.New(e.Error())
}
