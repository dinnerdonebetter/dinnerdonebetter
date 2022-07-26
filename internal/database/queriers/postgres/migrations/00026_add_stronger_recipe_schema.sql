CREATE TABLE IF NOT EXISTS valid_measurement_units (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "icon_path" TEXT NOT NULL,
    "created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "last_updated_on" BIGINT DEFAULT NULL,
    "archived_on" BIGINT DEFAULT NULL,
    UNIQUE("name")
);

-- introduce new temperature range into recipe steps and convert the value to a float
ALTER TABLE recipe_steps ALTER COLUMN "temperature_in_celsius" TYPE DOUBLE PRECISION USING temperature_in_celsius::double precision;
ALTER TABLE recipe_steps RENAME COLUMN "temperature_in_celsius" TO "min_temperature_in_celsius";
ALTER TABLE recipe_steps ADD COLUMN "max_temperature_in_celsius" DOUBLE PRECISION NOT NULL DEFAULT 0;

-- this is for, i.e. honey
ALTER TABLE valid_ingredients ADD COLUMN "animal_derived" BOOLEAN NOT NULL DEFAULT 'false';

-- add preference_rank to recipe_step_instruments

-- add type to recipe_step_products (to accommodate prepared instruments)

-- change quantity_type in recipe_step_ingredients to point to the new `valid_measurement_units` table.