ALTER TABLE api_clients ALTER COLUMN id TYPE TEXT;
ALTER TABLE api_clients ALTER COLUMN belongs_to_user TYPE TEXT;
ALTER TABLE household_invitations ALTER COLUMN id TYPE TEXT;
ALTER TABLE household_invitations ALTER COLUMN destination_household TYPE TEXT;
ALTER TABLE household_invitations ALTER COLUMN to_user TYPE TEXT;
ALTER TABLE household_invitations ALTER COLUMN from_user TYPE TEXT;
ALTER TABLE household_user_memberships ALTER COLUMN id TYPE TEXT;
ALTER TABLE household_user_memberships ALTER COLUMN belongs_to_household TYPE TEXT;
ALTER TABLE household_user_memberships ALTER COLUMN belongs_to_user TYPE TEXT;
ALTER TABLE households ALTER COLUMN id TYPE TEXT;
ALTER TABLE households ALTER COLUMN belongs_to_user TYPE TEXT;
ALTER TABLE meal_components ALTER COLUMN id TYPE TEXT;
ALTER TABLE meal_components ALTER COLUMN meal_id TYPE TEXT;
ALTER TABLE meal_components ALTER COLUMN recipe_id TYPE TEXT;
ALTER TABLE meal_plan_events ALTER COLUMN id TYPE TEXT;
ALTER TABLE meal_plan_events ALTER COLUMN belongs_to_meal_plan TYPE TEXT;
ALTER TABLE meal_plan_grocery_list_items ALTER COLUMN id TYPE TEXT;
ALTER TABLE meal_plan_grocery_list_items ALTER COLUMN valid_ingredient TYPE TEXT;
ALTER TABLE meal_plan_grocery_list_items ALTER COLUMN valid_measurement_unit TYPE TEXT;
ALTER TABLE meal_plan_grocery_list_items ALTER COLUMN purchased_measurement_unit TYPE TEXT;
ALTER TABLE meal_plan_grocery_list_items ALTER COLUMN belongs_to_meal_plan TYPE TEXT;
ALTER TABLE meal_plan_option_votes ALTER COLUMN id TYPE TEXT;
ALTER TABLE meal_plan_option_votes ALTER COLUMN by_user TYPE TEXT;
ALTER TABLE meal_plan_option_votes ALTER COLUMN belongs_to_meal_plan_option TYPE TEXT;
ALTER TABLE meal_plan_options ALTER COLUMN id TYPE TEXT;
ALTER TABLE meal_plan_options ALTER COLUMN meal_id TYPE TEXT;
ALTER TABLE meal_plan_options ALTER COLUMN assigned_cook TYPE TEXT;
ALTER TABLE meal_plan_options ALTER COLUMN assigned_dishwasher TYPE TEXT;
ALTER TABLE meal_plan_options ALTER COLUMN belongs_to_meal_plan_event TYPE TEXT;
ALTER TABLE meal_plan_tasks ALTER COLUMN id TYPE TEXT;
ALTER TABLE meal_plan_tasks ALTER COLUMN belongs_to_meal_plan_option TYPE TEXT;
ALTER TABLE meal_plan_tasks ALTER COLUMN belongs_to_recipe_prep_task TYPE TEXT;
ALTER TABLE meal_plan_tasks ALTER COLUMN assigned_to_user TYPE TEXT;
ALTER TABLE meal_plans ALTER COLUMN id TYPE TEXT;
ALTER TABLE meal_plans ALTER COLUMN belongs_to_household TYPE TEXT;
ALTER TABLE meals ALTER COLUMN id TYPE TEXT;
ALTER TABLE meals ALTER COLUMN created_by_user TYPE TEXT;
ALTER TABLE password_reset_tokens ALTER COLUMN id TYPE TEXT;
ALTER TABLE password_reset_tokens ALTER COLUMN belongs_to_user TYPE TEXT;
ALTER TABLE recipe_media ALTER COLUMN id TYPE TEXT;
ALTER TABLE recipe_media ALTER COLUMN belongs_to_recipe TYPE TEXT;
ALTER TABLE recipe_media ALTER COLUMN belongs_to_recipe_step TYPE TEXT;
ALTER TABLE recipe_prep_task_steps ALTER COLUMN id TYPE TEXT;
ALTER TABLE recipe_prep_task_steps ALTER COLUMN belongs_to_recipe_step TYPE TEXT;
ALTER TABLE recipe_prep_task_steps ALTER COLUMN belongs_to_recipe_prep_task TYPE TEXT;
ALTER TABLE recipe_prep_tasks ALTER COLUMN id TYPE TEXT;
ALTER TABLE recipe_prep_tasks ALTER COLUMN belongs_to_recipe TYPE TEXT;
ALTER TABLE recipe_step_ingredients ALTER COLUMN id TYPE TEXT;
ALTER TABLE recipe_step_ingredients ALTER COLUMN ingredient_id TYPE TEXT;
ALTER TABLE recipe_step_ingredients ALTER COLUMN belongs_to_recipe_step TYPE TEXT;
ALTER TABLE recipe_step_ingredients ALTER COLUMN recipe_step_product_id TYPE TEXT;
ALTER TABLE recipe_step_ingredients ALTER COLUMN measurement_unit TYPE TEXT;
ALTER TABLE recipe_step_instruments ALTER COLUMN id TYPE TEXT;
ALTER TABLE recipe_step_instruments ALTER COLUMN instrument_id TYPE TEXT;
ALTER TABLE recipe_step_instruments ALTER COLUMN belongs_to_recipe_step TYPE TEXT;
ALTER TABLE recipe_step_instruments ALTER COLUMN recipe_step_product_id TYPE TEXT;
ALTER TABLE recipe_step_products ALTER COLUMN id TYPE TEXT;
ALTER TABLE recipe_step_products ALTER COLUMN belongs_to_recipe_step TYPE TEXT;
ALTER TABLE recipe_step_products ALTER COLUMN measurement_unit TYPE TEXT;
ALTER TABLE recipe_steps ALTER COLUMN id TYPE TEXT;
ALTER TABLE recipe_steps ALTER COLUMN preparation_id TYPE TEXT;
ALTER TABLE recipe_steps ALTER COLUMN belongs_to_recipe TYPE TEXT;
ALTER TABLE recipes ALTER COLUMN id TYPE TEXT;
ALTER TABLE recipes ALTER COLUMN created_by_user TYPE TEXT;
ALTER TABLE users ALTER COLUMN id TYPE TEXT;
ALTER TABLE valid_ingredient_measurement_units ALTER COLUMN id TYPE TEXT;
ALTER TABLE valid_ingredient_measurement_units ALTER COLUMN valid_ingredient_id TYPE TEXT;
ALTER TABLE valid_ingredient_measurement_units ALTER COLUMN valid_measurement_unit_id TYPE TEXT;
ALTER TABLE valid_ingredient_preparations ALTER COLUMN id TYPE TEXT;
ALTER TABLE valid_ingredient_preparations ALTER COLUMN valid_preparation_id TYPE TEXT;
ALTER TABLE valid_ingredient_preparations ALTER COLUMN valid_ingredient_id TYPE TEXT;
ALTER TABLE valid_ingredient_states ALTER COLUMN id TYPE TEXT;
ALTER TABLE valid_ingredients ALTER COLUMN id TYPE TEXT;
ALTER TABLE valid_instruments ALTER COLUMN id TYPE TEXT;
ALTER TABLE valid_measurement_conversions ALTER COLUMN id TYPE TEXT;
ALTER TABLE valid_measurement_conversions ALTER COLUMN from_unit TYPE TEXT;
ALTER TABLE valid_measurement_conversions ALTER COLUMN to_unit TYPE TEXT;
ALTER TABLE valid_measurement_conversions ALTER COLUMN only_for_ingredient TYPE TEXT;
ALTER TABLE valid_measurement_units ALTER COLUMN id TYPE TEXT;
ALTER TABLE valid_preparation_instruments ALTER COLUMN id TYPE TEXT;
ALTER TABLE valid_preparation_instruments ALTER COLUMN valid_preparation_id TYPE TEXT;
ALTER TABLE valid_preparation_instruments ALTER COLUMN valid_instrument_id TYPE TEXT;
ALTER TABLE valid_preparations ALTER COLUMN id TYPE TEXT;
ALTER TABLE webhook_trigger_events ALTER COLUMN id TYPE TEXT;
ALTER TABLE webhook_trigger_events ALTER COLUMN belongs_to_webhook TYPE TEXT;
ALTER TABLE webhooks ALTER COLUMN id TYPE TEXT;
ALTER TABLE webhooks ALTER COLUMN belongs_to_household TYPE TEXT;

CREATE INDEX IF NOT EXISTS recipe_step_ingredients_product_of_recipe_step ON recipe_step_ingredients (recipe_step_product_id);
CREATE INDEX IF NOT EXISTS password_reset_token_belongs_to_user ON password_reset_tokens (belongs_to_user);
CREATE INDEX IF NOT EXISTS valid_preparation_instruments_valid_preparation_index ON valid_preparation_instruments (valid_preparation_id);
CREATE INDEX IF NOT EXISTS valid_preparation_instruments_valid_instrument_index ON valid_preparation_instruments (valid_instrument_id);
CREATE INDEX IF NOT EXISTS valid_ingredient_measurement_units_valid_ingredient_id_index ON valid_ingredient_measurement_units (valid_ingredient_id);
CREATE INDEX IF NOT EXISTS valid_ingredient_measurement_units_valid_measurement_unit_id_index ON valid_ingredient_measurement_units (valid_measurement_unit_id);
CREATE INDEX IF NOT EXISTS recipe_step_ingredients_measurement_unit_index ON recipe_step_ingredients (measurement_unit);
CREATE INDEX IF NOT EXISTS recipe_step_products_measurement_unit_index ON recipe_step_products (measurement_unit);
CREATE INDEX IF NOT EXISTS recipe_step_instruments_recipe_step_product_id_index ON recipe_step_instruments (recipe_step_product_id);
CREATE INDEX IF NOT EXISTS recipe_step_instruments_instrument_id_index ON recipe_step_instruments (instrument_id);
CREATE INDEX IF NOT EXISTS meal_plan_options_assigned_cook_index ON meal_plan_options (assigned_cook);
CREATE INDEX IF NOT EXISTS meal_plan_options_assigned_dishwasher_index ON meal_plan_options (assigned_dishwasher);
CREATE INDEX IF NOT EXISTS meal_plan_events_belongs_to_meal_pla_index ON meal_plan_events (belongs_to_meal_plan);
CREATE INDEX IF NOT EXISTS meal_plan_options_belongs_to_meal_plan_even_index ON meal_plan_options (belongs_to_meal_plan_event);
CREATE INDEX IF NOT EXISTS recipe_prep_tasks_belongs_to_recipe_index ON recipe_prep_tasks (belongs_to_recipe);
CREATE INDEX IF NOT EXISTS recipe_prep_task_steps_belongs_to_recipe_step_index ON recipe_prep_task_steps (belongs_to_recipe_step);
CREATE INDEX IF NOT EXISTS recipe_prep_task_steps_belongs_to_recipe_prep_task_index ON recipe_prep_task_steps (belongs_to_recipe_prep_task);
CREATE INDEX IF NOT EXISTS meal_plan_tasks_belongs_to_meal_plan_option_index ON meal_plan_tasks (belongs_to_meal_plan_option);
CREATE INDEX IF NOT EXISTS meal_plan_tasks_belongs_to_recipe_prep_task_index ON meal_plan_tasks (belongs_to_recipe_prep_task);
CREATE INDEX IF NOT EXISTS meal_plan_tasks_assigned_to_user_index ON meal_plan_tasks (assigned_to_user);
CREATE INDEX IF NOT EXISTS meal_plan_grocery_list_items_valid_ingredient_index ON meal_plan_grocery_list_items (valid_ingredient);
CREATE INDEX IF NOT EXISTS meal_plan_grocery_list_items_valid_measurement_unit_index ON meal_plan_grocery_list_items (valid_measurement_unit);
CREATE INDEX IF NOT EXISTS meal_plan_grocery_list_items_purchased_measurement_unit_index ON meal_plan_grocery_list_items (purchased_measurement_unit);
CREATE INDEX IF NOT EXISTS meal_plan_grocery_list_items_belongs_to_meal_pla_index ON meal_plan_grocery_list_items (belongs_to_meal_plan);
CREATE INDEX IF NOT EXISTS valid_measurement_conversions_from_unit_index ON valid_measurement_conversions (from_unit);
CREATE INDEX IF NOT EXISTS valid_measurement_conversions_to_unit_index ON valid_measurement_conversions (to_unit);
CREATE INDEX IF NOT EXISTS valid_measurement_conversions_only_for_ingredient_index ON valid_measurement_conversions (only_for_ingredient);
CREATE INDEX IF NOT EXISTS recipe_media_belongs_to_recipe_index ON recipe_media (belongs_to_recipe);
CREATE INDEX IF NOT EXISTS recipe_media_belongs_to_recipe_step_index ON recipe_media (belongs_to_recipe_step);
CREATE INDEX IF NOT EXISTS webhook_trigger_events_belongs_to_webhook_index ON webhook_trigger_events (belongs_to_webhook);

-- CREATE TYPE ingredient_attribute_type AS ENUM (
--     'texture',
--     'consistency',
--     'color',
--     'appearance',
--     'odor',
--     'taste',
--     'sound',
--     'other'
-- );

-- ALTER TABLE valid_ingredient_states ADD COLUMN attribute_type ingredient_attribute_type NOT NULL DEFAULT 'other';

CREATE TABLE IF NOT EXISTS recipe_step_conditions (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_recipe_step TEXT NOT NULL REFERENCES recipe_steps(id) ON DELETE CASCADE,
    ingredient_state TEXT NOT NULL REFERENCES valid_ingredient_states(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS recipe_step_condition_ingredients (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_recipe_step_condition TEXT NOT NULL REFERENCES recipe_step_conditions(id) ON DELETE CASCADE,
    recipe_step_ingredient TEXT NOT NULL REFERENCES recipe_step_ingredients(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);
