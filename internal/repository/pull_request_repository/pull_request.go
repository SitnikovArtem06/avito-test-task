package pull_request_repository

import (
	"context"
	"errors"
	"github.com/SitnikovArtem06/avito-test-task/internal/model"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PullRequestRepository struct {
	dbpool *pgxpool.Pool
}

func NewPullRequestRepository(dbpool *pgxpool.Pool) *PullRequestRepository {
	return &PullRequestRepository{dbpool: dbpool}
}

func (p *PullRequestRepository) CreatePullRequest(ctx context.Context, pullReq *model.PullRequestDb, reviewers []model.User) (*model.PullRequest, error) {

	tx, err := p.dbpool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	var created_at time.Time

	sqlInsertPull := `INSERT INTO pull_requests (id, name, author_id, status) VALUES($1,$2,$3,$4) RETURNING created_at;`

	if err = tx.QueryRow(ctx, sqlInsertPull, pullReq.Id, pullReq.Name, pullReq.AuthorId, pullReq.Status).Scan(&created_at); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, ErrAlreadyExists
			}
		}
		return nil, err
	}

	sqlInsertReviewers := `INSERT INTO pull_requests_reviewers (pull_request_id, reviewer_id) VALUES($1,$2);`

	for _, r := range reviewers {

		if _, err = tx.Exec(ctx, sqlInsertReviewers, pullReq.Id, r.Id); err != nil {
			return nil, err
		}

	}

	pullRequest := &model.PullRequest{
		Id:        pullReq.Id,
		Name:      pullReq.Name,
		AuthorId:  pullReq.AuthorId,
		Status:    pullReq.Status,
		Reviewers: reviewers,
		CreatedAt: created_at,
		MergedAt:  nil,
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return pullRequest, nil

}

func (r *PullRequestRepository) Merge(ctx context.Context, id string) error {

	sqlUpdate := `
        UPDATE pull_requests
        SET status = $2,
            merged_at = NOW()
        WHERE id = $1;
    `

	tag, err := r.dbpool.Exec(ctx, sqlUpdate, id, model.MERGED)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *PullRequestRepository) ReassignReviewer(ctx context.Context, prID string, oldReviewerID string, newReviewerID string) error {
	tx, err := r.dbpool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sqlDeleteReviewer := `
        DELETE FROM pull_requests_reviewers
        WHERE pull_request_id = $1 AND reviewer_id = $2;
    `
	tag, err := tx.Exec(ctx, sqlDeleteReviewer, prID, oldReviewerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotAssigned
	}

	sqlInsertReviewer := `
        INSERT INTO pull_requests_reviewers (pull_request_id, reviewer_id)
        VALUES ($1, $2);
    `
	if _, err = tx.Exec(ctx, sqlInsertReviewer, prID, newReviewerID); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PullRequestRepository) GetReviewers(ctx context.Context, prID string) ([]model.User, error) {

	sqlSelectReviewers := `
        SELECT u.id, u.name, u.is_active
        FROM pull_requests_reviewers pr
        JOIN users u ON u.id = pr.reviewer_id
        WHERE pr.pull_request_id = $1;
    `

	rows, err := r.dbpool.Query(ctx, sqlSelectReviewers, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []model.User
	for rows.Next() {
		var u model.User
		if err = rows.Scan(&u.Id, &u.Name, &u.IsActive); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviewers, nil
}

func (r *PullRequestRepository) GetByID(ctx context.Context, id string) (*model.PullRequestDb, error) {

	sqlSelectPR := `
        SELECT id, name, author_id, status, created_at, merged_at
        FROM pull_requests
        WHERE id = $1;
    `
	var prDb model.PullRequestDb
	if err := r.dbpool.QueryRow(ctx, sqlSelectPR, id).Scan(&prDb.Id, &prDb.Name, &prDb.AuthorId, &prDb.Status, &prDb.CreatedAt, &prDb.MergedAt); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &prDb, nil
}
