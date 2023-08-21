-- name: CheckMealPlanTaskExistence :one

SELECT EXISTS (
	SELECT meal_plan_tasks.id
	FROM meal_plan_tasks
		FULL OUTER JOIN meal_plan_options ON meal_plan_tasks.belongs_to_meal_plan_option=meal_plan_options.id
		FULL OUTER JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
		FULL OUTER JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
	WHERE meal_plan_tasks.completed_at IS NULL
		AND meal_plans.id = sqlc.arg(meal_plan_id)
		AND meal_plans.archived_at IS NULL
		AND meal_plan_tasks.id = sqlc.arg(meal_plan_task_id)
);
