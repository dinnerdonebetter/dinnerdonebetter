-- name: ArchiveMealPlan :exec
UPDATE meal_plans SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_household = $1 AND id = $2;