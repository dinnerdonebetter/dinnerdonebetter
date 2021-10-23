CREATE TABLE IF NOT EXISTS recipe_step_ingredients (
	"id" CHAR(27) NOT NULL PRIMARY KEY,
    "ingredient_id" TEXT,
	"quantity_type" TEXT NOT NULL,
	"quantity_value" DOUBLE PRECISION NOT NULL,
	"quantity_notes" TEXT NOT NULL,
	"product_of_recipe_step" BOOLEAN NOT NULL,
	"ingredient_notes" TEXT NOT NULL,
	"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	"last_updated_on" BIGINT DEFAULT NULL,
	"archived_on" BIGINT DEFAULT NULL,
	"belongs_to_recipe_step" CHAR(27) NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE
);