CREATE TABLE IF NOT EXISTS valid_ingredient_preparations (
	"id" CHAR(27) NOT NULL PRIMARY KEY,
	"notes" TEXT NOT NULL,
	"valid_preparation_id" CHAR(27) NOT NULL REFERENCES valid_preparations("id") ON DELETE CASCADE,
	"valid_ingredient_id" CHAR(27) NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
	"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	"last_updated_on" BIGINT DEFAULT NULL,
	"archived_on" BIGINT DEFAULT NULL,
    UNIQUE("valid_preparation_id", "valid_ingredient_id")
);