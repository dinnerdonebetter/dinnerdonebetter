CREATE TABLE IF NOT EXISTS recipes (
	id CHAR(27) NOT NULL PRIMARY KEY,
	name TEXT NOT NULL,
	source TEXT NOT NULL,
	description TEXT NOT NULL,
	inspired_by_recipe_id TEXT,
	created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	last_updated_on BIGINT DEFAULT NULL,
	archived_on BIGINT DEFAULT NULL,
	belongs_to_household CHAR(27) NOT NULL REFERENCES households(id) ON DELETE CASCADE
);