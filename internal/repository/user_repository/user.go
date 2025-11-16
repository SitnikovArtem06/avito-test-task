package user_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/SitnikovArtem06/avito-test-task/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {

	return &UserRepository{db: db}
}

func (u *UserRepository) UpdateActive(ctx context.Context, id string, active bool) (*model.User, error) {

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	sql := `UPDATE users SET is_active = $2 WHERE id = $1;`

	tag, err := tx.Exec(ctx, sql, id, active)
	if err != nil {
		return nil, err
	}

	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}

	sqlSelect := `
        SELECT u.id, u.name, tm.team_name, u.is_active
        FROM users u
        JOIN team_members tm ON tm.user_id = u.id
        WHERE u.id = $1;
    `

	var user model.User
	if err = tx.QueryRow(ctx, sqlSelect, id).Scan(&user.Id, &user.Name, &user.TeamName, &user.IsActive); err != nil {
		return nil, fmt.Errorf("db: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &user, nil

}

func (u *UserRepository) GetUserById(ctx context.Context, id string) (*model.User, error) {

	sqlSelect := `SELECT id,name,is_active from users WHERE id = $1;`

	var user model.User

	if err := u.db.QueryRow(ctx, sqlSelect, id).Scan(&user.Id, &user.Name, &user.IsActive); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil

}

func (u *UserRepository) GetReviewerPullRequestIDs(ctx context.Context, userID string) ([]string, error) {

	sqlSelect := `
        SELECT pull_request_id
        FROM pull_requests_reviewers
        WHERE reviewer_id = $1;
    `

	rows, err := u.db.Query(ctx, sqlSelect, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		prIDs = append(prIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return prIDs, nil
}
