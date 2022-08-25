-- name: MealPlanOptionExists :one
SELECT EXISTS ( SELECT meal_plan_options.id FROM meal_plan_options JOIN meal_plans ON meal_plan_options.belongs_to_meal_plan=meal_plans.id WHERE meal_plan_options.archived_on IS NULL AND meal_plan_options.belongs_to_meal_plan = $1 AND meal_plan_options.id = $2 AND meal_plans.archived_on IS NULL AND meal_plans.id = $3 );

-- name: GetMealPlanOptionQuery :one
SELECT
	meal_plan_options.id,
	meal_plan_options.day,
	meal_plan_options.meal_name,
	meal_plan_options.chosen,
	meal_plan_options.tiebroken,
	meal_plan_options.meal_id,
	meal_plan_options.notes,
	meal_plan_options.created_on,
	meal_plan_options.last_updated_on,
	meal_plan_options.archived_on,
	meal_plan_options.belongs_to_meal_plan,
	meals.id,
	meals.name,
	meals.description,
	meals.created_on,
	meals.last_updated_on,
	meals.archived_on,
	meals.created_by_user
FROM meal_plan_options
JOIN meal_plans ON meal_plan_options.belongs_to_meal_plan=meal_plans.id
JOIN meals ON meal_plan_options.meal_id=meals.id
WHERE meal_plan_options.archived_on IS NULL
AND meal_plan_options.belongs_to_meal_plan = $1
AND meal_plan_options.id = $2
AND meal_plans.archived_on IS NULL
AND meal_plans.id = $3;

-- name: GetTotalMealPlanOptionsCount :one
SELECT COUNT(meal_plan_options.id) FROM meal_plan_options WHERE meal_plan_options.archived_on IS NULL;

-- name: MealPlanOptionCreation :exec
INSERT INTO meal_plan_options (id,day,meal_name,meal_id,notes,belongs_to_meal_plan) VALUES ($1,$2,$3,$4,$5,$6);

-- name: UpdateMealPlanOption :exec
UPDATE meal_plan_options SET day = $1, meal_id = $2, meal_name = $3, notes = $4, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_meal_plan = $5 AND id = $6;

-- name: ArchiveMealPlanOption :exec
UPDATE meal_plan_options SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_meal_plan = $1 AND id = $2;

-- name: FinalizeMealPlanOption :exec
UPDATE meal_plan_options SET chosen = (belongs_to_meal_plan = $1 AND id = $2), tiebroken = $3 WHERE archived_on IS NULL AND belongs_to_meal_plan = $1 AND id = $2;
