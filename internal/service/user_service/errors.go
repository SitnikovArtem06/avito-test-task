package user_service

import "errors"

var (
	ErrNotFound error = errors.New("user not found")
)
