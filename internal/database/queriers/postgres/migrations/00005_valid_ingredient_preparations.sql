CREATE TABLE IF NOT EXISTS valid_ingredient_preparations (
	id CHAR(27) NOT NULL PRIMARY KEY,
	notes TEXT NOT NULL,
	valid_preparation_id TEXT NOT NULL,
	valid_ingredient_id TEXT NOT NULL,
	created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	last_updated_on BIGINT DEFAULT NULL,
	archived_on BIGINT DEFAULT NULL
);