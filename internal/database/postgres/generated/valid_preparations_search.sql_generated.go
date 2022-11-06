// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: valid_preparations_search.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const SearchForValidPreparations = `-- name: SearchForValidPreparations :many
SELECT
	valid_preparations.id,
	valid_preparations.name,
	valid_preparations.description,
	valid_preparations.icon_path,
	valid_preparations.yields_nothing,
	valid_preparations.restrict_to_ingredients,
	valid_preparations.zero_ingredients_allowable,
	valid_preparations.past_tense,
	valid_preparations.created_at,
	valid_preparations.last_updated_at,
	valid_preparations.archived_at
FROM valid_preparations
WHERE valid_preparations.archived_at IS NULL
	AND valid_preparations.name ILIKE $1
LIMIT 50
`

type SearchForValidPreparationsRow struct {
	ID                       string
	Name                     string
	Description              string
	IconPath                 string
	YieldsNothing            bool
	RestrictToIngredients    bool
	ZeroIngredientsAllowable bool
	PastTense                string
	CreatedAt                time.Time
	LastUpdatedAt            sql.NullTime
	ArchivedAt               sql.NullTime
}

func (q *Queries) SearchForValidPreparations(ctx context.Context, name string) ([]*SearchForValidPreparationsRow, error) {
	rows, err := q.db.QueryContext(ctx, SearchForValidPreparations, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidPreparationsRow{}
	for rows.Next() {
		var i SearchForValidPreparationsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.IconPath,
			&i.YieldsNothing,
			&i.RestrictToIngredients,
			&i.ZeroIngredientsAllowable,
			&i.PastTense,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
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
