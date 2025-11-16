package user_handler

type setActiveRequest struct {
	Id       string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type userResp struct {
	Id       string `json:"user_id"`
	Name     string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type errorBodyDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type userReviewsResponse struct {
	UserID       string                  `json:"user_id"`
	PullRequests []userReviewPullRequest `json:"pull_requests"`
}

type userReviewPullRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}
