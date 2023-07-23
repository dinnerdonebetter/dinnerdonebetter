ALTER TABLE recipe_step_instruments ADD COLUMN recipe_step_product_id CHAR(27) REFERENCES recipe_step_products("id");
ALTER TABLE recipe_step_instruments DROP COLUMN recipe_step_id;
ALTER TABLE recipe_step_instruments ADD COLUMN product_of_recipe_step BOOLEAN NOT NULL DEFAULT 'false';
ALTER TABLE recipe_step_instruments ADD COLUMN name TEXT NOT NULL DEFAULT '';
