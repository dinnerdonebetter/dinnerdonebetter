-- name: CheckValidIngredientGroupExistence :one

SELECT EXISTS ( SELECT valid_ingredient_groups.id FROM valid_ingredient_groups WHERE valid_ingredient_groups.archived_at IS NULL AND valid_ingredient_groups.id = $1 );
