-- name: CreateMealPlanEvent :exec
INSERT INTO
	meal_plan_events (id, notes, starts_at, ends_at, meal_name, belongs_to_meal_plan)
VALUES
	($1, $2, $3, $4, $5, $6);
