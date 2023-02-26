ALTER TABLE valid_ingredient_measurement_units ALTER COLUMN maximum_allowable_quantity DROP DEFAULT;
ALTER TABLE valid_ingredient_measurement_units ALTER COLUMN maximum_allowable_quantity DROP NOT NULL;

ALTER TABLE recipe_step_vessels ALTER COLUMN maximum_quantity DROP DEFAULT;
ALTER TABLE recipe_step_vessels ALTER COLUMN maximum_quantity DROP NOT NULL;

ALTER TABLE recipe_step_products ALTER COLUMN minimum_quantity_value DROP DEFAULT;
ALTER TABLE recipe_step_products ALTER COLUMN minimum_quantity_value DROP NOT NULL;
ALTER TABLE recipe_step_products ALTER COLUMN maximum_quantity_value DROP DEFAULT;
ALTER TABLE recipe_step_products ALTER COLUMN maximum_quantity_value DROP NOT NULL;
ALTER TABLE recipe_step_products ALTER COLUMN contained_in_vessel_index DROP DEFAULT;
ALTER TABLE recipe_step_products ALTER COLUMN contained_in_vessel_index DROP NOT NULL;

ALTER TABLE recipe_step_instruments ALTER COLUMN maximum_quantity DROP DEFAULT;
ALTER TABLE recipe_step_instruments ALTER COLUMN maximum_quantity DROP NOT NULL;

ALTER TABLE recipe_step_ingredients ALTER COLUMN maximum_quantity_value DROP DEFAULT;
ALTER TABLE recipe_step_ingredients ALTER COLUMN maximum_quantity_value DROP NOT NULL;
ALTER TABLE recipe_step_ingredients ALTER COLUMN vessel_index DROP DEFAULT;
ALTER TABLE recipe_step_ingredients ALTER COLUMN vessel_index DROP NOT NULL;
UPDATE recipe_step_ingredients SET vessel_index = NULL WHERE vessel_index = 0;

ALTER TABLE meal_plan_grocery_list_items ALTER COLUMN maximum_quantity_needed DROP DEFAULT;
ALTER TABLE meal_plan_grocery_list_items ALTER COLUMN maximum_quantity_needed DROP NOT NULL;








