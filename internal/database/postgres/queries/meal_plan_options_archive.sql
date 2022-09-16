UPDATE
    meal_plan_options
SET
    archived_at = NOW()
WHERE
    archived_at IS NULL
  AND belongs_to_meal_plan_event = $1
  AND id = $2;
