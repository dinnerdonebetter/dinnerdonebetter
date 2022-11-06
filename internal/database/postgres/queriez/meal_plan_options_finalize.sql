-- name: FinalizeMealPlanOption :exec
UPDATE meal_plan_options SET chosen = (belongs_to_meal_plan_event = $1 AND id = $2), tiebroken = $3 WHERE archived_at IS NULL AND belongs_to_meal_plan_event = $1 AND id = $2;
