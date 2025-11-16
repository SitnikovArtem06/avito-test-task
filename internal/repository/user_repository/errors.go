package user_repository

import "errors"

var (
	ErrNotFound error = errors.New("user not found")
)
