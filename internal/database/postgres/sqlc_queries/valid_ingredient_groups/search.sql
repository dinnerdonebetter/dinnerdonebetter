-- name: SearchForValidIngredientGroups :many

SELECT
	valid_ingredient_groups.id,
	valid_ingredient_groups.name,
	valid_ingredient_groups.description,
	valid_ingredient_groups.slug,
	valid_ingredient_groups.created_at,
	valid_ingredient_groups.last_updated_at,
	valid_ingredient_groups.archived_at,
	valid_ingredient_group_members.id,
    valid_ingredient_group_members.belongs_to_group,
    valid_ingredients.id,
    valid_ingredients.name,
    valid_ingredients.description,
    valid_ingredients.warning,
    valid_ingredients.contains_egg,
    valid_ingredients.contains_dairy,
    valid_ingredients.contains_peanut,
    valid_ingredients.contains_tree_nut,
    valid_ingredients.contains_soy,
    valid_ingredients.contains_wheat,
    valid_ingredients.contains_shellfish,
    valid_ingredients.contains_sesame,
    valid_ingredients.contains_fish,
    valid_ingredients.contains_gluten,
    valid_ingredients.animal_flesh,
    valid_ingredients.volumetric,
    valid_ingredients.is_liquid,
    valid_ingredients.icon_path,
    valid_ingredients.animal_derived,
    valid_ingredients.plural_name,
    valid_ingredients.restrict_to_preparations,
    valid_ingredients.minimum_ideal_storage_temperature_in_celsius,
    valid_ingredients.maximum_ideal_storage_temperature_in_celsius,
    valid_ingredients.storage_instructions,
    valid_ingredients.slug,
    valid_ingredients.contains_alcohol,
    valid_ingredients.shopping_suggestions,
    valid_ingredients.is_starch,
    valid_ingredients.is_protein,
    valid_ingredients.is_grain,
    valid_ingredients.is_fruit,
    valid_ingredients.is_salt,
    valid_ingredients.is_fat,
    valid_ingredients.is_acid,
    valid_ingredients.is_heat,
    valid_ingredients.created_at,
    valid_ingredients.last_updated_at,
    valid_ingredients.archived_at,
    valid_ingredient_group_members.created_at,
    valid_ingredient_group_members.archived_at
FROM valid_ingredient_groups
 JOIN valid_ingredient_group_members ON valid_ingredient_group_members.belongs_to_group=valid_ingredient_groups.id
  JOIN valid_ingredients ON valid_ingredients.id = valid_ingredient_group_members.valid_ingredient
WHERE valid_ingredient_groups.name ILIKE $1
AND valid_ingredient_groups.archived_at IS NULL
AND valid_ingredient_group_members.archived_at IS NULL
LIMIT 50;
