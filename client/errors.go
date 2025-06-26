package client

type ErrNotFound struct {
	Err error
}

func (e *ErrNotFound) Error() string {
	return e.Err.Error()
}

func (e *ErrNotFound) Unwrap() error {
	return e.Err
}

type ErrUnauthorized struct {
	Err error
}

func (e *ErrUnauthorized) Error() string {
	return e.Err.Error()
}

func (e *ErrUnauthorized) Unwrap() error {
	return e.Err
}
