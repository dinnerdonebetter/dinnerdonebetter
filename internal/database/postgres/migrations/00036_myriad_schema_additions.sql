ALTER TABLE valid_instruments DROP COLUMN IF EXISTS "variant";                                                                  -- #227
ALTER TABLE valid_instruments ADD COLUMN "plural_name" TEXT NOT NULL DEFAULT '';                                                -- #228
ALTER TABLE valid_ingredients ADD COLUMN "animal_derived" BOOLEAN NOT NULL DEFAULT 'false';                                     -- #189
ALTER TABLE valid_ingredients ADD COLUMN "plural_name" TEXT NOT NULL DEFAULT '';                                                -- #228
ALTER TABLE valid_ingredients ADD COLUMN "restrict_to_preparations" BOOLEAN NOT NULL DEFAULT 'false';                           -- #239
ALTER TABLE valid_ingredients ADD COLUMN "minimum_ideal_storage_temperature_in_celsius" DOUBLE PRECISION NOT NULL DEFAULT 0;    -- #247
ALTER TABLE valid_ingredients ADD COLUMN "maximum_ideal_storage_temperature_in_celsius" DOUBLE PRECISION NOT NULL DEFAULT 0;    -- #247
ALTER TABLE valid_preparations ADD COLUMN "yields_nothing" BOOLEAN NOT NULL DEFAULT 'false';                                    -- #195
ALTER TABLE valid_preparations ADD COLUMN "restrict_to_ingredients" BOOLEAN NOT NULL DEFAULT 'false';                           -- #219
ALTER TABLE valid_preparations ADD COLUMN "zero_ingredients_allowable" BOOLEAN NOT NULL DEFAULT 'false';                        -- #224
ALTER TABLE valid_preparations ADD COLUMN "past_tense" TEXT NOT NULL DEFAULT '';                                                -- #89
ALTER TABLE valid_measurement_units ADD COLUMN "universal" BOOLEAN NOT NULL DEFAULT 'false';                                    -- #212
ALTER TABLE valid_measurement_units ADD COLUMN "metric" BOOLEAN NOT NULL DEFAULT 'false';                                       -- #213
ALTER TABLE valid_measurement_units ADD COLUMN "imperial" BOOLEAN NOT NULL DEFAULT 'false';                                     -- #213
ALTER TABLE valid_measurement_units ADD COLUMN "plural_name" TEXT NOT NULL DEFAULT '';                                          -- #228
ALTER TABLE valid_ingredient_measurement_units ADD COLUMN "minimum_allowable_quantity" DOUBLE PRECISION NOT NULL DEFAULT 0;     -- #232
ALTER TABLE valid_ingredient_measurement_units ADD COLUMN "maximum_allowable_quantity" DOUBLE PRECISION NOT NULL DEFAULT 0;     -- #232
ALTER TABLE recipe_step_products ADD COLUMN "compostable" BOOLEAN NOT NULL DEFAULT 'false';                                     -- #252
ALTER TABLE recipe_step_products ADD COLUMN "maximum_storage_duration_in_seconds" INTEGER NOT NULL DEFAULT 0;                   -- #263
ALTER TABLE recipe_step_products ADD COLUMN "minimum_storage_temperature_in_celsius" DOUBLE PRECISION NOT NULL DEFAULT 0;       -- #263
ALTER TABLE recipe_step_products ADD COLUMN "maximum_storage_temperature_in_celsius" DOUBLE PRECISION NOT NULL DEFAULT 0;       -- #263
ALTER TABLE recipe_step_instruments ADD COLUMN "optional" BOOLEAN NOT NULL DEFAULT 'false';                                     -- #233
ALTER TABLE recipe_step_instruments ADD COLUMN "minimum_quantity" INTEGER NOT NULL DEFAULT 1;                                   -- #240
ALTER TABLE recipe_step_instruments ADD COLUMN "maximum_quantity" INTEGER NOT NULL DEFAULT 1;                                   -- #240
ALTER TABLE recipe_step_ingredients ADD COLUMN "optional" BOOLEAN NOT NULL DEFAULT 'false';                                     -- #233
ALTER TABLE meal_plan_options ADD COLUMN "assigned_cook" CHAR(27) REFERENCES users("id") ON DELETE CASCADE;                     -- #259
ALTER TABLE recipes ADD COLUMN "yields_portions" INTEGER NOT NULL DEFAULT 1;                                                    -- #253
ALTER TABLE recipes ADD COLUMN "seal_of_approval" BOOLEAN NOT NULL DEFAULT 'false';                                             -- #265
ALTER TABLE recipe_steps ADD COLUMN "explicit_instructions" TEXT NOT NULL DEFAULT '';                                           -- #243
CREATE TYPE time_zone AS ENUM ('UTC', 'US/Pacific', 'US/Mountain', 'US/Central', 'US/Eastern');                                                             -- #260
ALTER TABLE households ADD COLUMN "time_zone" time_zone NOT NULL DEFAULT 'US/Central';                                                  -- #260
-- CREATE TYPE component_type AS ENUM ('main', 'appetizer', 'side', 'garnish');                                                    -- #267
-- ALTER TABLE meal_recipes ADD COLUMN "component_type" component_type;                                                            -- #267
