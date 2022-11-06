-- name: GetFinalizedMealPlansForPlanning :exec
SELECT
	meal_plans.id as meal_plan_id,
	meal_plan_options.id as meal_plan_option_id,
	meals.id as meal_id,
	meal_plan_events.id as meal_plan_event_id,
	meal_recipes.recipe_id as recipe_id
FROM meal_plan_options
	FULL OUTER JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
	FULL OUTER JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
	FULL OUTER JOIN meal_recipes ON meal_plan_options.meal_id=meal_recipes.meal_id
	FULL OUTER JOIN meals ON meal_plan_options.meal_id=meals.id
WHERE meal_plans.archived_at IS NULL
	AND meal_plans.status = 'finalized'
	AND meal_plan_options.chosen IS TRUE
	AND meal_plans.tasks_created IS FALSE
	AND starts_at < NOW() + (1 * interval '1 week')
GROUP BY meal_plans.id, meal_plan_options.id, meals.id, meal_plan_events.id, meal_recipes.recipe_id
ORDER BY meal_plans.id;
