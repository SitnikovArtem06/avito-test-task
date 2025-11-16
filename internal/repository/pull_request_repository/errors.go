package pull_request_repository

import "errors"

var (
	ErrAlreadyExists error = errors.New("pr already exists")
	ErrNotFound      error = errors.New("pr not found")

	ErrNotAssigned = errors.New("reviewer wasnt assigned on this pr")
)
