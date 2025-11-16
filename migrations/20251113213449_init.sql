-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE teams(
    name TEXT PRIMARY KEY
);

CREATE TABLE users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL
);

CREATE TABLE team_members (
    user_id  TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_name  TEXT NOT NULL REFERENCES teams(name) ON DELETE CASCADE,
    PRIMARY KEY (user_id, team_name)
);

CREATE TABLE pull_requests(
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    merged_at TIMESTAMP

);

CREATE TABLE pull_requests_reviewers(
    pull_request_id TEXT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    PRIMARY KEY (pull_request_id, reviewer_id)

)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS pull_request_reviewers;
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
