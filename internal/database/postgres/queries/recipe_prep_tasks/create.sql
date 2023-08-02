-- name: CreateRecipePrepTask :exec

INSERT INTO recipe_prep_tasks (id,name,description,notes,optional,explicit_storage_instructions,minimum_time_buffer_before_recipe_in_seconds,maximum_time_buffer_before_recipe_in_seconds,storage_type,minimum_storage_temperature_in_celsius,maximum_storage_temperature_in_celsius,belongs_to_recipe)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);
