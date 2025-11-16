package pull_request_repository

import (
	"context"
	"github.com/SitnikovArtem06/avito-test-task/internal/model"
)

type PullRequestRepo interface {
	CreatePullRequest(ctx context.Context, pullReq *model.PullRequestDb, reviewers []model.User) (*model.PullRequest, error)
	Merge(ctx context.Context, id string) error

	GetByID(ctx context.Context, id string) (*model.PullRequestDb, error)

	ReassignReviewer(ctx context.Context, prID string, oldReviewerID string, newReviewerID string) error

	GetReviewers(ctx context.Context, prID string) ([]model.User, error)
}
