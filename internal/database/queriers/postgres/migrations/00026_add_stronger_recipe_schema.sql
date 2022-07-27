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

INSERT INTO valid_measurement_units (id,name,volumetric) VALUES
    ( '2CUumLd5Vxnbp79O7ewWFTeEG8b', 'gram', 'false' ),
    ( '2CUumRVyVpKSZ8bYf3MVnYcqFQf', 'milliliter', 'true' ),
    ( '2CUumOxYpP5N2I0Nb0nYqwbECI0', 'unit', 'false' ),
    ( '2CUumT5hzcIrvYpEDmrRmgrpi8J', 'clove', 'false' ),
    ( '2CUumOv3xTXNLIJpuhmrkJylhJ3', 'bunch', 'false' ),
    ( '2CUumPgusmUr9olCssJs8qUCv12', 'handful', 'false' ),
    ( '2CUumPXMn2nesduGGX7h6B8IV7e', 'teaspoon', 'true' ),
    ( '2CUumM34fwvbuzbPxb3z53GwANs', 'tablespoon', 'true' ),
    ( '2CUumO7M5l6dLlRNZ1iG28GH2Sw', 'can', 'false' ),
    ( '2CUumP2No8mYzti8SY0QHgAdgv7', 'cup', 'true' ),
    ( '2CUumS91nhzBU3RT1sebQEzUB5F', 'percent', 'false' );

CREATE TABLE IF NOT EXISTS valid_ingredient_measurement_units (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "notes" TEXT NOT NULL DEFAULT '',
    "valid_measurement_unit_id" CHAR(27) NOT NULL REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    "valid_ingredient_id" CHAR(27) NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    "created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "last_updated_on" BIGINT DEFAULT NULL,
    "archived_on" BIGINT DEFAULT NULL,
    UNIQUE("valid_measurement_unit_id", "valid_ingredient_id")
);


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

-- introduce new temperature range into recipe steps and convert the value to a float
ALTER TABLE recipe_steps ALTER COLUMN "temperature_in_celsius" TYPE DOUBLE PRECISION USING temperature_in_celsius::double precision;
ALTER TABLE recipe_steps RENAME COLUMN "temperature_in_celsius" TO "min_temperature_in_celsius";
ALTER TABLE recipe_steps ADD COLUMN "max_temperature_in_celsius" DOUBLE PRECISION NOT NULL DEFAULT 0;

-- this is for, i.e. honey
ALTER TABLE valid_ingredients ADD COLUMN "animal_derived" BOOLEAN NOT NULL DEFAULT 'false';

-- add preference_rank to recipe_step_instruments
ALTER TABLE recipe_step_instruments ADD COLUMN "preference_rank" INTEGER NOT NULL DEFAULT 0;

-- add type to recipe_step_products (to accommodate prepared instruments)
ALTER TABLE recipe_step_products ADD COLUMN "type" TEXT NOT NULL DEFAULT '';

-- add measurement_unit in recipe_step_ingredients to point to the new `valid_measurement_units` table.
ALTER TABLE recipe_step_ingredients ADD COLUMN "measurement_unit" CHAR(27) NOT NULL REFERENCES valid_measurement_units("id") ON DELETE RESTRICT;
ALTER TABLE recipe_step_ingredients DROP COLUMN IF EXISTS "quantity_type";
