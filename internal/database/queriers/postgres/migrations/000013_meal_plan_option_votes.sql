CREATE TABLE IF NOT EXISTS meal_plan_option_votes (
	id CHAR(27) NOT NULL PRIMARY KEY,
	meal_plan_option_id TEXT NOT NULL,
	day_of_week INTEGER NOT NULL,
	points INTEGER NOT NULL,
	abstain BOOLEAN NOT NULL,
	notes TEXT NOT NULL,
	created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	last_updated_on BIGINT DEFAULT NULL,
	archived_on BIGINT DEFAULT NULL,
	belongs_to_account CHAR(27) NOT NULL REFERENCES accounts(id) ON DELETE CASCADE
);