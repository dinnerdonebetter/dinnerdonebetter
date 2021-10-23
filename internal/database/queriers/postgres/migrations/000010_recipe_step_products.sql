CREATE TABLE IF NOT EXISTS recipe_step_products (
	"id" CHAR(27) NOT NULL PRIMARY KEY,
	"name" TEXT NOT NULL,
	"recipe_step_id" TEXT NOT NULL,
	"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	"last_updated_on" BIGINT DEFAULT NULL,
	"archived_on" BIGINT DEFAULT NULL,
	"belongs_to_recipe_step" CHAR(27) NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE
);