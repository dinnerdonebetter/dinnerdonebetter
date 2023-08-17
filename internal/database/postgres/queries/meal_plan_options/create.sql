-- name: CreateMealPlanOption :exec

INSERT INTO meal_plan_options (id,assigned_cook,assigned_dishwasher,meal_id,notes,meal_scale,belongs_to_meal_plan_event,chosen)
VALUES (
    $1, -- sqlc.arg(id),
    $2, -- sqlc.arg(assigned_cook),
    $3, -- sqlc.arg(assigned_dishwasher),
    $4, -- sqlc.arg(meal_id),
    $5, -- sqlc.arg(notes),
    $6, -- sqlc.arg(meal_scale)::float,
    $7, -- sqlc.arg(belongs_to_meal_plan_event),
    $8  -- sqlc.arg(chosen)::bool
);
