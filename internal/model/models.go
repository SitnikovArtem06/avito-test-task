package model

import "time"

type User struct {
	Id       string
	Name     string
	IsActive bool
	TeamName string
}

type Team struct {
	Name    string
	Members []User
}

type TeamCreateRequest struct {
	TeamName string
	Members  []User
}

type PullRequest struct {
	Id        string
	Name      string
	AuthorId  string
	Status    Status
	Reviewers []User
	CreatedAt time.Time
	MergedAt  *time.Time
}

type CreatePullRequestRequest struct {
	Id       string
	Name     string
	AuthorId string
}

type SetActiveRequest struct {
	Id       string
	IsActive bool
}

type ReassignResp struct {
	Pr          *PullRequest
	NewReviewer string
}
type Status string

const (
	OPEN   Status = "OPEN"
	MERGED Status = "MERGED"
)

func (s Status) String() string {
	return string(s)
}
