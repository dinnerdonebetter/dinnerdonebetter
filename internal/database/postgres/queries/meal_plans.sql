-- name: MealPlanExists :one
SELECT EXISTS ( SELECT meal_plans.id FROM meal_plans WHERE meal_plans.archived_on IS NULL AND meal_plans.id = $1 );

-- name: GetMealPlan :many
SELECT
    meal_plans.id,
    meal_plans.notes,
    meal_plans.status,
    meal_plans.voting_deadline,
    meal_plans.starts_at,
    meal_plans.ends_at,
    meal_plans.created_on,
    meal_plans.last_updated_on,
    meal_plans.archived_on,
    meal_plans.belongs_to_household,
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
    meal_plan_option_votes.id,
    meal_plan_option_votes.rank,
    meal_plan_option_votes.abstain,
    meal_plan_option_votes.notes,
    meal_plan_option_votes.by_user,
    meal_plan_option_votes.created_on,
    meal_plan_option_votes.last_updated_on,
    meal_plan_option_votes.archived_on,
    meal_plan_option_votes.belongs_to_meal_plan_option,
    meals.id,
    meals.name,
    meals.description,
    meals.created_on,
    meals.last_updated_on,
    meals.archived_on,
    meals.created_by_user
FROM meal_plans
         FULL OUTER JOIN meal_plan_options ON meal_plan_options.belongs_to_meal_plan=meal_plans.id
         FULL OUTER JOIN meal_plan_option_votes ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id
         FULL OUTER JOIN meals ON meal_plan_options.meal_id=meals.id
WHERE meal_plans.archived_on IS NULL
  AND meal_plans.id = $1
  AND meal_plans.belongs_to_household = $2
ORDER BY meal_plan_options.id;

-- name: GetMealPlanPastVotingDeadline :many
SELECT
    meal_plans.id,
    meal_plans.notes,
    meal_plans.status,
    meal_plans.voting_deadline,
    meal_plans.starts_at,
    meal_plans.ends_at,
    meal_plans.created_on,
    meal_plans.last_updated_on,
    meal_plans.archived_on,
    meal_plans.belongs_to_household,
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
    meal_plan_option_votes.id,
    meal_plan_option_votes.rank,
    meal_plan_option_votes.abstain,
    meal_plan_option_votes.notes,
    meal_plan_option_votes.by_user,
    meal_plan_option_votes.created_on,
    meal_plan_option_votes.last_updated_on,
    meal_plan_option_votes.archived_on,
    meal_plan_option_votes.belongs_to_meal_plan_option,
    meals.id,
    meals.name,
    meals.description,
    meals.created_on,
    meals.last_updated_on,
    meals.archived_on,
    meals.created_by_user
FROM meal_plans
 FULL OUTER JOIN meal_plan_options ON meal_plan_options.belongs_to_meal_plan=meal_plans.id
 FULL OUTER JOIN meal_plan_option_votes ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id
 FULL OUTER JOIN meals ON meal_plan_options.meal_id=meals.id
WHERE meal_plans.archived_on IS NULL
  AND meal_plans.id = $1
  AND meal_plans.belongs_to_household = $2
  AND meal_plans.status = 'awaiting_votes'
  AND extract(epoch from NOW()) > meal_plans.voting_deadline
ORDER BY meal_plan_options.id;

-- name: GetTotalMealPlansCount :one
SELECT COUNT(meal_plans.id) FROM meal_plans WHERE meal_plans.archived_on IS NULL;

-- name: CreateMealPlan :exec
INSERT INTO meal_plans (id,notes,status,voting_deadline,starts_at,ends_at,belongs_to_household) VALUES ($1,$2,$3,$4,$5,$6,$7);

-- name: UpdateMealPlan :exec
UPDATE meal_plans SET notes = $1, status = $2, voting_deadline = $3, starts_at = $4, ends_at = $5, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_household = $6 AND id = $7;

-- name: ArchiveMealPlan :exec
UPDATE meal_plans SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_household = $1 AND id = $2;

-- name: FinalizeMealPlan :exec
UPDATE meal_plans SET status = $1 WHERE archived_on IS NULL AND id = $2;

-- name: GetExpiredAndUnresolvedMealPlanIDs :many
SELECT
	meal_plans.id,
	meal_plans.notes,
	meal_plans.status,
	meal_plans.voting_deadline,
	meal_plans.starts_at,
	meal_plans.ends_at,
	meal_plans.created_on,
	meal_plans.last_updated_on,
	meal_plans.archived_on,
	meal_plans.belongs_to_household
FROM meal_plans
WHERE meal_plans.archived_on IS NULL
	AND meal_plans.status = 'awaiting_votes'
	AND to_timestamp(voting_deadline)::date < now()
GROUP BY meal_plans.id
ORDER BY meal_plans.id;
