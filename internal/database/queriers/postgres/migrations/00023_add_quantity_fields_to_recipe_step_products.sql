ALTER TABLE recipe_step_products DROP COLUMN IF EXISTS "recipe_step_id";

ALTER TABLE recipe_step_products ADD COLUMN "quantity_type" TEXT NOT NULL;
ALTER TABLE recipe_step_products ADD COLUMN "quantity_value" DOUBLE PRECISION NOT NULL;
ALTER TABLE recipe_step_products ADD COLUMN "quantity_notes" TEXT NOT NULL;
