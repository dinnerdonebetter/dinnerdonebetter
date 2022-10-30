--
-- PostgreSQL database dump
--

-- Dumped from database version 10.21 (Debian 10.21-1.pgdg90+1)
-- Dumped by pg_dump version 13.7 (Ubuntu 13.7-0ubuntu0.21.10.1)

CREATE TYPE grocery_list_item_status AS ENUM (
    'unknown',
    'already owned',
    'needs',
    'unavailable',
    'acquired'
);

CREATE TYPE invitation_state AS ENUM (
    'pending',
    'cancelled',
    'accepted',
    'rejected'
);

CREATE TYPE meal_name AS ENUM (
    'breakfast',
    'second_breakfast',
    'brunch',
    'lunch',
    'supper',
    'dinner'
);

CREATE TYPE meal_plan_status AS ENUM (
    'awaiting_votes',
    'finalized'
);

CREATE TYPE prep_step_status AS ENUM (
    'unfinished',
    'postponed',
    'ignored',
    'canceled',
    'finished'
);

CREATE TYPE recipe_step_product_type AS ENUM (
    'ingredient',
    'instrument'
);

CREATE TYPE storage_container_type AS ENUM (
    'uncovered',
    'covered',
    'on a wire rack',
    'in an airtight container'
);

CREATE TYPE time_zone AS ENUM (
    'UTC',
    'US/Pacific',
    'US/Mountain',
    'US/Central',
    'US/Eastern'
);

CREATE TYPE webhook_event AS ENUM (
    'webhook_created',
    'webhook_updated',
    'webhook_archived'
);

CREATE TABLE IF NOT EXISTS api_clients (
    id CHAR(20) NOT NULL PRIMARY KEY,
    name text DEFAULT ''::text,
    client_id text NOT NULL,
    secret_key bytea NOT NULL,
    permissions bigint DEFAULT 0 NOT NULL,
    admin_permissions bigint DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    belongs_to_user CHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS household_invitations (
    id CHAR(20) NOT NULL PRIMARY KEY,
    destination_household CHAR(20) NOT NULL REFERENCES households(id) ON DELETE CASCADE,
    to_email text NOT NULL,
    to_user CHAR(20),
    from_user CHAR(20) NOT NULL,
    status invitation_state DEFAULT 'pending'::invitation_state NOT NULL,
    note text DEFAULT ''::text NOT NULL,
    status_note text DEFAULT ''::text NOT NULL,
    token text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    UNIQUE(to_user, from_user, destination_household)
);

CREATE INDEX household_invitations_destination_household ON household_invitations USING btree (destination_household);

CREATE TABLE IF NOT EXISTS household_user_memberships (
    id CHAR(20) NOT NULL PRIMARY KEY,
    belongs_to_household CHAR(20) NOT NULL REFERENCES households(id) ON DELETE CASCADE,
    belongs_to_user CHAR(20) NOT NULL,
    default_household boolean DEFAULT false NOT NULL,
    household_roles text DEFAULT 'household_user'::text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    UNIQUE(belongs_to_household, belongs_to_user)
);

CREATE INDEX household_user_memberships_belongs_to_household ON household_user_memberships USING btree (belongs_to_household);
CREATE INDEX household_user_memberships_belongs_to_user ON household_user_memberships USING btree (belongs_to_user);

CREATE TABLE IF NOT EXISTS households (
    id CHAR(20) NOT NULL PRIMARY KEY,
    name text NOT NULL,
    billing_status text DEFAULT 'unpaid'::text NOT NULL,
    contact_email text DEFAULT ''::text NOT NULL,
    contact_phone text DEFAULT ''::text NOT NULL,
    payment_processor_customer_id text DEFAULT ''::text NOT NULL,
    subscription_plan_id text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    belongs_to_user CHAR(20) NOT NULL,
    time_zone time_zone DEFAULT 'US/Central'::time_zone NOT NULL
);

CREATE TABLE IF NOT EXISTS meal_plan_events (
    id CHAR(20) NOT NULL PRIMARY KEY,
    notes text DEFAULT ''::text NOT NULL,
    starts_at timestamp with time zone DEFAULT now() NOT NULL,
    ends_at timestamp with time zone DEFAULT now() NOT NULL,
    meal_name meal_name NOT NULL,
    belongs_to_meal_plan CHAR(20) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone
                        );

CREATE TABLE IF NOT EXISTS meal_plan_grocery_list_items (
    id CHAR(20) NOT NULL PRIMARY KEY,
    valid_ingredient CHAR(20) NOT NULL,
    valid_measurement_unit CHAR(20) NOT NULL,
    minimum_quantity_needed integer NOT NULL,
    maximum_quantity_needed integer NOT NULL,
    quantity_purchased integer,
    purchased_measurement_unit CHAR(20),
    purchased_upc text,
    purchase_price integer,
    status_explanation text DEFAULT ''::text NOT NULL,
    status grocery_list_item_status DEFAULT 'unknown'::grocery_list_item_status NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    completed_at timestamp with time zone,
    belongs_to_meal_plan CHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS meal_plan_option_votes (
    id CHAR(20) NOT NULL PRIMARY KEY,
    rank integer NOT NULL,
    abstain boolean NOT NULL,
    notes text NOT NULL,
    by_user CHAR(20) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    belongs_to_meal_plan_option CHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS meal_plan_options (
    id CHAR(20) NOT NULL PRIMARY KEY,
    meal_id CHAR(20) NOT NULL,
    notes text NOT NULL,
    chosen boolean DEFAULT false NOT NULL,
    tiebroken boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    assigned_cook CHAR(20),
    assigned_dishwasher CHAR(20),
    belongs_to_meal_plan_event CHAR(20)
);

CREATE TABLE IF NOT EXISTS meal_plan_tasks (
    id CHAR(20) NOT NULL PRIMARY KEY,
    belongs_to_meal_plan_option CHAR(20) NOT NULL,
    belongs_to_recipe_prep_task CHAR(20) NOT NULL,
    creation_explanation text DEFAULT ''::text NOT NULL,
    status_explanation text DEFAULT ''::text NOT NULL,
    status prep_step_status DEFAULT 'unfinished'::prep_step_status NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    assigned_to_user CHAR(20),
    completed_at timestamp with time zone
);

CREATE TABLE IF NOT EXISTS meal_plans (
    id CHAR(20) NOT NULL PRIMARY KEY,
    notes text NOT NULL,
    status meal_plan_status DEFAULT 'awaiting_votes'::meal_plan_status NOT NULL,
    voting_deadline timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    belongs_to_household CHAR(20) NOT NULL REFERENCES households(id) ON DELETE CASCADE,
    grocery_list_initialized boolean DEFAULT false NOT NULL,
    tasks_created boolean DEFAULT false NOT NULL
);

CREATE TABLE IF NOT EXISTS meal_recipes (
    id CHAR(20) NOT NULL PRIMARY KEY,
    meal_id CHAR(20) NOT NULL,
    recipe_id CHAR(20) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone
);

CREATE TABLE IF NOT EXISTS meals (
    id CHAR(20) NOT NULL PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    created_by_user CHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id CHAR(20) NOT NULL PRIMARY KEY,
    token text NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    redeemed_at timestamp with time zone,
    belongs_to_user CHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_media (
    id CHAR(20) NOT NULL PRIMARY KEY,
    belongs_to_recipe CHAR(20),
    belongs_to_recipe_step CHAR(20),
    mime_type text NOT NULL,
    internal_path text NOT NULL,
    external_path text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    index integer DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_prep_task_steps (
    id CHAR(20) NOT NULL PRIMARY KEY,
    satisfies_recipe_step boolean DEFAULT false NOT NULL,
    belongs_to_recipe_step CHAR(20) NOT NULL,
    belongs_to_recipe_prep_task CHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_prep_tasks (
    id CHAR(20) NOT NULL PRIMARY KEY,
    notes text DEFAULT ''::text NOT NULL,
    explicit_storage_instructions text DEFAULT ''::text NOT NULL,
    minimum_time_buffer_before_recipe_in_seconds integer NOT NULL,
    maximum_time_buffer_before_recipe_in_seconds integer,
    storage_type storage_container_type,
    minimum_storage_temperature_in_celsius integer,
    maximum_storage_temperature_in_celsius integer,
    belongs_to_recipe CHAR(20) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone
);

CREATE TABLE IF NOT EXISTS recipe_step_ingredients (
    id CHAR(20) NOT NULL PRIMARY KEY,
    ingredient_id CHAR(20),
    minimum_quantity_value double precision NOT NULL,
    quantity_notes text NOT NULL,
    product_of_recipe_step boolean NOT NULL,
    ingredient_notes text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    belongs_to_recipe_step CHAR(20) NOT NULL,
    name text NOT NULL,
    recipe_step_product_id CHAR(20),
    maximum_quantity_value double precision NOT NULL,
    measurement_unit CHAR(20),
    optional boolean DEFAULT false NOT NULL,
    CONSTRAINT valid_instrument_or_product CHECK (((recipe_step_product_id IS NOT NULL) OR (ingredient_id IS NOT NULL)))
);

CREATE TABLE IF NOT EXISTS recipe_step_instruments (
    id CHAR(20) NOT NULL PRIMARY KEY,
    instrument_id CHAR(20),
    notes text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    belongs_to_recipe_step CHAR(20) NOT NULL,
    preference_rank integer NOT NULL,
    recipe_step_product_id CHAR(20),
    product_of_recipe_step boolean DEFAULT false NOT NULL,
    name text DEFAULT ''::text NOT NULL,
    optional boolean DEFAULT false NOT NULL,
    minimum_quantity integer DEFAULT 1 NOT NULL,
    maximum_quantity integer DEFAULT 1 NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_step_products (
    id CHAR(20) NOT NULL PRIMARY KEY,
    name text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    belongs_to_recipe_step CHAR(20) NOT NULL,
    quantity_notes text NOT NULL,
    minimum_quantity_value double precision DEFAULT 0 NOT NULL,
    maximum_quantity_value double precision DEFAULT 0 NOT NULL,
    measurement_unit CHAR(20),
    type recipe_step_product_type NOT NULL,
    compostable boolean DEFAULT false NOT NULL,
    maximum_storage_duration_in_seconds integer,
    minimum_storage_temperature_in_celsius double precision,
    maximum_storage_temperature_in_celsius double precision,
    storage_instructions text DEFAULT ''::text NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_steps (
    id CHAR(20) NOT NULL PRIMARY KEY,
    index integer NOT NULL,
    preparation_id CHAR(20) NOT NULL,
    minimum_estimated_time_in_seconds bigint,
    maximum_estimated_time_in_seconds bigint,
    minimum_temperature_in_celsius integer,
    notes text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    belongs_to_recipe CHAR(20) NOT NULL,
    optional boolean DEFAULT false NOT NULL,
    maximum_temperature_in_celsius integer,
    explicit_instructions text DEFAULT ''::text NOT NULL
);

CREATE TABLE IF NOT EXISTS recipes (
    id CHAR(20) NOT NULL PRIMARY KEY,
    name text NOT NULL,
    source text NOT NULL,
    description text NOT NULL,
    inspired_by_recipe_id text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    created_by_user CHAR(20) NOT NULL,
    yields_portions integer DEFAULT 1 NOT NULL,
    seal_of_approval boolean DEFAULT false NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
   token text NOT NULL,
   data bytea NOT NULL,
   expiry timestamp with time zone NOT NULL,
   created_on bigint DEFAULT date_part('epoch'::text, now()) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id CHAR(20) NOT NULL PRIMARY KEY,
    username text NOT NULL,
    avatar_src text,
    email_address text NOT NULL,
    hashed_password text NOT NULL,
    password_last_changed_at timestamp with time zone,
    requires_password_change boolean DEFAULT false NOT NULL,
    two_factor_secret text NOT NULL,
    two_factor_secret_verified_at timestamp with time zone,
    service_roles text DEFAULT 'service_user'::text NOT NULL,
    user_account_status text DEFAULT 'unverified'::text NOT NULL,
    user_account_status_explanation text DEFAULT ''::text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    birth_day smallint,
    birth_month smallint
);

CREATE TABLE IF NOT EXISTS valid_ingredient_measurement_units (
    id CHAR(20) NOT NULL PRIMARY KEY,
    notes text DEFAULT ''::text NOT NULL,
    valid_ingredient_id CHAR(20) NOT NULL,
    valid_measurement_unit_id CHAR(20) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    minimum_allowable_quantity double precision DEFAULT 0 NOT NULL,
    maximum_allowable_quantity double precision DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS valid_ingredient_preparations (
    id CHAR(20) NOT NULL PRIMARY KEY,
    notes text NOT NULL,
    valid_preparation_id CHAR(20) NOT NULL,
    valid_ingredient_id CHAR(20) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone
);

CREATE TABLE IF NOT EXISTS valid_ingredients (
    id CHAR(20) NOT NULL PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL,
    warning text NOT NULL,
    contains_egg boolean NOT NULL,
    contains_dairy boolean NOT NULL,
    contains_peanut boolean NOT NULL,
    contains_tree_nut boolean NOT NULL,
    contains_soy boolean NOT NULL,
    contains_wheat boolean NOT NULL,
    contains_shellfish boolean NOT NULL,
    contains_sesame boolean NOT NULL,
    contains_fish boolean NOT NULL,
    contains_gluten boolean NOT NULL,
    animal_flesh boolean NOT NULL,
    volumetric boolean NOT NULL,
    icon_path text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    is_liquid boolean DEFAULT false,
    animal_derived boolean DEFAULT false NOT NULL,
    plural_name text DEFAULT ''::text NOT NULL,
    restrict_to_preparations boolean DEFAULT false NOT NULL,
    minimum_ideal_storage_temperature_in_celsius double precision,
    maximum_ideal_storage_temperature_in_celsius double precision,
    storage_instructions text DEFAULT ''::text NOT NULL
);

CREATE TABLE IF NOT EXISTS valid_instruments (
    id CHAR(20) NOT NULL PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL,
    icon_path text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    plural_name text DEFAULT ''::text NOT NULL,
    usable_for_storage boolean DEFAULT false NOT NULL
);

CREATE TABLE IF NOT EXISTS valid_measurement_conversions (
    id CHAR(20) NOT NULL PRIMARY KEY,
    from_unit CHAR(20) NOT NULL,
    to_unit CHAR(20) NOT NULL,
    only_for_ingredient CHAR(20),
    modifier bigint NOT NULL,
    notes text DEFAULT ''::text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone
);

CREATE TABLE IF NOT EXISTS valid_measurement_units (
    id CHAR(20) NOT NULL PRIMARY KEY,
    name text NOT NULL,
    description text DEFAULT ''::text NOT NULL,
    icon_path text DEFAULT ''::text NOT NULL,
    volumetric boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    universal boolean DEFAULT false NOT NULL,
    metric boolean DEFAULT false NOT NULL,
    imperial boolean DEFAULT false NOT NULL,
    plural_name text DEFAULT ''::text NOT NULL
);

CREATE TABLE IF NOT EXISTS valid_preparation_instruments (
    id CHAR(20) NOT NULL PRIMARY KEY,
    notes text DEFAULT ''::text NOT NULL,
    valid_preparation_id CHAR(20) NOT NULL,
    valid_instrument_id CHAR(20) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone
);

CREATE TABLE IF NOT EXISTS valid_preparations (
    id CHAR(20) NOT NULL PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL,
    icon_path text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    yields_nothing boolean DEFAULT false NOT NULL,
    restrict_to_ingredients boolean DEFAULT false NOT NULL,
    zero_ingredients_allowable boolean DEFAULT false NOT NULL,
    past_tense text DEFAULT ''::text NOT NULL
);

CREATE TABLE IF NOT EXISTS webhook_trigger_events (
    id CHAR(20) NOT NULL PRIMARY KEY,
    trigger_event webhook_event NOT NULL,
    belongs_to_webhook CHAR(20) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    archived_at timestamp with time zone
);

CREATE TABLE IF NOT EXISTS webhooks (
    id CHAR(20) NOT NULL PRIMARY KEY,
    name text NOT NULL,
    content_type text NOT NULL,
    url text NOT NULL,
    method text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_updated_at timestamp with time zone,
    archived_at timestamp with time zone,
    belongs_to_household CHAR(20) NOT NULL REFERENCES households(id) ON DELETE CASCADE
);

ALTER TABLE ONLY households
    ADD CONSTRAINT households_belongs_to_user_name_key UNIQUE (belongs_to_user, name);

ALTER TABLE ONLY meal_plan_option_votes
    ADD CONSTRAINT meal_plan_option_votes_by_user_belongs_to_meal_plan_option_key UNIQUE (by_user, belongs_to_meal_plan_option);

ALTER TABLE ONLY recipe_step_ingredients
    ADD CONSTRAINT recipe_step_ingredients_ingredient_id_belongs_to_recipe_ste_key UNIQUE (ingredient_id, belongs_to_recipe_step);

ALTER TABLE ONLY sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (token);

ALTER TABLE ONLY recipe_prep_task_steps
    ADD CONSTRAINT unique_recipe_step_and_prep_task UNIQUE (belongs_to_recipe_step, belongs_to_recipe_prep_task);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_username_key UNIQUE (username);

ALTER TABLE ONLY valid_ingredient_measurement_units
    ADD CONSTRAINT valid_ingredient_measurement__valid_ingredient_id_valid_mea_key UNIQUE (valid_ingredient_id, valid_measurement_unit_id, archived_at);

ALTER TABLE ONLY valid_ingredient_preparations
    ADD CONSTRAINT valid_ingredient_preparations_valid_preparation_id_valid_in_key UNIQUE (valid_preparation_id, valid_ingredient_id, archived_at);

ALTER TABLE ONLY valid_measurement_conversions
    ADD CONSTRAINT valid_measurement_conversions_from_unit_to_unit_only_for_in_key UNIQUE (from_unit, to_unit, only_for_ingredient);

ALTER TABLE ONLY valid_measurement_units
    ADD CONSTRAINT valid_measurement_units_name_key UNIQUE (name);

ALTER TABLE ONLY valid_preparation_instruments
    ADD CONSTRAINT valid_preparation_instruments_valid_preparation_id_valid_in_key UNIQUE (valid_preparation_id, valid_instrument_id, archived_at);

ALTER TABLE ONLY valid_preparations
    ADD CONSTRAINT valid_preparations_name_key UNIQUE (name);

ALTER TABLE ONLY webhook_trigger_events
    ADD CONSTRAINT webhook_trigger_events_trigger_event_belongs_to_webhook_key UNIQUE (trigger_event, belongs_to_webhook);

CREATE INDEX api_clients_belongs_to_user ON api_clients USING btree (belongs_to_user);

CREATE INDEX household_invitations_from_user ON household_invitations USING btree (from_user);

CREATE INDEX household_invitations_to_user ON household_invitations USING btree (to_user);


CREATE INDEX households_belongs_to_user ON households USING btree (belongs_to_user);

CREATE INDEX meal_plan_options_belongs_to_meal_plan_option ON meal_plan_option_votes USING btree (belongs_to_meal_plan_option);

CREATE INDEX meal_plans_belongs_to_household ON meal_plans USING btree (belongs_to_household);

CREATE INDEX meal_recipes_meal_id ON meal_recipes USING btree (meal_id);

CREATE INDEX meal_recipes_recipe_id ON meal_recipes USING btree (recipe_id);

CREATE INDEX meals_created_by_user ON meals USING btree (created_by_user);

CREATE INDEX recipe_step_products_belongs_to_recipe_step ON recipe_step_products USING btree (belongs_to_recipe_step);

CREATE INDEX recipe_steps_belongs_to_recipe ON recipe_steps USING btree (belongs_to_recipe);

CREATE INDEX recipes_created_by_user ON recipes USING btree (created_by_user);

CREATE INDEX sessions_expiry_idx ON sessions USING btree (expiry);

ALTER TABLE ONLY api_clients
    ADD CONSTRAINT api_clients_belongs_to_user_fkey FOREIGN KEY (belongs_to_user) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY recipe_step_instruments
    ADD CONSTRAINT fk_valid_recipe_step_instrument_ids FOREIGN KEY (instrument_id) REFERENCES valid_instruments(id);

ALTER TABLE ONLY household_invitations
    ADD CONSTRAINT household_invitations_from_user_fkey FOREIGN KEY (from_user) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY household_invitations
    ADD CONSTRAINT household_invitations_to_user_fkey FOREIGN KEY (to_user) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY household_user_memberships
    ADD CONSTRAINT household_user_memberships_belongs_to_user_fkey FOREIGN KEY (belongs_to_user) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY households
    ADD CONSTRAINT households_belongs_to_user_fkey FOREIGN KEY (belongs_to_user) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY meal_plan_events
    ADD CONSTRAINT meal_plan_events_belongs_to_meal_plan_fkey FOREIGN KEY (belongs_to_meal_plan) REFERENCES meal_plans(id) ON DELETE CASCADE;

ALTER TABLE ONLY meal_plan_grocery_list_items
    ADD CONSTRAINT meal_plan_grocery_list_items_belongs_to_meal_plan_fkey FOREIGN KEY (belongs_to_meal_plan) REFERENCES meal_plans(id) ON DELETE CASCADE;

ALTER TABLE ONLY meal_plan_grocery_list_items
    ADD CONSTRAINT meal_plan_grocery_list_items_purchased_measurement_unit_fkey FOREIGN KEY (purchased_measurement_unit) REFERENCES valid_measurement_units(id) ON DELETE CASCADE;

ALTER TABLE ONLY meal_plan_grocery_list_items
    ADD CONSTRAINT meal_plan_grocery_list_items_valid_ingredient_fkey FOREIGN KEY (valid_ingredient) REFERENCES valid_ingredients(id) ON DELETE CASCADE;

ALTER TABLE ONLY meal_plan_grocery_list_items
    ADD CONSTRAINT meal_plan_grocery_list_items_valid_measurement_unit_fkey FOREIGN KEY (valid_measurement_unit) REFERENCES valid_measurement_units(id) ON DELETE CASCADE;

ALTER TABLE ONLY meal_plan_option_votes
    ADD CONSTRAINT meal_plan_option_votes_belongs_to_meal_plan_option_fkey FOREIGN KEY (belongs_to_meal_plan_option) REFERENCES meal_plan_options(id) ON DELETE CASCADE;

ALTER TABLE ONLY meal_plan_option_votes
    ADD CONSTRAINT meal_plan_option_votes_by_user_fkey FOREIGN KEY (by_user) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY meal_plan_options
    ADD CONSTRAINT meal_plan_options_assigned_cook_fkey FOREIGN KEY (assigned_cook) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY meal_plan_options
    ADD CONSTRAINT meal_plan_options_assigned_dishwasher_fkey FOREIGN KEY (assigned_dishwasher) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY meal_plan_options
    ADD CONSTRAINT meal_plan_options_belongs_to_meal_plan_event_fkey FOREIGN KEY (belongs_to_meal_plan_event) REFERENCES meal_plan_events(id) ON DELETE CASCADE;

ALTER TABLE ONLY meal_plan_options
    ADD CONSTRAINT meal_plan_options_meal_id_fkey FOREIGN KEY (meal_id) REFERENCES meals(id) ON DELETE CASCADE;

ALTER TABLE ONLY meal_plan_tasks
    ADD CONSTRAINT meal_plan_tasks_assigned_to_user_fkey FOREIGN KEY (assigned_to_user) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY meal_plan_tasks
    ADD CONSTRAINT meal_plan_tasks_belongs_to_meal_plan_option_fkey FOREIGN KEY (belongs_to_meal_plan_option) REFERENCES meal_plan_options(id) ON DELETE CASCADE;

ALTER TABLE ONLY meal_plan_tasks
    ADD CONSTRAINT meal_plan_tasks_belongs_to_recipe_prep_task_fkey FOREIGN KEY (belongs_to_recipe_prep_task) REFERENCES recipe_prep_tasks(id) ON DELETE CASCADE;

ALTER TABLE ONLY meal_recipes
    ADD CONSTRAINT meal_recipes_meal_id_fkey FOREIGN KEY (meal_id) REFERENCES meals(id) ON DELETE CASCADE;

ALTER TABLE ONLY meal_recipes
    ADD CONSTRAINT meal_recipes_recipe_id_fkey FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE;

ALTER TABLE ONLY meals
    ADD CONSTRAINT meals_created_by_user_fkey FOREIGN KEY (created_by_user) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY password_reset_tokens
    ADD CONSTRAINT password_reset_tokens_belongs_to_user_fkey FOREIGN KEY (belongs_to_user) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY recipe_media
    ADD CONSTRAINT recipe_media_belongs_to_recipe_fkey FOREIGN KEY (belongs_to_recipe) REFERENCES recipes(id) ON DELETE CASCADE;

ALTER TABLE ONLY recipe_media
    ADD CONSTRAINT recipe_media_belongs_to_recipe_step_fkey FOREIGN KEY (belongs_to_recipe_step) REFERENCES recipe_steps(id) ON DELETE CASCADE;

ALTER TABLE ONLY recipe_prep_task_steps
    ADD CONSTRAINT recipe_prep_task_steps_belongs_to_recipe_prep_task_fkey FOREIGN KEY (belongs_to_recipe_prep_task) REFERENCES recipe_prep_tasks(id) ON DELETE CASCADE;

ALTER TABLE ONLY recipe_prep_task_steps
    ADD CONSTRAINT recipe_prep_task_steps_belongs_to_recipe_step_fkey FOREIGN KEY (belongs_to_recipe_step) REFERENCES recipe_steps(id) ON DELETE CASCADE;

ALTER TABLE ONLY recipe_prep_tasks
    ADD CONSTRAINT recipe_prep_tasks_belongs_to_recipe_fkey FOREIGN KEY (belongs_to_recipe) REFERENCES recipes(id) ON DELETE CASCADE;

ALTER TABLE ONLY recipe_step_ingredients
    ADD CONSTRAINT recipe_step_ingredients_belongs_to_recipe_step_fkey FOREIGN KEY (belongs_to_recipe_step) REFERENCES recipe_steps(id) ON DELETE CASCADE;

ALTER TABLE ONLY recipe_step_ingredients
    ADD CONSTRAINT recipe_step_ingredients_ingredient_id_fkey FOREIGN KEY (ingredient_id) REFERENCES valid_ingredients(id) ON DELETE CASCADE;

ALTER TABLE ONLY recipe_step_ingredients
    ADD CONSTRAINT recipe_step_ingredients_measurement_unit_fkey FOREIGN KEY (measurement_unit) REFERENCES valid_measurement_units(id);

ALTER TABLE ONLY recipe_step_ingredients
    ADD CONSTRAINT recipe_step_ingredients_recipe_step_product_id_fkey FOREIGN KEY (recipe_step_product_id) REFERENCES recipe_step_products(id) ON DELETE RESTRICT;

ALTER TABLE ONLY recipe_step_instruments
    ADD CONSTRAINT recipe_step_instruments_belongs_to_recipe_step_fkey FOREIGN KEY (belongs_to_recipe_step) REFERENCES recipe_steps(id) ON DELETE CASCADE;

ALTER TABLE ONLY recipe_step_instruments
    ADD CONSTRAINT recipe_step_instruments_recipe_step_product_id_fkey FOREIGN KEY (recipe_step_product_id) REFERENCES recipe_step_products(id);

ALTER TABLE ONLY recipe_step_products
    ADD CONSTRAINT recipe_step_products_belongs_to_recipe_step_fkey FOREIGN KEY (belongs_to_recipe_step) REFERENCES recipe_steps(id) ON DELETE CASCADE;

ALTER TABLE ONLY recipe_step_products
    ADD CONSTRAINT recipe_step_products_measurement_unit_fkey FOREIGN KEY (measurement_unit) REFERENCES valid_measurement_units(id);

ALTER TABLE ONLY recipe_steps
    ADD CONSTRAINT recipe_steps_belongs_to_recipe_fkey FOREIGN KEY (belongs_to_recipe) REFERENCES recipes(id) ON DELETE CASCADE;

ALTER TABLE ONLY recipe_steps
    ADD CONSTRAINT recipe_steps_preparation_id_fkey FOREIGN KEY (preparation_id) REFERENCES valid_preparations(id) ON DELETE CASCADE;

ALTER TABLE ONLY recipes
    ADD CONSTRAINT recipes_created_by_user_fkey FOREIGN KEY (created_by_user) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY valid_ingredient_measurement_units
    ADD CONSTRAINT valid_ingredient_measurement_uni_valid_measurement_unit_id_fkey FOREIGN KEY (valid_measurement_unit_id) REFERENCES valid_measurement_units(id) ON DELETE CASCADE;

ALTER TABLE ONLY valid_ingredient_measurement_units
    ADD CONSTRAINT valid_ingredient_measurement_units_valid_ingredient_id_fkey FOREIGN KEY (valid_ingredient_id) REFERENCES valid_ingredients(id) ON DELETE CASCADE;

ALTER TABLE ONLY valid_ingredient_preparations
    ADD CONSTRAINT valid_ingredient_preparations_valid_ingredient_id_fkey FOREIGN KEY (valid_ingredient_id) REFERENCES valid_ingredients(id) ON DELETE CASCADE;

ALTER TABLE ONLY valid_ingredient_preparations
    ADD CONSTRAINT valid_ingredient_preparations_valid_preparation_id_fkey FOREIGN KEY (valid_preparation_id) REFERENCES valid_preparations(id) ON DELETE CASCADE;

ALTER TABLE ONLY valid_measurement_conversions
    ADD CONSTRAINT valid_measurement_conversions_from_unit_fkey FOREIGN KEY (from_unit) REFERENCES valid_measurement_units(id) ON DELETE CASCADE;

ALTER TABLE ONLY valid_measurement_conversions
    ADD CONSTRAINT valid_measurement_conversions_only_for_ingredient_fkey FOREIGN KEY (only_for_ingredient) REFERENCES valid_ingredients(id) ON DELETE CASCADE;

ALTER TABLE ONLY valid_measurement_conversions
    ADD CONSTRAINT valid_measurement_conversions_to_unit_fkey FOREIGN KEY (to_unit) REFERENCES valid_measurement_units(id) ON DELETE CASCADE;

ALTER TABLE ONLY valid_preparation_instruments
    ADD CONSTRAINT valid_preparation_instruments_valid_instrument_id_fkey FOREIGN KEY (valid_instrument_id) REFERENCES valid_instruments(id) ON DELETE CASCADE;

ALTER TABLE ONLY valid_preparation_instruments
    ADD CONSTRAINT valid_preparation_instruments_valid_preparation_id_fkey FOREIGN KEY (valid_preparation_id) REFERENCES valid_preparations(id) ON DELETE CASCADE;

ALTER TABLE ONLY webhook_trigger_events
    ADD CONSTRAINT webhook_trigger_events_belongs_to_webhook_fkey FOREIGN KEY (belongs_to_webhook) REFERENCES webhooks(id) ON DELETE CASCADE;

--
-- PostgreSQL database dump complete
--
