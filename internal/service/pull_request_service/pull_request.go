package pull_request_service

import (
	"context"
	"errors"
	"github.com/SitnikovArtem06/avito-test-task/internal/model"
	"github.com/SitnikovArtem06/avito-test-task/internal/repository/pull_request_repository"
	"github.com/SitnikovArtem06/avito-test-task/internal/repository/team_repository"
	"github.com/SitnikovArtem06/avito-test-task/internal/repository/user_repository"
	"math/rand"
	"time"
)

type PullRequestService struct {
	tr  team_repository.TeamRepo
	ur  user_repository.UserRepo
	rr  pull_request_repository.PullRequestRepo
	rnd *rand.Rand
}

func NewPullRequestService(tr team_repository.TeamRepo, ur user_repository.UserRepo, rr pull_request_repository.PullRequestRepo) *PullRequestService {
	return &PullRequestService{
		tr:  tr,
		ur:  ur,
		rr:  rr,
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (p *PullRequestService) CreatePullRequest(ctx context.Context, req *model.CreatePullRequestRequest) (*model.PullRequest, error) {

	if _, err := p.ur.GetUserById(ctx, req.AuthorId); err != nil {
		if errors.Is(err, user_repository.ErrNotFound) {
			return nil, ErrNotFoundAuthor
		}
		return nil, err
	}

	name, err := p.tr.GetTeamNameByMemberId(ctx, req.AuthorId)
	if err != nil {
		if errors.Is(err, team_repository.ErrNotFound) {
			return nil, ErrNotFoundTeam
		}
		return nil, err
	}

	team, err := p.tr.GetTeam(ctx, name)
	if err != nil {
		if errors.Is(err, team_repository.ErrNotFound) {
			return nil, ErrNotFoundTeam
		}
		return nil, err
	}

	users := team.Members

	var potentialReviewers []model.User

	for _, u := range users {
		if u.Id != req.AuthorId && u.IsActive {
			potentialReviewers = append(potentialReviewers, u)
		}
	}

	n := len(potentialReviewers)

	k := 2
	if n < k {
		k = n
	}

	reviewers := potentialReviewers[:k]

	reqDb := &model.PullRequestDb{
		Id:        req.Id,
		Name:      req.Name,
		AuthorId:  req.AuthorId,
		Status:    model.OPEN,
		CreatedAt: time.Time{},
		MergedAt:  nil,
	}

	pr, err := p.rr.CreatePullRequest(ctx, reqDb, reviewers)
	if err != nil {
		if errors.Is(err, pull_request_repository.ErrAlreadyExists) {
			return nil, ErrAlreadyExist
		}

		return nil, err
	}

	return pr, nil

}

func (p *PullRequestService) Merge(ctx context.Context, id string) (*model.PullRequest, error) {

	pr, err := p.rr.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pull_request_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	users, err := p.rr.GetReviewers(ctx, id)
	if err != nil {
		return nil, err
	}

	if pr.Status == model.MERGED {
		return &model.PullRequest{
			Id:        pr.Id,
			Name:      pr.Name,
			AuthorId:  pr.AuthorId,
			Status:    pr.Status,
			Reviewers: users,
			CreatedAt: pr.CreatedAt,
			MergedAt:  pr.MergedAt,
		}, nil
	}

	if err = p.rr.Merge(ctx, id); err != nil {
		if errors.Is(err, pull_request_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	pr2, err := p.rr.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pull_request_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &model.PullRequest{
		Id:        pr2.Id,
		Name:      pr2.Name,
		AuthorId:  pr2.AuthorId,
		Status:    pr2.Status,
		Reviewers: users,
		CreatedAt: pr2.CreatedAt,
		MergedAt:  pr2.MergedAt,
	}, nil

}

func (p *PullRequestService) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*model.ReassignResp, error) {

	if _, err := p.ur.GetUserById(ctx, oldReviewerID); err != nil {
		if errors.Is(err, user_repository.ErrNotFound) {
			return nil, ErrNotFoundUser
		}
		return nil, err
	}

	prDb, err := p.rr.GetByID(ctx, prID)
	if err != nil {
		if errors.Is(err, pull_request_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if prDb.Status == model.MERGED {
		return nil, ErrMerged
	}

	reviewers, err := p.rr.GetReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}

	check := make(map[string]struct{})

	for _, r := range reviewers {
		check[r.Id] = struct{}{}
	}

	teamName, err := p.tr.GetTeamNameByMemberId(ctx, oldReviewerID)
	if err != nil {
		if errors.Is(err, team_repository.ErrNotFound) {
			return nil, ErrNotFoundTeam
		}
		return nil, err
	}

	team, err := p.tr.GetTeam(ctx, teamName)
	if err != nil {
		if errors.Is(err, team_repository.ErrNotFound) {
			return nil, ErrNotFoundTeam
		}
		return nil, err
	}

	members := team.Members

	var candidates []model.User

	for _, m := range members {
		_, ok := check[m.Id]
		if m.IsActive && !ok && m.Id != prDb.AuthorId && m.Id != oldReviewerID {
			candidates = append(candidates, m)
		}
	}

	if len(candidates) == 0 {
		return nil, ErrNoCandidate
	}

	idx := p.rnd.Intn(len(candidates))
	candidateID := candidates[idx].Id

	if err = p.rr.ReassignReviewer(ctx, prID, oldReviewerID, candidateID); err != nil {
		if errors.Is(err, pull_request_repository.ErrNotAssigned) {
			return nil, ErrNotAssigned
		}
		return nil, err
	}

	reviewers, err = p.rr.GetReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}

	pr := &model.PullRequest{
		Id:        prDb.Id,
		Name:      prDb.Name,
		AuthorId:  prDb.AuthorId,
		Status:    prDb.Status,
		Reviewers: reviewers,
		CreatedAt: prDb.CreatedAt,
		MergedAt:  prDb.MergedAt,
	}
	return &model.ReassignResp{
		pr,
		candidateID,
	}, nil
}
