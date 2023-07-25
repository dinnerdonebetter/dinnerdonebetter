--
-- Name: component_type; Type: TYPE;
--

CREATE TYPE component_type AS ENUM (
    'unspecified',
    'amuse-bouche',
    'appetizer',
    'soup',
    'main',
    'salad',
    'beverage',
    'side',
    'dessert'
);

--
-- Name: grocery_list_item_status; Type: TYPE;
--

CREATE TYPE grocery_list_item_status AS ENUM (
    'unknown',
    'already owned',
    'needs',
    'unavailable',
    'acquired'
);

--
-- Name: ingredient_attribute_type; Type: TYPE;
--

CREATE TYPE ingredient_attribute_type AS ENUM (
    'texture',
    'consistency',
    'color',
    'appearance',
    'odor',
    'taste',
    'sound',
    'temperature',
    'other'
);

--
-- Name: invitation_state; Type: TYPE;
--

CREATE TYPE invitation_state AS ENUM (
    'pending',
    'cancelled',
    'accepted',
    'rejected'
);

--
-- Name: meal_name; Type: TYPE;
--

CREATE TYPE meal_name AS ENUM (
    'breakfast',
    'second_breakfast',
    'brunch',
    'lunch',
    'supper',
    'dinner'
);

--
-- Name: meal_plan_status; Type: TYPE;
--

CREATE TYPE meal_plan_status AS ENUM (
    'awaiting_votes',
    'finalized'
);

--
-- Name: oauth2_client_token_scopes; Type: TYPE;
--

CREATE TYPE oauth2_client_token_scopes AS ENUM (
    'unknown',
    'household_member',
    'household_admin',
    'service_admin'
);

--
-- Name: prep_step_status; Type: TYPE;
--

CREATE TYPE prep_step_status AS ENUM (
    'unfinished',
    'postponed',
    'ignored',
    'canceled',
    'finished'
);

--
-- Name: recipe_ingredient_scale; Type: TYPE;
--

CREATE TYPE recipe_ingredient_scale AS ENUM (
    'linear',
    'logarithmic',
    'exponential'
);

--
-- Name: recipe_step_product_type; Type: TYPE;
--

CREATE TYPE recipe_step_product_type AS ENUM (
    'ingredient',
    'instrument',
    'vessel'
);

--
-- Name: setting_type; Type: TYPE;
--

CREATE TYPE setting_type AS ENUM (
    'user',
    'household',
    'membership'
);

--
-- Name: storage_container_type; Type: TYPE;
--

CREATE TYPE storage_container_type AS ENUM (
    'uncovered',
    'covered',
    'on a wire rack',
    'in an airtight container'
);

--
-- Name: time_zone; Type: TYPE;
--

CREATE TYPE time_zone AS ENUM (
    'UTC',
    'US/Pacific',
    'US/Mountain',
    'US/Central',
    'US/Eastern'
);

--
-- Name: valid_election_method; Type: TYPE;
--

CREATE TYPE valid_election_method AS ENUM (
    'schulze',
    'instant-runoff'
);

--
-- Name: vessel_shape; Type: TYPE;
--

CREATE TYPE vessel_shape AS ENUM (
    'hemisphere',
    'rectangle',
    'cone',
    'pyramid',
    'cylinder',
    'sphere',
    'cube',
    'other'
);

--
-- Name: webhook_event; Type: TYPE;
--

CREATE TYPE webhook_event AS ENUM (
    'webhook_created',
    'webhook_updated',
    'webhook_archived'
);

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: darwin_migrations; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS darwin_migrations (
    id integer NOT NULL,
    version real NOT NULL,
    description character varying(255) NOT NULL,
    checksum character varying(32) NOT NULL,
    applied_at integer NOT NULL,
    execution_time real NOT NULL
);

--
-- Name: household_instrument_ownerships; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS household_instrument_ownerships (
    id TEXT NOT NULL,
    notes text DEFAULT ''::TEXT NOT NULL,
    quantity integer DEFAULT 0 NOT NULL,
    valid_instrument_id TEXT NOT NULL,
    belongs_to_household TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: household_invitations; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS household_invitations (
    id TEXT NOT NULL,
    destination_household TEXT NOT NULL,
    to_email TEXT NOT NULL,
    to_user text,
    from_user TEXT NOT NULL,
    status invitation_state DEFAULT 'pending'::invitation_state NOT NULL,
    note text DEFAULT ''::TEXT NOT NULL,
    status_note text DEFAULT ''::TEXT NOT NULL,
    token TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT (now() + '7 days'::interval) NOT NULL,
    to_name text DEFAULT ''::TEXT NOT NULL
);

--
-- Name: household_user_memberships; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS household_user_memberships (
    id TEXT NOT NULL,
    belongs_to_household TEXT NOT NULL,
    belongs_to_user TEXT NOT NULL,
    default_household boolean DEFAULT false NOT NULL,
    household_role text DEFAULT 'household_user'::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: households; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS households (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    billing_status text DEFAULT 'unpaid'::TEXT NOT NULL,
    contact_phone text DEFAULT ''::TEXT NOT NULL,
    payment_processor_customer_id text DEFAULT ''::TEXT NOT NULL,
    subscription_plan_id text,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_user TEXT NOT NULL,
    time_zone time_zone DEFAULT 'US/Central'::time_zone NOT NULL,
    address_line_1 text DEFAULT ''::TEXT NOT NULL,
    address_line_2 text DEFAULT ''::TEXT NOT NULL,
    city text DEFAULT ''::TEXT NOT NULL,
    state text DEFAULT ''::TEXT NOT NULL,
    zip_code text DEFAULT ''::TEXT NOT NULL,
    country text DEFAULT ''::TEXT NOT NULL,
    latitude NUMERIC(14,11),
    longitude NUMERIC(14,11),
    last_payment_provider_sync_occurred_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: meal_components; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS meal_components (
    id TEXT NOT NULL,
    meal_id TEXT NOT NULL,
    recipe_id TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    meal_component_type component_type DEFAULT 'unspecified'::component_type NOT NULL,
    recipe_scale NUMERIC(14,2) DEFAULT 1.0 NOT NULL
);

--
-- Name: meal_plan_events; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS meal_plan_events (
    id TEXT NOT NULL,
    notes text DEFAULT ''::TEXT NOT NULL,
    starts_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    ends_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    meal_name meal_name NOT NULL,
    belongs_to_meal_plan TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: meal_plan_grocery_list_items; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS meal_plan_grocery_list_items (
    id TEXT NOT NULL,
    valid_ingredient TEXT NOT NULL,
    valid_measurement_unit TEXT NOT NULL,
    minimum_quantity_needed NUMERIC(14,2) NOT NULL,
    maximum_quantity_needed NUMERIC(14,2),
    quantity_purchased NUMERIC(14,2),
    purchased_measurement_unit text,
    purchased_upc text,
    purchase_price NUMERIC(14,2),
    status_explanation text DEFAULT ''::TEXT NOT NULL,
    status grocery_list_item_status DEFAULT 'unknown'::grocery_list_item_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_meal_plan TEXT NOT NULL
);

--
-- Name: meal_plan_option_votes; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS meal_plan_option_votes (
    id TEXT NOT NULL,
    rank integer NOT NULL,
    abstain boolean NOT NULL,
    notes TEXT NOT NULL,
    by_user TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_meal_plan_option TEXT NOT NULL
);

--
-- Name: meal_plan_options; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS meal_plan_options (
    id TEXT NOT NULL,
    meal_id TEXT NOT NULL,
    notes TEXT NOT NULL,
    chosen boolean DEFAULT false NOT NULL,
    tiebroken boolean DEFAULT false NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    assigned_cook text,
    assigned_dishwasher text,
    belongs_to_meal_plan_event text,
    meal_scale NUMERIC(14,2) DEFAULT 1.0 NOT NULL
);

--
-- Name: meal_plan_tasks; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS meal_plan_tasks (
    id TEXT NOT NULL,
    belongs_to_meal_plan_option TEXT NOT NULL,
    belongs_to_recipe_prep_task TEXT NOT NULL,
    creation_explanation text DEFAULT ''::TEXT NOT NULL,
    status_explanation text DEFAULT ''::TEXT NOT NULL,
    status prep_step_status DEFAULT 'unfinished'::prep_step_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    assigned_to_user text,
    completed_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: meal_plans; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS meal_plans (
    id TEXT NOT NULL,
    notes TEXT NOT NULL,
    status meal_plan_status DEFAULT 'awaiting_votes'::meal_plan_status NOT NULL,
    voting_deadline TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_household TEXT NOT NULL,
    grocery_list_initialized boolean DEFAULT false NOT NULL,
    tasks_created boolean DEFAULT false NOT NULL,
    election_method valid_election_method DEFAULT 'schulze'::valid_election_method NOT NULL,
    created_by_user TEXT NOT NULL
);

--
-- Name: meals; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS meals (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    created_by_user TEXT NOT NULL,
    min_estimated_portions NUMERIC(14,2) DEFAULT 1.0 NOT NULL,
    max_estimated_portions NUMERIC(14,2),
    eligible_for_meal_plans boolean DEFAULT true NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: oauth2_client_tokens; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS oauth2_client_tokens (
    id TEXT NOT NULL,
    client_id TEXT NOT NULL,
    belongs_to_user TEXT NOT NULL,
    redirect_uri text DEFAULT ''::TEXT NOT NULL,
    scope oauth2_client_token_scopes DEFAULT 'unknown'::oauth2_client_token_scopes NOT NULL,
    code text DEFAULT ''::TEXT NOT NULL,
    code_challenge text DEFAULT ''::TEXT NOT NULL,
    code_challenge_method text DEFAULT ''::TEXT NOT NULL,
    code_created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    code_expires_at TIMESTAMP WITH TIME ZONE DEFAULT (now() + '01:00:00'::interval) NOT NULL,
    access text DEFAULT ''::TEXT NOT NULL,
    access_created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    access_expires_at TIMESTAMP WITH TIME ZONE DEFAULT (now() + '01:00:00'::interval) NOT NULL,
    refresh text DEFAULT ''::TEXT NOT NULL,
    refresh_created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    refresh_expires_at TIMESTAMP WITH TIME ZONE DEFAULT (now() + '01:00:00'::interval) NOT NULL
);

--
-- Name: oauth2_clients; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS oauth2_clients (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    description text DEFAULT ''::TEXT NOT NULL,
    client_id TEXT NOT NULL,
    client_secret TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: password_reset_tokens; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id TEXT NOT NULL,
    token TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    redeemed_at TIMESTAMP WITH TIME ZONE,
    belongs_to_user TEXT NOT NULL
);

--
-- Name: recipe_media; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS recipe_media (
    id TEXT NOT NULL,
    belongs_to_recipe text,
    belongs_to_recipe_step text,
    mime_type TEXT NOT NULL,
    internal_path TEXT NOT NULL,
    external_path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    index integer DEFAULT 0 NOT NULL
);

--
-- Name: recipe_prep_task_steps; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS recipe_prep_task_steps (
    id TEXT NOT NULL,
    satisfies_recipe_step boolean DEFAULT false NOT NULL,
    belongs_to_recipe_step TEXT NOT NULL,
    belongs_to_recipe_prep_task TEXT NOT NULL
);

--
-- Name: recipe_prep_tasks; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS recipe_prep_tasks (
    id TEXT NOT NULL,
    notes text DEFAULT ''::TEXT NOT NULL,
    explicit_storage_instructions text DEFAULT ''::TEXT NOT NULL,
    minimum_time_buffer_before_recipe_in_seconds integer NOT NULL,
    maximum_time_buffer_before_recipe_in_seconds integer,
    storage_type storage_container_type,
    minimum_storage_temperature_in_celsius NUMERIC(14,2),
    maximum_storage_temperature_in_celsius NUMERIC(14,2),
    belongs_to_recipe TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    name text DEFAULT ''::TEXT NOT NULL,
    description text DEFAULT ''::TEXT NOT NULL,
    optional boolean DEFAULT true NOT NULL
);

--
-- Name: recipe_ratings; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS recipe_ratings (
    id TEXT NOT NULL,
    recipe_id TEXT NOT NULL,
    taste NUMERIC(14,2),
    difficulty NUMERIC(14,2),
    cleanup NUMERIC(14,2),
    instructions NUMERIC(14,2),
    overall NUMERIC(14,2),
    notes text DEFAULT ''::TEXT NOT NULL,
    by_user TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: recipe_step_completion_condition_ingredients; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS recipe_step_completion_condition_ingredients (
    id TEXT NOT NULL,
    belongs_to_recipe_step_completion_condition TEXT NOT NULL,
    recipe_step_ingredient TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: recipe_step_completion_conditions; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS recipe_step_completion_conditions (
    id TEXT NOT NULL,
    belongs_to_recipe_step TEXT NOT NULL,
    ingredient_state TEXT NOT NULL,
    notes text DEFAULT ''::TEXT NOT NULL,
    optional boolean DEFAULT false NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: recipe_step_ingredients; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS recipe_step_ingredients (
    id TEXT NOT NULL,
    ingredient_id text,
    minimum_quantity_value NUMERIC(14,2) NOT NULL,
    quantity_notes TEXT NOT NULL,
    ingredient_notes TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_recipe_step TEXT NOT NULL,
    name TEXT NOT NULL,
    recipe_step_product_id text,
    maximum_quantity_value NUMERIC(14,2),
    measurement_unit text,
    optional boolean DEFAULT false NOT NULL,
    option_index integer DEFAULT 0 NOT NULL,
    vessel_index integer,
    to_taste boolean DEFAULT false NOT NULL,
    product_percentage_to_use NUMERIC(14,2),
    recipe_step_product_recipe_id text,
    CONSTRAINT valid_instrument_or_product CHECK (((recipe_step_product_id IS NOT NULL) OR (ingredient_id IS NOT NULL)))
);

--
-- Name: recipe_step_instruments; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS recipe_step_instruments (
    id TEXT NOT NULL,
    instrument_id text,
    notes TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_recipe_step TEXT NOT NULL,
    preference_rank integer NOT NULL,
    recipe_step_product_id text,
    name text DEFAULT ''::TEXT NOT NULL,
    optional boolean DEFAULT false NOT NULL,
    minimum_quantity integer DEFAULT 1 NOT NULL,
    maximum_quantity integer,
    option_index integer DEFAULT 0 NOT NULL
);

--
-- Name: recipe_step_products; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS recipe_step_products (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_recipe_step TEXT NOT NULL,
    quantity_notes TEXT NOT NULL,
    minimum_quantity_value NUMERIC(14,2),
    maximum_quantity_value NUMERIC(14,2),
    measurement_unit text,
    type recipe_step_product_type NOT NULL,
    compostable boolean DEFAULT false NOT NULL,
    maximum_storage_duration_in_seconds integer,
    minimum_storage_temperature_in_celsius NUMERIC(14,2),
    maximum_storage_temperature_in_celsius NUMERIC(14,2),
    storage_instructions text DEFAULT ''::TEXT NOT NULL,
    is_liquid boolean DEFAULT false NOT NULL,
    is_waste boolean DEFAULT false NOT NULL,
    index integer DEFAULT 0 NOT NULL,
    contained_in_vessel_index integer
);

--
-- Name: recipe_step_vessels; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS recipe_step_vessels (
    id TEXT NOT NULL,
    name text DEFAULT ''::TEXT NOT NULL,
    notes text DEFAULT ''::TEXT NOT NULL,
    belongs_to_recipe_step TEXT NOT NULL,
    recipe_step_product_id text,
    vessel_predicate text DEFAULT ''::TEXT NOT NULL,
    minimum_quantity integer DEFAULT 0 NOT NULL,
    maximum_quantity integer,
    unavailable_after_step boolean DEFAULT false NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    valid_vessel_id text
);

--
-- Name: recipe_steps; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS recipe_steps (
    id TEXT NOT NULL,
    index integer NOT NULL,
    preparation_id TEXT NOT NULL,
    minimum_estimated_time_in_seconds bigint,
    maximum_estimated_time_in_seconds bigint,
    minimum_temperature_in_celsius NUMERIC(14,2),
    notes TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_recipe TEXT NOT NULL,
    optional boolean DEFAULT false NOT NULL,
    maximum_temperature_in_celsius NUMERIC(14,2),
    explicit_instructions text DEFAULT ''::TEXT NOT NULL,
    condition_expression text DEFAULT ''::TEXT NOT NULL,
    start_timer_automatically boolean DEFAULT false NOT NULL
);

--
-- Name: recipes; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS recipes (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    source TEXT NOT NULL,
    description TEXT NOT NULL,
    inspired_by_recipe_id text,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    created_by_user TEXT NOT NULL,
    min_estimated_portions NUMERIC(14,2) DEFAULT 1 NOT NULL,
    seal_of_approval boolean DEFAULT false NOT NULL,
    slug text DEFAULT ''::TEXT NOT NULL,
    portion_name text DEFAULT 'portion'::TEXT NOT NULL,
    plural_portion_name text DEFAULT 'portions'::TEXT NOT NULL,
    max_estimated_portions NUMERIC(14,2),
    eligible_for_meals boolean DEFAULT true NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE,
    last_validated_at TIMESTAMP WITH TIME ZONE,
    yields_component_type component_type DEFAULT 'unspecified'::component_type NOT NULL
);

--
-- Name: service_setting_configurations; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS service_setting_configurations (
    id TEXT NOT NULL,
    value text DEFAULT ''::TEXT NOT NULL,
    notes text DEFAULT ''::TEXT NOT NULL,
    service_setting_id TEXT NOT NULL,
    belongs_to_user TEXT NOT NULL,
    belongs_to_household TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: service_settings; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS service_settings (
    id TEXT NOT NULL,
    name text DEFAULT ''::TEXT NOT NULL,
    type setting_type DEFAULT 'user'::setting_type NOT NULL,
    description text DEFAULT ''::TEXT NOT NULL,
    default_value text,
    enumeration text DEFAULT ''::TEXT NOT NULL,
    admins_only boolean DEFAULT true NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: sessions; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS sessions (
    token TEXT NOT NULL,
    data bytea NOT NULL,
    expiry TIMESTAMP WITH TIME ZONE NOT NULL,
    created_on bigint DEFAULT date_part('epoch'::text, now()) NOT NULL
);

--
-- Name: user_ingredient_preferences; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS user_ingredient_preferences (
    id TEXT NOT NULL,
    ingredient TEXT NOT NULL,
    rating smallint DEFAULT 0 NOT NULL,
    notes text DEFAULT ''::TEXT NOT NULL,
    allergy boolean DEFAULT false NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_user TEXT NOT NULL
);

--
-- Name: users; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS users (
    id TEXT NOT NULL,
    username TEXT NOT NULL,
    avatar_src text,
    email_address TEXT NOT NULL,
    hashed_password TEXT NOT NULL,
    password_last_changed_at TIMESTAMP WITH TIME ZONE,
    requires_password_change boolean DEFAULT false NOT NULL,
    two_factor_secret TEXT NOT NULL,
    two_factor_secret_verified_at TIMESTAMP WITH TIME ZONE,
    service_role text DEFAULT 'service_user'::TEXT NOT NULL,
    user_account_status text DEFAULT 'unverified'::TEXT NOT NULL,
    user_account_status_explanation text DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    birthday TIMESTAMP WITH TIME ZONE,
    email_address_verification_token text DEFAULT ''::text,
    email_address_verified_at TIMESTAMP WITH TIME ZONE,
    first_name text DEFAULT ''::TEXT NOT NULL,
    last_name text DEFAULT ''::TEXT NOT NULL,
    last_accepted_terms_of_service TIMESTAMP WITH TIME ZONE,
    last_accepted_privacy_policy TIMESTAMP WITH TIME ZONE
);

--
-- Name: valid_ingredient_group_members; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS valid_ingredient_group_members (
    id TEXT NOT NULL,
    belongs_to_group TEXT NOT NULL,
    valid_ingredient TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: valid_ingredient_groups; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS valid_ingredient_groups (
    id TEXT NOT NULL,
    name text DEFAULT ''::TEXT NOT NULL,
    slug text DEFAULT ''::TEXT NOT NULL,
    description text DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: valid_ingredient_measurement_units; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS valid_ingredient_measurement_units (
    id TEXT NOT NULL,
    notes text DEFAULT ''::TEXT NOT NULL,
    valid_ingredient_id TEXT NOT NULL,
    valid_measurement_unit_id TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    minimum_allowable_quantity NUMERIC(14,2) DEFAULT 0 NOT NULL,
    maximum_allowable_quantity NUMERIC(14,2)
);

--
-- Name: valid_ingredient_preparations; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS valid_ingredient_preparations (
    id TEXT NOT NULL,
    notes TEXT NOT NULL,
    valid_preparation_id TEXT NOT NULL,
    valid_ingredient_id TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: valid_ingredient_state_ingredients; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS valid_ingredient_state_ingredients (
    id TEXT NOT NULL,
    valid_ingredient TEXT NOT NULL,
    valid_ingredient_state TEXT NOT NULL,
    notes text DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: valid_ingredient_states; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS valid_ingredient_states (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    past_tense text DEFAULT ''::TEXT NOT NULL,
    slug text DEFAULT ''::TEXT NOT NULL,
    description text DEFAULT ''::TEXT NOT NULL,
    icon_path text DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    attribute_type ingredient_attribute_type DEFAULT 'other'::ingredient_attribute_type NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: valid_ingredients; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS valid_ingredients (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    warning TEXT NOT NULL,
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
    icon_path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    is_liquid boolean DEFAULT false,
    animal_derived boolean DEFAULT false NOT NULL,
    plural_name text DEFAULT ''::TEXT NOT NULL,
    restrict_to_preparations boolean DEFAULT false NOT NULL,
    minimum_ideal_storage_temperature_in_celsius NUMERIC(14,2),
    maximum_ideal_storage_temperature_in_celsius NUMERIC(14,2),
    storage_instructions text DEFAULT ''::TEXT NOT NULL,
    contains_alcohol boolean DEFAULT false NOT NULL,
    slug text DEFAULT ''::TEXT NOT NULL,
    shopping_suggestions text DEFAULT ''::TEXT NOT NULL,
    is_starch boolean DEFAULT false NOT NULL,
    is_protein boolean DEFAULT false NOT NULL,
    is_grain boolean DEFAULT false NOT NULL,
    is_fruit boolean DEFAULT false NOT NULL,
    is_salt boolean DEFAULT false NOT NULL,
    is_fat boolean DEFAULT false NOT NULL,
    is_acid boolean DEFAULT false NOT NULL,
    is_heat boolean DEFAULT false NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: valid_instruments; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS valid_instruments (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    icon_path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    plural_name text DEFAULT ''::TEXT NOT NULL,
    usable_for_storage boolean DEFAULT false NOT NULL,
    slug text DEFAULT ''::TEXT NOT NULL,
    display_in_summary_lists boolean DEFAULT true NOT NULL,
    include_in_generated_instructions boolean DEFAULT true NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: valid_measurement_conversions; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS valid_measurement_conversions (
    id TEXT NOT NULL,
    from_unit TEXT NOT NULL,
    to_unit TEXT NOT NULL,
    only_for_ingredient text,
    modifier NUMERIC(14,5) NOT NULL,
    notes text DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: valid_measurement_units; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS valid_measurement_units (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    description text DEFAULT ''::TEXT NOT NULL,
    icon_path text DEFAULT ''::TEXT NOT NULL,
    volumetric boolean DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    universal boolean DEFAULT false NOT NULL,
    metric boolean DEFAULT false NOT NULL,
    imperial boolean DEFAULT false NOT NULL,
    plural_name text DEFAULT ''::TEXT NOT NULL,
    slug text DEFAULT ''::TEXT NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: valid_preparation_instruments; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS valid_preparation_instruments (
    id TEXT NOT NULL,
    notes text DEFAULT ''::TEXT NOT NULL,
    valid_preparation_id TEXT NOT NULL,
    valid_instrument_id TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: valid_preparation_vessels; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS valid_preparation_vessels (
    id TEXT NOT NULL,
    valid_preparation_id TEXT NOT NULL,
    valid_vessel_id TEXT NOT NULL,
    notes text DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: valid_preparations; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS valid_preparations (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    icon_path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    yields_nothing boolean DEFAULT false NOT NULL,
    restrict_to_ingredients boolean DEFAULT false NOT NULL,
    past_tense text DEFAULT ''::TEXT NOT NULL,
    slug text DEFAULT ''::TEXT NOT NULL,
    minimum_ingredient_count integer DEFAULT 1 NOT NULL,
    maximum_ingredient_count integer,
    minimum_instrument_count integer DEFAULT 1 NOT NULL,
    maximum_instrument_count integer,
    temperature_required boolean DEFAULT false NOT NULL,
    time_estimate_required boolean DEFAULT false NOT NULL,
    condition_expression_required boolean DEFAULT false NOT NULL,
    consumes_vessel boolean DEFAULT false NOT NULL,
    only_for_vessels boolean DEFAULT false NOT NULL,
    minimum_vessel_count integer DEFAULT 0 NOT NULL,
    maximum_vessel_count integer,
    last_indexed_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: valid_vessels; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS valid_vessels (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    plural_name text DEFAULT ''::TEXT NOT NULL,
    description text DEFAULT ''::TEXT NOT NULL,
    icon_path TEXT NOT NULL,
    usable_for_storage boolean DEFAULT false NOT NULL,
    slug text DEFAULT ''::TEXT NOT NULL,
    display_in_summary_lists boolean DEFAULT true NOT NULL,
    include_in_generated_instructions boolean DEFAULT true NOT NULL,
    capacity NUMERIC(14,2) DEFAULT 0 NOT NULL,
    capacity_unit text,
    width_in_millimeters NUMERIC(14,2),
    length_in_millimeters NUMERIC(14,2),
    height_in_millimeters NUMERIC(14,2),
    shape vessel_shape DEFAULT 'other'::vessel_shape NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: webhook_trigger_events; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS webhook_trigger_events (
    id TEXT NOT NULL,
    trigger_event webhook_event NOT NULL,
    belongs_to_webhook TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    archived_at TIMESTAMP WITH TIME ZONE
);

--
-- Name: webhooks; Type: TABLE;
--

CREATE TABLE IF NOT EXISTS webhooks (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    content_type TEXT NOT NULL,
    url TEXT NOT NULL,
    method TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_household TEXT NOT NULL
);

--
-- Name: darwin_migrations id; Type: DEFAULT;
--

ALTER TABLE ONLY darwin_migrations ALTER COLUMN id SET DEFAULT nextval('darwin_migrations_id_seq'::regclass);

--
-- Name: darwin_migrations darwin_migrations_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY darwin_migrations ADD CONSTRAINT darwin_migrations_pkey PRIMARY KEY (id);

--
-- Name: darwin_migrations darwin_migrations_version_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY darwin_migrations ADD CONSTRAINT darwin_migrations_version_key UNIQUE (version);

--
-- Name: household_instrument_ownerships household_instrument_ownershi_valid_instrument_id_belongs_t_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY household_instrument_ownerships ADD CONSTRAINT household_instrument_ownershi_valid_instrument_id_belongs_t_key UNIQUE (valid_instrument_id, belongs_to_household, archived_at);

--
-- Name: household_instrument_ownerships household_instrument_ownerships_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY household_instrument_ownerships ADD CONSTRAINT household_instrument_ownerships_pkey PRIMARY KEY (id);

--
-- Name: household_invitations household_invitations_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY household_invitations ADD CONSTRAINT household_invitations_pkey PRIMARY KEY (id);

--
-- Name: household_invitations household_invitations_to_user_from_user_destination_househo_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY household_invitations ADD CONSTRAINT household_invitations_to_user_from_user_destination_househo_key UNIQUE (to_user, from_user, destination_household);

--
-- Name: household_user_memberships household_user_memberships_belongs_to_household_belongs_to__key; Type: CONSTRAINT;
--

ALTER TABLE ONLY household_user_memberships ADD CONSTRAINT household_user_memberships_belongs_to_household_belongs_to__key UNIQUE (belongs_to_household, belongs_to_user);

--
-- Name: household_user_memberships household_user_memberships_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY household_user_memberships ADD CONSTRAINT household_user_memberships_pkey PRIMARY KEY (id);

--
-- Name: households households_belongs_to_user_name_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY households ADD CONSTRAINT households_belongs_to_user_name_key UNIQUE (belongs_to_user, name);

--
-- Name: households households_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY households ADD CONSTRAINT households_pkey PRIMARY KEY (id);

--
-- Name: meal_plan_events meal_plan_events_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_events ADD CONSTRAINT meal_plan_events_pkey PRIMARY KEY (id);

--
-- Name: meal_plan_grocery_list_items meal_plan_grocery_list_items_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_grocery_list_items ADD CONSTRAINT meal_plan_grocery_list_items_pkey PRIMARY KEY (id);

--
-- Name: meal_plan_option_votes meal_plan_option_votes_by_user_belongs_to_meal_plan_option_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_option_votes ADD CONSTRAINT meal_plan_option_votes_by_user_belongs_to_meal_plan_option_key UNIQUE (by_user, belongs_to_meal_plan_option);

--
-- Name: meal_plan_option_votes meal_plan_option_votes_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_option_votes ADD CONSTRAINT meal_plan_option_votes_pkey PRIMARY KEY (id);

--
-- Name: meal_plan_options meal_plan_options_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_options ADD CONSTRAINT meal_plan_options_pkey PRIMARY KEY (id);

--
-- Name: meal_plan_tasks meal_plan_tasks_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_tasks ADD CONSTRAINT meal_plan_tasks_pkey PRIMARY KEY (id);

--
-- Name: meal_plans meal_plans_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY meal_plans ADD CONSTRAINT meal_plans_pkey PRIMARY KEY (id);

--
-- Name: meal_components meal_recipes_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY meal_components ADD CONSTRAINT meal_recipes_pkey PRIMARY KEY (id);

--
-- Name: meals meals_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY meals ADD CONSTRAINT meals_pkey PRIMARY KEY (id);

--
-- Name: oauth2_client_tokens oauth2_client_tokens_belongs_to_user_client_id_scope_code_e_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY oauth2_client_tokens ADD CONSTRAINT oauth2_client_tokens_belongs_to_user_client_id_scope_code_e_key UNIQUE (belongs_to_user, client_id, scope, code_expires_at, access_expires_at, refresh_expires_at);

--
-- Name: oauth2_client_tokens oauth2_client_tokens_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY oauth2_client_tokens ADD CONSTRAINT oauth2_client_tokens_pkey PRIMARY KEY (id);

--
-- Name: oauth2_clients oauth2_clients_client_id_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY oauth2_clients ADD CONSTRAINT oauth2_clients_client_id_key UNIQUE (client_id);

--
-- Name: oauth2_clients oauth2_clients_client_secret_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY oauth2_clients ADD CONSTRAINT oauth2_clients_client_secret_key UNIQUE (client_secret);

--
-- Name: oauth2_clients oauth2_clients_name_archived_at_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY oauth2_clients ADD CONSTRAINT oauth2_clients_name_archived_at_key UNIQUE (name, archived_at);

--
-- Name: oauth2_clients oauth2_clients_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY oauth2_clients ADD CONSTRAINT oauth2_clients_pkey PRIMARY KEY (id);

--
-- Name: password_reset_tokens password_reset_tokens_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY password_reset_tokens ADD CONSTRAINT password_reset_tokens_pkey PRIMARY KEY (id);

--
-- Name: recipe_media recipe_media_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY recipe_media ADD CONSTRAINT recipe_media_pkey PRIMARY KEY (id);

--
-- Name: recipe_prep_task_steps recipe_prep_task_steps_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY recipe_prep_task_steps ADD CONSTRAINT recipe_prep_task_steps_pkey PRIMARY KEY (id);

--
-- Name: recipe_prep_tasks recipe_prep_tasks_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY recipe_prep_tasks ADD CONSTRAINT recipe_prep_tasks_pkey PRIMARY KEY (id);

--
-- Name: recipe_ratings recipe_ratings_by_user_recipe_id_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY recipe_ratings ADD CONSTRAINT recipe_ratings_by_user_recipe_id_key UNIQUE (by_user, recipe_id);

--
-- Name: recipe_ratings recipe_ratings_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY recipe_ratings ADD CONSTRAINT recipe_ratings_pkey PRIMARY KEY (id);

--
-- Name: recipe_step_completion_conditions recipe_step_completion_condit_belongs_to_recipe_step_ingred_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_completion_conditions ADD CONSTRAINT recipe_step_completion_condit_belongs_to_recipe_step_ingred_key UNIQUE (belongs_to_recipe_step, ingredient_state);

--
-- Name: recipe_step_completion_condition_ingredients recipe_step_completion_condition_ingredients_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_completion_condition_ingredients ADD CONSTRAINT recipe_step_completion_condition_ingredients_pkey PRIMARY KEY (id);

--
-- Name: recipe_step_completion_conditions recipe_step_completion_conditions_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_completion_conditions ADD CONSTRAINT recipe_step_completion_conditions_pkey PRIMARY KEY (id);

--
-- Name: recipe_step_ingredients recipe_step_ingredients_ingredient_id_belongs_to_recipe_ste_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_ingredients ADD CONSTRAINT recipe_step_ingredients_ingredient_id_belongs_to_recipe_ste_key UNIQUE (ingredient_id, belongs_to_recipe_step);

--
-- Name: recipe_step_ingredients recipe_step_ingredients_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_ingredients ADD CONSTRAINT recipe_step_ingredients_pkey PRIMARY KEY (id);

--
-- Name: recipe_step_instruments recipe_step_instruments_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_instruments ADD CONSTRAINT recipe_step_instruments_pkey PRIMARY KEY (id);

--
-- Name: recipe_step_products recipe_step_products_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_products ADD CONSTRAINT recipe_step_products_pkey PRIMARY KEY (id);

--
-- Name: recipe_step_vessels recipe_step_vessels_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_vessels ADD CONSTRAINT recipe_step_vessels_pkey PRIMARY KEY (id);

--
-- Name: recipe_steps recipe_steps_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY recipe_steps ADD CONSTRAINT recipe_steps_pkey PRIMARY KEY (id);

--
-- Name: recipes recipes_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY recipes ADD CONSTRAINT recipes_pkey PRIMARY KEY (id);

--
-- Name: service_setting_configurations service_setting_configuration_belongs_to_user_belongs_to_ho_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY service_setting_configurations ADD CONSTRAINT service_setting_configuration_belongs_to_user_belongs_to_ho_key UNIQUE (belongs_to_user, belongs_to_household, service_setting_id);

--
-- Name: service_setting_configurations service_setting_configurations_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY service_setting_configurations ADD CONSTRAINT service_setting_configurations_pkey PRIMARY KEY (id);

--
-- Name: service_settings service_settings_name_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY service_settings ADD CONSTRAINT service_settings_name_key UNIQUE (name);

--
-- Name: service_settings service_settings_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY service_settings ADD CONSTRAINT service_settings_pkey PRIMARY KEY (id);

--
-- Name: sessions sessions_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY sessions ADD CONSTRAINT sessions_pkey PRIMARY KEY (token);

--
-- Name: recipe_prep_task_steps unique_recipe_step_and_prep_task; Type: CONSTRAINT;
--

ALTER TABLE ONLY recipe_prep_task_steps ADD CONSTRAINT unique_recipe_step_and_prep_task UNIQUE (belongs_to_recipe_step, belongs_to_recipe_prep_task);

--
-- Name: user_ingredient_preferences user_ingredient_preferences_belongs_to_user_ingredient_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY user_ingredient_preferences ADD CONSTRAINT user_ingredient_preferences_belongs_to_user_ingredient_key UNIQUE (belongs_to_user, ingredient);

--
-- Name: user_ingredient_preferences user_ingredient_preferences_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY user_ingredient_preferences ADD CONSTRAINT user_ingredient_preferences_pkey PRIMARY KEY (id);

--
-- Name: users users_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY users ADD CONSTRAINT users_pkey PRIMARY KEY (id);

--
-- Name: users users_username_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY users ADD CONSTRAINT users_username_key UNIQUE (username);

--
-- Name: valid_ingredient_group_members valid_ingredient_group_members_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_group_members ADD CONSTRAINT valid_ingredient_group_members_pkey PRIMARY KEY (id);

--
-- Name: valid_ingredient_groups valid_ingredient_groups_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_groups ADD CONSTRAINT valid_ingredient_groups_pkey PRIMARY KEY (id);

--
-- Name: valid_ingredient_measurement_units valid_ingredient_measurement__valid_ingredient_id_valid_mea_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_measurement_units ADD CONSTRAINT valid_ingredient_measurement__valid_ingredient_id_valid_mea_key UNIQUE (valid_ingredient_id, valid_measurement_unit_id, archived_at);

--
-- Name: valid_ingredient_measurement_units valid_ingredient_measurement_units_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_measurement_units ADD CONSTRAINT valid_ingredient_measurement_units_pkey PRIMARY KEY (id);

--
-- Name: valid_ingredient_preparations valid_ingredient_preparations_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_preparations ADD CONSTRAINT valid_ingredient_preparations_pkey PRIMARY KEY (id);

--
-- Name: valid_ingredient_preparations valid_ingredient_preparations_valid_preparation_id_valid_in_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_preparations ADD CONSTRAINT valid_ingredient_preparations_valid_preparation_id_valid_in_key UNIQUE (valid_preparation_id, valid_ingredient_id, archived_at);

--
-- Name: valid_ingredient_state_ingredients valid_ingredient_state_ingred_valid_ingredient_valid_ingred_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_state_ingredients ADD CONSTRAINT valid_ingredient_state_ingred_valid_ingredient_valid_ingred_key UNIQUE (valid_ingredient, valid_ingredient_state, archived_at);

--
-- Name: valid_ingredient_state_ingredients valid_ingredient_state_ingredients_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_state_ingredients ADD CONSTRAINT valid_ingredient_state_ingredients_pkey PRIMARY KEY (id);

--
-- Name: valid_ingredient_states valid_ingredient_states_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_states ADD CONSTRAINT valid_ingredient_states_pkey PRIMARY KEY (id);

--
-- Name: valid_ingredients valid_ingredients_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredients ADD CONSTRAINT valid_ingredients_pkey PRIMARY KEY (id);

--
-- Name: valid_preparation_instruments valid_instrument_preparation_pair; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_preparation_instruments ADD CONSTRAINT valid_instrument_preparation_pair UNIQUE (valid_instrument_id, valid_preparation_id, archived_at);

--
-- Name: valid_instruments valid_instruments_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_instruments ADD CONSTRAINT valid_instruments_pkey PRIMARY KEY (id);

--
-- Name: valid_measurement_conversions valid_measurement_conversions_from_unit_to_unit_only_for_in_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_measurement_conversions ADD CONSTRAINT valid_measurement_conversions_from_unit_to_unit_only_for_in_key UNIQUE (from_unit, to_unit, only_for_ingredient);

--
-- Name: valid_measurement_conversions valid_measurement_conversions_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_measurement_conversions ADD CONSTRAINT valid_measurement_conversions_pkey PRIMARY KEY (id);

--
-- Name: valid_measurement_units valid_measurement_units_name_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_measurement_units ADD CONSTRAINT valid_measurement_units_name_key UNIQUE (name);

--
-- Name: valid_measurement_units valid_measurement_units_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_measurement_units ADD CONSTRAINT valid_measurement_units_pkey PRIMARY KEY (id);

--
-- Name: valid_preparation_instruments valid_preparation_instruments_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_preparation_instruments ADD CONSTRAINT valid_preparation_instruments_pkey PRIMARY KEY (id);

--
-- Name: valid_preparation_instruments valid_preparation_instruments_valid_preparation_id_valid_in_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_preparation_instruments ADD CONSTRAINT valid_preparation_instruments_valid_preparation_id_valid_in_key UNIQUE (valid_preparation_id, valid_instrument_id, archived_at);

--
-- Name: valid_preparation_vessels valid_preparation_vessels_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_preparation_vessels ADD CONSTRAINT valid_preparation_vessels_pkey PRIMARY KEY (id);

--
-- Name: valid_preparation_vessels valid_preparation_vessels_valid_preparation_id_valid_vessel_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_preparation_vessels ADD CONSTRAINT valid_preparation_vessels_valid_preparation_id_valid_vessel_key UNIQUE (valid_preparation_id, valid_vessel_id, archived_at);

--
-- Name: valid_preparations valid_preparations_name_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_preparations ADD CONSTRAINT valid_preparations_name_key UNIQUE (name);

--
-- Name: valid_preparations valid_preparations_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_preparations ADD CONSTRAINT valid_preparations_pkey PRIMARY KEY (id);

--
-- Name: valid_vessels valid_vessels_name_archived_at_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_vessels ADD CONSTRAINT valid_vessels_name_archived_at_key UNIQUE (name, archived_at);

--
-- Name: valid_vessels valid_vessels_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_vessels ADD CONSTRAINT valid_vessels_pkey PRIMARY KEY (id);

--
-- Name: valid_vessels valid_vessels_slug_archived_at_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY valid_vessels ADD CONSTRAINT valid_vessels_slug_archived_at_key UNIQUE (slug, archived_at);

--
-- Name: webhook_trigger_events webhook_trigger_events_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY webhook_trigger_events ADD CONSTRAINT webhook_trigger_events_pkey PRIMARY KEY (id);

--
-- Name: webhook_trigger_events webhook_trigger_events_trigger_event_belongs_to_webhook_key; Type: CONSTRAINT;
--

ALTER TABLE ONLY webhook_trigger_events ADD CONSTRAINT webhook_trigger_events_trigger_event_belongs_to_webhook_key UNIQUE (trigger_event, belongs_to_webhook);

--
-- Name: webhooks webhooks_pkey; Type: CONSTRAINT;
--

ALTER TABLE ONLY webhooks ADD CONSTRAINT webhooks_pkey PRIMARY KEY (id);

--
-- Name: household_invitations_destination_household; Type: INDEX;
--

CREATE INDEX household_invitations_destination_household ON household_invitations USING btree (destination_household);

--
-- Name: household_invitations_from_user; Type: INDEX;
--

CREATE INDEX household_invitations_from_user ON household_invitations USING btree (from_user);

--
-- Name: household_invitations_to_user; Type: INDEX;
--

CREATE INDEX household_invitations_to_user ON household_invitations USING btree (to_user);

--
-- Name: household_user_memberships_belongs_to_household; Type: INDEX;
--

CREATE INDEX household_user_memberships_belongs_to_household ON household_user_memberships USING btree (belongs_to_household);

--
-- Name: household_user_memberships_belongs_to_user; Type: INDEX;
--

CREATE INDEX household_user_memberships_belongs_to_user ON household_user_memberships USING btree (belongs_to_user);

--
-- Name: households_belongs_to_user; Type: INDEX;
--

CREATE INDEX households_belongs_to_user ON households USING btree (belongs_to_user);

--
-- Name: meal_plan_events_belongs_to_meal_pla_index; Type: INDEX;
--

CREATE INDEX meal_plan_events_belongs_to_meal_pla_index ON meal_plan_events USING btree (belongs_to_meal_plan);

--
-- Name: meal_plan_grocery_list_items_belongs_to_meal_pla_index; Type: INDEX;
--

CREATE INDEX meal_plan_grocery_list_items_belongs_to_meal_pla_index ON meal_plan_grocery_list_items USING btree (belongs_to_meal_plan);

--
-- Name: meal_plan_grocery_list_items_purchased_measurement_unit_index; Type: INDEX;
--

CREATE INDEX meal_plan_grocery_list_items_purchased_measurement_unit_index ON meal_plan_grocery_list_items USING btree (purchased_measurement_unit);

--
-- Name: meal_plan_grocery_list_items_valid_ingredient_index; Type: INDEX;
--

CREATE INDEX meal_plan_grocery_list_items_valid_ingredient_index ON meal_plan_grocery_list_items USING btree (valid_ingredient);

--
-- Name: meal_plan_grocery_list_items_valid_measurement_unit_index; Type: INDEX;
--

CREATE INDEX meal_plan_grocery_list_items_valid_measurement_unit_index ON meal_plan_grocery_list_items USING btree (valid_measurement_unit);

--
-- Name: meal_plan_options_assigned_cook_index; Type: INDEX;
--

CREATE INDEX meal_plan_options_assigned_cook_index ON meal_plan_options USING btree (assigned_cook);

--
-- Name: meal_plan_options_assigned_dishwasher_index; Type: INDEX;
--

CREATE INDEX meal_plan_options_assigned_dishwasher_index ON meal_plan_options USING btree (assigned_dishwasher);

--
-- Name: meal_plan_options_belongs_to_meal_plan_even_index; Type: INDEX;
--

CREATE INDEX meal_plan_options_belongs_to_meal_plan_even_index ON meal_plan_options USING btree (belongs_to_meal_plan_event);

--
-- Name: meal_plan_options_belongs_to_meal_plan_option; Type: INDEX;
--

CREATE INDEX meal_plan_options_belongs_to_meal_plan_option ON meal_plan_option_votes USING btree (belongs_to_meal_plan_option);

--
-- Name: meal_plan_tasks_assigned_to_user_index; Type: INDEX;
--

CREATE INDEX meal_plan_tasks_assigned_to_user_index ON meal_plan_tasks USING btree (assigned_to_user);

--
-- Name: meal_plan_tasks_belongs_to_meal_plan_option_index; Type: INDEX;
--

CREATE INDEX meal_plan_tasks_belongs_to_meal_plan_option_index ON meal_plan_tasks USING btree (belongs_to_meal_plan_option);

--
-- Name: meal_plan_tasks_belongs_to_recipe_prep_task_index; Type: INDEX;
--

CREATE INDEX meal_plan_tasks_belongs_to_recipe_prep_task_index ON meal_plan_tasks USING btree (belongs_to_recipe_prep_task);

--
-- Name: meal_plans_belongs_to_household; Type: INDEX;
--

CREATE INDEX meal_plans_belongs_to_household ON meal_plans USING btree (belongs_to_household);

--
-- Name: meal_recipes_meal_id; Type: INDEX;
--

CREATE INDEX meal_recipes_meal_id ON meal_components USING btree (meal_id);

--
-- Name: meal_recipes_recipe_id; Type: INDEX;
--

CREATE INDEX meal_recipes_recipe_id ON meal_components USING btree (recipe_id);

--
-- Name: meals_created_by_user; Type: INDEX;
--

CREATE INDEX meals_created_by_user ON meals USING btree (created_by_user);

--
-- Name: password_reset_token_belongs_to_user; Type: INDEX;
--

CREATE INDEX password_reset_token_belongs_to_user ON password_reset_tokens USING btree (belongs_to_user);

--
-- Name: recipe_media_belongs_to_recipe_index; Type: INDEX;
--

CREATE INDEX recipe_media_belongs_to_recipe_index ON recipe_media USING btree (belongs_to_recipe);

--
-- Name: recipe_media_belongs_to_recipe_step_index; Type: INDEX;
--

CREATE INDEX recipe_media_belongs_to_recipe_step_index ON recipe_media USING btree (belongs_to_recipe_step);

--
-- Name: recipe_prep_task_steps_belongs_to_recipe_prep_task_index; Type: INDEX;
--

CREATE INDEX recipe_prep_task_steps_belongs_to_recipe_prep_task_index ON recipe_prep_task_steps USING btree (belongs_to_recipe_prep_task);

--
-- Name: recipe_prep_task_steps_belongs_to_recipe_step_index; Type: INDEX;
--

CREATE INDEX recipe_prep_task_steps_belongs_to_recipe_step_index ON recipe_prep_task_steps USING btree (belongs_to_recipe_step);

--
-- Name: recipe_prep_tasks_belongs_to_recipe_index; Type: INDEX;
--

CREATE INDEX recipe_prep_tasks_belongs_to_recipe_index ON recipe_prep_tasks USING btree (belongs_to_recipe);

--
-- Name: recipe_step_ingredients_measurement_unit_index; Type: INDEX;
--

CREATE INDEX recipe_step_ingredients_measurement_unit_index ON recipe_step_ingredients USING btree (measurement_unit);

--
-- Name: recipe_step_ingredients_product_of_recipe_step; Type: INDEX;
--

CREATE INDEX recipe_step_ingredients_product_of_recipe_step ON recipe_step_ingredients USING btree (recipe_step_product_id);

--
-- Name: recipe_step_instruments_instrument_id_index; Type: INDEX;
--

CREATE INDEX recipe_step_instruments_instrument_id_index ON recipe_step_instruments USING btree (instrument_id);

--
-- Name: recipe_step_instruments_recipe_step_product_id_index; Type: INDEX;
--

CREATE INDEX recipe_step_instruments_recipe_step_product_id_index ON recipe_step_instruments USING btree (recipe_step_product_id);

--
-- Name: recipe_step_products_belongs_to_recipe_step; Type: INDEX;
--

CREATE INDEX recipe_step_products_belongs_to_recipe_step ON recipe_step_products USING btree (belongs_to_recipe_step);

--
-- Name: recipe_step_products_measurement_unit_index; Type: INDEX;
--

CREATE INDEX recipe_step_products_measurement_unit_index ON recipe_step_products USING btree (measurement_unit);

--
-- Name: recipe_steps_belongs_to_recipe; Type: INDEX;
--

CREATE INDEX recipe_steps_belongs_to_recipe ON recipe_steps USING btree (belongs_to_recipe);

--
-- Name: recipes_created_by_user; Type: INDEX;
--

CREATE INDEX recipes_created_by_user ON recipes USING btree (created_by_user);

--
-- Name: sessions_expiry_idx; Type: INDEX;
--

CREATE INDEX sessions_expiry_idx ON sessions USING btree (expiry);

--
-- Name: valid_ingredient_measurement_units_valid_ingredient_id_index; Type: INDEX;
--

CREATE INDEX valid_ingredient_measurement_units_valid_ingredient_id_index ON valid_ingredient_measurement_units USING btree (valid_ingredient_id);

--
-- Name: valid_ingredient_measurement_units_valid_measurement_unit_id_in; Type: INDEX;
--

CREATE INDEX valid_ingredient_measurement_units_valid_measurement_unit_id_in ON valid_ingredient_measurement_units USING btree (valid_measurement_unit_id);

--
-- Name: valid_ingredient_state_ingredients_referncing_valid_ingredient_; Type: INDEX;
--

CREATE INDEX valid_ingredient_state_ingredients_referncing_valid_ingredient_ ON valid_ingredient_state_ingredients USING btree (valid_ingredient);

--
-- Name: valid_measurement_conversions_from_unit_index; Type: INDEX;
--

CREATE INDEX valid_measurement_conversions_from_unit_index ON valid_measurement_conversions USING btree (from_unit);

--
-- Name: valid_measurement_conversions_only_for_ingredient_index; Type: INDEX;
--

CREATE INDEX valid_measurement_conversions_only_for_ingredient_index ON valid_measurement_conversions USING btree (only_for_ingredient);

--
-- Name: valid_measurement_conversions_to_unit_index; Type: INDEX;
--

CREATE INDEX valid_measurement_conversions_to_unit_index ON valid_measurement_conversions USING btree (to_unit);

--
-- Name: valid_preparation_instruments_valid_instrument_index; Type: INDEX;
--

CREATE INDEX valid_preparation_instruments_valid_instrument_index ON valid_preparation_instruments USING btree (valid_instrument_id);

--
-- Name: valid_preparation_instruments_valid_preparation_index; Type: INDEX;
--

CREATE INDEX valid_preparation_instruments_valid_preparation_index ON valid_preparation_instruments USING btree (valid_preparation_id);

--
-- Name: valid_preparation_vessels_referencing_valid_preparations_idx; Type: INDEX;
--

CREATE INDEX valid_preparation_vessels_referencing_valid_preparations_idx ON valid_preparation_vessels USING btree (valid_preparation_id);

--
-- Name: valid_preparation_vessels_referencing_valid_vessels_idx; Type: INDEX;
--

CREATE INDEX valid_preparation_vessels_referencing_valid_vessels_idx ON valid_preparation_vessels USING btree (valid_vessel_id);

--
-- Name: webhook_trigger_events_belongs_to_webhook_index; Type: INDEX;
--

CREATE INDEX webhook_trigger_events_belongs_to_webhook_index ON webhook_trigger_events USING btree (belongs_to_webhook);

--
-- Name: recipe_step_instruments fk_valid_recipe_step_instrument_ids; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_instruments ADD CONSTRAINT fk_valid_recipe_step_instrument_ids FOREIGN KEY (instrument_id) REFERENCES valid_instruments(id);

--
-- Name: household_instrument_ownerships household_instrument_ownerships_belongs_to_household_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY household_instrument_ownerships ADD CONSTRAINT household_instrument_ownerships_belongs_to_household_fkey FOREIGN KEY (belongs_to_household) REFERENCES households(id) ON DELETE CASCADE;

--
-- Name: household_instrument_ownerships household_instrument_ownerships_valid_instrument_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY household_instrument_ownerships ADD CONSTRAINT household_instrument_ownerships_valid_instrument_id_fkey FOREIGN KEY (valid_instrument_id) REFERENCES valid_instruments(id) ON DELETE CASCADE;

--
-- Name: household_invitations household_invitations_destination_household_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY household_invitations ADD CONSTRAINT household_invitations_destination_household_fkey FOREIGN KEY (destination_household) REFERENCES households(id) ON DELETE CASCADE;

--
-- Name: household_invitations household_invitations_from_user_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY household_invitations ADD CONSTRAINT household_invitations_from_user_fkey FOREIGN KEY (from_user) REFERENCES users(id) ON DELETE CASCADE;

--
-- Name: household_invitations household_invitations_to_user_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY household_invitations ADD CONSTRAINT household_invitations_to_user_fkey FOREIGN KEY (to_user) REFERENCES users(id) ON DELETE CASCADE;

--
-- Name: household_user_memberships household_user_memberships_belongs_to_household_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY household_user_memberships ADD CONSTRAINT household_user_memberships_belongs_to_household_fkey FOREIGN KEY (belongs_to_household) REFERENCES households(id) ON DELETE CASCADE;

--
-- Name: household_user_memberships household_user_memberships_belongs_to_user_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY household_user_memberships ADD CONSTRAINT household_user_memberships_belongs_to_user_fkey FOREIGN KEY (belongs_to_user) REFERENCES users(id) ON DELETE CASCADE;

--
-- Name: households households_belongs_to_user_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY households ADD CONSTRAINT households_belongs_to_user_fkey FOREIGN KEY (belongs_to_user) REFERENCES users(id) ON DELETE CASCADE;

--
-- Name: meal_plan_events meal_plan_events_belongs_to_meal_plan_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_events ADD CONSTRAINT meal_plan_events_belongs_to_meal_plan_fkey FOREIGN KEY (belongs_to_meal_plan) REFERENCES meal_plans(id) ON DELETE CASCADE;

--
-- Name: meal_plan_grocery_list_items meal_plan_grocery_list_items_belongs_to_meal_plan_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_grocery_list_items ADD CONSTRAINT meal_plan_grocery_list_items_belongs_to_meal_plan_fkey FOREIGN KEY (belongs_to_meal_plan) REFERENCES meal_plans(id) ON DELETE CASCADE;

--
-- Name: meal_plan_grocery_list_items meal_plan_grocery_list_items_purchased_measurement_unit_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_grocery_list_items ADD CONSTRAINT meal_plan_grocery_list_items_purchased_measurement_unit_fkey FOREIGN KEY (purchased_measurement_unit) REFERENCES valid_measurement_units(id) ON DELETE CASCADE;

--
-- Name: meal_plan_grocery_list_items meal_plan_grocery_list_items_valid_ingredient_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_grocery_list_items ADD CONSTRAINT meal_plan_grocery_list_items_valid_ingredient_fkey FOREIGN KEY (valid_ingredient) REFERENCES valid_ingredients(id) ON DELETE CASCADE;

--
-- Name: meal_plan_grocery_list_items meal_plan_grocery_list_items_valid_measurement_unit_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_grocery_list_items ADD CONSTRAINT meal_plan_grocery_list_items_valid_measurement_unit_fkey FOREIGN KEY (valid_measurement_unit) REFERENCES valid_measurement_units(id) ON DELETE CASCADE;

--
-- Name: meal_plan_option_votes meal_plan_option_votes_belongs_to_meal_plan_option_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_option_votes ADD CONSTRAINT meal_plan_option_votes_belongs_to_meal_plan_option_fkey FOREIGN KEY (belongs_to_meal_plan_option) REFERENCES meal_plan_options(id) ON DELETE CASCADE;

--
-- Name: meal_plan_option_votes meal_plan_option_votes_by_user_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_option_votes ADD CONSTRAINT meal_plan_option_votes_by_user_fkey FOREIGN KEY (by_user) REFERENCES users(id) ON DELETE CASCADE;

--
-- Name: meal_plan_options meal_plan_options_assigned_cook_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_options ADD CONSTRAINT meal_plan_options_assigned_cook_fkey FOREIGN KEY (assigned_cook) REFERENCES users(id) ON DELETE CASCADE;

--
-- Name: meal_plan_options meal_plan_options_assigned_dishwasher_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_options ADD CONSTRAINT meal_plan_options_assigned_dishwasher_fkey FOREIGN KEY (assigned_dishwasher) REFERENCES users(id) ON DELETE CASCADE;

--
-- Name: meal_plan_options meal_plan_options_belongs_to_meal_plan_event_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_options ADD CONSTRAINT meal_plan_options_belongs_to_meal_plan_event_fkey FOREIGN KEY (belongs_to_meal_plan_event) REFERENCES meal_plan_events(id) ON DELETE CASCADE;

--
-- Name: meal_plan_options meal_plan_options_meal_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_options ADD CONSTRAINT meal_plan_options_meal_id_fkey FOREIGN KEY (meal_id) REFERENCES meals(id) ON DELETE CASCADE;

--
-- Name: meal_plan_tasks meal_plan_tasks_assigned_to_user_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_tasks ADD CONSTRAINT meal_plan_tasks_assigned_to_user_fkey FOREIGN KEY (assigned_to_user) REFERENCES users(id) ON DELETE CASCADE;

--
-- Name: meal_plan_tasks meal_plan_tasks_belongs_to_meal_plan_option_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_tasks ADD CONSTRAINT meal_plan_tasks_belongs_to_meal_plan_option_fkey FOREIGN KEY (belongs_to_meal_plan_option) REFERENCES meal_plan_options(id) ON DELETE CASCADE;

--
-- Name: meal_plan_tasks meal_plan_tasks_belongs_to_recipe_prep_task_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_plan_tasks ADD CONSTRAINT meal_plan_tasks_belongs_to_recipe_prep_task_fkey FOREIGN KEY (belongs_to_recipe_prep_task) REFERENCES recipe_prep_tasks(id) ON DELETE CASCADE;

--
-- Name: meal_plans meal_plans_belongs_to_household_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_plans ADD CONSTRAINT meal_plans_belongs_to_household_fkey FOREIGN KEY (belongs_to_household) REFERENCES households(id) ON DELETE CASCADE;

--
-- Name: meal_plans meal_plans_created_by_user_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_plans ADD CONSTRAINT meal_plans_created_by_user_fkey FOREIGN KEY (created_by_user) REFERENCES users(id);

--
-- Name: meal_components meal_recipes_meal_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_components ADD CONSTRAINT meal_recipes_meal_id_fkey FOREIGN KEY (meal_id) REFERENCES meals(id) ON DELETE CASCADE;

--
-- Name: meal_components meal_recipes_recipe_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meal_components ADD CONSTRAINT meal_recipes_recipe_id_fkey FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE;

--
-- Name: meals meals_created_by_user_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY meals ADD CONSTRAINT meals_created_by_user_fkey FOREIGN KEY (created_by_user) REFERENCES users(id) ON DELETE CASCADE;

--
-- Name: oauth2_client_tokens oauth2_client_tokens_belongs_to_user_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY oauth2_client_tokens ADD CONSTRAINT oauth2_client_tokens_belongs_to_user_fkey FOREIGN KEY (belongs_to_user) REFERENCES users(id) ON DELETE CASCADE;

--
-- Name: oauth2_client_tokens oauth2_client_tokens_client_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY oauth2_client_tokens ADD CONSTRAINT oauth2_client_tokens_client_id_fkey FOREIGN KEY (client_id) REFERENCES oauth2_clients(client_id) ON DELETE CASCADE;

--
-- Name: password_reset_tokens password_reset_tokens_belongs_to_user_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY password_reset_tokens ADD CONSTRAINT password_reset_tokens_belongs_to_user_fkey FOREIGN KEY (belongs_to_user) REFERENCES users(id) ON DELETE CASCADE;

--
-- Name: recipe_media recipe_media_belongs_to_recipe_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_media ADD CONSTRAINT recipe_media_belongs_to_recipe_fkey FOREIGN KEY (belongs_to_recipe) REFERENCES recipes(id) ON DELETE CASCADE;

--
-- Name: recipe_media recipe_media_belongs_to_recipe_step_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_media ADD CONSTRAINT recipe_media_belongs_to_recipe_step_fkey FOREIGN KEY (belongs_to_recipe_step) REFERENCES recipe_steps(id) ON DELETE CASCADE;

--
-- Name: recipe_prep_task_steps recipe_prep_task_steps_belongs_to_recipe_prep_task_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_prep_task_steps ADD CONSTRAINT recipe_prep_task_steps_belongs_to_recipe_prep_task_fkey FOREIGN KEY (belongs_to_recipe_prep_task) REFERENCES recipe_prep_tasks(id) ON DELETE CASCADE;

--
-- Name: recipe_prep_task_steps recipe_prep_task_steps_belongs_to_recipe_step_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_prep_task_steps ADD CONSTRAINT recipe_prep_task_steps_belongs_to_recipe_step_fkey FOREIGN KEY (belongs_to_recipe_step) REFERENCES recipe_steps(id) ON DELETE CASCADE;

--
-- Name: recipe_prep_tasks recipe_prep_tasks_belongs_to_recipe_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_prep_tasks ADD CONSTRAINT recipe_prep_tasks_belongs_to_recipe_fkey FOREIGN KEY (belongs_to_recipe) REFERENCES recipes(id) ON DELETE CASCADE;

--
-- Name: recipe_ratings recipe_ratings_by_user_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_ratings ADD CONSTRAINT recipe_ratings_by_user_fkey FOREIGN KEY (by_user) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: recipe_ratings recipe_ratings_recipe_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_ratings ADD CONSTRAINT recipe_ratings_recipe_id_fkey FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: recipe_step_completion_condition_ingredients recipe_step_completion_condit_belongs_to_recipe_step_compl_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_completion_condition_ingredients ADD CONSTRAINT recipe_step_completion_condit_belongs_to_recipe_step_compl_fkey FOREIGN KEY (belongs_to_recipe_step_completion_condition) REFERENCES recipe_step_completion_conditions(id) ON DELETE CASCADE;

--
-- Name: recipe_step_completion_condition_ingredients recipe_step_completion_condition_in_recipe_step_ingredient_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_completion_condition_ingredients ADD CONSTRAINT recipe_step_completion_condition_in_recipe_step_ingredient_fkey FOREIGN KEY (recipe_step_ingredient) REFERENCES recipe_step_ingredients(id) ON DELETE CASCADE;

--
-- Name: recipe_step_completion_conditions recipe_step_completion_conditions_belongs_to_recipe_step_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_completion_conditions ADD CONSTRAINT recipe_step_completion_conditions_belongs_to_recipe_step_fkey FOREIGN KEY (belongs_to_recipe_step) REFERENCES recipe_steps(id) ON DELETE CASCADE;

--
-- Name: recipe_step_completion_conditions recipe_step_completion_conditions_ingredient_state_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_completion_conditions ADD CONSTRAINT recipe_step_completion_conditions_ingredient_state_fkey FOREIGN KEY (ingredient_state) REFERENCES valid_ingredient_states(id) ON DELETE CASCADE;

--
-- Name: recipe_step_ingredients recipe_step_ingredients_belongs_to_recipe_step_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_ingredients ADD CONSTRAINT recipe_step_ingredients_belongs_to_recipe_step_fkey FOREIGN KEY (belongs_to_recipe_step) REFERENCES recipe_steps(id) ON DELETE CASCADE;

--
-- Name: recipe_step_ingredients recipe_step_ingredients_ingredient_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_ingredients ADD CONSTRAINT recipe_step_ingredients_ingredient_id_fkey FOREIGN KEY (ingredient_id) REFERENCES valid_ingredients(id) ON DELETE CASCADE;

--
-- Name: recipe_step_ingredients recipe_step_ingredients_measurement_unit_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_ingredients ADD CONSTRAINT recipe_step_ingredients_measurement_unit_fkey FOREIGN KEY (measurement_unit) REFERENCES valid_measurement_units(id);

--
-- Name: recipe_step_ingredients recipe_step_ingredients_recipe_step_product_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_ingredients ADD CONSTRAINT recipe_step_ingredients_recipe_step_product_id_fkey FOREIGN KEY (recipe_step_product_id) REFERENCES recipe_step_products(id) ON DELETE RESTRICT;

--
-- Name: recipe_step_ingredients recipe_step_ingredients_recipe_step_product_recipe_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_ingredients ADD CONSTRAINT recipe_step_ingredients_recipe_step_product_recipe_id_fkey FOREIGN KEY (recipe_step_product_recipe_id) REFERENCES recipes(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: recipe_step_instruments recipe_step_instruments_belongs_to_recipe_step_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_instruments ADD CONSTRAINT recipe_step_instruments_belongs_to_recipe_step_fkey FOREIGN KEY (belongs_to_recipe_step) REFERENCES recipe_steps(id) ON DELETE CASCADE;

--
-- Name: recipe_step_instruments recipe_step_instruments_recipe_step_product_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_instruments ADD CONSTRAINT recipe_step_instruments_recipe_step_product_id_fkey FOREIGN KEY (recipe_step_product_id) REFERENCES recipe_step_products(id);

--
-- Name: recipe_step_products recipe_step_products_belongs_to_recipe_step_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_products ADD CONSTRAINT recipe_step_products_belongs_to_recipe_step_fkey FOREIGN KEY (belongs_to_recipe_step) REFERENCES recipe_steps(id) ON DELETE CASCADE;

--
-- Name: recipe_step_products recipe_step_products_measurement_unit_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_products ADD CONSTRAINT recipe_step_products_measurement_unit_fkey FOREIGN KEY (measurement_unit) REFERENCES valid_measurement_units(id);

--
-- Name: recipe_step_vessels recipe_step_vessels_belongs_to_recipe_step_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_vessels ADD CONSTRAINT recipe_step_vessels_belongs_to_recipe_step_fkey FOREIGN KEY (belongs_to_recipe_step) REFERENCES recipe_steps(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: recipe_step_vessels recipe_step_vessels_recipe_step_product_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_vessels ADD CONSTRAINT recipe_step_vessels_recipe_step_product_id_fkey FOREIGN KEY (recipe_step_product_id) REFERENCES recipe_step_products(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: recipe_step_vessels recipe_step_vessels_valid_vessel_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_step_vessels ADD CONSTRAINT recipe_step_vessels_valid_vessel_id_fkey FOREIGN KEY (valid_vessel_id) REFERENCES valid_vessels(id) ON DELETE CASCADE;

--
-- Name: recipe_steps recipe_steps_belongs_to_recipe_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_steps ADD CONSTRAINT recipe_steps_belongs_to_recipe_fkey FOREIGN KEY (belongs_to_recipe) REFERENCES recipes(id) ON DELETE CASCADE;

--
-- Name: recipe_steps recipe_steps_preparation_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipe_steps ADD CONSTRAINT recipe_steps_preparation_id_fkey FOREIGN KEY (preparation_id) REFERENCES valid_preparations(id) ON DELETE CASCADE;

--
-- Name: recipes recipes_created_by_user_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY recipes ADD CONSTRAINT recipes_created_by_user_fkey FOREIGN KEY (created_by_user) REFERENCES users(id) ON DELETE CASCADE;

--
-- Name: service_setting_configurations service_setting_configurations_belongs_to_household_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY service_setting_configurations ADD CONSTRAINT service_setting_configurations_belongs_to_household_fkey FOREIGN KEY (belongs_to_household) REFERENCES households(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: service_setting_configurations service_setting_configurations_belongs_to_user_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY service_setting_configurations ADD CONSTRAINT service_setting_configurations_belongs_to_user_fkey FOREIGN KEY (belongs_to_user) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: service_setting_configurations service_setting_configurations_service_setting_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY service_setting_configurations ADD CONSTRAINT service_setting_configurations_service_setting_id_fkey FOREIGN KEY (service_setting_id) REFERENCES service_settings(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: user_ingredient_preferences user_ingredient_preferences_belongs_to_user_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY user_ingredient_preferences ADD CONSTRAINT user_ingredient_preferences_belongs_to_user_fkey FOREIGN KEY (belongs_to_user) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: user_ingredient_preferences user_ingredient_preferences_ingredient_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY user_ingredient_preferences ADD CONSTRAINT user_ingredient_preferences_ingredient_fkey FOREIGN KEY (ingredient) REFERENCES valid_ingredients(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: valid_ingredient_group_members valid_ingredient_group_members_belongs_to_group_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_group_members ADD CONSTRAINT valid_ingredient_group_members_belongs_to_group_fkey FOREIGN KEY (belongs_to_group) REFERENCES valid_ingredient_groups(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: valid_ingredient_group_members valid_ingredient_group_members_valid_ingredient_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_group_members ADD CONSTRAINT valid_ingredient_group_members_valid_ingredient_fkey FOREIGN KEY (valid_ingredient) REFERENCES valid_ingredients(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: valid_ingredient_measurement_units valid_ingredient_measurement_uni_valid_measurement_unit_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_measurement_units ADD CONSTRAINT valid_ingredient_measurement_uni_valid_measurement_unit_id_fkey FOREIGN KEY (valid_measurement_unit_id) REFERENCES valid_measurement_units(id) ON DELETE CASCADE;

--
-- Name: valid_ingredient_measurement_units valid_ingredient_measurement_units_valid_ingredient_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_measurement_units ADD CONSTRAINT valid_ingredient_measurement_units_valid_ingredient_id_fkey FOREIGN KEY (valid_ingredient_id) REFERENCES valid_ingredients(id) ON DELETE CASCADE;

--
-- Name: valid_ingredient_preparations valid_ingredient_preparations_valid_ingredient_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_preparations ADD CONSTRAINT valid_ingredient_preparations_valid_ingredient_id_fkey FOREIGN KEY (valid_ingredient_id) REFERENCES valid_ingredients(id) ON DELETE CASCADE;

--
-- Name: valid_ingredient_preparations valid_ingredient_preparations_valid_preparation_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_preparations ADD CONSTRAINT valid_ingredient_preparations_valid_preparation_id_fkey FOREIGN KEY (valid_preparation_id) REFERENCES valid_preparations(id) ON DELETE CASCADE;

--
-- Name: valid_ingredient_state_ingredients valid_ingredient_state_ingredients_valid_ingredient_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_state_ingredients ADD CONSTRAINT valid_ingredient_state_ingredients_valid_ingredient_fkey FOREIGN KEY (valid_ingredient) REFERENCES valid_ingredients(id) ON DELETE CASCADE;

--
-- Name: valid_ingredient_state_ingredients valid_ingredient_state_ingredients_valid_ingredient_state_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY valid_ingredient_state_ingredients ADD CONSTRAINT valid_ingredient_state_ingredients_valid_ingredient_state_fkey FOREIGN KEY (valid_ingredient_state) REFERENCES valid_ingredient_states(id) ON DELETE CASCADE;

--
-- Name: valid_measurement_conversions valid_measurement_conversions_from_unit_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY valid_measurement_conversions ADD CONSTRAINT valid_measurement_conversions_from_unit_fkey FOREIGN KEY (from_unit) REFERENCES valid_measurement_units(id) ON DELETE CASCADE;

--
-- Name: valid_measurement_conversions valid_measurement_conversions_only_for_ingredient_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY valid_measurement_conversions ADD CONSTRAINT valid_measurement_conversions_only_for_ingredient_fkey FOREIGN KEY (only_for_ingredient) REFERENCES valid_ingredients(id) ON DELETE CASCADE;

--
-- Name: valid_measurement_conversions valid_measurement_conversions_to_unit_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY valid_measurement_conversions ADD CONSTRAINT valid_measurement_conversions_to_unit_fkey FOREIGN KEY (to_unit) REFERENCES valid_measurement_units(id) ON DELETE CASCADE;

--
-- Name: valid_preparation_instruments valid_preparation_instruments_valid_instrument_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY valid_preparation_instruments ADD CONSTRAINT valid_preparation_instruments_valid_instrument_id_fkey FOREIGN KEY (valid_instrument_id) REFERENCES valid_instruments(id) ON DELETE CASCADE;

--
-- Name: valid_preparation_instruments valid_preparation_instruments_valid_preparation_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY valid_preparation_instruments ADD CONSTRAINT valid_preparation_instruments_valid_preparation_id_fkey FOREIGN KEY (valid_preparation_id) REFERENCES valid_preparations(id) ON DELETE CASCADE;

--
-- Name: valid_preparation_vessels valid_preparation_vessels_valid_preparation_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY valid_preparation_vessels ADD CONSTRAINT valid_preparation_vessels_valid_preparation_id_fkey FOREIGN KEY (valid_preparation_id) REFERENCES valid_preparations(id) ON DELETE CASCADE;

--
-- Name: valid_preparation_vessels valid_preparation_vessels_valid_vessel_id_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY valid_preparation_vessels ADD CONSTRAINT valid_preparation_vessels_valid_vessel_id_fkey FOREIGN KEY (valid_vessel_id) REFERENCES valid_vessels(id) ON DELETE CASCADE;

--
-- Name: valid_vessels valid_vessels_capacity_unit_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY valid_vessels ADD CONSTRAINT valid_vessels_capacity_unit_fkey FOREIGN KEY (capacity_unit) REFERENCES valid_measurement_units(id) ON DELETE CASCADE;

--
-- Name: webhook_trigger_events webhook_trigger_events_belongs_to_webhook_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY webhook_trigger_events ADD CONSTRAINT webhook_trigger_events_belongs_to_webhook_fkey FOREIGN KEY (belongs_to_webhook) REFERENCES webhooks(id) ON DELETE CASCADE;

--
-- Name: webhooks webhooks_belongs_to_household_fkey; Type: FK CONSTRAINT;
--

ALTER TABLE ONLY webhooks ADD CONSTRAINT webhooks_belongs_to_household_fkey FOREIGN KEY (belongs_to_household) REFERENCES households(id) ON DELETE CASCADE;

--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: cloudsqlsuperuser
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO PUBLIC;

--
-- Name: FUNCTION pg_replication_origin_advance(text, pg_lsn); Type: ACL; Schema: pg_catalog; Owner: cloudsqladmin
--

GRANT ALL ON FUNCTION pg_catalog.pg_replication_origin_advance(text, pg_lsn) TO cloudsqlsuperuser;

--
-- Name: FUNCTION pg_replication_origin_create(text); Type: ACL; Schema: pg_catalog; Owner: cloudsqladmin
--

GRANT ALL ON FUNCTION pg_catalog.pg_replication_origin_create(text) TO cloudsqlsuperuser;

--
-- Name: FUNCTION pg_replication_origin_drop(text); Type: ACL; Schema: pg_catalog; Owner: cloudsqladmin
--

GRANT ALL ON FUNCTION pg_catalog.pg_replication_origin_drop(text) TO cloudsqlsuperuser;

--
-- Name: FUNCTION pg_replication_origin_oid(text); Type: ACL; Schema: pg_catalog; Owner: cloudsqladmin
--

GRANT ALL ON FUNCTION pg_catalog.pg_replication_origin_oid(text) TO cloudsqlsuperuser;

--
-- Name: FUNCTION pg_replication_origin_progress(text, boolean); Type: ACL; Schema: pg_catalog; Owner: cloudsqladmin
--

GRANT ALL ON FUNCTION pg_catalog.pg_replication_origin_progress(text, boolean) TO cloudsqlsuperuser;

--
-- Name: FUNCTION pg_replication_origin_session_is_setup(); Type: ACL; Schema: pg_catalog; Owner: cloudsqladmin
--

GRANT ALL ON FUNCTION pg_catalog.pg_replication_origin_session_is_setup() TO cloudsqlsuperuser;

--
-- Name: FUNCTION pg_replication_origin_session_progress(boolean); Type: ACL; Schema: pg_catalog; Owner: cloudsqladmin
--

GRANT ALL ON FUNCTION pg_catalog.pg_replication_origin_session_progress(boolean) TO cloudsqlsuperuser;

--
-- Name: FUNCTION pg_replication_origin_session_reset(); Type: ACL; Schema: pg_catalog; Owner: cloudsqladmin
--

GRANT ALL ON FUNCTION pg_catalog.pg_replication_origin_session_reset() TO cloudsqlsuperuser;

--
-- Name: FUNCTION pg_replication_origin_session_setup(text); Type: ACL; Schema: pg_catalog; Owner: cloudsqladmin
--

GRANT ALL ON FUNCTION pg_catalog.pg_replication_origin_session_setup(text) TO cloudsqlsuperuser;

--
-- Name: FUNCTION pg_replication_origin_xact_reset(); Type: ACL; Schema: pg_catalog; Owner: cloudsqladmin
--

GRANT ALL ON FUNCTION pg_catalog.pg_replication_origin_xact_reset() TO cloudsqlsuperuser;

--
-- Name: FUNCTION pg_replication_origin_xact_setup(pg_lsn, TIMESTAMP WITH TIME ZONE); Type: ACL; Schema: pg_catalog; Owner: cloudsqladmin
--

GRANT ALL ON FUNCTION pg_catalog.pg_replication_origin_xact_setup(pg_lsn, TIMESTAMP WITH TIME ZONE) TO cloudsqlsuperuser;

--
-- Name: FUNCTION pg_show_replication_origin_status(OUT local_id oid, OUT external_id text, OUT remote_lsn pg_lsn, OUT local_lsn pg_lsn); Type: ACL; Schema: pg_catalog; Owner: cloudsqladmin
--

GRANT ALL ON FUNCTION pg_catalog.pg_show_replication_origin_status(OUT local_id oid, OUT external_id text, OUT remote_lsn pg_lsn, OUT local_lsn pg_lsn) TO cloudsqlsuperuser;

--
-- PostgreSQL database dump complete
--

