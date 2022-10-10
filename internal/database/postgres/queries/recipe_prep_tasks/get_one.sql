SELECT
    recipe_prep_tasks.id,
    recipe_prep_tasks.notes,
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
    recipe_prep_task_steps.satisfies_recipe_step
FROM recipe_prep_tasks
     FULL OUTER JOIN meal_plan_task_recipe_steps ON recipe_prep_task_steps.belongs_to_recipe_prep_task=recipe_prep_tasks.id
WHERE recipe_prep_tasks.archived_at IS NULL
    AND recipe_prep_tasks.id = $1
    AND meal_plan_task_recipe_steps.archived_at IS NULL;
