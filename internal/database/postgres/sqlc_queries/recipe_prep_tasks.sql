-- name: ArchiveRecipePrepTask :execrows

UPDATE recipe_prep_tasks SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;

-- name: CreateRecipePrepTask :exec

INSERT INTO recipe_prep_tasks (id,name,description,notes,optional,explicit_storage_instructions,minimum_time_buffer_before_recipe_in_seconds,maximum_time_buffer_before_recipe_in_seconds,storage_type,minimum_storage_temperature_in_celsius,maximum_storage_temperature_in_celsius,belongs_to_recipe)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);

-- name: CheckRecipePrepTaskExistence :one

SELECT EXISTS (
	SELECT recipe_prep_tasks.id
	FROM recipe_prep_tasks
	JOIN recipes ON recipe_prep_tasks.belongs_to_recipe=recipes.id
	WHERE recipe_prep_tasks.archived_at IS NULL
	  AND recipe_prep_tasks.belongs_to_recipe = $1
	  AND recipe_prep_tasks.id = $2
	  AND recipes.archived_at IS NULL
	  AND recipes.id = $1
);

-- name: GetRecipePrepTask :many

SELECT
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
	recipe_prep_task_steps.id as task_step_id,
	recipe_prep_task_steps.belongs_to_recipe_step as task_step_belongs_to_recipe_step,
	recipe_prep_task_steps.belongs_to_recipe_prep_task as task_step_belongs_to_recipe_prep_task,
	recipe_prep_task_steps.satisfies_recipe_step as task_step_satisfies_recipe_step
FROM recipe_prep_tasks
    JOIN recipe_prep_task_steps ON recipe_prep_tasks.id=recipe_prep_task_steps.belongs_to_recipe_prep_task
WHERE recipe_prep_tasks.archived_at IS NULL
    AND recipe_prep_tasks.id = sqlc.arg(recipe_prep_task_id)
    AND recipe_prep_tasks.archived_at IS NULL;

-- name: ListAllRecipePrepTasksByRecipe :many

SELECT
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
    recipe_prep_task_steps.id as task_step_id,
    recipe_prep_task_steps.belongs_to_recipe_step as task_step_belongs_to_recipe_step,
    recipe_prep_task_steps.belongs_to_recipe_prep_task as task_step_belongs_to_recipe_prep_task,
    recipe_prep_task_steps.satisfies_recipe_step as task_step_satisfies_recipe_step
FROM recipe_prep_tasks
         JOIN recipe_prep_task_steps ON recipe_prep_task_steps.belongs_to_recipe_prep_task=recipe_prep_tasks.id
         JOIN recipe_steps ON recipe_prep_task_steps.belongs_to_recipe_step=recipe_steps.id
         JOIN recipes ON recipe_prep_tasks.belongs_to_recipe=recipes.id
WHERE recipe_prep_tasks.archived_at IS NULL
  AND recipe_steps.archived_at IS NULL
  AND recipes.archived_at IS NULL
  AND recipes.id = sqlc.arg(recipe_id);

-- name: UpdateRecipePrepTask :execrows

UPDATE recipe_prep_tasks SET
	 name = $1,
	 description = $2,
	 notes = $3,
	 optional = $4,
	 explicit_storage_instructions = $5,
	 minimum_time_buffer_before_recipe_in_seconds = $6,
	 maximum_time_buffer_before_recipe_in_seconds = $7,
	 storage_type = $8,
	 minimum_storage_temperature_in_celsius = $9,
	 maximum_storage_temperature_in_celsius = $10,
	 belongs_to_recipe = $11,
	 last_updated_at = NOW()
WHERE archived_at IS NULL AND id = $12;
