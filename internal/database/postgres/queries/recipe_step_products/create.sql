INSERT INTO recipe_step_products
(id,name,type,measurement_unit,minimum_quantity_value,maximum_quantity_value,quantity_notes,compostable,maximum_storage_duration_in_seconds,minimum_storage_temperature_in_celsius,maximum_storage_temperature_in_celsius,storage_instructions,belongs_to_recipe_step,is_liquid,is_waste)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15);
