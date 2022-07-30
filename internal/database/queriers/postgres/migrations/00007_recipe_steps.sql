CREATE TABLE IF NOT EXISTS recipe_steps (
	"id" CHAR(27) NOT NULL PRIMARY KEY,
	"index" INTEGER NOT NULL,
	"preparation_id" CHAR(27) NOT NULL REFERENCES valid_preparations("id") ON DELETE CASCADE,
	"minimum_estimated_time_in_seconds" BIGINT NOT NULL,
	"maximum_estimated_time_in_seconds" BIGINT NOT NULL,
    "minimum_temperature_in_celsius" INTEGER,
    "maximum_temperature_in_celsius" INTEGER,
    "notes" TEXT NOT NULL,
    "optional" BOOLEAN NOT NULL DEFAULT 'false',
	"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	"last_updated_on" BIGINT DEFAULT NULL,
	"archived_on" BIGINT DEFAULT NULL,
	"belongs_to_recipe" CHAR(27) NOT NULL REFERENCES recipes("id") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS recipe_steps_belongs_to_recipe ON recipe_steps (belongs_to_recipe);
