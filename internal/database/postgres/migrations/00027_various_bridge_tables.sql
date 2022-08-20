CREATE TABLE IF NOT EXISTS valid_preparation_instruments (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "notes" TEXT NOT NULL DEFAULT '',
    "valid_preparation_id" CHAR(27) NOT NULL REFERENCES valid_preparations("id") ON DELETE CASCADE,
    "valid_instrument_id" CHAR(27) NOT NULL REFERENCES valid_instruments("id") ON DELETE CASCADE,
    "created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "last_updated_on" BIGINT DEFAULT NULL,
    "archived_on" BIGINT DEFAULT NULL,
    UNIQUE("valid_preparation_id", "valid_instrument_id")
);

CREATE TABLE IF NOT EXISTS valid_ingredient_measurement_units (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "notes" TEXT NOT NULL DEFAULT '',
    "valid_ingredient_id" CHAR(27) NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    "valid_measurement_unit_id" CHAR(27) NOT NULL REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    "created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "last_updated_on" BIGINT DEFAULT NULL,
    "archived_on" BIGINT DEFAULT NULL,
    UNIQUE("valid_ingredient_id", "valid_measurement_unit_id")
);
