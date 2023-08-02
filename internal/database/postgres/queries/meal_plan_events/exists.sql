-- name: CheckMealPlanEventExistence :one

SELECT EXISTS ( SELECT meal_plan_events.id FROM meal_plan_events WHERE meal_plan_events.archived_at IS NULL AND meal_plan_events.id = $1 );
