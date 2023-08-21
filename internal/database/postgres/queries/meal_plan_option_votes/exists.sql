-- name: CheckMealPlanOptionVoteExistence :one

SELECT EXISTS (
	SELECT
	 meal_plan_option_votes.id
	FROM
	 meal_plan_option_votes
		JOIN meal_plan_options ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id
		JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
		JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
	WHERE meal_plan_option_votes.archived_at IS NULL
	AND meal_plan_option_votes.belongs_to_meal_plan_option = sqlc.arg(meal_plan_option_id)
	AND meal_plan_option_votes.id = sqlc.arg(meal_plan_option_vote_id)
	AND meal_plan_options.archived_at IS NULL
	AND meal_plan_options.belongs_to_meal_plan_event = sqlc.arg(meal_plan_event_id)
	AND meal_plan_events.archived_at IS NULL
	AND meal_plan_events.belongs_to_meal_plan = sqlc.arg(meal_plan_id)
	AND meal_plan_options.id = sqlc.arg(meal_plan_option_id)
	AND meal_plans.archived_at IS NULL
	AND meal_plans.id = sqlc.arg(meal_plan_id)
);