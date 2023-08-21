-- name: CreateValidIngredientPreparation :exec

INSERT INTO valid_ingredient_preparations (id,notes,valid_preparation_id,valid_ingredient_id) VALUES ($1,$2,$3,$4);
