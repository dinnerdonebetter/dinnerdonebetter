CREATE TYPE vessel_shape AS ENUM ('hemisphere', 'rectangle', 'cone', 'pyramid', 'cylinder', 'sphere', 'cube', 'other');

CREATE TABLE IF NOT EXISTS valid_vessels (
    "id" TEXT NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "plural_name" TEXT NOT NULL DEFAULT '',
    "description" TEXT NOT NULL DEFAULT '',
    "icon_path" TEXT NOT NULL,
    "usable_for_storage" BOOLEAN NOT NULL DEFAULT 'false',
    "slug" TEXT NOT NULL DEFAULT '',
    "display_in_summary_lists" BOOLEAN NOT NULL DEFAULT 'true',
    "include_in_generated_instructions" BOOLEAN NOT NULL DEFAULT 'true',
    "capacity" NUMERIC(14, 2) NOT NULL DEFAULT 0,
    "capacity_unit" TEXT NOT NULL REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    "width_in_millimeters" NUMERIC(14, 2),
    "length_in_millimeters" NUMERIC(14, 2),
    "height_in_millimeters" NUMERIC(14, 2),
    "shape" vessel_shape NOT NULL DEFAULT 'other',
    "last_indexed_at" TIMESTAMP WITH TIME ZONE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    UNIQUE("name", "archived_at"),
    UNIQUE("slug", "archived_at")
);

-- ALTER TABLE valid_instruments DROP COLUMN IF EXISTS "is_vessel";
-- ALTER TABLE valid_instruments DROP COLUMN IF EXISTS "is_exclusively_vessel";

ALTER TABLE recipe_step_vessels ADD COLUMN "valid_vessel_id" TEXT REFERENCES valid_vessels("id") ON DELETE CASCADE;
