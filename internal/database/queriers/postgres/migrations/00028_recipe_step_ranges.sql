ALTER TABLE recipe_steps RENAME COLUMN temperature_in_celsius TO minimum_temperature_in_celsius;
ALTER TABLE recipe_steps ADD COLUMN "maximum_temperature_in_celsius" INTEGER;
ALTER TABLE recipe_steps RENAME COLUMN min_estimated_time_in_seconds TO minimum_estimated_time_in_seconds;
ALTER TABLE recipe_steps RENAME COLUMN max_estimated_time_in_seconds TO maximum_estimated_time_in_seconds;

ALTER TABLE recipe_step_ingredients RENAME COLUMN quantity_value TO minimum_quantity_value;
ALTER TABLE recipe_step_ingredients ADD COLUMN "maximum_quantity_value" DOUBLE PRECISION NOT NULL;
