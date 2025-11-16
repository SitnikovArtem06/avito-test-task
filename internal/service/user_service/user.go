package user_service

import (
	"context"
	"errors"
	"github.com/SitnikovArtem06/avito-test-task/internal/model"
	"github.com/SitnikovArtem06/avito-test-task/internal/repository/pull_request_repository"
	"github.com/SitnikovArtem06/avito-test-task/internal/repository/user_repository"
)

type UserService struct {
	userRepo user_repository.UserRepo
	prRepo   pull_request_repository.PullRequestRepo
}

func NewUserService(userRepo user_repository.UserRepo, prRepo pull_request_repository.PullRequestRepo) *UserService {
	return &UserService{userRepo: userRepo, prRepo: prRepo}
}

func (u *UserService) UpdateActive(ctx context.Context, req *model.SetActiveRequest) (*model.User, error) {

	user, err := u.userRepo.UpdateActive(ctx, req.Id, req.IsActive)

	if err != nil {
		if errors.Is(err, user_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil

}

func (u *UserService) GetPRsByUser(ctx context.Context, id string) ([]model.PullRequest, error) {

	if _, err := u.userRepo.GetUserById(ctx, id); err != nil {
		if errors.Is(err, user_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	prIds, err := u.userRepo.GetReviewerPullRequestIDs(ctx, id)
	if err != nil {
		return nil, err
	}

	var prsDb []model.PullRequestDb

	for _, p := range prIds {

		pr, err := u.prRepo.GetByID(ctx, p)
		if err != nil {
			return nil, err
		}
		prsDb = append(prsDb, *pr)

	}

	var prs []model.PullRequest

	for _, p := range prsDb {

		prs = append(prs, model.PullRequest{
			Id:        p.Id,
			Name:      p.Name,
			AuthorId:  p.AuthorId,
			Status:    p.Status,
			Reviewers: nil,
			CreatedAt: p.CreatedAt,
			MergedAt:  p.MergedAt,
		})

	}

	return prs, nil

}
