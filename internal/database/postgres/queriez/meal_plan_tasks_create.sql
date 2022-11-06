-- name: CreateMealPlanTask :exec
INSERT INTO meal_plan_tasks (id,status,status_explanation,creation_explanation,belongs_to_meal_plan_option,belongs_to_recipe_prep_task,assigned_to_user)
VALUES ($1,$2,$3,$4,$5,$6,$7);
