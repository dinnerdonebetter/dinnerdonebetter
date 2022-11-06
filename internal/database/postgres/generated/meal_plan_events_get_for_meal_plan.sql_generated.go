// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: meal_plan_events_get_for_meal_plan.sql

package generated

import (
	"context"
)

const GetMealPlanEventsForMealPlan = `-- name: GetMealPlanEventsForMealPlan :many
SELECT
	meal_plan_events.id,
	meal_plan_events.notes,
	meal_plan_events.starts_at,
	meal_plan_events.ends_at,
	meal_plan_events.meal_name,
	meal_plan_events.belongs_to_meal_plan,
	meal_plan_events.created_at,
	meal_plan_events.last_updated_at,
	meal_plan_events.archived_at
FROM meal_plan_events
WHERE meal_plan_events.archived_at IS NULL
	AND meal_plan_events.belongs_to_meal_plan = $1
`

func (q *Queries) GetMealPlanEventsForMealPlan(ctx context.Context, belongsToMealPlan string) ([]*MealPlanEvents, error) {
	rows, err := q.db.QueryContext(ctx, GetMealPlanEventsForMealPlan, belongsToMealPlan)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*MealPlanEvents{}
	for rows.Next() {
		var i MealPlanEvents
		if err := rows.Scan(
			&i.ID,
			&i.Notes,
			&i.StartsAt,
			&i.EndsAt,
			&i.MealName,
			&i.BelongsToMealPlan,
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
