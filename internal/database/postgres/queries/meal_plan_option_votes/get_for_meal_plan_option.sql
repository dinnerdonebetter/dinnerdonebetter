SELECT
	meal_plan_option_votes.id,
	meal_plan_option_votes.rank,
	meal_plan_option_votes.abstain,
	meal_plan_option_votes.notes,
	meal_plan_option_votes.by_user,
	meal_plan_option_votes.created_at,
	meal_plan_option_votes.last_updated_at,
	meal_plan_option_votes.archived_at,
	meal_plan_option_votes.belongs_to_meal_plan_option
FROM meal_plan_option_votes
	     JOIN meal_plan_options ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id
	     JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
	     JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
WHERE meal_plan_option_votes.archived_at IS NULL
  AND meal_plan_option_votes.belongs_to_meal_plan_option = $3
  AND meal_plan_options.archived_at IS NULL
  AND meal_plan_options.belongs_to_meal_plan_event = $2
  AND meal_plan_options.id = $3
  AND meal_plan_events.archived_at IS NULL
  AND meal_plan_events.belongs_to_meal_plan = $1
  AND meal_plan_events.id = $2
  AND meal_plans.archived_at IS NULL
  AND meal_plans.id = $1;
