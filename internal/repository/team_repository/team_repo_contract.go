package team_repository

import (
	"context"
	"github.com/SitnikovArtem06/avito-test-task/internal/model"
)

type TeamRepo interface {
	CreateTeam(ctx context.Context, team *model.TeamDb, users []*model.UserDb) error
	GetTeam(ctx context.Context, teamName string) (*model.Team, error)

	GetTeamNameByMemberId(ctx context.Context, id string) (string, error)
}
