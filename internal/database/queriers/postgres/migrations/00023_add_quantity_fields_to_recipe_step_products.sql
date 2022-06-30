ALTER TABLE recipe_step_products DROP COLUMN IF EXISTS "recipe_step_id";
ALTER TABLE recipe_steps DROP COLUMN IF EXISTS "prerequisite_step";

ALTER TABLE recipe_step_ingredients ADD COLUMN "name" TEXT NOT NULL;
ALTER TABLE recipe_step_ingredients ADD COLUMN "recipe_step_product_id" CHAR(27) REFERENCES recipe_step_products("id") ON DELETE RESTRICT;
ALTER TABLE recipe_step_products ADD COLUMN "quantity_type" TEXT NOT NULL;
ALTER TABLE recipe_step_products ADD COLUMN "quantity_value" DOUBLE PRECISION NOT NULL;
ALTER TABLE recipe_step_products ADD COLUMN "quantity_notes" TEXT NOT NULL;
ALTER TABLE recipe_step_ingredients ADD CONSTRAINT valid_instrument_or_product check (recipe_step_product_id is not null or ingredient_id is not null)
