-- name: ValidIngredientPreparationExists :one
SELECT EXISTS ( SELECT valid_ingredient_preparations.id FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL AND valid_ingredient_preparations.id = $1 );

-- name: GetValidIngredientPreparation :one
SELECT
    valid_ingredient_preparations.id,
    valid_ingredient_preparations.notes,
    valid_preparations.id,
    valid_preparations.name,
    valid_preparations.description,
    valid_preparations.icon_path,
    valid_preparations.created_on,
    valid_preparations.last_updated_on,
    valid_preparations.archived_on,
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
    valid_ingredients.archived_on,
    valid_ingredient_preparations.created_on,
    valid_ingredient_preparations.last_updated_on,
    valid_ingredient_preparations.archived_on
FROM valid_ingredient_preparations
         JOIN valid_ingredients ON valid_ingredient_preparations.valid_ingredient_id = valid_ingredients.id
         JOIN valid_preparations ON valid_ingredient_preparations.valid_preparation_id = valid_preparations.id
WHERE valid_ingredient_preparations.archived_on IS NULL
  AND valid_ingredient_preparations.id = $1;

-- name: GetTotalValidIngredientPreparationsCount :one
SELECT COUNT(valid_ingredient_preparations.id) FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL;

-- name: CreateValidIngredientPreparation :exec
INSERT INTO valid_ingredient_preparations (id,notes,valid_preparation_id,valid_ingredient_id) VALUES ($1,$2,$3,$4);

-- name: UpdateValidIngredientPreparation :exec
UPDATE valid_ingredient_preparations SET notes = $1, valid_preparation_id = $2, valid_ingredient_id = $3, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $4;

-- name: ArchiveValidIngredientPreparation :exec
UPDATE valid_ingredient_preparations SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1;
