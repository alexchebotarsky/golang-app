package database

type ErrNotFound struct {
	Err error
}

func (e *ErrNotFound) Error() string {
	return e.Err.Error()
}

func (e *ErrNotFound) Unwrap() error {
	return e.Err
}
