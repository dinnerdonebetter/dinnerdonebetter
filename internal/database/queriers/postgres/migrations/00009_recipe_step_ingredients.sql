CREATE TABLE IF NOT EXISTS recipe_step_ingredients (
	"id" CHAR(27) NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "ingredient_id" CHAR(27) REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    "measurement_unit" CHAR(27) REFERENCES valid_measurement_units("id"),
    "minimum_quantity_value" DOUBLE PRECISION NOT NULL,
    "maximum_quantity_value" DOUBLE PRECISION NOT NULL,
	"quantity_notes" TEXT NOT NULL,
	"product_of_recipe_step" BOOLEAN NOT NULL,
	"ingredient_notes" TEXT NOT NULL,
	"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	"last_updated_on" BIGINT DEFAULT NULL,
	"archived_on" BIGINT DEFAULT NULL,
    "recipe_step_product_id" CHAR(27) REFERENCES recipe_step_products("id") ON DELETE RESTRICT,
    "belongs_to_recipe_step" CHAR(27) NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    UNIQUE("ingredient_id", "belongs_to_recipe_step")
);

ALTER TABLE recipe_step_ingredients ADD CONSTRAINT valid_instrument_or_product check (recipe_step_product_id is not null or ingredient_id is not null)
