-- name: ArchiveRecipe :exec

UPDATE recipes SET archived_at = NOW() WHERE archived_at IS NULL AND created_by_user = $1 AND id = $2;

-- name: CreateRecipe :exec

INSERT INTO recipes (id,"name",slug,"source",description,inspired_by_recipe_id,min_estimated_portions,max_estimated_portions,portion_name,plural_portion_name,seal_of_approval,eligible_for_meals,yields_component_type,created_by_user) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14);

-- name: CheckRecipeExistence :one

SELECT EXISTS ( SELECT recipes.id FROM recipes WHERE recipes.archived_at IS NULL AND recipes.id = $1 );

-- name: GetRecipeByID :one

SELECT
	recipes.id,
	recipes.name,
	recipes.slug,
	recipes.source,
	recipes.description,
	recipes.inspired_by_recipe_id,
	recipes.min_estimated_portions,
	recipes.max_estimated_portions,
	recipes.portion_name,
	recipes.plural_portion_name,
	recipes.seal_of_approval,
	recipes.eligible_for_meals,
	recipes.yields_component_type,
	recipes.created_at,
	recipes.last_updated_at,
	recipes.archived_at,
	recipes.created_by_user,
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
FROM recipes
	FULL OUTER JOIN recipe_steps ON recipes.id=recipe_steps.belongs_to_recipe
	FULL OUTER JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipes.archived_at IS NULL
	AND recipes.id = $1
ORDER BY recipe_steps.index;

-- name: GetRecipeByIDAndAuthorID :many

SELECT
	recipes.id,
	recipes.name,
	recipes.slug,
	recipes.source,
	recipes.description,
	recipes.inspired_by_recipe_id,
	recipes.min_estimated_portions,
	recipes.max_estimated_portions,
	recipes.portion_name,
	recipes.plural_portion_name,
	recipes.seal_of_approval,
	recipes.eligible_for_meals,
	recipes.yields_component_type,
	recipes.created_at,
	recipes.last_updated_at,
	recipes.archived_at,
	recipes.created_by_user,
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
FROM recipes
	FULL OUTER JOIN recipe_steps ON recipes.id=recipe_steps.belongs_to_recipe
	FULL OUTER JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipes.archived_at IS NULL
	AND recipes.id = $1
	AND recipes.created_by_user = $2
ORDER BY recipe_steps.index;

-- name: GetRecipes :many

SELECT
    recipes.id,
    recipes.name,
    recipes.slug,
    recipes.source,
    recipes.description,
    recipes.inspired_by_recipe_id,
    recipes.min_estimated_portions,
    recipes.max_estimated_portions,
    recipes.portion_name,
    recipes.plural_portion_name,
    recipes.seal_of_approval,
    recipes.eligible_for_meals,
    recipes.yields_component_type,
    recipes.created_at,
    recipes.last_updated_at,
    recipes.archived_at,
    recipes.created_by_user,
    (
        SELECT
            COUNT(recipes.id)
        FROM
            recipes
        WHERE
            recipes.archived_at IS NULL
          AND recipes.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
          AND recipes.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
          AND (
                recipes.last_updated_at IS NULL
                OR recipes.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
            )
          AND (
                recipes.last_updated_at IS NULL
                OR recipes.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
            )
        OFFSET sqlc.narg(query_offset)
    ) AS filtered_count,
    (
        SELECT
            COUNT(recipes.id)
        FROM
            recipes
        WHERE
            recipes.archived_at IS NULL
    ) AS total_count
FROM recipes
    WHERE recipes.archived_at IS NULL
    AND recipes.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
    AND recipes.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    AND (
        recipes.last_updated_at IS NULL
        OR recipes.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
    )
    AND (
        recipes.last_updated_at IS NULL
        OR recipes.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
    )
OFFSET sqlc.narg(query_offset)
LIMIT sqlc.narg(query_limit);

-- name: GetRecipesNeedingIndexing :many

SELECT recipes.id
  FROM recipes
 WHERE (recipes.archived_at IS NULL)
       AND (
			(recipes.last_indexed_at IS NULL)
			OR recipes.last_indexed_at
				< now() - '24 hours'::INTERVAL
		);

-- name: GetRecipeIDsForMeal :many

SELECT
	recipes.id
FROM
	recipes
	 JOIN meal_components ON meal_components.recipe_id = recipes.id
	 JOIN meals ON meal_components.meal_id = meals.id
WHERE
	recipes.archived_at IS NULL
	AND meals.id = $1
GROUP BY
	recipes.id
ORDER BY
	recipes.id;

-- name: UpdateRecipe :exec

UPDATE recipes SET
    name = $1,
    slug = $2,
    source = $3,
    description = $4,
    inspired_by_recipe_id = $5,
	min_estimated_portions = $6,
	max_estimated_portions = $7,
    portion_name = $8,
    plural_portion_name = $9,
    seal_of_approval = $10,
    eligible_for_meals = $11,
	yields_component_type = $12,
    last_updated_at = NOW()
WHERE archived_at IS NULL
  AND created_by_user = $13
  AND id = $14;

-- name: UpdateRecipeLastIndexedAt :exec

UPDATE recipes SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL;
