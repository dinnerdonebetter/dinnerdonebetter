-- name: ArchiveMealPlanOption :exec

UPDATE
	meal_plan_options
SET
	archived_at = NOW()
WHERE
	archived_at IS NULL
	AND belongs_to_meal_plan_event = $1
	AND id = $2;

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

-- name: CheckMealPlanOptionExistence :one

SELECT EXISTS (
	SELECT
	 meal_plan_options.id
	FROM
	 meal_plan_options
		JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event = meal_plan_events.id
		JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan = meal_plans.id
	WHERE
	 meal_plan_options.archived_at IS NULL
	AND meal_plan_options.belongs_to_meal_plan_event = sqlc.arg(meal_plan_event_id)
	AND meal_plan_options.id = sqlc.arg(meal_plan_option_id)
	AND meal_plan_events.archived_at IS NULL
	AND meal_plan_events.belongs_to_meal_plan = sqlc.arg(meal_plan_id)
	AND meal_plan_events.id = sqlc.arg(meal_plan_event_id)
	AND meal_plans.archived_at IS NULL
	AND meal_plans.id = sqlc.arg(meal_plan_id)
);

-- name: FinalizeMealPlanOption :exec

UPDATE meal_plan_options SET chosen = (belongs_to_meal_plan_event = $1 AND id = $2), tiebroken = $3 WHERE archived_at IS NULL AND belongs_to_meal_plan_event = $1 AND id = $2;

-- name: GetMealPlanOptionsForMealPlanEvent :many

SELECT
	meal_plan_options.id,
	meal_plan_options.assigned_cook,
	meal_plan_options.assigned_dishwasher,
	meal_plan_options.chosen,
	meal_plan_options.tiebroken,
    meal_plan_options.meal_scale,
	meal_plan_options.meal_id,
	meal_plan_options.notes,
	meal_plan_options.created_at,
	meal_plan_options.last_updated_at,
	meal_plan_options.archived_at,
	meal_plan_options.belongs_to_meal_plan_event,
	meals.id,
	meals.name,
	meals.description,
	meals.min_estimated_portions,
	meals.max_estimated_portions,
    meals.eligible_for_meal_plans,
	meals.created_at,
	meals.last_updated_at,
	meals.archived_at,
	meals.created_by_user
FROM meal_plan_options
	JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event = meal_plan_events.id
	JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan = meal_plans.id
	JOIN meals ON meal_plan_options.meal_id = meals.id
WHERE
	meal_plan_options.archived_at IS NULL
	AND meal_plan_options.belongs_to_meal_plan_event = $1
	AND meal_plan_events.id = $1
	AND meal_plan_events.belongs_to_meal_plan = $2
	AND meal_plans.archived_at IS NULL
	AND meal_plans.id = $2;

-- name: GetMealPlanOption :one

SELECT
	meal_plan_options.id,
	meal_plan_options.assigned_cook,
	meal_plan_options.assigned_dishwasher,
	meal_plan_options.chosen,
	meal_plan_options.tiebroken,
    meal_plan_options.meal_scale,
	meal_plan_options.meal_id,
	meal_plan_options.notes,
	meal_plan_options.created_at,
	meal_plan_options.last_updated_at,
	meal_plan_options.archived_at,
	meal_plan_options.belongs_to_meal_plan_event,
	meals.id,
	meals.name,
	meals.description,
	meals.min_estimated_portions,
	meals.max_estimated_portions,
    meals.eligible_for_meal_plans,
	meals.created_at,
	meals.last_updated_at,
	meals.archived_at,
	meals.created_by_user
FROM meal_plan_options
	JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event = meal_plan_events.id
	JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan = meal_plans.id
	JOIN meals ON meal_plan_options.meal_id = meals.id
WHERE meal_plan_options.archived_at IS NULL
	AND meal_plan_options.belongs_to_meal_plan_event = $2
	AND meal_plan_options.id = $3
	AND meal_plan_events.id = $2
	AND meal_plan_events.belongs_to_meal_plan = $1
	AND meal_plans.archived_at IS NULL
	AND meal_plans.id = $1;

-- name: GetMealPlanOptionByID :one

SELECT
	meal_plan_options.id,
	meal_plan_options.assigned_cook,
	meal_plan_options.assigned_dishwasher,
	meal_plan_options.chosen,
	meal_plan_options.tiebroken,
    meal_plan_options.meal_scale,
	meal_plan_options.meal_id,
	meal_plan_options.notes,
	meal_plan_options.created_at,
	meal_plan_options.last_updated_at,
	meal_plan_options.archived_at,
	meal_plan_options.belongs_to_meal_plan_event,
	meals.id,
	meals.name,
	meals.description,
	meals.min_estimated_portions,
	meals.max_estimated_portions,
    meals.eligible_for_meal_plans,
	meals.created_at,
	meals.last_updated_at,
	meals.archived_at,
	meals.created_by_user
FROM meal_plan_options
	JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event = meal_plan_events.id
	JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan = meal_plans.id
	JOIN meals ON meal_plan_options.meal_id = meals.id
WHERE meal_plan_options.archived_at IS NULL
	AND meal_plan_options.id = $1;

-- name: UpdateMealPlanOption :exec

UPDATE meal_plan_options
SET
	assigned_cook = $1,
	assigned_dishwasher = $2,
	meal_id = $3,
	notes = $4,
	meal_scale = $5,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_meal_plan_event = $6
	AND id = $7;
