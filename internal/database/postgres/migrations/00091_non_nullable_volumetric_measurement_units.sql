ALTER TABLE valid_measurement_units ALTER COLUMN volumetric DROP NOT NULL;
ALTER TABLE valid_ingredients ALTER COLUMN is_liquid DROP NOT NULL;

-- ALTER TABLE valid_measurement_conversions RENAME TO valid_measurement_unit_conversions;
