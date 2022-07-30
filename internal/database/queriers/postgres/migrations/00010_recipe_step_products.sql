CREATE TABLE IF NOT EXISTS recipe_step_products (
	"id" CHAR(27) NOT NULL PRIMARY KEY,
	"name" TEXT NOT NULL,
    "quantity_type" TEXT NOT NULL,
    "quantity_value" DOUBLE PRECISION NOT NULL,
    "quantity_notes" TEXT NOT NULL,
	"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	"last_updated_on" BIGINT DEFAULT NULL,
	"archived_on" BIGINT DEFAULT NULL,
    "measurement_unit" CHAR(27) REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
	"belongs_to_recipe_step" CHAR(27) NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS recipe_step_products_belongs_to_recipe_step ON recipe_step_products (belongs_to_recipe_step);