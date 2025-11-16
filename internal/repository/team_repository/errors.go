package team_repository

import "errors"

var (
	ErrAlreadyExists error = errors.New("team already exist")
	ErrNotFound      error = errors.New("team not found")
)
