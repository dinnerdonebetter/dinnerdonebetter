SELECT
	meal_plan_tasks.id,
	meal_plan_options.id,
	meal_plan_options.assigned_cook,
	meal_plan_options.assigned_dishwasher,
	meal_plan_options.chosen,
	meal_plan_options.tiebroken,
    meal_plan_options.meal_scale,
	meal_plan_options.meal_id,
	meal_plan_options.notes,
	meal_plan_options.created_at,
	meal_plan_options.last_updated_at,
	meal_plan_options.archived_at,
	meal_plan_options.belongs_to_meal_plan_event,
	recipe_prep_tasks.id,
	recipe_prep_tasks.name,
	recipe_prep_tasks.description,
	recipe_prep_tasks.notes,
	recipe_prep_tasks.optional,
	recipe_prep_tasks.explicit_storage_instructions,
	recipe_prep_tasks.minimum_time_buffer_before_recipe_in_seconds,
	recipe_prep_tasks.maximum_time_buffer_before_recipe_in_seconds,
	recipe_prep_tasks.storage_type,
	recipe_prep_tasks.minimum_storage_temperature_in_celsius,
	recipe_prep_tasks.maximum_storage_temperature_in_celsius,
	recipe_prep_tasks.belongs_to_recipe,
	recipe_prep_tasks.created_at,
	recipe_prep_tasks.last_updated_at,
	recipe_prep_tasks.archived_at,
	recipe_prep_task_steps.id,
	recipe_prep_task_steps.belongs_to_recipe_step,
	recipe_prep_task_steps.belongs_to_recipe_prep_task,
	recipe_prep_task_steps.satisfies_recipe_step,
	meal_plan_tasks.created_at,
	meal_plan_tasks.last_updated_at,
	meal_plan_tasks.completed_at,
	meal_plan_tasks.status,
	meal_plan_tasks.creation_explanation,
	meal_plan_tasks.status_explanation,
	meal_plan_tasks.assigned_to_user
FROM meal_plan_tasks
	FULL OUTER JOIN meal_plan_options ON meal_plan_tasks.belongs_to_meal_plan_option=meal_plan_options.id
	FULL OUTER JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
	FULL OUTER JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
	FULL OUTER JOIN meals ON meal_plan_options.meal_id=meals.id
	JOIN recipe_prep_tasks ON meal_plan_tasks.belongs_to_recipe_prep_task=recipe_prep_tasks.id
	JOIN recipe_prep_task_steps ON recipe_prep_task_steps.belongs_to_recipe_prep_task=recipe_prep_tasks.id
	JOIN recipe_steps ON recipe_prep_task_steps.belongs_to_recipe_step=recipe_steps.id
WHERE meal_plan_options.archived_at IS NULL
	AND meal_plan_events.archived_at IS NULL
	AND meal_plans.archived_at IS NULL
	AND meals.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND meal_plan_tasks.id = $1;
