DROP TABLE recipe_step_products;

ALTER TABLE recipe_steps ADD COLUMN "optional" BOOLEAN NOT NULL DEFAULT 'false';
ALTER TABLE recipe_steps ADD COLUMN "yields" TEXT NOT NULL DEFAULT '';

ALTER TABLE recipe_steps DROP COLUMN IF EXISTS "why";
