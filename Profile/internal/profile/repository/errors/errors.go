package repoerr

import "errors"

var ErrConflict = errors.New("conflict")
var ErrFrobidden = errors.New("forbidden")
var ErrNotFound = errors.New("not found")

type RepoError struct {
	Field string
	Err   error
}

func (e *RepoError) Error() string {
	return e.Err.Error()
}

func (e *RepoError) Unwrap() error {
	return e.Err
}
