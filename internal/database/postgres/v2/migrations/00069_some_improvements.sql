ALTER TABLE valid_preparations ADD COLUMN "universal" BOOLEAN NOT NULL DEFAULT 'false';
ALTER TABLE valid_instruments ADD COLUMN "include_in_generated_instructions" BOOLEAN NOT NULL DEFAULT 'true';
ALTER TABLE recipe_step_ingredients ADD COLUMN "to_taste" BOOLEAN NOT NULL DEFAULT 'false';
ALTER TABLE recipe_step_ingredients ADD COLUMN "product_percentage_to_use" NUMERIC(14, 2);
ALTER TABLE recipes ADD COLUMN "slug" TEXT NOT NULL DEFAULT '';
ALTER TABLE recipes ADD COLUMN "portion_name" TEXT NOT NULL DEFAULT 'portion';
ALTER TABLE recipes ADD COLUMN "plural_portion_name" TEXT NOT NULL DEFAULT 'portions';
