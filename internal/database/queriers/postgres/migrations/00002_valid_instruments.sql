CREATE TABLE IF NOT EXISTS valid_instruments (
	id CHAR(27) NOT NULL PRIMARY KEY,
	name TEXT NOT NULL,
	variant TEXT NOT NULL,
	description TEXT NOT NULL,
	icon TEXT NOT NULL,
	created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	last_updated_on BIGINT DEFAULT NULL,
	archived_on BIGINT DEFAULT NULL
);