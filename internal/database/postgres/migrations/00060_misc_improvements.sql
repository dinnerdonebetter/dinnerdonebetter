ALTER TABLE household_user_memberships RENAME COLUMN household_roles TO household_role;
ALTER TABLE household_invitations ADD COLUMN expires_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() + interval '1 week';                -- #24
ALTER TABLE recipe_step_instruments ADD COLUMN option_index INTEGER NOT NULL DEFAULT 0;                                                     -- #230/#295
ALTER TABLE recipe_step_ingredients ADD COLUMN option_index INTEGER NOT NULL DEFAULT 0;                                                     -- #230/#295
ALTER TABLE recipe_step_ingredients ADD COLUMN requires_defrost BOOLEAN NOT NULL DEFAULT 'false';                                           -- #318
ALTER TABLE users RENAME COLUMN service_roles TO service_role;
ALTER TABLE users DROP COLUMN birth_day;                                                                                                    -- #364
ALTER TABLE users DROP COLUMN birth_month;                                                                                                  -- #364
ALTER TABLE users ADD COLUMN birthday TIMESTAMP WITH TIME ZONE;                                                                             -- #364
ALTER TABLE users ADD COLUMN email_address_verification_token TEXT NOT NULL DEFAULT 'replaceme';                                            -- #156
ALTER TABLE users ADD COLUMN email_address_verified_at TIMESTAMP WITH TIME ZONE;                                                            -- #156
ALTER TABLE valid_ingredients ADD COLUMN contains_alcohol BOOLEAN NOT NULL DEFAULT 'false';                                                 -- #363
ALTER TABLE valid_ingredients ADD COLUMN slug TEXT NOT NULL DEFAULT '';                                                                     -- #184
ALTER TABLE valid_instruments ADD COLUMN slug TEXT NOT NULL DEFAULT '';                                                                     -- #184
ALTER TABLE valid_measurement_units ADD COLUMN slug TEXT NOT NULL DEFAULT '';                                                               -- #184
ALTER TABLE valid_preparations ADD COLUMN slug TEXT NOT NULL DEFAULT '';                                                                    -- #184
ALTER TABLE valid_ingredients ADD COLUMN shopping_suggestions TEXT NOT NULL DEFAULT '';                                                     -- #320
ALTER TABLE valid_instruments ADD COLUMN display_in_summary_lists BOOLEAN NOT NULL DEFAULT 'true';                                          -- #241
CREATE TYPE valid_election_method AS ENUM ('schulze', 'instant-runoff');                                                                    -- #188
ALTER TABLE meal_plans ADD COLUMN election_method valid_election_method NOT NULL DEFAULT 'schulze';                                         -- #188
ALTER TABLE recipe_step_products ADD COLUMN is_liquid BOOLEAN NOT NULL DEFAULT 'false';                                                     -- #160
ALTER TABLE recipe_step_products ADD COLUMN is_waste BOOLEAN NOT NULL DEFAULT 'false';                                                      -- #384
CREATE TYPE component_type AS ENUM ('unspecified', 'amuse-bouche', 'appetizer', 'soup', 'main', 'salad', 'beverage', 'side', 'dessert');    -- #267
ALTER TABLE meal_recipes ADD COLUMN "meal_component_type" component_type NOT NULL DEFAULT 'unspecified';                                    -- #267
CREATE TYPE recipe_ingredient_scale AS ENUM ('linear', 'logarithmic', 'exponential');                                                       -- #294
ALTER TABLE meal_recipes RENAME TO meal_components;
