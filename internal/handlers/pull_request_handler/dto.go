package pull_request_handler

import "time"

type createPRDto struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type errorBodyDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type pullRequestDTO struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}

type mergeRequest struct {
	PullRequestId string `json:"pull_request_id"`
}

type mergeResp struct {
	PullRequestID     string    `json:"pull_request_id"`
	PullRequestName   string    `json:"pull_request_name"`
	AuthorID          string    `json:"author_id"`
	Status            string    `json:"status"`
	AssignedReviewers []string  `json:"assigned_reviewers"`
	MergedAt          time.Time `json:"mergedAt"`
}

type reassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

type reassignRespDto struct {
	Pr          pullRequestDTO `json:"pr"`
	NewReviewer string         `json:"replaced_by"`
}
