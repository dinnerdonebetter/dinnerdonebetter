SELECT
	meal_plan_options.id,
	meal_plan_options.assigned_cook,
	meal_plan_options.assigned_dishwasher,
	meal_plan_options.chosen,
	meal_plan_options.tiebroken,
	meal_plan_options.meal_id,
	meal_plan_options.notes,
	meal_plan_options.created_at,
	meal_plan_options.last_updated_at,
	meal_plan_options.archived_at,
	meal_plan_options.belongs_to_meal_plan_event,
	meals.id,
	meals.name,
	meals.description,
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