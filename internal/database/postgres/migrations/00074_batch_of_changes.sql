ALTER TABLE valid_ingredients ADD COLUMN "is_starch" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE valid_ingredients ADD COLUMN "is_protein" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE valid_ingredients ADD COLUMN "is_grain" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE valid_ingredients ADD COLUMN "is_fruit" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE valid_ingredients ADD COLUMN "is_salt" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE valid_ingredients ADD COLUMN "is_fat" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE valid_ingredients ADD COLUMN "is_acid" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE valid_ingredients ADD COLUMN "is_heat" BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE recipes ADD COLUMN "eligible_for_meals" BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE meals ADD COLUMN "eligible_for_meals_plans" BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE recipe_steps ADD COLUMN "start_timer_automatically" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE recipe_step_ingredients ADD COLUMN "product_of_recipe_id" TEXT REFERENCES recipes("id") ON DELETE CASCADE;

ALTER TABLE recipe_prep_tasks ALTER COLUMN minimum_storage_temperature_in_celsius DROP DEFAULT;
ALTER TABLE recipe_prep_tasks ALTER COLUMN minimum_storage_temperature_in_celsius DROP NOT NULL;

ALTER TABLE recipe_prep_tasks ALTER COLUMN maximum_storage_temperature_in_celsius DROP DEFAULT;
ALTER TABLE recipe_prep_tasks ALTER COLUMN maximum_storage_temperature_in_celsius DROP NOT NULL;

ALTER TABLE recipe_prep_tasks ALTER COLUMN maximum_time_buffer_before_recipe_in_seconds DROP DEFAULT;
ALTER TABLE recipe_prep_tasks ALTER COLUMN maximum_time_buffer_before_recipe_in_seconds DROP NOT NULL;

CREATE TYPE setting_type AS ENUM (
    'user',
    'household',
    'membership'
);

CREATE TABLE IF NOT EXISTS valid_settings (
    "id" TEXT NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL DEFAULT '',
    "type" setting_type NOT NULL DEFAULT 'user',
    "description" TEXT NOT NULL DEFAULT '',
    "default_value" TEXT,
    "admins_only" BOOLEAN NOT NULL DEFAULT 'true',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    UNIQUE("name")
);

CREATE TABLE IF NOT EXISTS setting_configurations (
    "id" TEXT NOT NULL PRIMARY KEY,
    "value" TEXT NOT NULL DEFAULT '',
    "notes" TEXT NOT NULL DEFAULT '',
    "setting_id" TEXT NOT NULL REFERENCES valid_settings("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "user_id" TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "household_id" TEXT NOT NULL REFERENCES households("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    UNIQUE("user_id", "household_id", "setting_id")
);