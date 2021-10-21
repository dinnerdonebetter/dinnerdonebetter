CREATE TABLE IF NOT EXISTS recipe_step_instruments (
	id CHAR(27) NOT NULL PRIMARY KEY,
	instrument_id TEXT,
	recipe_step_id TEXT NOT NULL,
	notes TEXT NOT NULL,
	created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	last_updated_on BIGINT DEFAULT NULL,
	archived_on BIGINT DEFAULT NULL,
	"belongs_to_recipe_step" CHAR(27) NOT NULL REFERENCES recipe_steps(id) ON DELETE CASCADE
);