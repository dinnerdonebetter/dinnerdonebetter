ALTER TABLE valid_preparations ADD COLUMN minimum_ingredient_count INTEGER NOT NULL DEFAULT 1;
ALTER TABLE valid_preparations ADD COLUMN maximum_ingredient_count INTEGER;
ALTER TABLE valid_preparations ADD COLUMN minimum_instrument_count INTEGER NOT NULL DEFAULT 1;
ALTER TABLE valid_preparations ADD COLUMN maximum_instrument_count INTEGER;
ALTER TABLE valid_preparations ADD COLUMN temperature_required BOOLEAN NOT NULL DEFAULT 'false';
ALTER TABLE valid_preparations ADD COLUMN time_estimate_required BOOLEAN NOT NULL DEFAULT 'false';
ALTER TABLE valid_preparations DROP COLUMN IF EXISTS "zero_ingredients_allowable";
