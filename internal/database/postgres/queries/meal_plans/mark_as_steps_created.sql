UPDATE meal_plans
SET
    tasks_created = 'true',
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $1;
