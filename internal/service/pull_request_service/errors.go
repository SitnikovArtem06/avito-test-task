package pull_request_service

import "errors"

var (
	ErrNotFoundAuthor error = errors.New("author not found")

	ErrNotFoundUser error = errors.New("user not found")

	ErrNotFoundTeam error = errors.New("team not found")

	ErrAlreadyExist error = errors.New("PR already exists")

	ErrNoCandidate error = errors.New("no one of member can be reviewer")

	ErrMerged   error = errors.New("pr already merged")
	ErrNotFound error = errors.New("pr not found")

	ErrNotAssigned = errors.New("reviewer wasnt assigned on this pr")
)
