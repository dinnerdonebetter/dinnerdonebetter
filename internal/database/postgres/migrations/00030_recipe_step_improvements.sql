CREATE TYPE recipe_step_product_type AS ENUM ('ingredient', 'instrument');
ALTER TABLE recipe_step_products ADD COLUMN type recipe_step_product_type NOT NULL;
