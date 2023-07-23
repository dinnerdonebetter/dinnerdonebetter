CREATE TABLE IF NOT EXISTS recipe_steps (
	"id" CHAR(27) NOT NULL PRIMARY KEY,
	"index" INTEGER NOT NULL,
	"preparation_id" CHAR(27) NOT NULL REFERENCES valid_preparations("id") ON DELETE CASCADE,
	"prerequisite_step" BIGINT NOT NULL,
	"min_estimated_time_in_seconds" BIGINT NOT NULL,
	"max_estimated_time_in_seconds" BIGINT NOT NULL,
	"temperature_in_celsius" INTEGER,
    "why" TEXT NOT NULL,
    "notes" TEXT NOT NULL,
	"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	"last_updated_on" BIGINT DEFAULT NULL,
	"archived_on" BIGINT DEFAULT NULL,
	"belongs_to_recipe" CHAR(27) NOT NULL REFERENCES recipes("id") ON DELETE CASCADE
);