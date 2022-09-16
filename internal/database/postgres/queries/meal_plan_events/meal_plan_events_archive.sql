UPDATE meal_plan_events SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1 AND belongs_to_meal_plan = $2;
