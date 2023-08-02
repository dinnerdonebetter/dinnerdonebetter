-- name: CheckMealPlanExistence :one

SELECT EXISTS ( SELECT meal_plans.id FROM meal_plans WHERE meal_plans.archived_at IS NULL AND meal_plans.id = $1 );