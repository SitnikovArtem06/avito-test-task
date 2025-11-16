package model

import "time"

type UserDb struct {
	Id       string `db:"id"`
	Name     string `db:"name"`
	IsActive bool   `db:"is_active"`
}

type TeamDb struct {
	Name string `db:"name"`
}

type MemberTeamDb struct {
	UserId   string `db:"user_id"`
	TeamName string `db:"team_name"`
}

type PullRequestDb struct {
	Id        string     `db:"id"`
	Name      string     `db:"name"`
	AuthorId  string     `db:"author_id"`
	Status    Status     `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
	MergedAt  *time.Time `db:"merged_at"`
}

type PullRequestReviewerDB struct {
	PullRequestId string `db:"pull_request_id"`
	ReviewerId    string `db:"reviewer_id"`
}
