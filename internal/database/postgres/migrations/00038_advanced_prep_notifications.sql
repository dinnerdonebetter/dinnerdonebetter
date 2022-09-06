CREATE TABLE IF NOT EXISTS advanced_prep_steps (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "belongs_to_meal_plan_option" CHAR(27) NOT NULL REFERENCES meal_plan_options("id") ON DELETE CASCADE,
    "satisfies_recipe_step" CHAR(27) NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    "cannot_complete_before" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "cannot_complete_after" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "completed_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS advanced_prep_notifications (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "advanced_prep_step_id" CHAR(27) NOT NULL REFERENCES advanced_prep_steps("id") ON DELETE CASCADE,
    "notification_sent_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "completed_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

ALTER TABLE valid_instruments ADD COLUMN "usable_for_storage" BOOLEAN NOT NULL DEFAULT 'false';
ALTER TABLE valid_ingredients ADD COLUMN "storage_instructions" TEXT NOT NULL DEFAULT '';
ALTER TABLE recipe_step_products ADD COLUMN "storage_instructions" TEXT NOT NULL DEFAULT '';
ALTER TABLE meal_plan_options ADD COLUMN "prep_steps_created" BOOLEAN NOT NULL DEFAULT 'false';
ALTER TABLE meal_plan_options ADD COLUMN "assigned_dishwasher" CHAR(27) REFERENCES users("id") ON DELETE CASCADE;

-- api_clients
ALTER TABLE api_clients ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE api_clients ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE api_clients ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE api_clients ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE api_clients ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- household_invitations
ALTER TABLE household_invitations ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE household_invitations ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE household_invitations ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE household_invitations ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE household_invitations ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- households
ALTER TABLE households ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE households ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE households ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE households ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE households ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- household_user_memberships
ALTER TABLE household_user_memberships ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE household_user_memberships ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE household_user_memberships ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE household_user_memberships ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE household_user_memberships ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- meal_plan_options
ALTER TABLE meal_plan_options ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE meal_plan_options ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE meal_plan_options ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE meal_plan_options ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE meal_plan_options ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- meal_plan_option_votes
ALTER TABLE meal_plan_option_votes ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE meal_plan_option_votes ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE meal_plan_option_votes ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE meal_plan_option_votes ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE meal_plan_option_votes ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- meal_plans
ALTER TABLE meal_plans ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE meal_plans ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE meal_plans ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE meal_plans ALTER COLUMN starts_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(starts_at);
ALTER TABLE meal_plans ALTER COLUMN ends_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(ends_at);
ALTER TABLE meal_plans ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE meal_plans ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- meal_recipes
ALTER TABLE meal_recipes ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE meal_recipes ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE meal_recipes ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE meal_recipes ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE meal_recipes ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- meals
ALTER TABLE meals ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE meals ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE meals ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE meals ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE meals ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- password_reset_tokens
ALTER TABLE password_reset_tokens ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE password_reset_tokens ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE password_reset_tokens ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE password_reset_tokens ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE password_reset_tokens ALTER COLUMN redeemed_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(redeemed_at);

-- recipes
ALTER TABLE recipes ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE recipes ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE recipes ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE recipes ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE recipes ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- recipe_step_ingredients
ALTER TABLE recipe_step_ingredients ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE recipe_step_ingredients ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE recipe_step_ingredients ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE recipe_step_ingredients ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE recipe_step_ingredients ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- recipe_step_instruments
ALTER TABLE recipe_step_instruments ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE recipe_step_instruments ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE recipe_step_instruments ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE recipe_step_instruments ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE recipe_step_instruments ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- recipe_step_products
ALTER TABLE recipe_step_products ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE recipe_step_products ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE recipe_step_products ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE recipe_step_products ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE recipe_step_products ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- recipe_steps
ALTER TABLE recipe_steps ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE recipe_steps ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE recipe_steps ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE recipe_steps ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE recipe_steps ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- users
ALTER TABLE users ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE users ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE users ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE users ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE users ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);
ALTER TABLE users ALTER COLUMN password_last_changed_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(password_last_changed_at);
ALTER TABLE users ALTER COLUMN two_factor_secret_verified_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(two_factor_secret_verified_at);

-- valid_ingredient_measurement_units
ALTER TABLE valid_ingredient_measurement_units ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE valid_ingredient_measurement_units ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE valid_ingredient_measurement_units ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE valid_ingredient_measurement_units ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE valid_ingredient_measurement_units ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- valid_ingredient_preparations
ALTER TABLE valid_ingredient_preparations ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE valid_ingredient_preparations ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE valid_ingredient_preparations ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE valid_ingredient_preparations ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE valid_ingredient_preparations ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- valid_ingredients
ALTER TABLE valid_ingredients ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE valid_ingredients ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE valid_ingredients ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE valid_ingredients ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE valid_ingredients ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- valid_instruments
ALTER TABLE valid_instruments ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE valid_instruments ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE valid_instruments ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE valid_instruments ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE valid_instruments ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- valid_measurement_units
ALTER TABLE valid_measurement_units ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE valid_measurement_units ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE valid_measurement_units ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE valid_measurement_units ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE valid_measurement_units ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- valid_preparation_instruments
ALTER TABLE valid_preparation_instruments ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE valid_preparation_instruments ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE valid_preparation_instruments ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE valid_preparation_instruments ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE valid_preparation_instruments ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- valid_preparations
ALTER TABLE valid_preparations ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE valid_preparations ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE valid_preparations ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE valid_preparations ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE valid_preparations ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);

-- webhooks
ALTER TABLE webhooks ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE webhooks ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(created_at);
ALTER TABLE webhooks ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE webhooks ALTER COLUMN last_updated_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(last_updated_at);
ALTER TABLE webhooks ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING to_timestamp(archived_at);
