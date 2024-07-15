-- name: ArchiveMealPlan :execrows

UPDATE meal_plans SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_household = sqlc.arg(belongs_to_household) AND id = sqlc.arg(id);
-- name: CreateMealPlan :exec

INSERT INTO meal_plans (
	id,
	notes,
	status,
	voting_deadline,
	belongs_to_household,
	created_by_user
) VALUES (
	sqlc.arg(id),
	sqlc.arg(notes),
	sqlc.arg(status),
	sqlc.arg(voting_deadline),
	sqlc.arg(belongs_to_household),
	sqlc.arg(created_by_user)
);
-- name: CheckMealPlanExistence :one

SELECT EXISTS (
	SELECT meal_plans.id
	FROM meal_plans
	WHERE meal_plans.archived_at IS NULL
		AND meal_plans.id = sqlc.arg(meal_plan_id)
		AND meal_plans.belongs_to_household = sqlc.arg(belongs_to_household)
);
-- name: FinalizeMealPlan :exec

UPDATE meal_plans SET status = sqlc.arg(status) WHERE archived_at IS NULL AND id = sqlc.arg(id);
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
	AND voting_deadline < NOW()
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
	JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event = meal_plan_events.id
	JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan = meal_plans.id
	JOIN meal_components ON meal_plan_options.meal_id = meal_components.meal_id
	JOIN meals ON meal_plan_options.meal_id = meals.id
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
  AND meal_plans.id = sqlc.arg(id)
  AND meal_plans.belongs_to_household = sqlc.arg(belongs_to_household);
-- name: GetMealPlans :many

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
	meal_plans.created_by_user,
	(
		SELECT COUNT(meal_plans.id)
		FROM meal_plans
		WHERE meal_plans.archived_at IS NULL
			AND meal_plans.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND meal_plans.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				meal_plans.last_updated_at IS NULL
				OR meal_plans.last_updated_at > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				meal_plans.last_updated_at IS NULL
				OR meal_plans.last_updated_at < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + '999 years'::INTERVAL))
			)
			AND meal_plans.belongs_to_household = sqlc.arg(household_id)
	) AS filtered_count,
	(
		SELECT COUNT(meal_plans.id)
		FROM meal_plans
		WHERE meal_plans.archived_at IS NULL
			AND meal_plans.belongs_to_household = sqlc.arg(household_id)
	) AS total_count
FROM meal_plans
WHERE meal_plans.archived_at IS NULL
	AND meal_plans.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
	AND meal_plans.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		meal_plans.last_updated_at IS NULL
		OR meal_plans.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		meal_plans.last_updated_at IS NULL
		OR meal_plans.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + '999 years'::INTERVAL))
	)
	AND meal_plans.belongs_to_household = sqlc.arg(household_id)
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);
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

UPDATE meal_plans SET
	grocery_list_initialized = TRUE,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = sqlc.arg(id);
-- name: MarkMealPlanAsPrepTasksCreated :exec

UPDATE meal_plans SET
	tasks_created = TRUE,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = sqlc.arg(id);
-- name: UpdateMealPlan :execrows

UPDATE meal_plans SET
	notes = sqlc.arg(notes),
	status = sqlc.arg(status),
	voting_deadline = sqlc.arg(voting_deadline),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_household = sqlc.arg(belongs_to_household)
	AND id = sqlc.arg(id);
