ALTER TABLE users RENAME COLUMN service_roles TO service_role;
ALTER TABLE household_user_memberships RENAME COLUMN household_roles TO household_role;
CREATE TYPE recipe_ingredient_scale AS ENUM ('linear', 'logarithmic', 'exponential');                                                       -- #294
ALTER TABLE household_invitations ADD COLUMN expires_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() + interval '1 week';                -- #24
ALTER TABLE recipe_step_instruments ADD COLUMN option_index INTEGER NOT NULL DEFAULT 0;                                                  -- #230/#295
ALTER TABLE recipe_step_ingredients ADD COLUMN option_index INTEGER NOT NULL DEFAULT 0;                                                  -- #230/#295
ALTER TABLE recipe_step_ingredients ADD COLUMN requires_defrost BOOLEAN NOT NULL DEFAULT 'false';                                        -- #318
-- ALTER TABLE users DROP COLUMN birth_year;                                                                                                -- #364
-- ALTER TABLE users DROP COLUMN birth_month;                                                                                               -- #364
-- ALTER TABLE users ADD COLUMN birthday TIMESTAMP WITH TIME ZONE;                                                                          -- #364
-- ALTER TABLE valid_ingredients ADD COLUMN contains_alcohol BOOLEAN NOT NULL DEFAULT 'false';                                              -- #363
-- ALTER TABLE valid_ingredients ADD COLUMN slug TEXT NOT NULL DEFAULT '';                                                                  -- #184
-- ALTER TABLE valid_instruments ADD COLUMN slug TEXT NOT NULL DEFAULT '';                                                                  -- #184
-- ALTER TABLE valid_ingredients ADD COLUMN shopping_suggestions TEXT NOT NULL DEFAULT '';                                                  -- #320
-- ALTER TABLE valid_measurement_units ADD COLUMN slug TEXT NOT NULL DEFAULT '';                                                            -- #184
-- ALTER TABLE valid_preparations ADD COLUMN slug TEXT NOT NULL DEFAULT '';                                                                 -- #184
-- CREATE TYPE valid_election_method AS ENUM ('schulze', 'instant-runoff');                                                                 -- #188
-- ALTER TABLE meal_plans ADD COLUMN election_method valid_election_method NOT NULL DEFAULT 'schulze';                                      -- #188
-- ALTER TABLE recipe_step_products ADD COLUMN is_liquid NOT NULL DEFAULT 'false';                                                          -- #160
-- ALTER TABLE recipe_step_products ADD COLUMN is_waste NOT NULL DEFAULT 'false';                                                           -- #384
-- ALTER TABLE users RENAME COLUMN email_address_verified_at TIMESTAMP WITH TIME ZONE;                                                      -- #156
-- ALTER TABLE valid_instruments ADD COLUMN display_in_summary_lists BOOLEAN NOT NULL DEFAULT 'true';                                       -- #241
-- CREATE TYPE component_type AS ENUM ('unspecified', 'amuse-bouche', 'appetizer', 'soup', 'main', 'salad', 'beverage', 'side', 'dessert'); -- #267
-- ALTER TABLE meal_recipes ADD COLUMN "component_type" component_type NOT NULL DEFAULT 'unspecified';                                      -- #267

-- CREATE TABLE IF NOT EXISTS valid_ingredient_statuses (
--     id CHAR(27) NOT NULL PRIMARY KEY,
--     "name" text NOT NULL,
--     plural_name text DEFAULT ''::text NOT NULL,
--     slug text DEFAULT ''::text NOT NULL,
--     description text text DEFAULT ''::text NOT NULL,
--     icon_path text text DEFAULT ''::text NOT NULL,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
--     last_updated_at TIMESTAMP WITH TIME ZONE,
--     archived_at TIMESTAMP WITH TIME ZONE
-- );

-- CREATE TABLE IF NOT EXISTS valid_ingredient_status_ingredients (
--     id CHAR(27) NOT NULL PRIMARY KEY,
--     notes text DEFAULT ''::text NOT NULL,
--     valid_ingredient_id CHAR(27) NOT NULL REFERENCES valid_ingredients(id) ON DELETE CASCADE,
--     valid_ingredient_status_id CHAR(27) NOT NULL REFERENCES valid_ingredient_statuses(id) ON DELETE CASCADE,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
--     last_updated_at TIMESTAMP WITH TIME ZONE,
--     archived_at TIMESTAMP WITH TIME ZONE,
--     UNIQUE (valid_ingredient_id, valid_ingredient_status_id, archived_at)
-- );

-- -- #194
-- CREATE TABLE IF NOT EXISTS valid_instrument_ownerships (
--     id CHAR(27) NOT NULL PRIMARY KEY,
--     notes text DEFAULT ''::text NOT NULL,
--     valid_instrument_id CHAR(27) NOT NULL REFERENCES valid_instruments(id) ON DELETE CASCADE,
--     household_id CHAR(27) NOT NULL REFERENCES households(id) ON DELETE CASCADE,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
--     last_updated_at TIMESTAMP WITH TIME ZONE,
--     archived_at TIMESTAMP WITH TIME ZONE,
--     UNIQUE (household_id, valid_instrument_id, archived_at)
-- );

-- -- #152
-- CREATE TABLE IF NOT EXISTS preparation_conditions (
--     id CHAR(27) NOT NULL PRIMARY KEY,
--     notes text DEFAULT ''::text NOT NULL,
--     valid_preparation_id CHAR(27) NOT NULL REFERENCES valid_preparations(id) ON DELETE CASCADE,
--     liquid_count INTEGER NOT NULL DEFAULT 0,
--     notice TEXT DEFAULT ''::text NOT NULL,
--     minimum_temperature NUMERIC(15, 2),
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
--     last_updated_at TIMESTAMP WITH TIME ZONE,
--     archived_at TIMESTAMP WITH TIME ZONE,
--     UNIQUE (household_id, valid_instrument_id, archived_at)
-- );

-- -- #231
-- CREATE TABLE IF NOT EXISTS recipe_lists (
--     id CHAR(27) NOT NULL PRIMARY KEY,
--     name text,
--     notes text DEFAULT ''::text NOT NULL,
--     belongs_to_user CHAR(27) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
--     "public" boolean NOT NULL DEFAULT 'false',
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
--     last_updated_at TIMESTAMP WITH TIME ZONE,
--     archived_at TIMESTAMP WITH TIME ZONE
--     UNIQUE (name, belongs_to_user, archived_at)
-- );

-- -- #231
-- CREATE TABLE IF NOT EXISTS recipe_list_entries (
--     id CHAR(27) NOT NULL PRIMARY KEY,
--     notes text DEFAULT ''::text NOT NULL,
--     recipe_id CHAR(27) NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
--     belongs_to_recipe_collection CHAR(27) NOT NULL REFERENCES recipe_lists(id) ON DELETE CASCADE,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
--     last_updated_at TIMESTAMP WITH TIME ZONE,
--     archived_at TIMESTAMP WITH TIME ZONE
--     UNIQUE (name, recipe_id, belongs_to_recipe_collection, archived_at)
-- );

-- -- #296
-- CREATE TABLE IF NOT EXISTS meal_lists (
--     id CHAR(27) NOT NULL PRIMARY KEY,
--     name text,
--     notes text DEFAULT ''::text NOT NULL,
--     belongs_to_user CHAR(27) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
--     "public" boolean NOT NULL DEFAULT 'false',
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
--     last_updated_at TIMESTAMP WITH TIME ZONE,
--     archived_at TIMESTAMP WITH TIME ZONE
--     UNIQUE (name, belongs_to_user, archived_at)
-- );

-- -- #296
-- CREATE TABLE IF NOT EXISTS meal_list_entries (
--     id CHAR(27) NOT NULL PRIMARY KEY,
--     notes text DEFAULT ''::text NOT NULL,
--     meal_id CHAR(27) NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
--     belongs_to_meal_collection CHAR(27) NOT NULL REFERENCES meal_lists(id) ON DELETE CASCADE,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
--     last_updated_at TIMESTAMP WITH TIME ZONE,
--     archived_at TIMESTAMP WITH TIME ZONE
--     UNIQUE (name, meal_id, belongs_to_meal_collection, archived_at)
-- );

-- -- #278, #279, #280
-- CREATE TYPE us_state AS ENUM ('AL', 'AK', 'AZ', 'AR', 'CA', 'CO', 'CT', 'DE', 'FL', 'GA', 'HI', 'ID', 'IL', 'IN', 'IA', 'KS', 'KY', 'LA', 'ME', 'MD', 'MA', 'MI', 'MN', 'MS', 'MO', 'MT', 'NE', 'NV', 'NH', 'NJ', 'NM', 'NY', 'NC', 'ND', 'OH', 'OK', 'PA', 'RI', 'SC', 'SD', 'TN', 'TX', 'UT', 'VT', 'VA', 'WA', 'WV', 'WY');
-- CREATE TABLE IF NOT EXISTS grocery_stores (
--     "id" CHAR(27) NOT NULL PRIMARY KEY,
--     "brand" TEXT NOT NULL DEFAULT '',
--     "street_address_1" TEXT NOT NULL DEFAULT '',
--     "street_address_2" TEXT NOT NULL DEFAULT '',
--     "city" TEXT NOT NULL DEFAULT '',
--     "state" us_state NOT NULL,
--     "zip_code" TEXT NOT NULL DEFAULT '',
--     "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
--     "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
--     "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
-- );

-- CREATE TABLE IF NOT EXISTS valid_grocery_sections (
--     "id" CHAR(27) NOT NULL PRIMARY KEY,
--     "name" TEXT NOT NULL,
--     "description" TEXT NOT NULL,
--     "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
--     "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
--     "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
--     "created_by_user" CHAR(27) NOT NULL REFERENCES users("id") ON DELETE CASCADE,
--     UNIQUE("name", "belongs_to_store", "archived_at")
-- );

-- CREATE TABLE IF NOT EXISTS valid_ingredient_grocery_sections (
--     "id" CHAR(27) NOT NULL PRIMARY KEY,
--     "notes" TEXT NOT NULL DEFAULT '',
--     "valid_ingredient_id" CHAR(27) NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
--     "valid_grocery_section_id" CHAR(27) NOT NULL REFERENCES valid_grocery_sections("id") ON DELETE CASCADE,
--     "grocery_store_id" CHAR(27) NOT NULL REFERENCES grocery_stores("id") ON DELETE CASCADE,
--     "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
--     "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
--     "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
--     UNIQUE("valid_ingredient_id", "grocery_store_id", "valid_grocery_section_id","archived_at")
-- );

-- -- #335
-- CREATE TABLE IF NOT EXISTS valid_allergens (
--     "id" CHAR(27) NOT NULL PRIMARY KEY,
--     "name" TEXT NOT NULL,
--     "description" TEXT NOT NULL,
--     "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
--     "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
--     "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
--     UNIQUE(name, archived_at)
-- );

-- CREATE TABLE IF NOT EXISTS household_user_allergens (
--     "id" CHAR(27) NOT NULL PRIMARY KEY,
--     "household_member_id" CHAR(27) NOT NULL REFERENCES household_user_memberships("id") ON DELETE CASCADE,
--     "valid_allergen_id" CHAR(27) NOT NULL REFERENCES valid_allergens("id") ON DELETE CASCADE,
--     "created_on" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
--     "last_updated_on" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
--     "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
--     UNIQUE(household_member_id, valid_allergen_id, archived_at)
-- );

-- -- #394
-- CREATE TABLE IF NOT EXISTS valid_package_quantities (
--     id CHAR(27) NOT NULL PRIMARY KEY,
--     name TEXT NOT NULL,
--     description TEXT NOT NULL,
--     valid_measurement_unit_id CHAR(27) NOT NULL REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
--     quantity NUMERIC(15, 2),
--     created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
--     last_updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
--     archived_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
--     created_by_user CHAR(27) NOT NULL REFERENCES users("id") ON DELETE CASCADE,
--     UNIQUE("name", "belongs_to_store", "archived_at")
-- );

-- -- #394
-- CREATE TABLE IF NOT EXISTS valid_ingredient_package_quantities (
--     id CHAR(27) NOT NULL PRIMARY KEY,
--     notes text DEFAULT ''::text NOT NULL,
--     valid_ingredient_id CHAR(27) NOT NULL REFERENCES valid_ingredients(id) ON DELETE CASCADE,
--     valid_package_quantity_id CHAR(27) NOT NULL REFERENCES valid_package_quantities(id) ON DELETE CASCADE,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
--     last_updated_at TIMESTAMP WITH TIME ZONE,
--     archived_at TIMESTAMP WITH TIME ZONE,
--     UNIQUE (valid_ingredient_id, valid_package_quantity_id, archived_at)
-- );

