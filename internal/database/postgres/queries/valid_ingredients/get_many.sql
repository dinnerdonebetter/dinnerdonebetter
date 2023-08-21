-- name: GetValidIngredients :many

SELECT
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
    (
        SELECT
            COUNT(valid_ingredients.id)
        FROM
            valid_ingredients
        WHERE
            valid_ingredients.archived_at IS NULL
          AND valid_ingredients.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
          AND valid_ingredients.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
          AND (valid_ingredients.last_updated_at IS NULL OR valid_ingredients.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
          AND (valid_ingredients.last_updated_at IS NULL OR valid_ingredients.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    ) as filtered_count,
    (
        SELECT
            COUNT(valid_ingredients.id)
        FROM
            valid_ingredients
        WHERE
            valid_ingredients.archived_at IS NULL
    ) as total_count
FROM
  valid_ingredients
WHERE
  valid_ingredients.archived_at IS NULL
  AND valid_ingredients.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
  AND valid_ingredients.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
  AND (valid_ingredients.last_updated_at IS NULL OR valid_ingredients.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
  AND (valid_ingredients.last_updated_at IS NULL OR valid_ingredients.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
GROUP BY
  valid_ingredients.id
ORDER BY
  valid_ingredients.id
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);
