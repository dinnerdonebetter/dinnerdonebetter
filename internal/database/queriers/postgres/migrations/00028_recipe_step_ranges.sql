ALTER TABLE recipe_steps RENAME COLUMN temperature_in_celsius TO minimum_temperature_in_celsius;
ALTER TABLE recipe_steps ADD COLUMN "maximum_temperature_in_celsius" INTEGER;
