ALTER TABLE valid_measurement_units ALTER COLUMN volumetric DROP NOT NULL;
ALTER TABLE valid_ingredients ALTER COLUMN is_liquid DROP NOT NULL;
ALTER TABLE meal_plan_options ALTER COLUMN belongs_to_meal_plan_event DROP NOT NULL;
-- ALTER TABLE valid_ingredients ALTER COLUMN minimum_ideal_storage_temperature_in_celsius TYPE real;
-- ALTER TABLE valid_ingredients ALTER COLUMN maximum_ideal_storage_temperature_in_celsius TYPE real;

-- ALTER TABLE valid_measurement_conversions RENAME TO valid_measurement_unit_conversions;
