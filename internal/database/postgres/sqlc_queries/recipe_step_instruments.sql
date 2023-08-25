-- name: ArchiveRecipeStepInstrument :exec

UPDATE recipe_step_instruments SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_recipe_step = $1 AND id = $2;

-- name: CreateRecipeStepInstrument :exec

INSERT INTO recipe_step_instruments
(id,instrument_id,recipe_step_product_id,"name",notes,preference_rank,optional,option_index,minimum_quantity,maximum_quantity,belongs_to_recipe_step)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);

-- name: CheckRecipeStepInstrumentExistence :one

SELECT EXISTS ( SELECT recipe_step_instruments.id FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_instruments.archived_at IS NULL AND recipe_step_instruments.belongs_to_recipe_step = sqlc.arg(recipe_step_id) AND recipe_step_instruments.id = sqlc.arg(recipe_step_instrument_id) AND recipe_steps.archived_at IS NULL AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id) AND recipe_steps.id = sqlc.arg(recipe_step_id) AND recipes.archived_at IS NULL AND recipes.id = sqlc.arg(recipe_id) );

-- name: GetRecipeStepInstrumentsForRecipe :many

SELECT
	recipe_step_instruments.id,
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.plural_name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.usable_for_storage,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.slug,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at,
	recipe_step_instruments.recipe_step_product_id,
	recipe_step_instruments.name,
	recipe_step_instruments.notes,
	recipe_step_instruments.preference_rank,
	recipe_step_instruments.optional,
	recipe_step_instruments.minimum_quantity,
	recipe_step_instruments.maximum_quantity,
	recipe_step_instruments.option_index,
	recipe_step_instruments.created_at,
	recipe_step_instruments.last_updated_at,
	recipe_step_instruments.archived_at,
	recipe_step_instruments.belongs_to_recipe_step
FROM recipe_step_instruments
	LEFT JOIN valid_instruments ON recipe_step_instruments.instrument_id=valid_instruments.id
	JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_instruments.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = $1
	AND recipes.archived_at IS NULL
	AND recipes.id = $1;

-- name: GetRecipeStepInstrument :one

SELECT
	recipe_step_instruments.id,
	valid_instruments.id as valid_instrument_id,
	valid_instruments.name as valid_instrument_name,
	valid_instruments.plural_name as valid_instrument_plural_name,
	valid_instruments.description as valid_instrument_description,
	valid_instruments.icon_path as valid_instrument_icon_path,
	valid_instruments.usable_for_storage as valid_instrument_usable_for_storage,
    valid_instruments.display_in_summary_lists as valid_instrument_display_in_summary_lists,
	valid_instruments.include_in_generated_instructions as valid_instrument_include_in_generated_instructions,
	valid_instruments.slug as valid_instrument_slug,
	valid_instruments.created_at as valid_instrument_created_at,
	valid_instruments.last_updated_at as valid_instrument_last_updated_at,
	valid_instruments.archived_at as valid_instrument_archived_at,
	recipe_step_instruments.recipe_step_product_id,
	recipe_step_instruments.name,
	recipe_step_instruments.notes,
	recipe_step_instruments.preference_rank,
	recipe_step_instruments.optional,
	recipe_step_instruments.minimum_quantity,
	recipe_step_instruments.maximum_quantity,
	recipe_step_instruments.option_index,
	recipe_step_instruments.created_at,
	recipe_step_instruments.last_updated_at,
	recipe_step_instruments.archived_at,
	recipe_step_instruments.belongs_to_recipe_step
FROM recipe_step_instruments
	LEFT JOIN valid_instruments ON recipe_step_instruments.instrument_id=valid_instruments.id
	JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_instruments.archived_at IS NULL
	AND recipe_step_instruments.belongs_to_recipe_step = sqlc.arg(recipe_step_id)
	AND recipe_step_instruments.id = sqlc.arg(recipe_step_instrument_id)
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id)
	AND recipe_steps.id = sqlc.arg(recipe_step_id)
	AND recipes.archived_at IS NULL
	AND recipes.id = sqlc.arg(recipe_id);

-- name: GetRecipeStepInstruments :many

SELECT
    recipe_step_instruments.id,
    valid_instruments.id as valid_instrument_id,
    valid_instruments.name as valid_instrument_name,
    valid_instruments.plural_name as valid_instrument_plural_name,
    valid_instruments.description as valid_instrument_description,
    valid_instruments.icon_path as valid_instrument_icon_path,
    valid_instruments.usable_for_storage as valid_instrument_usable_for_storage,
    valid_instruments.display_in_summary_lists as valid_instrument_display_in_summary_lists,
    valid_instruments.include_in_generated_instructions as valid_instrument_include_in_generated_instructions,
    valid_instruments.slug as valid_instrument_slug,
    valid_instruments.created_at as valid_instrument_created_at,
    valid_instruments.last_updated_at as valid_instrument_last_updated_at,
    valid_instruments.archived_at as valid_instrument_archived_at,
    recipe_step_instruments.recipe_step_product_id,
    recipe_step_instruments.name,
    recipe_step_instruments.notes,
    recipe_step_instruments.preference_rank,
    recipe_step_instruments.optional,
    recipe_step_instruments.minimum_quantity,
    recipe_step_instruments.maximum_quantity,
    recipe_step_instruments.option_index,
    recipe_step_instruments.created_at,
    recipe_step_instruments.last_updated_at,
    recipe_step_instruments.archived_at,
    recipe_step_instruments.belongs_to_recipe_step,
    (
        SELECT
            COUNT(recipe_step_instruments.id)
        FROM
            recipe_step_instruments
        WHERE
            recipe_step_instruments.archived_at IS NULL
            AND recipe_step_instruments.belongs_to_recipe_step = sqlc.arg(recipe_step_id)
            AND recipe_step_instruments.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
            AND recipe_step_instruments.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
            AND (recipe_step_instruments.last_updated_at IS NULL OR recipe_step_instruments.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
            AND (recipe_step_instruments.last_updated_at IS NULL OR recipe_step_instruments.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    ) as filtered_count,
    (
        SELECT
            COUNT(recipe_step_instruments.id)
        FROM
            recipe_step_instruments
        WHERE
            recipe_step_instruments.archived_at IS NULL
    ) as total_count
FROM recipe_step_instruments
    LEFT JOIN valid_instruments ON recipe_step_instruments.instrument_id=valid_instruments.id
    JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id
    JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE
    recipe_step_instruments.archived_at IS NULL
    AND recipe_step_instruments.belongs_to_recipe_step = sqlc.arg(recipe_step_id)
    AND recipe_step_instruments.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
    AND recipe_step_instruments.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    AND (recipe_step_instruments.last_updated_at IS NULL OR recipe_step_instruments.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
    AND (recipe_step_instruments.last_updated_at IS NULL OR recipe_step_instruments.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    AND recipe_steps.archived_at IS NULL
    AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id)
    AND recipe_steps.id = sqlc.arg(recipe_step_id)
    AND recipes.archived_at IS NULL
    AND recipes.id = sqlc.arg(recipe_id)
    OFFSET sqlc.arg(query_offset)
    LIMIT sqlc.arg(query_limit);

-- name: UpdateRecipeStepInstrument :exec

UPDATE recipe_step_instruments SET
	instrument_id = $1,
	recipe_step_product_id = $2,
	name = $3,
	notes = $4,
	preference_rank = $5,
	optional = $6,
	option_index = $7,
	minimum_quantity = $8,
	maximum_quantity = $9,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_recipe_step = $10
	AND id = $11;