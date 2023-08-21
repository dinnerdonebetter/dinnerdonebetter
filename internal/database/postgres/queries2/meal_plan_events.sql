-- name: ArchiveMealPlanEvent :exec

UPDATE meal_plan_events SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1 AND belongs_to_meal_plan = $2;


-- name: CreateMealPlanEvent :exec

INSERT INTO
	meal_plan_events (id, notes, starts_at, ends_at, meal_name, belongs_to_meal_plan)
VALUES
	($1, $2, $3, $4, $5, $6);


-- name: MealPlanEventIsEligibleForVoting :one

SELECT
  EXISTS (
    SELECT
      meal_plan_events.id
    FROM
      meal_plan_events
      JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan = meal_plans.id
    WHERE
      meal_plan_events.archived_at IS NULL
      AND meal_plans.id = sqlc.arg(meal_plan_id)
      AND meal_plans.status = 'awaiting_votes'
      AND meal_plans.archived_at IS NULL
      AND meal_plan_events.id = sqlc.arg(meal_plan_event_id)
      AND meal_plan_events.archived_at IS NULL
  );

-- name: CheckMealPlanEventExistence :one

SELECT EXISTS ( SELECT meal_plan_events.id FROM meal_plan_events WHERE meal_plan_events.archived_at IS NULL AND meal_plan_events.id = sqlc.arg(meal_plan_event_id) AND meal_plan_events.belongs_to_meal_plan = sqlc.arg(meal_plan_id));


-- name: GetMealPlanEventsForMealPlan :many

SELECT
	meal_plan_events.id,
	meal_plan_events.notes,
	meal_plan_events.starts_at,
	meal_plan_events.ends_at,
	meal_plan_events.meal_name,
	meal_plan_events.belongs_to_meal_plan,
	meal_plan_events.created_at,
	meal_plan_events.last_updated_at,
	meal_plan_events.archived_at
FROM meal_plan_events
WHERE meal_plan_events.archived_at IS NULL
	AND meal_plan_events.belongs_to_meal_plan = $1;


-- name: GetMealPlanEvent :one

SELECT
	meal_plan_events.id,
	meal_plan_events.notes,
	meal_plan_events.starts_at,
	meal_plan_events.ends_at,
	meal_plan_events.meal_name,
	meal_plan_events.belongs_to_meal_plan,
	meal_plan_events.created_at,
	meal_plan_events.last_updated_at,
	meal_plan_events.archived_at
FROM meal_plan_events
WHERE meal_plan_events.archived_at IS NULL
	AND meal_plan_events.id = $1
	AND meal_plan_events.belongs_to_meal_plan = $2;


-- name: UpdateMealPlanEvent :exec

UPDATE meal_plan_events
SET notes = $1,
	starts_at = $2,
	ends_at = $3,
	meal_name = $4,
	belongs_to_meal_plan = $5,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $6;
