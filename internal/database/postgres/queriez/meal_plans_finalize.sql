-- name: FinalizeMealPlan :exec
UPDATE meal_plans SET status = $1 WHERE archived_at IS NULL AND id = $2;