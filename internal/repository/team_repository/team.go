package team_repository

import (
	"context"
	"errors"
	"github.com/SitnikovArtem06/avito-test-task/internal/model"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	dbpool *pgxpool.Pool
}

func NewTeamRepository(dbpool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{dbpool: dbpool}
}

func (t *TeamRepository) CreateTeam(ctx context.Context, team *model.TeamDb, users []*model.UserDb) error {

	tx, err := t.dbpool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	sqlAddTeam := `INSERT INTO teams(name) VALUES($1);`

	_, err = tx.Exec(ctx, sqlAddTeam, team.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return ErrAlreadyExists
			}
		}
		return err
	}

	sqlAddUser := `INSERT INTO users (id,name,is_active) VALUES ($1,$2,$3) ON CONFLICT(id) DO UPDATE SET
    name = EXCLUDED.name, is_active = EXCLUDED.is_active;`

	for _, u := range users {
		_, err = tx.Exec(ctx, sqlAddUser, u.Id, u.Name, u.IsActive)
		if err != nil {
			return err
		}
	}

	sqlAddMember := `INSERT INTO team_members(user_id, team_name) VALUES($1,$2);`
	for _, u := range users {
		_, err = tx.Exec(ctx, sqlAddMember, u.Id, team.Name)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil

}

func (t *TeamRepository) GetTeam(ctx context.Context, teamName string) (*model.Team, error) {

	tx, err := t.dbpool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	sqlGetTeam := `SELECT name FROM teams WHERE name = $1;`

	var teamNameDb string

	if err = tx.QueryRow(ctx, sqlGetTeam, teamName).Scan(&teamNameDb); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	sqlGetMembers := `SELECT u.id, u.name, u.is_active  FROM team_members tm join users u on u.id = tm.user_id WHERE tm.team_name = $1;`

	team := &model.Team{
		Name: teamName,
	}

	var user model.User

	rows, err := tx.Query(ctx, sqlGetMembers, teamName)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		if err = rows.Scan(&user.Id, &user.Name, &user.IsActive); err != nil {
			return nil, err
		}

		team.Members = append(team.Members, user)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return team, nil

}

func (t *TeamRepository) GetTeamNameByMemberId(ctx context.Context, id string) (string, error) {

	sqlSelect := `SELECT team_name from team_members WHERE user_id = $1;`

	var name string

	if err := t.dbpool.QueryRow(ctx, sqlSelect, id).Scan(&name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}

	return name, nil

}
