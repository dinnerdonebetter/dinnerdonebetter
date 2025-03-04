// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: audit_logs.sql

package generated

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

const createAuditLogEntry = `-- name: CreateAuditLogEntry :exec
INSERT INTO audit_log_entries (
	id,
	resource_type,
	relevant_id,
	event_type,
	changes,
	belongs_to_user,
	belongs_to_household
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7
)
`

type CreateAuditLogEntryParams struct {
	ID                 string
	ResourceType       string
	RelevantID         string
	EventType          AuditLogEventType
	Changes            json.RawMessage
	BelongsToUser      sql.NullString
	BelongsToHousehold sql.NullString
}

func (q *Queries) CreateAuditLogEntry(ctx context.Context, db DBTX, arg *CreateAuditLogEntryParams) error {
	_, err := db.ExecContext(ctx, createAuditLogEntry,
		arg.ID,
		arg.ResourceType,
		arg.RelevantID,
		arg.EventType,
		arg.Changes,
		arg.BelongsToUser,
		arg.BelongsToHousehold,
	)
	return err
}

const getAuditLogEntriesForHousehold = `-- name: GetAuditLogEntriesForHousehold :many
SELECT
	audit_log_entries.id,
	audit_log_entries.resource_type,
	audit_log_entries.relevant_id,
	audit_log_entries.event_type,
	audit_log_entries.changes,
	audit_log_entries.belongs_to_user,
	audit_log_entries.belongs_to_household,
	audit_log_entries.created_at,
	(
		SELECT COUNT(audit_log_entries.id)
		FROM audit_log_entries
		WHERE
			audit_log_entries.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND audit_log_entries.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND audit_log_entries.belongs_to_household = $3
	) AS filtered_count,
	(
		SELECT COUNT(audit_log_entries.id)
		FROM audit_log_entries
		WHERE
			audit_log_entries.belongs_to_household = $3
	) AS total_count
FROM audit_log_entries
WHERE audit_log_entries.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND audit_log_entries.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND audit_log_entries.belongs_to_household = $3
LIMIT $5
OFFSET $4
`

type GetAuditLogEntriesForHouseholdParams struct {
	CreatedAfter       sql.NullTime
	CreatedBefore      sql.NullTime
	BelongsToHousehold sql.NullString
	QueryOffset        sql.NullInt32
	QueryLimit         sql.NullInt32
}

type GetAuditLogEntriesForHouseholdRow struct {
	CreatedAt          time.Time
	ID                 string
	ResourceType       string
	RelevantID         string
	EventType          AuditLogEventType
	Changes            json.RawMessage
	BelongsToUser      sql.NullString
	BelongsToHousehold sql.NullString
	FilteredCount      int64
	TotalCount         int64
}

func (q *Queries) GetAuditLogEntriesForHousehold(ctx context.Context, db DBTX, arg *GetAuditLogEntriesForHouseholdParams) ([]*GetAuditLogEntriesForHouseholdRow, error) {
	rows, err := db.QueryContext(ctx, getAuditLogEntriesForHousehold,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.BelongsToHousehold,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetAuditLogEntriesForHouseholdRow{}
	for rows.Next() {
		var i GetAuditLogEntriesForHouseholdRow
		if err := rows.Scan(
			&i.ID,
			&i.ResourceType,
			&i.RelevantID,
			&i.EventType,
			&i.Changes,
			&i.BelongsToUser,
			&i.BelongsToHousehold,
			&i.CreatedAt,
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

const getAuditLogEntriesForHouseholdAndResourceType = `-- name: GetAuditLogEntriesForHouseholdAndResourceType :many
SELECT
	audit_log_entries.id,
	audit_log_entries.resource_type,
	audit_log_entries.relevant_id,
	audit_log_entries.event_type,
	audit_log_entries.changes,
	audit_log_entries.belongs_to_user,
	audit_log_entries.belongs_to_household,
	audit_log_entries.created_at,
	(
		SELECT COUNT(audit_log_entries.id)
		FROM audit_log_entries
		WHERE
			audit_log_entries.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND audit_log_entries.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND audit_log_entries.belongs_to_household = $3
			AND audit_log_entries.resource_type = ANY($4::text[])
	) AS filtered_count,
	(
		SELECT COUNT(audit_log_entries.id)
		FROM audit_log_entries
		WHERE
			audit_log_entries.belongs_to_household = $3
			AND audit_log_entries.resource_type = ANY($4::text[])
	) AS total_count
FROM audit_log_entries
WHERE audit_log_entries.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND audit_log_entries.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND audit_log_entries.belongs_to_household = $3
	AND audit_log_entries.resource_type = ANY($4::text[])
LIMIT $6
OFFSET $5
`

type GetAuditLogEntriesForHouseholdAndResourceTypeParams struct {
	CreatedAfter       sql.NullTime
	CreatedBefore      sql.NullTime
	BelongsToHousehold sql.NullString
	Resources          []string
	QueryOffset        sql.NullInt32
	QueryLimit         sql.NullInt32
}

type GetAuditLogEntriesForHouseholdAndResourceTypeRow struct {
	CreatedAt          time.Time
	ID                 string
	ResourceType       string
	RelevantID         string
	EventType          AuditLogEventType
	Changes            json.RawMessage
	BelongsToUser      sql.NullString
	BelongsToHousehold sql.NullString
	FilteredCount      int64
	TotalCount         int64
}

func (q *Queries) GetAuditLogEntriesForHouseholdAndResourceType(ctx context.Context, db DBTX, arg *GetAuditLogEntriesForHouseholdAndResourceTypeParams) ([]*GetAuditLogEntriesForHouseholdAndResourceTypeRow, error) {
	rows, err := db.QueryContext(ctx, getAuditLogEntriesForHouseholdAndResourceType,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.BelongsToHousehold,
		pq.Array(arg.Resources),
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetAuditLogEntriesForHouseholdAndResourceTypeRow{}
	for rows.Next() {
		var i GetAuditLogEntriesForHouseholdAndResourceTypeRow
		if err := rows.Scan(
			&i.ID,
			&i.ResourceType,
			&i.RelevantID,
			&i.EventType,
			&i.Changes,
			&i.BelongsToUser,
			&i.BelongsToHousehold,
			&i.CreatedAt,
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

const getAuditLogEntriesForUser = `-- name: GetAuditLogEntriesForUser :many
SELECT
	audit_log_entries.id,
	audit_log_entries.resource_type,
	audit_log_entries.relevant_id,
	audit_log_entries.event_type,
	audit_log_entries.changes,
	audit_log_entries.belongs_to_user,
	audit_log_entries.belongs_to_household,
	audit_log_entries.created_at,
	(
		SELECT COUNT(audit_log_entries.id)
		FROM audit_log_entries
		WHERE
			audit_log_entries.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND audit_log_entries.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND audit_log_entries.belongs_to_user = $3
	) AS filtered_count,
	(
		SELECT COUNT(audit_log_entries.id)
		FROM audit_log_entries
		WHERE
			audit_log_entries.belongs_to_user = $3
	) AS total_count
FROM audit_log_entries
WHERE audit_log_entries.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND audit_log_entries.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND audit_log_entries.belongs_to_user = $3
LIMIT $5
OFFSET $4
`

type GetAuditLogEntriesForUserParams struct {
	CreatedAfter  sql.NullTime
	CreatedBefore sql.NullTime
	BelongsToUser sql.NullString
	QueryOffset   sql.NullInt32
	QueryLimit    sql.NullInt32
}

type GetAuditLogEntriesForUserRow struct {
	CreatedAt          time.Time
	ID                 string
	ResourceType       string
	RelevantID         string
	EventType          AuditLogEventType
	Changes            json.RawMessage
	BelongsToUser      sql.NullString
	BelongsToHousehold sql.NullString
	FilteredCount      int64
	TotalCount         int64
}

func (q *Queries) GetAuditLogEntriesForUser(ctx context.Context, db DBTX, arg *GetAuditLogEntriesForUserParams) ([]*GetAuditLogEntriesForUserRow, error) {
	rows, err := db.QueryContext(ctx, getAuditLogEntriesForUser,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.BelongsToUser,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetAuditLogEntriesForUserRow{}
	for rows.Next() {
		var i GetAuditLogEntriesForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.ResourceType,
			&i.RelevantID,
			&i.EventType,
			&i.Changes,
			&i.BelongsToUser,
			&i.BelongsToHousehold,
			&i.CreatedAt,
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

const getAuditLogEntriesForUserAndResourceType = `-- name: GetAuditLogEntriesForUserAndResourceType :many
SELECT
	audit_log_entries.id,
	audit_log_entries.resource_type,
	audit_log_entries.relevant_id,
	audit_log_entries.event_type,
	audit_log_entries.changes,
	audit_log_entries.belongs_to_user,
	audit_log_entries.belongs_to_household,
	audit_log_entries.created_at,
	(
		SELECT COUNT(audit_log_entries.id)
		FROM audit_log_entries
		WHERE
			audit_log_entries.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND audit_log_entries.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND audit_log_entries.belongs_to_user = $3
			AND audit_log_entries.resource_type = ANY($4::text[])
	) AS filtered_count,
	(
		SELECT COUNT(audit_log_entries.id)
		FROM audit_log_entries
		WHERE
			audit_log_entries.belongs_to_user = $3
			AND audit_log_entries.resource_type = ANY($4::text[])
	) AS total_count
FROM audit_log_entries
WHERE audit_log_entries.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND audit_log_entries.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND audit_log_entries.belongs_to_user = $3
	AND audit_log_entries.resource_type = ANY($4::text[])
LIMIT $6
OFFSET $5
`

type GetAuditLogEntriesForUserAndResourceTypeParams struct {
	CreatedAfter  sql.NullTime
	CreatedBefore sql.NullTime
	BelongsToUser sql.NullString
	Resources     []string
	QueryOffset   sql.NullInt32
	QueryLimit    sql.NullInt32
}

type GetAuditLogEntriesForUserAndResourceTypeRow struct {
	CreatedAt          time.Time
	ID                 string
	ResourceType       string
	RelevantID         string
	EventType          AuditLogEventType
	Changes            json.RawMessage
	BelongsToUser      sql.NullString
	BelongsToHousehold sql.NullString
	FilteredCount      int64
	TotalCount         int64
}

func (q *Queries) GetAuditLogEntriesForUserAndResourceType(ctx context.Context, db DBTX, arg *GetAuditLogEntriesForUserAndResourceTypeParams) ([]*GetAuditLogEntriesForUserAndResourceTypeRow, error) {
	rows, err := db.QueryContext(ctx, getAuditLogEntriesForUserAndResourceType,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.BelongsToUser,
		pq.Array(arg.Resources),
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetAuditLogEntriesForUserAndResourceTypeRow{}
	for rows.Next() {
		var i GetAuditLogEntriesForUserAndResourceTypeRow
		if err := rows.Scan(
			&i.ID,
			&i.ResourceType,
			&i.RelevantID,
			&i.EventType,
			&i.Changes,
			&i.BelongsToUser,
			&i.BelongsToHousehold,
			&i.CreatedAt,
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

const getAuditLogEntry = `-- name: GetAuditLogEntry :one
SELECT
	audit_log_entries.id,
	audit_log_entries.resource_type,
	audit_log_entries.relevant_id,
	audit_log_entries.event_type,
	audit_log_entries.changes,
	audit_log_entries.belongs_to_user,
	audit_log_entries.belongs_to_household,
	audit_log_entries.created_at
FROM audit_log_entries
WHERE audit_log_entries.id = $1
`

type GetAuditLogEntryRow struct {
	CreatedAt          time.Time
	ID                 string
	ResourceType       string
	RelevantID         string
	EventType          AuditLogEventType
	Changes            json.RawMessage
	BelongsToUser      sql.NullString
	BelongsToHousehold sql.NullString
}

func (q *Queries) GetAuditLogEntry(ctx context.Context, db DBTX, id string) (*GetAuditLogEntryRow, error) {
	row := db.QueryRowContext(ctx, getAuditLogEntry, id)
	var i GetAuditLogEntryRow
	err := row.Scan(
		&i.ID,
		&i.ResourceType,
		&i.RelevantID,
		&i.EventType,
		&i.Changes,
		&i.BelongsToUser,
		&i.BelongsToHousehold,
		&i.CreatedAt,
	)
	return &i, err
}
