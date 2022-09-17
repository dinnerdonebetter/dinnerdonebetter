INSERT INTO recipe_step_instruments
(id,instrument_id,recipe_step_product_id,name,product_of_recipe_step,notes,preference_rank,optional,minimum_quantity,maximum_quantity,belongs_to_recipe_step)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);
