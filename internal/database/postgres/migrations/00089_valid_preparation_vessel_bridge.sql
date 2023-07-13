CREATE TABLE IF NOT EXISTS valid_preparation_vessels (
    "id" TEXT NOT NULL PRIMARY KEY,
    "valid_preparation_id" TEXT NOT NULL REFERENCES valid_preparations("id") ON DELETE CASCADE,
    "valid_vessel_id" TEXT NOT NULL REFERENCES valid_vessels("id") ON DELETE CASCADE,
    "notes" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE,
    "archived_at" TIMESTAMP WITH TIME ZONE,
    UNIQUE("valid_preparation_id", "valid_vessel_id", "archived_at")
);

CREATE INDEX IF NOT EXISTS valid_preparation_vessels_referencing_valid_preparations_idx ON valid_preparation_vessels (valid_preparation_id);
CREATE INDEX IF NOT EXISTS valid_preparation_vessels_referencing_valid_vessels_idx ON valid_preparation_vessels (valid_vessel_id);

ALTER TABLE recipe_step_vessels DROP COLUMN valid_instrument_id;
