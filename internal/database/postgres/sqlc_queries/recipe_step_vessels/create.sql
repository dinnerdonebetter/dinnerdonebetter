-- name: CreateRecipeStepVessel :exec

INSERT INTO recipe_step_vessels
(id,"name",notes,belongs_to_recipe_step,recipe_step_product_id,valid_vessel_id,vessel_predicate,minimum_quantity,maximum_quantity,unavailable_after_step)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);
