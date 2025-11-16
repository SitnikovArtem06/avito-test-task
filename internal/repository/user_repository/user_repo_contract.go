package user_repository

import (
	"context"
	"github.com/SitnikovArtem06/avito-test-task/internal/model"
)

type UserRepo interface {
	UpdateActive(ctx context.Context, id string, active bool) (*model.User, error)
	GetUserById(ctx context.Context, id string) (*model.User, error)

	GetReviewerPullRequestIDs(ctx context.Context, userID string) ([]string, error)
}
