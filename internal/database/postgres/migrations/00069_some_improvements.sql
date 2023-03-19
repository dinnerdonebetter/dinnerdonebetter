ALTER TABLE valid_preparations ADD COLUMN "applies_to_all_ingredients" BOOLEAN NOT NULL DEFAULT 'false';
ALTER TABLE recipe_step_ingredients ADD COLUMN "to_taste" BOOLEAN NOT NULL DEFAULT 'false';
ALTER TABLE recipe_step_ingredients ADD COLUMN "product_percentage_to_use" NUMERIC(14, 2);
