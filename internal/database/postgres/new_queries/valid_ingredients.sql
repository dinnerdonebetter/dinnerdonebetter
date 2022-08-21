-- name: ValidIngredientExists :one
SELECT EXISTS ( SELECT valid_ingredients.id FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL AND valid_ingredients.id = $1 );

-- name: GetValidIngredient :one
SELECT
    valid_ingredients.id,
    valid_ingredients.name,
    valid_ingredients.plural_name,
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
    valid_ingredients.animal_derived,
    valid_ingredients.volumetric,
    valid_ingredients.restrict_to_preparations,
    valid_ingredients.minimum_ideal_storage_temperature_in_celsius,
    valid_ingredients.maximum_ideal_storage_temperature_in_celsius,
    valid_ingredients.is_liquid,
    valid_ingredients.icon_path,
    valid_ingredients.created_on,
    valid_ingredients.last_updated_on,
    valid_ingredients.archived_on
FROM valid_ingredients
WHERE valid_ingredients.archived_on IS NULL
  AND valid_ingredients.id = $1;

-- name: GetRandomValidIngredient :one
SELECT
    valid_ingredients.id,
    valid_ingredients.name,
    valid_ingredients.plural_name,
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
    valid_ingredients.animal_derived,
    valid_ingredients.volumetric,
    valid_ingredients.restrict_to_preparations,
    valid_ingredients.minimum_ideal_storage_temperature_in_celsius,
    valid_ingredients.maximum_ideal_storage_temperature_in_celsius,
    valid_ingredients.is_liquid,
    valid_ingredients.icon_path,
    valid_ingredients.created_on,
    valid_ingredients.last_updated_on,
    valid_ingredients.archived_on
FROM valid_ingredients
WHERE valid_ingredients.archived_on IS NULL
ORDER BY random() LIMIT 1;

-- name: SearchForValidIngredients :many
SELECT
    valid_ingredients.id,
    valid_ingredients.name,
    valid_ingredients.plural_name,
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
    valid_ingredients.animal_derived,
    valid_ingredients.volumetric,
    valid_ingredients.restrict_to_preparations,
    valid_ingredients.minimum_ideal_storage_temperature_in_celsius,
    valid_ingredients.maximum_ideal_storage_temperature_in_celsius,
    valid_ingredients.is_liquid,
    valid_ingredients.icon_path,
    valid_ingredients.created_on,
    valid_ingredients.last_updated_on,
    valid_ingredients.archived_on
FROM valid_ingredients
WHERE valid_ingredients.name ILIKE $1
AND valid_ingredients.archived_on IS NULL
    LIMIT 50;

-- name: GetTotalValidIngredientCount :one
SELECT COUNT(valid_ingredients.id) FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL;

-- name: CreateValidIngredient :exec
INSERT INTO valid_ingredients (
    id,
    name,
    plural_name,
    description,
    warning,
    contains_egg,
    contains_dairy,
    contains_peanut,
    contains_tree_nut,
    contains_soy,
    contains_wheat,
    contains_shellfish,
    contains_sesame,
    contains_fish,
    contains_gluten,
    animal_flesh,
    animal_derived,
    volumetric,
    restrict_to_preparations,
    minimum_ideal_storage_temperature_in_celsius,
    maximum_ideal_storage_temperature_in_celsius,
    is_liquid,
    icon_path
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23);

-- name: UpdateValidIngredient :exec
UPDATE valid_ingredients SET
 name = $1,
 plural_name = $2,
 description = $3,
 warning = $4,
 contains_egg = $5,
 contains_dairy = $6,
 contains_peanut = $7,
 contains_tree_nut = $8,
 contains_soy = $9,
 contains_wheat = $10,
 contains_shellfish = $11,
 contains_sesame = $12,
 contains_fish = $13,
 contains_gluten = $14,
 animal_flesh = $15,
 animal_derived = $16,
 volumetric = $17,
 restrict_to_preparations = $18,
 minimum_ideal_storage_temperature_in_celsius = $19,
 maximum_ideal_storage_temperature_in_celsius = $20,
 is_liquid = $21,
 icon_path = $22,
 last_updated_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL AND id = $23;

-- name: ArchiveValidIngredient :exec
UPDATE valid_ingredients SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1;