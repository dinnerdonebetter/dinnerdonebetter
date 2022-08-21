-- name: ValidIngredientExists :one
SELECT EXISTS ( SELECT valid_ingredients.id FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL AND valid_ingredients.id = $1 );

-- name: GetValidIngredient :one
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
    valid_ingredients.created_on,
    valid_ingredients.last_updated_on,
    valid_ingredients.archived_on
FROM valid_ingredients
WHERE valid_ingredients.archived_on IS NULL
ORDER BY random() LIMIT 1;

-- name: SearchForValidIngredients :many
SELECT valid_ingredients.id, valid_ingredients.name, valid_ingredients.description, valid_ingredients.warning, valid_ingredients.contains_egg, valid_ingredients.contains_dairy, valid_ingredients.contains_peanut, valid_ingredients.contains_tree_nut, valid_ingredients.contains_soy, valid_ingredients.contains_wheat, valid_ingredients.contains_shellfish, valid_ingredients.contains_sesame, valid_ingredients.contains_fish, valid_ingredients.contains_gluten, valid_ingredients.animal_flesh, valid_ingredients.volumetric, valid_ingredients.is_liquid, valid_ingredients.icon_path, valid_ingredients.created_on, valid_ingredients.last_updated_on, valid_ingredients.archived_on FROM valid_ingredients WHERE valid_ingredients.name ILIKE $1 AND valid_ingredients.archived_on IS NULL LIMIT 50;

-- name: GetTotalValidIngredientCount :one
SELECT COUNT(valid_ingredients.id) FROM valid_ingredients WHERE valid_ingredients.archived_on IS NULL;

-- name: CreateValidIngredient :exec
INSERT INTO valid_ingredients (id,name,description,warning,contains_egg,contains_dairy,contains_peanut,contains_tree_nut,contains_soy,contains_wheat,contains_shellfish,contains_sesame,contains_fish,contains_gluten,animal_flesh,volumetric,is_liquid,icon_path) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18);

-- name: UpdateValidIngredient :exec

UPDATE valid_ingredients SET
 name = $1,
 description = $2,
 warning = $3,
 contains_egg = $4,
 contains_dairy = $5,
 contains_peanut = $6,
 contains_tree_nut = $7,
 contains_soy = $8,
 contains_wheat = $9,
 contains_shellfish = $10,
 contains_sesame = $11,
 contains_fish = $12,
 contains_gluten = $13,
 animal_flesh = $14,
 volumetric = $15,
 is_liquid = $16,
 icon_path = $17,
 last_updated_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL AND id = $18;

-- name: ArchiveValidIngredient :exec
UPDATE valid_ingredients SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1;