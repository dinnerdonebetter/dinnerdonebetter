CREATE TABLE IF NOT EXISTS valid_measurement_conversions (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "from_unit" CHAR(27) NOT NULL REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    "to_unit" CHAR(27) NOT NULL REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    "only_for_ingredient" CHAR(27) REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    "modifier" BIGINT NOT NULL,
    "notes" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    UNIQUE("from_unit", "to_unit", "only_for_ingredient")
);