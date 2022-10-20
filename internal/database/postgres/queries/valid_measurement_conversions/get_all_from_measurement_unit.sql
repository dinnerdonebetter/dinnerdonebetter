SELECT
    valid_measurement_conversions.id,
    valid_measurement_units_from.id,
    valid_measurement_units_from.name,
    valid_measurement_units_from.description,
    valid_measurement_units_from.volumetric,
    valid_measurement_units_from.icon_path,
    valid_measurement_units_from.universal,
    valid_measurement_units_from.metric,
    valid_measurement_units_from.imperial,
    valid_measurement_units_from.plural_name,
    valid_measurement_units_from.created_at,
    valid_measurement_units_from.last_updated_at,
    valid_measurement_units_from.archived_at,
    valid_measurement_units_to.id,
    valid_measurement_units_to.name,
    valid_measurement_units_to.description,
    valid_measurement_units_to.volumetric,
    valid_measurement_units_to.icon_path,
    valid_measurement_units_to.universal,
    valid_measurement_units_to.metric,
    valid_measurement_units_to.imperial,
    valid_measurement_units_to.plural_name,
    valid_measurement_units_to.created_at,
    valid_measurement_units_to.last_updated_at,
    valid_measurement_units_to.archived_at,
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
    valid_ingredients.created_at,
    valid_ingredients.last_updated_at,
    valid_ingredients.archived_at,
    valid_measurement_conversions.modifier,
    valid_measurement_conversions.notes,
    valid_measurement_conversions.created_at,
    valid_measurement_conversions.last_updated_at,
    valid_measurement_conversions.archived_at
FROM valid_measurement_conversions
         LEFT JOIN valid_ingredients ON valid_measurement_conversions.only_for_ingredient = valid_ingredients.id
         JOIN valid_measurement_units AS valid_measurement_units_from ON valid_measurement_conversions.from_unit = valid_measurement_units_from.id
         JOIN valid_measurement_units AS valid_measurement_units_to ON valid_measurement_conversions.to_unit = valid_measurement_units_to.id
WHERE valid_measurement_conversions.archived_at IS NULL
  AND valid_measurement_units_from.archived_at IS NULL
  AND valid_measurement_units_from.id = $1
  AND valid_measurement_units_to.archived_at IS NULL
