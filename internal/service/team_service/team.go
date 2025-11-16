package team_service

import (
	"context"
	"errors"
	"github.com/SitnikovArtem06/avito-test-task/internal/model"
	"github.com/SitnikovArtem06/avito-test-task/internal/repository/team_repository"
)

type TeamService struct {
	teamRepo team_repository.TeamRepo
}

func NewTeamService(teamRepo team_repository.TeamRepo) *TeamService {
	return &TeamService{teamRepo: teamRepo}
}

func (t *TeamService) CreateTeam(ctx context.Context, req *model.TeamCreateRequest) (*model.Team, error) {

	teamDb := &model.TeamDb{
		Name: req.TeamName,
	}

	usersDb := make([]*model.UserDb, 0, len(req.Members))
	for _, m := range req.Members {
		usersDb = append(usersDb, &model.UserDb{
			Id:       m.Id,
			Name:     m.Name,
			IsActive: m.IsActive,
		})
	}

	err := t.teamRepo.CreateTeam(ctx, teamDb, usersDb)
	if err != nil {
		if errors.Is(err, team_repository.ErrAlreadyExists) {
			return nil, ErrAlreadyExists
		}
		return nil, err
	}

	team := &model.Team{
		Name:    req.TeamName,
		Members: req.Members,
	}

	return team, nil
}

func (t *TeamService) GetTeam(ctx context.Context, teamName string) (*model.Team, error) {

	team, err := t.teamRepo.GetTeam(ctx, teamName)
	if err != nil {
		if errors.Is(err, team_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return team, nil

}
