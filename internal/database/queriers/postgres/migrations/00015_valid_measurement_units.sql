CREATE TABLE IF NOT EXISTS valid_measurement_units (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL DEFAULT '',
    "icon_path" TEXT NOT NULL DEFAULT '',
    "volumetric" BOOLEAN DEFAULT 'false',
    "created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "last_updated_on" BIGINT DEFAULT NULL,
    "archived_on" BIGINT DEFAULT NULL,
    UNIQUE("name")
);
