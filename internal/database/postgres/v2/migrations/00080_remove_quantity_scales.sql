ALTER TABLE recipe_steps DROP COLUMN IF EXISTS "time_scale_factor";
ALTER TABLE recipe_step_products DROP COLUMN IF EXISTS "quantity_scale_factor";
ALTER TABLE recipe_step_ingredients DROP COLUMN IF EXISTS "quantity_scale_factor";
ALTER TABLE recipe_step_vessels DROP COLUMN IF EXISTS "quantity_scale_factor";
ALTER TABLE recipe_step_instruments DROP COLUMN IF EXISTS "quantity_scale_factor";
