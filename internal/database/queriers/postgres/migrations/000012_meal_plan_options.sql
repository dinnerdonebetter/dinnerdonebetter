CREATE TABLE IF NOT EXISTS meal_plan_options (
	id CHAR(27) NOT NULL PRIMARY KEY,
	meal_plan_id TEXT NOT NULL,
	day_of_week INTEGER NOT NULL,
	recipe_id TEXT NOT NULL,
	notes TEXT NOT NULL,
	created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	last_updated_on BIGINT DEFAULT NULL,
	archived_on BIGINT DEFAULT NULL,
	belongs_to_household CHAR(27) NOT NULL REFERENCES households(id) ON DELETE CASCADE
);