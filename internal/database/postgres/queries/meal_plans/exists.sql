-- name: CheckMealPlanExistence :one

SELECT EXISTS ( SELECT meal_plans.id FROM meal_plans WHERE meal_plans.archived_at IS NULL AND meal_plans.id = sqlc.arg(meal_plan_id) AND meal_plans.belongs_to_household = sqlc.arg(household_id) );