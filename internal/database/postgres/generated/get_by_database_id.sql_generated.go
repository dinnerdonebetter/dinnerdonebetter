// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: get_by_database_id.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const GetOAuth2ClientByDatabaseID = `-- name: GetOAuth2ClientByDatabaseID :one

SELECT
	oauth2_clients.id,
	oauth2_clients.name,
	oauth2_clients.client_id,
	oauth2_clients.client_secret,
	oauth2_clients.created_at,
	oauth2_clients.archived_at
FROM oauth2_clients
WHERE oauth2_clients.archived_at IS NULL
	AND oauth2_clients.id = $1
`

type GetOAuth2ClientByDatabaseIDRow struct {
	ID           string       `db:"id"`
	Name         string       `db:"name"`
	ClientID     string       `db:"client_id"`
	ClientSecret string       `db:"client_secret"`
	CreatedAt    time.Time    `db:"created_at"`
	ArchivedAt   sql.NullTime `db:"archived_at"`
}

func (q *Queries) GetOAuth2ClientByDatabaseID(ctx context.Context, db DBTX, id string) (*GetOAuth2ClientByDatabaseIDRow, error) {
	row := db.QueryRowContext(ctx, GetOAuth2ClientByDatabaseID, id)
	var i GetOAuth2ClientByDatabaseIDRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ClientID,
		&i.ClientSecret,
		&i.CreatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}