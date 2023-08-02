-- name: UpdateMealPlanEvent :exec

UPDATE meal_plan_events
SET notes = $1,
	starts_at = $2,
	ends_at = $3,
	meal_name = $4,
	belongs_to_meal_plan = $5,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $6;
