ALTER TABLE recipe_step_ingredients ALTER COLUMN minimum_quantity_value TYPE INTEGER;
ALTER TABLE recipe_step_ingredients ALTER COLUMN maximum_quantity_value TYPE INTEGER;
ALTER TABLE recipe_step_products ALTER COLUMN minimum_quantity_value TYPE INTEGER;
ALTER TABLE recipe_step_products ALTER COLUMN maximum_quantity_value TYPE INTEGER;
ALTER TABLE recipe_step_products ALTER COLUMN minimum_storage_temperature_in_celsius TYPE INTEGER;
ALTER TABLE recipe_step_products ALTER COLUMN maximum_storage_temperature_in_celsius TYPE INTEGER;
-- ALTER TABLE valid_ingredient_measurement_units ALTER COLUMN minimum_allowable_quantity TYPE INTEGER;
-- ALTER TABLE valid_ingredient_measurement_units ALTER COLUMN maximum_allowable_quantity TYPE INTEGER;
-- ALTER TABLE valid_ingredients ALTER COLUMN minimum_ideal_storage_temperature_in_celsius TYPE INTEGER;
-- ALTER TABLE valid_ingredients ALTER COLUMN maximum_ideal_storage_temperature_in_celsius TYPE INTEGER;
