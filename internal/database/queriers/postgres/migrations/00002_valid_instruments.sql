CREATE TABLE IF NOT EXISTS valid_instruments (
	"id" CHAR(27) NOT NULL PRIMARY KEY,
	"name" TEXT NOT NULL,
	"variant" TEXT NOT NULL, -- TODO: Dump me
	"description" TEXT NOT NULL,
    "icon_path" TEXT NOT NULL,
	"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	"last_updated_on" BIGINT DEFAULT NULL,
	"archived_on" BIGINT DEFAULT NULL,
    UNIQUE("name", "variant")
);