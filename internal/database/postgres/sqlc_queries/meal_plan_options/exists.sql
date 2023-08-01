SELECT EXISTS (
	SELECT
	 meal_plan_options.id
	FROM
	 meal_plan_options
		JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event = meal_plan_events.id
		JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan = meal_plans.id
	WHERE
	 meal_plan_options.archived_at IS NULL
	AND meal_plan_options.belongs_to_meal_plan_event = $2
	AND meal_plan_options.id = $3
	AND meal_plan_events.archived_at IS NULL
	AND meal_plan_events.belongs_to_meal_plan = $1
	AND meal_plan_events.id = $2
	AND meal_plans.archived_at IS NULL
	AND meal_plans.id = $1
);
