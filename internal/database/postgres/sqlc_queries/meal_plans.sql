-- name: ArchiveMealPlan :exec

UPDATE meal_plans SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_household = $1 AND id = $2;

-- name: CreateMealPlan :exec

INSERT INTO meal_plans (id,notes,status,voting_deadline,belongs_to_household,created_by_user) VALUES ($1,$2,$3,$4,$5,$6);

-- name: CheckMealPlanExistence :one

SELECT EXISTS ( SELECT meal_plans.id FROM meal_plans WHERE meal_plans.archived_at IS NULL AND meal_plans.id = sqlc.arg(meal_plan_id) AND meal_plans.belongs_to_household = sqlc.arg(household_id) );

-- name: FinalizeMealPlan :exec

UPDATE meal_plans SET status = $1 WHERE archived_at IS NULL AND id = $2;

-- name: GetExpiredAndUnresolvedMealPlans :many

SELECT
	meal_plans.id,
	meal_plans.notes,
	meal_plans.status,
	meal_plans.voting_deadline,
	meal_plans.grocery_list_initialized,
	meal_plans.tasks_created,
	meal_plans.election_method,
	meal_plans.created_at,
	meal_plans.last_updated_at,
	meal_plans.archived_at,
	meal_plans.belongs_to_household,
	meal_plans.created_by_user
FROM meal_plans
WHERE meal_plans.archived_at IS NULL
	AND meal_plans.status = 'awaiting_votes'
	AND voting_deadline < now()
GROUP BY meal_plans.id
ORDER BY meal_plans.id;
-- name: GetFinalizedMealPlansForPlanning :many

SELECT
  meal_plans.id as meal_plan_id,
  meal_plan_options.id as meal_plan_option_id,
  meals.id as meal_id,
  meal_plan_events.id as meal_plan_event_id,
  meal_components.recipe_id as recipe_id
FROM
  meal_plan_options
  FULL OUTER JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event = meal_plan_events.id
  FULL OUTER JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan = meal_plans.id
  FULL OUTER JOIN meal_components ON meal_plan_options.meal_id = meal_components.meal_id
  FULL OUTER JOIN meals ON meal_plan_options.meal_id = meals.id
WHERE
  meal_plans.archived_at IS NULL
  AND meal_plans.status = 'finalized'
  AND meal_plan_options.chosen IS TRUE
  AND meal_plans.tasks_created IS FALSE
GROUP BY
  meal_plans.id,
  meal_plan_options.id,
  meals.id,
  meal_plan_events.id,
  meal_components.recipe_id
ORDER BY
  meal_plans.id;

-- name: GetFinalizedMealPlansWithoutGroceryListInit :many

SELECT
	meal_plans.id,
	meal_plans.belongs_to_household
FROM meal_plans
WHERE meal_plans.archived_at IS NULL
	AND meal_plans.status = 'finalized'
	AND meal_plans.grocery_list_initialized IS FALSE;

-- name: GetMealPlan :one

SELECT
	meal_plans.id,
	meal_plans.notes,
	meal_plans.status,
	meal_plans.voting_deadline,
	meal_plans.grocery_list_initialized,
	meal_plans.tasks_created,
	meal_plans.election_method,
	meal_plans.created_at,
	meal_plans.last_updated_at,
	meal_plans.archived_at,
	meal_plans.belongs_to_household,
	meal_plans.created_by_user
FROM meal_plans
WHERE meal_plans.archived_at IS NULL
	AND meal_plans.id = $1
	AND meal_plans.belongs_to_household = $2;

-- name: GetMealPlanPastVotingDeadline :one

SELECT
	meal_plans.id,
	meal_plans.notes,
	meal_plans.status,
	meal_plans.voting_deadline,
	meal_plans.grocery_list_initialized,
	meal_plans.tasks_created,
	meal_plans.election_method,
	meal_plans.created_at,
	meal_plans.last_updated_at,
	meal_plans.archived_at,
	meal_plans.belongs_to_household,
	meal_plans.created_by_user
FROM meal_plans
WHERE meal_plans.archived_at IS NULL
	AND meal_plans.id = sqlc.arg(meal_plan_id)
	AND meal_plans.belongs_to_household = sqlc.arg(household_id)
	AND meal_plans.status = 'awaiting_votes'
	AND NOW() > meal_plans.voting_deadline;

-- name: MarkMealPlanAsGroceryListInitialized :exec

UPDATE meal_plans
SET
	grocery_list_initialized = 'true',
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $1;

-- name: MarkMealPlanAsPrepTasksCreated :exec

UPDATE meal_plans
SET
	tasks_created = 'true',
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $1;

-- name: UpdateMealPlan :exec

UPDATE meal_plans SET notes = $1, status = $2, voting_deadline = $3, last_updated_at = NOW() WHERE archived_at IS NULL AND belongs_to_household = $4 AND id = $5;