SELECT
    recipes.id,
    recipes.name,
    recipes.source,
    recipes.description,
    recipes.inspired_by_recipe_id,
    recipes.yields_portions,
    recipes.seal_of_approval,
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
    valid_preparations.zero_ingredients_allowable,
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
    recipe_steps.optional,
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