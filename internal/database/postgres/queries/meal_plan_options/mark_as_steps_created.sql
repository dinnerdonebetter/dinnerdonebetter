UPDATE meal_plan_options
SET
    prep_steps_created = 'true',
    last_updated_at = NOW()
WHERE archived_at IS NULL
  AND id = $1;
