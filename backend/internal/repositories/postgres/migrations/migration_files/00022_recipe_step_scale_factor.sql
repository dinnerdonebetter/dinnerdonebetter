-- Add scale_factor to recipe step instruments, vessels, and ingredients.
-- Default 1.0: when scaling a recipe (e.g. 2x), effective quantity scale = recipe_scale * scale_factor
-- (e.g. vessel with scale_factor 0.5 still counts as one when recipe is doubled).
ALTER TABLE recipe_step_instruments
	ADD COLUMN IF NOT EXISTS scale_factor NUMERIC(14,4) DEFAULT 1.0 NOT NULL;

ALTER TABLE recipe_step_vessels
	ADD COLUMN IF NOT EXISTS scale_factor NUMERIC(14,4) DEFAULT 1.0 NOT NULL;

ALTER TABLE recipe_step_ingredients
	ADD COLUMN IF NOT EXISTS scale_factor NUMERIC(14,4) DEFAULT 1.0 NOT NULL;
