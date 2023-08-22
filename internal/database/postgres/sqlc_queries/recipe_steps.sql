-- name: ArchiveRecipeStep :exec

UPDATE recipe_steps SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_recipe = $1 AND id = $2;

-- name: CreateRecipeStep :exec

INSERT INTO recipe_steps
(id,index,preparation_id,minimum_estimated_time_in_seconds,maximum_estimated_time_in_seconds,minimum_temperature_in_celsius,maximum_temperature_in_celsius,notes,explicit_instructions,condition_expression,optional,start_timer_automatically,belongs_to_recipe)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13);

-- name: CheckRecipeStepExistence :one

SELECT EXISTS (
	SELECT recipe_steps.id
	FROM recipe_steps
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	WHERE recipe_steps.archived_at IS NULL
	  AND recipe_steps.belongs_to_recipe = $1
	  AND recipe_steps.id = $2
	  AND recipes.archived_at IS NULL
	  AND recipes.id = $1
);

-- name: GetRecipeStep :one

SELECT
	recipe_steps.id,
	recipe_steps.index,
	valid_preparations.id,
	valid_preparations.name,
	valid_preparations.description,
	valid_preparations.icon_path,
	valid_preparations.yields_nothing,
	valid_preparations.restrict_to_ingredients,
	valid_preparations.minimum_ingredient_count,
	valid_preparations.maximum_ingredient_count,
	valid_preparations.minimum_instrument_count,
	valid_preparations.maximum_instrument_count,
	valid_preparations.temperature_required,
	valid_preparations.time_estimate_required,
	valid_preparations.condition_expression_required,
    valid_preparations.consumes_vessel,
    valid_preparations.only_for_vessels,
    valid_preparations.minimum_vessel_count,
    valid_preparations.maximum_vessel_count,
	valid_preparations.slug,
	valid_preparations.past_tense,
	valid_preparations.created_at,
	valid_preparations.last_updated_at,
	valid_preparations.archived_at,
	recipe_steps.minimum_estimated_time_in_seconds,
	recipe_steps.maximum_estimated_time_in_seconds,
	recipe_steps.minimum_temperature_in_celsius,
	recipe_steps.maximum_temperature_in_celsius,
	recipe_steps.notes,
	recipe_steps.explicit_instructions,
	recipe_steps.condition_expression,
	recipe_steps.optional,
	recipe_steps.start_timer_automatically,
	recipe_steps.created_at,
	recipe_steps.last_updated_at,
	recipe_steps.archived_at,
	recipe_steps.belongs_to_recipe
FROM recipe_steps
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = $1
	AND recipe_steps.id = $2
	AND recipes.archived_at IS NULL
	AND recipes.id = $1;

-- name: GetRecipeStepByRecipeID :one

SELECT
	recipe_steps.id,
	recipe_steps.index,
	valid_preparations.id,
	valid_preparations.name,
	valid_preparations.description,
	valid_preparations.icon_path,
	valid_preparations.yields_nothing,
	valid_preparations.restrict_to_ingredients,
	valid_preparations.minimum_ingredient_count,
	valid_preparations.maximum_ingredient_count,
	valid_preparations.minimum_instrument_count,
	valid_preparations.maximum_instrument_count,
	valid_preparations.temperature_required,
	valid_preparations.time_estimate_required,
	valid_preparations.condition_expression_required,
    valid_preparations.consumes_vessel,
    valid_preparations.only_for_vessels,
    valid_preparations.minimum_vessel_count,
    valid_preparations.maximum_vessel_count,
	valid_preparations.slug,
	valid_preparations.past_tense,
	valid_preparations.created_at,
	valid_preparations.last_updated_at,
	valid_preparations.archived_at,
	recipe_steps.minimum_estimated_time_in_seconds,
	recipe_steps.maximum_estimated_time_in_seconds,
	recipe_steps.minimum_temperature_in_celsius,
	recipe_steps.maximum_temperature_in_celsius,
	recipe_steps.notes,
	recipe_steps.explicit_instructions,
	recipe_steps.condition_expression,
	recipe_steps.optional,
	recipe_steps.start_timer_automatically,
	recipe_steps.created_at,
	recipe_steps.last_updated_at,
	recipe_steps.archived_at,
	recipe_steps.belongs_to_recipe
FROM recipe_steps
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipe_steps.archived_at IS NULL
	AND recipe_steps.id = $1;

-- name: UpdateRecipeStep :exec

UPDATE recipe_steps SET
	index = $1,
	preparation_id = $2,
	minimum_estimated_time_in_seconds = $3,
	maximum_estimated_time_in_seconds = $4,
	minimum_temperature_in_celsius = $5,
	maximum_temperature_in_celsius = $6,
	notes = $7,
	explicit_instructions = $8,
	condition_expression = $9,
	optional = $10,
	start_timer_automatically = $11,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_recipe = $12
	AND id = $13;
