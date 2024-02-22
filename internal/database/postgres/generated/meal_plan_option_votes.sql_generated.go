// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: meal_plan_option_votes.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const archiveMealPlanOptionVote = `-- name: ArchiveMealPlanOptionVote :execrows

UPDATE meal_plan_option_votes SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_meal_plan_option = $1 AND id = $2
`

type ArchiveMealPlanOptionVoteParams struct {
	BelongsToMealPlanOption string
	ID                      string
}

func (q *Queries) ArchiveMealPlanOptionVote(ctx context.Context, db DBTX, arg *ArchiveMealPlanOptionVoteParams) (int64, error) {
	result, err := db.ExecContext(ctx, archiveMealPlanOptionVote, arg.BelongsToMealPlanOption, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const checkMealPlanOptionVoteExistence = `-- name: CheckMealPlanOptionVoteExistence :one

SELECT EXISTS (
	SELECT meal_plan_option_votes.id
	FROM meal_plan_option_votes
		JOIN meal_plan_options ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id
		JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
		JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
	WHERE meal_plan_option_votes.archived_at IS NULL
		AND meal_plan_option_votes.belongs_to_meal_plan_option = $1
		AND meal_plan_option_votes.id = $2
		AND meal_plan_options.archived_at IS NULL
		AND meal_plan_options.belongs_to_meal_plan_event = $3
		AND meal_plan_events.archived_at IS NULL
		AND meal_plan_events.belongs_to_meal_plan = $4
		AND meal_plan_options.id = $1
		AND meal_plans.archived_at IS NULL
		AND meal_plans.id = $4
)
`

type CheckMealPlanOptionVoteExistenceParams struct {
	MealPlanOptionID     string
	MealPlanOptionVoteID string
	MealPlanID           string
	MealPlanEventID      sql.NullString
}

func (q *Queries) CheckMealPlanOptionVoteExistence(ctx context.Context, db DBTX, arg *CheckMealPlanOptionVoteExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkMealPlanOptionVoteExistence,
		arg.MealPlanOptionID,
		arg.MealPlanOptionVoteID,
		arg.MealPlanEventID,
		arg.MealPlanID,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createMealPlanOptionVote = `-- name: CreateMealPlanOptionVote :exec

INSERT INTO meal_plan_option_votes (
	id,
	rank,
	abstain,
	notes,
	by_user,
	belongs_to_meal_plan_option
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
)
`

type CreateMealPlanOptionVoteParams struct {
	ID                      string
	Notes                   string
	ByUser                  string
	BelongsToMealPlanOption string
	Rank                    int32
	Abstain                 bool
}

func (q *Queries) CreateMealPlanOptionVote(ctx context.Context, db DBTX, arg *CreateMealPlanOptionVoteParams) error {
	_, err := db.ExecContext(ctx, createMealPlanOptionVote,
		arg.ID,
		arg.Rank,
		arg.Abstain,
		arg.Notes,
		arg.ByUser,
		arg.BelongsToMealPlanOption,
	)
	return err
}

const getMealPlanOptionVote = `-- name: GetMealPlanOptionVote :one

SELECT
	meal_plan_option_votes.id,
	meal_plan_option_votes.rank,
	meal_plan_option_votes.abstain,
	meal_plan_option_votes.notes,
	meal_plan_option_votes.by_user,
	meal_plan_option_votes.created_at,
	meal_plan_option_votes.last_updated_at,
	meal_plan_option_votes.archived_at,
	meal_plan_option_votes.belongs_to_meal_plan_option
FROM meal_plan_option_votes
	JOIN meal_plan_options ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id
	JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
	JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
WHERE meal_plan_option_votes.archived_at IS NULL
	AND meal_plan_option_votes.belongs_to_meal_plan_option = $1
	AND meal_plan_option_votes.id = $2
	AND meal_plan_options.archived_at IS NULL
	AND meal_plan_options.belongs_to_meal_plan_event = $3
	AND meal_plan_events.archived_at IS NULL
	AND meal_plan_events.belongs_to_meal_plan = $4
	AND meal_plan_options.id = $1
	AND meal_plans.archived_at IS NULL
	AND meal_plans.id = $4
`

type GetMealPlanOptionVoteParams struct {
	MealPlanOptionID     string
	MealPlanOptionVoteID string
	MealPlanID           string
	MealPlanEventID      sql.NullString
}

func (q *Queries) GetMealPlanOptionVote(ctx context.Context, db DBTX, arg *GetMealPlanOptionVoteParams) (*MealPlanOptionVotes, error) {
	row := db.QueryRowContext(ctx, getMealPlanOptionVote,
		arg.MealPlanOptionID,
		arg.MealPlanOptionVoteID,
		arg.MealPlanEventID,
		arg.MealPlanID,
	)
	var i MealPlanOptionVotes
	err := row.Scan(
		&i.ID,
		&i.Rank,
		&i.Abstain,
		&i.Notes,
		&i.ByUser,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
		&i.BelongsToMealPlanOption,
	)
	return &i, err
}

const getMealPlanOptionVotes = `-- name: GetMealPlanOptionVotes :many

SELECT
	meal_plan_option_votes.id,
	meal_plan_option_votes.rank,
	meal_plan_option_votes.abstain,
	meal_plan_option_votes.notes,
	meal_plan_option_votes.by_user,
	meal_plan_option_votes.created_at,
	meal_plan_option_votes.last_updated_at,
	meal_plan_option_votes.archived_at,
	meal_plan_option_votes.belongs_to_meal_plan_option,
	(
		SELECT COUNT(meal_plan_option_votes.id)
		FROM meal_plan_option_votes
		WHERE meal_plan_option_votes.archived_at IS NULL
			AND meal_plan_option_votes.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND meal_plan_option_votes.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				meal_plan_option_votes.last_updated_at IS NULL
				OR meal_plan_option_votes.last_updated_at > COALESCE($3, (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				meal_plan_option_votes.last_updated_at IS NULL
				OR meal_plan_option_votes.last_updated_at < COALESCE($4, (SELECT NOW() + '999 years'::INTERVAL))
			)
			AND meal_plan_option_votes.belongs_to_meal_plan_option = $5
	) AS filtered_count,
	(
		SELECT COUNT(meal_plan_option_votes.id)
		FROM meal_plan_option_votes
		WHERE meal_plan_option_votes.archived_at IS NULL
	) AS total_count
FROM meal_plan_option_votes
	JOIN meal_plan_options ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id
	JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
	JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
WHERE meal_plan_option_votes.archived_at IS NULL
	AND meal_plan_option_votes.belongs_to_meal_plan_option = $5
	AND meal_plan_options.archived_at IS NULL
	AND meal_plan_options.belongs_to_meal_plan_event = $6
	AND meal_plan_options.id = $5
	AND meal_plan_events.archived_at IS NULL
	AND meal_plan_events.belongs_to_meal_plan = $7
	AND meal_plan_events.id = $6
	AND meal_plans.archived_at IS NULL
	AND meal_plans.id = $7
	AND meal_plan_option_votes.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND meal_plan_option_votes.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		meal_plan_option_votes.last_updated_at IS NULL
		OR meal_plan_option_votes.last_updated_at > COALESCE($4, (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		meal_plan_option_votes.last_updated_at IS NULL
		OR meal_plan_option_votes.last_updated_at < COALESCE($3, (SELECT NOW() + '999 years'::INTERVAL))
	)
GROUP BY
	meal_plan_option_votes.id,
	meal_plan_options.id,
	meal_plan_events.id,
	meal_plans.id
LIMIT $9
OFFSET $8
`

type GetMealPlanOptionVotesParams struct {
	CreatedAfter     sql.NullTime
	CreatedBefore    sql.NullTime
	UpdatedBefore    sql.NullTime
	UpdatedAfter     sql.NullTime
	MealPlanOptionID string
	MealPlanID       string
	MealPlanEventID  sql.NullString
	QueryOffset      sql.NullInt32
	QueryLimit       sql.NullInt32
}

type GetMealPlanOptionVotesRow struct {
	CreatedAt               time.Time
	LastUpdatedAt           sql.NullTime
	ArchivedAt              sql.NullTime
	ID                      string
	Notes                   string
	ByUser                  string
	BelongsToMealPlanOption string
	FilteredCount           int64
	TotalCount              int64
	Rank                    int32
	Abstain                 bool
}

func (q *Queries) GetMealPlanOptionVotes(ctx context.Context, db DBTX, arg *GetMealPlanOptionVotesParams) ([]*GetMealPlanOptionVotesRow, error) {
	rows, err := db.QueryContext(ctx, getMealPlanOptionVotes,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.UpdatedBefore,
		arg.UpdatedAfter,
		arg.MealPlanOptionID,
		arg.MealPlanEventID,
		arg.MealPlanID,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetMealPlanOptionVotesRow{}
	for rows.Next() {
		var i GetMealPlanOptionVotesRow
		if err := rows.Scan(
			&i.ID,
			&i.Rank,
			&i.Abstain,
			&i.Notes,
			&i.ByUser,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.BelongsToMealPlanOption,
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

const getMealPlanOptionVotesForMealPlanOption = `-- name: GetMealPlanOptionVotesForMealPlanOption :many

SELECT
	meal_plan_option_votes.id,
	meal_plan_option_votes.rank,
	meal_plan_option_votes.abstain,
	meal_plan_option_votes.notes,
	meal_plan_option_votes.by_user,
	meal_plan_option_votes.created_at,
	meal_plan_option_votes.last_updated_at,
	meal_plan_option_votes.archived_at,
	meal_plan_option_votes.belongs_to_meal_plan_option
FROM meal_plan_option_votes
	JOIN meal_plan_options ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id
	JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
	JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
WHERE meal_plan_option_votes.archived_at IS NULL
	AND meal_plan_option_votes.belongs_to_meal_plan_option = $1
	AND meal_plan_options.archived_at IS NULL
	AND meal_plan_options.belongs_to_meal_plan_event = $2
	AND meal_plan_options.id = $1
	AND meal_plan_events.archived_at IS NULL
	AND meal_plan_events.belongs_to_meal_plan = $3
	AND meal_plan_events.id = $2
	AND meal_plans.archived_at IS NULL
	AND meal_plans.id = $3
`

type GetMealPlanOptionVotesForMealPlanOptionParams struct {
	MealPlanOptionID string
	MealPlanID       string
	MealPlanEventID  sql.NullString
}

func (q *Queries) GetMealPlanOptionVotesForMealPlanOption(ctx context.Context, db DBTX, arg *GetMealPlanOptionVotesForMealPlanOptionParams) ([]*MealPlanOptionVotes, error) {
	rows, err := db.QueryContext(ctx, getMealPlanOptionVotesForMealPlanOption, arg.MealPlanOptionID, arg.MealPlanEventID, arg.MealPlanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*MealPlanOptionVotes{}
	for rows.Next() {
		var i MealPlanOptionVotes
		if err := rows.Scan(
			&i.ID,
			&i.Rank,
			&i.Abstain,
			&i.Notes,
			&i.ByUser,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.BelongsToMealPlanOption,
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

const updateMealPlanOptionVote = `-- name: UpdateMealPlanOptionVote :execrows

UPDATE meal_plan_option_votes SET
	rank = $1,
	abstain = $2,
	notes = $3,
	by_user = $4,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_meal_plan_option = $5
	AND id = $6
`

type UpdateMealPlanOptionVoteParams struct {
	Notes                   string
	ByUser                  string
	BelongsToMealPlanOption string
	ID                      string
	Rank                    int32
	Abstain                 bool
}

func (q *Queries) UpdateMealPlanOptionVote(ctx context.Context, db DBTX, arg *UpdateMealPlanOptionVoteParams) (int64, error) {
	result, err := db.ExecContext(ctx, updateMealPlanOptionVote,
		arg.Rank,
		arg.Abstain,
		arg.Notes,
		arg.ByUser,
		arg.BelongsToMealPlanOption,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
