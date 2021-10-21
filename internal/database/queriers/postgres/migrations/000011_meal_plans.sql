CREATE TABLE IF NOT EXISTS meal_plans (
	id CHAR(27) NOT NULL PRIMARY KEY,
	state TEXT NOT NULL,
	starts_at BIGINT NOT NULL,
	ends_at BIGINT NOT NULL,
	created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	last_updated_on BIGINT DEFAULT NULL,
	archived_on BIGINT DEFAULT NULL,
	belongs_to_account CHAR(27) NOT NULL REFERENCES accounts(id) ON DELETE CASCADE
);