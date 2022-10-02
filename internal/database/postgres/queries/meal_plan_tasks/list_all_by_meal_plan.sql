SELECT
    meal_plan_tasks.id,
    meal_plan_options.id,
    meal_plan_options.assigned_cook,
    meal_plan_options.assigned_dishwasher,
    meal_plan_options.chosen,
    meal_plan_options.tiebroken,
    meal_plan_options.meal_id,
    meal_plan_options.notes,
    meal_plan_options.prep_steps_created,
    meal_plan_options.created_at,
    meal_plan_options.last_updated_at,
    meal_plan_options.archived_at,
    meal_plan_options.belongs_to_meal_plan_event,
    meal_plan_tasks.cannot_complete_before,
    meal_plan_tasks.cannot_complete_after,
    meal_plan_tasks.created_at,
    meal_plan_tasks.completed_at,
    meal_plan_tasks.status,
    meal_plan_tasks.creation_explanation,
    meal_plan_tasks.status_explanation,
    meal_plan_tasks.assigned_to_user,
    meal_plan_task_recipe_steps.satisfies_recipe_step
FROM meal_plan_tasks
    FULL OUTER JOIN meal_plan_options ON meal_plan_tasks.belongs_to_meal_plan_option=meal_plan_options.id
    FULL OUTER JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
    FULL OUTER JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
    FULL OUTER JOIN meals ON meal_plan_options.meal_id=meals.id
    JOIN meal_plan_task_recipe_steps ON meal_plan_task_recipe_steps.belongs_to_meal_plan_task=meal_plan_tasks.id
    JOIN recipe_steps ON meal_plan_task_recipe_steps.satisfies_recipe_step=recipe_steps.id
WHERE meal_plan_options.archived_at IS NULL
    AND meal_plan_events.archived_at IS NULL
    AND meal_plans.archived_at IS NULL
    AND meals.archived_at IS NULL
    AND recipe_steps.archived_at IS NULL
    AND meal_plans.id = $1;
