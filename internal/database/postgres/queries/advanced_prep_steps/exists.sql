SELECT EXISTS (
	SELECT advanced_prep_steps.id
	FROM advanced_prep_steps
		FULL OUTER JOIN meal_plan_options ON advanced_prep_steps.belongs_to_meal_plan_option=meal_plan_options.id
		FULL OUTER JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
		FULL OUTER JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
	WHERE advanced_prep_steps.settled_at IS NULL
		AND meal_plans.id = $1
		AND meal_plans.archived_at IS NULL
		AND advanced_prep_steps.id = $2
);
