INSERT INTO recipe_steps
(id,index,preparation_id,minimum_estimated_time_in_seconds,maximum_estimated_time_in_seconds,minimum_temperature_in_celsius,maximum_temperature_in_celsius,notes,explicit_instructions,optional,belongs_to_recipe)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);
