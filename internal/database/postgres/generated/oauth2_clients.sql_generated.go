// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: oauth2_clients.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const archiveOAuth2Client = `-- name: ArchiveOAuth2Client :execrows

UPDATE oauth2_clients SET
	archived_at = NOW()
WHERE archived_at IS NULL
	AND id = $1
`

func (q *Queries) ArchiveOAuth2Client(ctx context.Context, db DBTX, id string) (int64, error) {
	result, err := db.ExecContext(ctx, archiveOAuth2Client, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const createOAuth2Client = `-- name: CreateOAuth2Client :exec

INSERT INTO oauth2_clients (
	id,
	name,
	description,
	client_id,
	client_secret
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
)
`

type CreateOAuth2ClientParams struct {
	ID           string
	Name         string
	Description  string
	ClientID     string
	ClientSecret string
}

func (q *Queries) CreateOAuth2Client(ctx context.Context, db DBTX, arg *CreateOAuth2ClientParams) error {
	_, err := db.ExecContext(ctx, createOAuth2Client,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.ClientID,
		arg.ClientSecret,
	)
	return err
}

const getOAuth2ClientByClientID = `-- name: GetOAuth2ClientByClientID :one

SELECT
	oauth2_clients.id,
	oauth2_clients.name,
	oauth2_clients.description,
	oauth2_clients.client_id,
	oauth2_clients.client_secret,
	oauth2_clients.created_at,
	oauth2_clients.archived_at
FROM oauth2_clients
WHERE oauth2_clients.archived_at IS NULL
	AND oauth2_clients.client_id = $1
`

func (q *Queries) GetOAuth2ClientByClientID(ctx context.Context, db DBTX, clientID string) (*Oauth2Clients, error) {
	row := db.QueryRowContext(ctx, getOAuth2ClientByClientID, clientID)
	var i Oauth2Clients
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.ClientID,
		&i.ClientSecret,
		&i.CreatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}

const getOAuth2ClientByDatabaseID = `-- name: GetOAuth2ClientByDatabaseID :one

SELECT
	oauth2_clients.id,
	oauth2_clients.name,
	oauth2_clients.description,
	oauth2_clients.client_id,
	oauth2_clients.client_secret,
	oauth2_clients.created_at,
	oauth2_clients.archived_at
FROM oauth2_clients
WHERE oauth2_clients.archived_at IS NULL
	AND oauth2_clients.id = $1
`

func (q *Queries) GetOAuth2ClientByDatabaseID(ctx context.Context, db DBTX, id string) (*Oauth2Clients, error) {
	row := db.QueryRowContext(ctx, getOAuth2ClientByDatabaseID, id)
	var i Oauth2Clients
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.ClientID,
		&i.ClientSecret,
		&i.CreatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}

const getOAuth2Clients = `-- name: GetOAuth2Clients :many

SELECT
	oauth2_clients.id,
	oauth2_clients.name,
	oauth2_clients.description,
	oauth2_clients.client_id,
	oauth2_clients.client_secret,
	oauth2_clients.created_at,
	oauth2_clients.archived_at,
	(
		SELECT COUNT(oauth2_clients.id)
		FROM oauth2_clients
		WHERE oauth2_clients.archived_at IS NULL
			AND oauth2_clients.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND oauth2_clients.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	) as filtered_count,
	(
		SELECT COUNT(users.id)
		FROM users
		WHERE users.archived_at IS NULL
	) AS total_count
FROM oauth2_clients
WHERE oauth2_clients.archived_at IS NULL
	AND oauth2_clients.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND oauth2_clients.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
LIMIT $4
OFFSET $3
`

type GetOAuth2ClientsParams struct {
	CreatedAfter  sql.NullTime
	CreatedBefore sql.NullTime
	QueryOffset   sql.NullInt32
	QueryLimit    sql.NullInt32
}

type GetOAuth2ClientsRow struct {
	ID            string
	Name          string
	Description   string
	ClientID      string
	ClientSecret  string
	CreatedAt     time.Time
	ArchivedAt    sql.NullTime
	FilteredCount int64
	TotalCount    int64
}

func (q *Queries) GetOAuth2Clients(ctx context.Context, db DBTX, arg *GetOAuth2ClientsParams) ([]*GetOAuth2ClientsRow, error) {
	rows, err := db.QueryContext(ctx, getOAuth2Clients,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetOAuth2ClientsRow{}
	for rows.Next() {
		var i GetOAuth2ClientsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.ClientID,
			&i.ClientSecret,
			&i.CreatedAt,
			&i.ArchivedAt,
			&i.FilteredCount,
			&i.TotalCount,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
