UPDATE meal_plan_options
SET
	assigned_cook = $1,
	assigned_dishwasher = $2,
	meal_id = $3,
	notes = $4,
	prep_steps_created = $5,
	last_updated_at = NOW()
WHERE archived_at IS NULL
  AND belongs_to_meal_plan_event = $6
  AND id = $7;
