SELECT
  EXISTS (
    SELECT
      meal_plan_events.id
    FROM
      meal_plan_events
      JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan = meal_plans.id
    WHERE
      meal_plan_events.archived_at IS NULL
      AND meal_plans.id = $1
      AND meal_plans.status = 'awaiting_votes'
      AND meal_plans.archived_at IS NULL
      AND meal_plan_events.id = $2
      AND meal_plan_events.archived_at IS NULL
  );