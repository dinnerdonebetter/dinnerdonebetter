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
 ( '2CUumPXMn2nesduGGX7h6B8IV7e', 'teaspoon', 'true' ),
 ( '2CUumM34fwvbuzbPxb3z53GwANs', 'tablespoon', 'true' ),
 ( '2CUumO7M5l6dLlRNZ1iG28GH2Sw', 'can', 'false' ),
 ( '2CUumP2No8mYzti8SY0QHgAdgv7', 'cup', 'true' ),
 ( '2CUumS91nhzBU3RT1sebQEzUB5F', 'percent', 'false' );
