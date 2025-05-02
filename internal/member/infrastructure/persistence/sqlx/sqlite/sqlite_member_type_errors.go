package sqlite

import "fmt"

type DBError struct {
	CustomError error
	RawError    error
}

func (e *DBError) Error() string {
	return fmt.Sprintf("%v: %v", e.CustomError, e.RawError)
}

func (e *DBError) Unwrap() error {
	return e.CustomError
}
