-- Meal Planning Domain Migration
-- All recipe, meal, and meal planning functionality

-- =============================================================================
-- ENUMERATED TYPES
-- =============================================================================

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

CREATE TYPE grocery_list_item_status AS ENUM (
    'unknown',
    'already owned',
    'needs',
    'unavailable',
    'acquired'
);

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

CREATE TYPE recipe_ingredient_scale AS ENUM (
    'linear',
    'logarithmic',
    'exponential'
);

CREATE TYPE recipe_step_product_type AS ENUM (
    'ingredient',
    'instrument',
    'vessel'
);

CREATE TYPE storage_container_type AS ENUM (
    'uncovered',
    'covered',
    'on a wire rack',
    'in an airtight container'
);

CREATE TYPE valid_election_method AS ENUM (
    'schulze',
    'instant-runoff'
);

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

CREATE TYPE recipe_status AS ENUM (
    'submitted',
    'approved',
    'needs revision'
);

-- =============================================================================
-- VALID ENTITY TABLES
-- =============================================================================

CREATE TABLE IF NOT EXISTS valid_ingredient_states (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    past_tense TEXT DEFAULT ''::TEXT NOT NULL,
    slug TEXT DEFAULT ''::TEXT NOT NULL,
    description TEXT DEFAULT ''::TEXT NOT NULL,
    icon_path TEXT DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    attribute_type ingredient_attribute_type DEFAULT 'other'::ingredient_attribute_type NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS valid_ingredients (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    warning TEXT NOT NULL,
    contains_egg BOOLEAN NOT NULL,
    contains_dairy BOOLEAN NOT NULL,
    contains_peanut BOOLEAN NOT NULL,
    contains_tree_nut BOOLEAN NOT NULL,
    contains_soy BOOLEAN NOT NULL,
    contains_wheat BOOLEAN NOT NULL,
    contains_shellfish BOOLEAN NOT NULL,
    contains_sesame BOOLEAN NOT NULL,
    contains_fish BOOLEAN NOT NULL,
    contains_gluten BOOLEAN NOT NULL,
    animal_flesh BOOLEAN NOT NULL,
    icon_path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    is_liquid BOOLEAN DEFAULT FALSE,
    animal_derived BOOLEAN DEFAULT FALSE NOT NULL,
    plural_name TEXT DEFAULT ''::TEXT NOT NULL,
    restrict_to_preparations BOOLEAN DEFAULT FALSE NOT NULL,
    minimum_ideal_storage_temperature_in_celsius NUMERIC(14,2),
    maximum_ideal_storage_temperature_in_celsius NUMERIC(14,2),
    storage_instructions TEXT DEFAULT ''::TEXT NOT NULL,
    contains_alcohol BOOLEAN DEFAULT FALSE NOT NULL,
    slug TEXT DEFAULT ''::TEXT NOT NULL,
    shopping_suggestions TEXT DEFAULT ''::TEXT NOT NULL,
    is_starch BOOLEAN DEFAULT FALSE NOT NULL,
    is_protein BOOLEAN DEFAULT FALSE NOT NULL,
    is_grain BOOLEAN DEFAULT FALSE NOT NULL,
    is_fruit BOOLEAN DEFAULT FALSE NOT NULL,
    is_salt BOOLEAN DEFAULT FALSE NOT NULL,
    is_fat BOOLEAN DEFAULT FALSE NOT NULL,
    is_acid BOOLEAN DEFAULT FALSE NOT NULL,
    is_heat BOOLEAN DEFAULT FALSE NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE,
    contaminates_equipment BOOLEAN DEFAULT FALSE NOT NULL
);

CREATE TABLE IF NOT EXISTS valid_ingredient_groups (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT DEFAULT ''::TEXT NOT NULL,
    slug TEXT DEFAULT ''::TEXT NOT NULL,
    description TEXT DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS valid_ingredient_group_members (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_group TEXT NOT NULL REFERENCES valid_ingredient_groups("id") ON DELETE CASCADE,
    valid_ingredient TEXT NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS valid_measurement_units (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT ''::TEXT NOT NULL,
    icon_path TEXT DEFAULT ''::TEXT NOT NULL,
    volumetric BOOLEAN DEFAULT FALSE,
    universal BOOLEAN DEFAULT FALSE NOT NULL,
    metric BOOLEAN DEFAULT FALSE NOT NULL,
    imperial BOOLEAN DEFAULT FALSE NOT NULL,
    plural_name TEXT DEFAULT ''::TEXT NOT NULL,
    slug TEXT DEFAULT ''::TEXT NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(name, archived_at)
);

CREATE TABLE IF NOT EXISTS valid_ingredient_measurement_units (
    id TEXT NOT NULL PRIMARY KEY,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    valid_ingredient_id TEXT NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    valid_measurement_unit_id TEXT NOT NULL REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    minimum_allowable_quantity NUMERIC(14,2) DEFAULT 0 NOT NULL,
    maximum_allowable_quantity NUMERIC(14,2),
    UNIQUE(valid_ingredient_id, valid_measurement_unit_id, archived_at)
);

CREATE TABLE IF NOT EXISTS valid_measurement_unit_conversions (
    id TEXT NOT NULL PRIMARY KEY,
    from_unit TEXT NOT NULL REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    to_unit TEXT NOT NULL REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    only_for_ingredient TEXT REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    modifier NUMERIC(14,5) NOT NULL,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE NULLS NOT DISTINCT (from_unit, to_unit, only_for_ingredient),
    CHECK (from_unit < to_unit)
);

CREATE TABLE IF NOT EXISTS valid_ingredient_state_ingredients (
    id TEXT NOT NULL PRIMARY KEY,
    valid_ingredient TEXT NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    valid_ingredient_state TEXT NOT NULL REFERENCES valid_ingredient_states("id") ON DELETE CASCADE,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(valid_ingredient, valid_ingredient_state, archived_at)
);

CREATE TABLE IF NOT EXISTS valid_preparations (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    icon_path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    yields_nothing BOOLEAN DEFAULT FALSE NOT NULL,
    restrict_to_ingredients BOOLEAN DEFAULT FALSE NOT NULL,
    past_tense TEXT DEFAULT ''::TEXT NOT NULL,
    slug TEXT DEFAULT ''::TEXT NOT NULL,
    minimum_ingredient_count INTEGER DEFAULT 1 NOT NULL,
    maximum_ingredient_count INTEGER,
    minimum_instrument_count INTEGER DEFAULT 1 NOT NULL,
    maximum_instrument_count INTEGER,
    temperature_required BOOLEAN DEFAULT FALSE NOT NULL,
    time_estimate_required BOOLEAN DEFAULT FALSE NOT NULL,
    condition_expression_required BOOLEAN DEFAULT FALSE NOT NULL,
    consumes_vessel BOOLEAN DEFAULT FALSE NOT NULL,
    only_for_vessels BOOLEAN DEFAULT FALSE NOT NULL,
    minimum_vessel_count INTEGER DEFAULT 0 NOT NULL,
    maximum_vessel_count INTEGER,
    last_indexed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(name, archived_at)
);

CREATE TABLE IF NOT EXISTS valid_ingredient_preparations (
    id TEXT NOT NULL PRIMARY KEY,
    notes TEXT NOT NULL,
    valid_preparation_id TEXT NOT NULL REFERENCES valid_preparations("id") ON DELETE CASCADE,
    valid_ingredient_id TEXT NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(valid_preparation_id, valid_ingredient_id, archived_at)
);

CREATE TABLE IF NOT EXISTS valid_prep_task_configs (
    id TEXT NOT NULL PRIMARY KEY,
    valid_ingredient_id TEXT NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    valid_preparation_id TEXT NOT NULL REFERENCES valid_preparations("id") ON DELETE CASCADE,
    minimum_storage_duration_in_seconds INTEGER NOT NULL,
    maximum_storage_duration_in_seconds INTEGER,
    storage_container_type storage_container_type NOT NULL DEFAULT 'covered',
    minimum_storage_temperature_in_celsius NUMERIC(14,2),
    maximum_storage_temperature_in_celsius NUMERIC(14,2),
    storage_instructions TEXT DEFAULT ''::TEXT NOT NULL,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    source TEXT DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE NULLS NOT DISTINCT (
        valid_ingredient_id,
        valid_preparation_id,
        storage_container_type,
        minimum_storage_temperature_in_celsius,
        maximum_storage_temperature_in_celsius,
        archived_at
    )
);

CREATE TABLE IF NOT EXISTS valid_instruments (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    icon_path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    plural_name TEXT DEFAULT ''::TEXT NOT NULL,
    usable_for_storage BOOLEAN DEFAULT FALSE NOT NULL,
    slug TEXT DEFAULT ''::TEXT NOT NULL,
    display_in_summary_lists BOOLEAN DEFAULT true NOT NULL,
    include_in_generated_instructions BOOLEAN DEFAULT true NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS valid_preparation_instruments (
    id TEXT NOT NULL PRIMARY KEY,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    valid_preparation_id TEXT NOT NULL REFERENCES valid_preparations("id") ON DELETE CASCADE,
    valid_instrument_id TEXT NOT NULL REFERENCES valid_instruments("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(valid_instrument_id, valid_preparation_id, archived_at)
);

CREATE TABLE IF NOT EXISTS valid_vessels (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    plural_name TEXT DEFAULT ''::TEXT NOT NULL,
    description TEXT DEFAULT ''::TEXT NOT NULL,
    icon_path TEXT NOT NULL,
    usable_for_storage BOOLEAN DEFAULT FALSE NOT NULL,
    slug TEXT DEFAULT ''::TEXT NOT NULL,
    display_in_summary_lists BOOLEAN DEFAULT true NOT NULL,
    include_in_generated_instructions BOOLEAN DEFAULT true NOT NULL,
    capacity NUMERIC(14,2) DEFAULT 0 NOT NULL,
    capacity_unit TEXT REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    width_in_millimeters NUMERIC(14,2),
    length_in_millimeters NUMERIC(14,2),
    height_in_millimeters NUMERIC(14,2),
    shape vessel_shape DEFAULT 'other'::vessel_shape NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(name, archived_at),
    UNIQUE(slug, archived_at)
);

CREATE TABLE IF NOT EXISTS valid_preparation_vessels (
    id TEXT NOT NULL PRIMARY KEY,
    valid_preparation_id TEXT NOT NULL REFERENCES valid_preparations("id") ON DELETE CASCADE,
    valid_vessel_id TEXT NOT NULL REFERENCES valid_vessels("id") ON DELETE CASCADE,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(valid_preparation_id, valid_vessel_id, archived_at)
);

-- =============================================================================
-- USER PREFERENCES AND ACCOUNT OWNERSHIPS
-- =============================================================================

CREATE TABLE IF NOT EXISTS user_ingredient_preferences (
    id TEXT NOT NULL PRIMARY KEY,
    ingredient TEXT NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    rating smallint DEFAULT 0 NOT NULL,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    allergy BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    UNIQUE(belongs_to_user, ingredient)
);

CREATE TABLE IF NOT EXISTS account_instrument_ownerships (
    id TEXT NOT NULL PRIMARY KEY,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    quantity INTEGER DEFAULT 0 NOT NULL,
    valid_instrument_id TEXT NOT NULL,
    belongs_to_account TEXT NOT NULL REFERENCES accounts("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(valid_instrument_id, belongs_to_account, archived_at)
);

-- =============================================================================
-- RECIPE TABLES
-- =============================================================================



CREATE TABLE IF NOT EXISTS recipes (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    source TEXT NOT NULL,
    description TEXT NOT NULL,
    status recipe_status NOT NULL DEFAULT 'submitted',
    inspired_by_recipe_id TEXT REFERENCES recipes("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    created_by_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    min_estimated_portions NUMERIC(14,2) DEFAULT 1 NOT NULL,
    slug TEXT DEFAULT ''::TEXT NOT NULL,
    portion_name TEXT DEFAULT 'portion'::TEXT NOT NULL,
    plural_portion_name TEXT DEFAULT 'portions'::TEXT NOT NULL,
    max_estimated_portions NUMERIC(14,2),
    eligible_for_meals BOOLEAN DEFAULT true NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE,
    last_validated_at TIMESTAMP WITH TIME ZONE,
    yields_component_type component_type DEFAULT 'unspecified'::component_type NOT NULL,
    source_isbn TEXT DEFAULT ''::TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_steps (
    id TEXT NOT NULL PRIMARY KEY,
    index INTEGER NOT NULL,
    preparation_id TEXT NOT NULL REFERENCES valid_preparations("id") ON DELETE CASCADE,
    minimum_estimated_time_in_seconds bigint,
    maximum_estimated_time_in_seconds bigint,
    minimum_temperature_in_celsius NUMERIC(14,2),
    notes TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_recipe TEXT NOT NULL REFERENCES recipes("id") ON DELETE CASCADE,
    optional BOOLEAN DEFAULT FALSE NOT NULL,
    maximum_temperature_in_celsius NUMERIC(14,2),
    explicit_instructions TEXT DEFAULT ''::TEXT NOT NULL,
    condition_expression TEXT DEFAULT ''::TEXT NOT NULL,
    start_timer_automatically BOOLEAN DEFAULT FALSE NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_step_products (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_recipe_step TEXT NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    quantity_notes TEXT NOT NULL,
    minimum_measurement_quantity_value NUMERIC(14,2),
    maximum_measurement_quantity_value NUMERIC(14,2),
    minimum_item_quantity_value NUMERIC(14,2),
    maximum_item_quantity_value NUMERIC(14,2),
    measurement_unit TEXT REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    type recipe_step_product_type NOT NULL,
    compostable BOOLEAN DEFAULT FALSE NOT NULL,
    maximum_storage_duration_in_seconds INTEGER,
    minimum_storage_temperature_in_celsius NUMERIC(14,2),
    maximum_storage_temperature_in_celsius NUMERIC(14,2),
    storage_instructions TEXT DEFAULT ''::TEXT NOT NULL,
    is_liquid BOOLEAN DEFAULT FALSE NOT NULL,
    is_waste BOOLEAN DEFAULT FALSE NOT NULL,
    index INTEGER DEFAULT 0 NOT NULL,
    contained_in_vessel_index INTEGER
);

CREATE TABLE IF NOT EXISTS recipe_step_ingredients (
    id TEXT NOT NULL PRIMARY KEY,
    ingredient_id TEXT REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    minimum_quantity_value NUMERIC(14,2) NOT NULL,
    quantity_notes TEXT NOT NULL,
    ingredient_notes TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_recipe_step TEXT NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    name TEXT NOT NULL,
    recipe_step_product_id TEXT REFERENCES recipe_step_products("id") ON DELETE CASCADE,
    maximum_quantity_value NUMERIC(14,2),
    measurement_unit TEXT REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    optional BOOLEAN DEFAULT FALSE NOT NULL,
    index INTEGER DEFAULT 0 NOT NULL,
    option_index INTEGER DEFAULT 0 NOT NULL,
    vessel_index INTEGER,
    to_taste BOOLEAN DEFAULT FALSE NOT NULL,
    product_percentage_to_use NUMERIC(14,2),
    recipe_step_product_recipe_id TEXT REFERENCES recipes("id") ON DELETE CASCADE,
    scale_factor NUMERIC(14,4) DEFAULT 1.0 NOT NULL,
    CONSTRAINT valid_instrument_or_product CHECK (((recipe_step_product_id IS NOT NULL) OR (ingredient_id IS NOT NULL))),
    UNIQUE(ingredient_id, belongs_to_recipe_step)
);

CREATE TABLE IF NOT EXISTS recipe_step_instruments (
    id TEXT NOT NULL PRIMARY KEY,
    instrument_id TEXT REFERENCES valid_instruments("id") ON DELETE CASCADE,
    notes TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_recipe_step TEXT NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    preference_rank INTEGER NOT NULL,
    recipe_step_product_id TEXT REFERENCES recipe_step_products("id") ON DELETE CASCADE,
    name TEXT DEFAULT ''::TEXT NOT NULL,
    optional BOOLEAN DEFAULT FALSE NOT NULL,
    minimum_quantity INTEGER DEFAULT 1 NOT NULL,
    maximum_quantity INTEGER,
    index INTEGER DEFAULT 0 NOT NULL,
    option_index INTEGER DEFAULT 0 NOT NULL,
    scale_factor NUMERIC(14,4) DEFAULT 1.0 NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_step_vessels (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT DEFAULT ''::TEXT NOT NULL,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    belongs_to_recipe_step TEXT NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    recipe_step_product_id TEXT REFERENCES recipe_step_products("id") ON DELETE CASCADE,
    vessel_predicate TEXT DEFAULT ''::TEXT NOT NULL,
    minimum_quantity INTEGER DEFAULT 0 NOT NULL,
    maximum_quantity INTEGER,
    unavailable_after_step BOOLEAN DEFAULT FALSE NOT NULL,
    index INTEGER DEFAULT 0 NOT NULL,
    option_index INTEGER DEFAULT 0 NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    valid_vessel_id TEXT REFERENCES valid_vessels("id") ON DELETE CASCADE,
    scale_factor NUMERIC(14,4) DEFAULT 1.0 NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_step_completion_conditions (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_recipe_step TEXT NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    ingredient_state TEXT NOT NULL REFERENCES valid_ingredient_states("id") ON DELETE CASCADE,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    optional BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(belongs_to_recipe_step, ingredient_state)
);

CREATE TABLE IF NOT EXISTS recipe_step_completion_condition_ingredients (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_recipe_step_completion_condition TEXT NOT NULL REFERENCES recipe_step_completion_conditions("id") ON DELETE CASCADE,
    recipe_step_ingredient TEXT NOT NULL REFERENCES recipe_step_ingredients("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS recipe_media (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_recipe TEXT REFERENCES recipes("id") ON DELETE CASCADE,
    belongs_to_recipe_step TEXT REFERENCES recipe_steps("id") ON DELETE CASCADE,
    mime_type TEXT NOT NULL,
    internal_path TEXT NOT NULL,
    external_path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    index INTEGER DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_prep_tasks (
    id TEXT NOT NULL PRIMARY KEY,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    explicit_storage_instructions TEXT DEFAULT ''::TEXT NOT NULL,
    minimum_time_buffer_before_recipe_in_seconds INTEGER NOT NULL,
    maximum_time_buffer_before_recipe_in_seconds INTEGER,
    storage_type storage_container_type,
    minimum_storage_temperature_in_celsius NUMERIC(14,2),
    maximum_storage_temperature_in_celsius NUMERIC(14,2),
    belongs_to_recipe TEXT NOT NULL REFERENCES recipes("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    name TEXT DEFAULT ''::TEXT NOT NULL,
    description TEXT DEFAULT ''::TEXT NOT NULL,
    optional BOOLEAN DEFAULT true NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_prep_task_steps (
    id TEXT NOT NULL PRIMARY KEY,
    satisfies_recipe_step BOOLEAN DEFAULT FALSE NOT NULL,
    belongs_to_recipe_step TEXT NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    belongs_to_recipe_prep_task TEXT NOT NULL REFERENCES recipe_prep_tasks("id") ON DELETE CASCADE,
    UNIQUE(belongs_to_recipe_step, belongs_to_recipe_prep_task)
);

CREATE TABLE IF NOT EXISTS recipe_ratings (
    id TEXT NOT NULL PRIMARY KEY,
    recipe_id TEXT NOT NULL REFERENCES recipes("id") ON DELETE CASCADE,
    taste NUMERIC(14,2),
    difficulty NUMERIC(14,2),
    cleanup NUMERIC(14,2),
    instructions NUMERIC(14,2),
    overall NUMERIC(14,2),
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    by_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(by_user, recipe_id)
);

-- =============================================================================
-- MEAL TABLES
-- =============================================================================

CREATE TABLE IF NOT EXISTS meals (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    created_by_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    min_estimated_portions NUMERIC(14,2) DEFAULT 1.0 NOT NULL,
    max_estimated_portions NUMERIC(14,2),
    eligible_for_meal_plans BOOLEAN DEFAULT true NOT NULL,
    last_indexed_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS meal_components (
    id TEXT NOT NULL PRIMARY KEY,
    meal_id TEXT NOT NULL REFERENCES meals("id") ON DELETE CASCADE,
    recipe_id TEXT NOT NULL REFERENCES recipes("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    meal_component_type component_type DEFAULT 'unspecified'::component_type NOT NULL,
    recipe_scale NUMERIC(14,2) DEFAULT 1.0 NOT NULL
);

-- =============================================================================
-- MEAL PLAN TABLES
-- =============================================================================

CREATE TABLE IF NOT EXISTS meal_plans (
    id TEXT NOT NULL PRIMARY KEY,
    notes TEXT NOT NULL,
    status meal_plan_status DEFAULT 'awaiting_votes'::meal_plan_status NOT NULL,
    voting_deadline TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_account TEXT NOT NULL REFERENCES accounts("id") ON DELETE CASCADE,
    grocery_list_initialized BOOLEAN DEFAULT FALSE NOT NULL,
    tasks_created BOOLEAN DEFAULT FALSE NOT NULL,
    election_method valid_election_method DEFAULT 'schulze'::valid_election_method NOT NULL,
    created_by_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS meal_plan_events (
    id TEXT NOT NULL PRIMARY KEY,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    starts_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    ends_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    meal_name meal_name NOT NULL,
    belongs_to_meal_plan TEXT NOT NULL REFERENCES meal_plans("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS meal_plan_options (
    id TEXT NOT NULL PRIMARY KEY,
    meal_id TEXT NOT NULL REFERENCES meals("id") ON DELETE CASCADE,
    notes TEXT NOT NULL,
    chosen BOOLEAN DEFAULT FALSE NOT NULL,
    tiebroken BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    assigned_cook TEXT REFERENCES users("id") ON DELETE CASCADE,
    assigned_dishwasher TEXT REFERENCES users("id") ON DELETE CASCADE,
    belongs_to_meal_plan_event TEXT REFERENCES meal_plan_events("id") ON DELETE CASCADE,
    meal_scale NUMERIC(14,2) DEFAULT 1.0 NOT NULL
);

CREATE TABLE IF NOT EXISTS meal_plan_grocery_list_items (
    id TEXT NOT NULL PRIMARY KEY,
    valid_ingredient TEXT NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    valid_measurement_unit TEXT NOT NULL REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    minimum_quantity_needed NUMERIC(14,2) NOT NULL,
    maximum_quantity_needed NUMERIC(14,2),
    quantity_purchased NUMERIC(14,2),
    purchased_measurement_unit TEXT REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    purchased_upc TEXT,
    purchase_price NUMERIC(14,2),
    status_explanation TEXT DEFAULT ''::TEXT NOT NULL,
    status grocery_list_item_status DEFAULT 'unknown'::grocery_list_item_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_meal_plan TEXT NOT NULL REFERENCES meal_plans("id") ON DELETE CASCADE,
    belongs_to_meal_plan_option TEXT REFERENCES meal_plan_options("id") ON DELETE CASCADE,
    recipe_id TEXT REFERENCES recipes("id") ON DELETE CASCADE,
    recipe_step_id TEXT REFERENCES recipe_steps("id") ON DELETE CASCADE,
    ingredient_index INTEGER,
    option_index INTEGER
);

CREATE TABLE IF NOT EXISTS meal_plan_option_votes (
    id TEXT NOT NULL PRIMARY KEY,
    rank INTEGER NOT NULL,
    abstain BOOLEAN NOT NULL,
    notes TEXT NOT NULL,
    by_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_meal_plan_option TEXT NOT NULL REFERENCES meal_plan_options("id") ON DELETE CASCADE,
    UNIQUE(by_user, belongs_to_meal_plan_option)
);

CREATE TABLE IF NOT EXISTS meal_plan_recipe_option_selections (
    id TEXT PRIMARY KEY,
    belongs_to_meal_plan_option TEXT NOT NULL REFERENCES meal_plan_options("id") ON DELETE CASCADE,
    recipe_id TEXT NOT NULL REFERENCES recipes("id") ON DELETE CASCADE,
    recipe_step_id TEXT NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    ingredient_index INTEGER NOT NULL,
    selected_option_index INTEGER NOT NULL,
    selection_type TEXT NOT NULL CHECK (selection_type IN ('ingredient', 'instrument', 'vessel')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(belongs_to_meal_plan_option, recipe_step_id, ingredient_index, selection_type)
);

CREATE TABLE IF NOT EXISTS meal_plan_tasks (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_meal_plan_option TEXT NOT NULL REFERENCES meal_plan_options("id") ON DELETE CASCADE,
    belongs_to_recipe_prep_task TEXT NOT NULL REFERENCES recipe_prep_tasks("id") ON DELETE CASCADE,
    creation_explanation TEXT DEFAULT ''::TEXT NOT NULL,
    status_explanation TEXT DEFAULT ''::TEXT NOT NULL,
    status prep_step_status DEFAULT 'unfinished'::prep_step_status NOT NULL,
    notification_sent_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    assigned_to_user TEXT REFERENCES users("id") ON DELETE CASCADE,
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS meal_lists (
    "id" TEXT NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL DEFAULT '',
    "description" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "belongs_to_user" TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS meal_list_items (
    "id" TEXT NOT NULL PRIMARY KEY,
    "meal_id" TEXT NOT NULL REFERENCES meals("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "notes" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "belongs_to_meal_list" TEXT NOT NULL REFERENCES meal_lists("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS recipe_lists (
    "id" TEXT NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL DEFAULT '',
    "description" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "belongs_to_user" TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS recipe_list_items (
    "id" TEXT NOT NULL PRIMARY KEY,
    "recipe_id" TEXT NOT NULL REFERENCES recipes("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "notes" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "belongs_to_recipe_list" TEXT NOT NULL REFERENCES recipe_lists("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS recipe_images (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_recipe TEXT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    uploaded_media_id TEXT NOT NULL REFERENCES uploaded_media(id) ON DELETE CASCADE,
    uploaded_by_user TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS meal_images (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_meal TEXT NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
    uploaded_media_id TEXT NOT NULL REFERENCES uploaded_media(id) ON DELETE CASCADE,
    uploaded_by_user TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);

-- Add user_temperature_unit service setting (celsius vs fahrenheit, default fahrenheit)
INSERT INTO service_settings (id, name, type, description, default_value, enumeration, admins_only)
VALUES (
   'd6me6i4n9qd3gcf5j1p0',
   'user_temperature_unit',
   'user',
   'Preferred unit for displaying temperatures (e.g. oven, storage)',
   'fahrenheit',
   'celsius|fahrenheit',
   false
);


-- =============================================================================
-- INDEXES FOR MEAL PLANNING TABLES
-- =============================================================================

-- Valid ingredient states indexes
CREATE INDEX idx_valid_ingredient_states_archived_at ON valid_ingredient_states (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_ingredient_states_name ON valid_ingredient_states (name) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_ingredient_states_slug ON valid_ingredient_states (slug) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_ingredient_states_type ON valid_ingredient_states (attribute_type) WHERE archived_at IS NULL;

-- Valid ingredients indexes
CREATE INDEX idx_valid_ingredients_archived_at ON valid_ingredients (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_ingredients_name ON valid_ingredients (name) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_ingredients_slug ON valid_ingredients (slug) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_ingredients_allergens ON valid_ingredients (contains_egg, contains_dairy, contains_peanut, contains_tree_nut, contains_soy, contains_wheat, contains_shellfish, contains_sesame, contains_fish, contains_gluten) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_ingredients_properties ON valid_ingredients (is_liquid, animal_derived, animal_flesh, is_starch, is_protein, is_grain, is_fruit, is_salt, is_fat, is_acid, is_heat) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_ingredients_indexing ON valid_ingredients (last_indexed_at) WHERE archived_at IS NULL;

-- Valid ingredient groups indexes
CREATE INDEX idx_valid_ingredient_groups_archived_at ON valid_ingredient_groups (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_ingredient_groups_name ON valid_ingredient_groups (name) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_ingredient_groups_slug ON valid_ingredient_groups (slug) WHERE archived_at IS NULL;

-- Valid ingredient group members indexes
CREATE INDEX idx_ingredient_group_members_group ON valid_ingredient_group_members (belongs_to_group) WHERE archived_at IS NULL;
CREATE INDEX idx_ingredient_group_members_ingredient ON valid_ingredient_group_members (valid_ingredient) WHERE archived_at IS NULL;

-- Valid measurement units indexes
CREATE INDEX idx_valid_measurement_units_archived_at ON valid_measurement_units (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_measurement_units_name ON valid_measurement_units (name) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_measurement_units_slug ON valid_measurement_units (slug) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_measurement_units_properties ON valid_measurement_units (volumetric, universal, metric, imperial) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_measurement_units_indexing ON valid_measurement_units (last_indexed_at) WHERE archived_at IS NULL;

-- Valid ingredient measurement units indexes
CREATE INDEX idx_ingredient_measurement_units_ingredient ON valid_ingredient_measurement_units (valid_ingredient_id) WHERE archived_at IS NULL;
CREATE INDEX idx_ingredient_measurement_units_unit ON valid_ingredient_measurement_units (valid_measurement_unit_id) WHERE archived_at IS NULL;
CREATE UNIQUE INDEX idx_valid_ingredient_measurement_units_ing_unit_active
    ON valid_ingredient_measurement_units (valid_ingredient_id, valid_measurement_unit_id)
    WHERE archived_at IS NULL;

-- Valid measurement unit conversions indexes
CREATE INDEX idx_measurement_conversions_from_unit ON valid_measurement_unit_conversions (from_unit) WHERE archived_at IS NULL;
CREATE INDEX idx_measurement_conversions_to_unit ON valid_measurement_unit_conversions (to_unit) WHERE archived_at IS NULL;
CREATE INDEX idx_measurement_conversions_ingredient ON valid_measurement_unit_conversions (only_for_ingredient) WHERE archived_at IS NULL;
CREATE INDEX idx_measurement_conversions_from_to ON valid_measurement_unit_conversions (from_unit, to_unit) WHERE archived_at IS NULL;

-- Valid ingredient state ingredients indexes
CREATE INDEX idx_ingredient_state_ingredients_ingredient ON valid_ingredient_state_ingredients (valid_ingredient) WHERE archived_at IS NULL;
CREATE INDEX idx_ingredient_state_ingredients_state ON valid_ingredient_state_ingredients (valid_ingredient_state) WHERE archived_at IS NULL;
CREATE UNIQUE INDEX idx_valid_ingredient_state_ingredients_ing_state_active
    ON valid_ingredient_state_ingredients (valid_ingredient, valid_ingredient_state)
    WHERE archived_at IS NULL;

-- Valid preparations indexes
CREATE INDEX idx_valid_preparations_archived_at ON valid_preparations (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_preparations_name ON valid_preparations (name) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_preparations_slug ON valid_preparations (slug) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_preparations_properties ON valid_preparations (yields_nothing, restrict_to_ingredients, temperature_required, time_estimate_required, consumes_vessel, only_for_vessels) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_preparations_indexing ON valid_preparations (last_indexed_at) WHERE archived_at IS NULL;

-- Valid ingredient preparations indexes
CREATE INDEX idx_ingredient_preparations_preparation ON valid_ingredient_preparations (valid_preparation_id) WHERE archived_at IS NULL;
CREATE INDEX idx_ingredient_preparations_ingredient ON valid_ingredient_preparations (valid_ingredient_id) WHERE archived_at IS NULL;
CREATE UNIQUE INDEX idx_valid_ingredient_preparations_prep_ing_active
    ON valid_ingredient_preparations (valid_preparation_id, valid_ingredient_id)
    WHERE archived_at IS NULL;

-- Valid instruments indexes
CREATE INDEX idx_valid_instruments_archived_at ON valid_instruments (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_instruments_name ON valid_instruments (name) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_instruments_slug ON valid_instruments (slug) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_instruments_properties ON valid_instruments (usable_for_storage, display_in_summary_lists, include_in_generated_instructions) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_instruments_indexing ON valid_instruments (last_indexed_at) WHERE archived_at IS NULL;

-- Valid preparation instruments indexes
CREATE INDEX idx_preparation_instruments_preparation ON valid_preparation_instruments (valid_preparation_id) WHERE archived_at IS NULL;
CREATE INDEX idx_preparation_instruments_instrument ON valid_preparation_instruments (valid_instrument_id) WHERE archived_at IS NULL;
CREATE UNIQUE INDEX idx_valid_preparation_instruments_prep_instrument_active
    ON valid_preparation_instruments (valid_preparation_id, valid_instrument_id)
    WHERE archived_at IS NULL;

-- Valid vessels indexes
CREATE INDEX idx_valid_vessels_archived_at ON valid_vessels (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_vessels_name ON valid_vessels (name) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_vessels_slug ON valid_vessels (slug) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_vessels_properties ON valid_vessels (usable_for_storage, display_in_summary_lists, include_in_generated_instructions) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_vessels_capacity ON valid_vessels (capacity_unit) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_vessels_shape ON valid_vessels (shape) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_vessels_indexing ON valid_vessels (last_indexed_at) WHERE archived_at IS NULL;

-- Valid preparation vessels indexes
CREATE INDEX idx_preparation_vessels_preparation ON valid_preparation_vessels (valid_preparation_id) WHERE archived_at IS NULL;
CREATE INDEX idx_preparation_vessels_vessel ON valid_preparation_vessels (valid_vessel_id) WHERE archived_at IS NULL;
CREATE UNIQUE INDEX idx_valid_preparation_vessels_prep_vessel_active
    ON valid_preparation_vessels (valid_preparation_id, valid_vessel_id)
    WHERE archived_at IS NULL;

-- User ingredient preferences indexes
CREATE INDEX idx_user_ingredient_preferences_user ON user_ingredient_preferences (belongs_to_user) WHERE archived_at IS NULL;
CREATE INDEX idx_user_ingredient_preferences_ingredient ON user_ingredient_preferences (ingredient) WHERE archived_at IS NULL;
CREATE INDEX idx_user_ingredient_preferences_allergy ON user_ingredient_preferences (belongs_to_user, allergy) WHERE archived_at IS NULL AND allergy = TRUE;
CREATE INDEX idx_user_ingredient_preferences_rating ON user_ingredient_preferences (belongs_to_user, rating) WHERE archived_at IS NULL;

-- Account instrument ownerships indexes
CREATE INDEX idx_instrument_ownerships_account ON account_instrument_ownerships (belongs_to_account) WHERE archived_at IS NULL;
CREATE INDEX idx_instrument_ownerships_instrument ON account_instrument_ownerships (valid_instrument_id) WHERE archived_at IS NULL;


-- Recipes indexes
CREATE INDEX idx_recipes_archived_at ON recipes (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_recipes_created_by_user ON recipes (created_by_user) WHERE archived_at IS NULL;
CREATE INDEX idx_recipes_name ON recipes (name) WHERE archived_at IS NULL;
CREATE INDEX idx_recipes_slug ON recipes (slug) WHERE archived_at IS NULL;
CREATE INDEX idx_recipes_inspired_by ON recipes (inspired_by_recipe_id) WHERE archived_at IS NULL;
CREATE INDEX idx_recipes_eligible_meals ON recipes (eligible_for_meals) WHERE archived_at IS NULL AND eligible_for_meals = TRUE;
CREATE INDEX idx_recipes_component_type ON recipes (yields_component_type) WHERE archived_at IS NULL;
CREATE INDEX idx_recipes_indexing ON recipes (last_indexed_at) WHERE archived_at IS NULL;
CREATE INDEX idx_recipes_validation ON recipes (last_validated_at) WHERE archived_at IS NULL;
CREATE INDEX idx_recipes_user_created_at ON recipes (created_by_user, created_at) WHERE archived_at IS NULL;
CREATE INDEX idx_recipes_created_at_id ON recipes (created_at, id) WHERE archived_at IS NULL;
CREATE INDEX idx_recipes_validation_needed ON recipes (last_validated_at) WHERE archived_at IS NULL AND last_validated_at IS NULL;
CREATE INDEX idx_recipes_indexing_needed ON recipes (last_indexed_at) WHERE archived_at IS NULL;

-- Recipe steps indexes
CREATE INDEX idx_recipe_steps_recipe ON recipe_steps (belongs_to_recipe) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_steps_recipe_all ON recipe_steps (belongs_to_recipe); -- Non-partial for edge cases
CREATE INDEX idx_recipe_steps_preparation ON recipe_steps (preparation_id) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_steps_recipe_index ON recipe_steps (belongs_to_recipe, index) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_steps_optional ON recipe_steps (belongs_to_recipe, optional) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_steps_recipe_preparation_index ON recipe_steps (belongs_to_recipe, preparation_id, index) WHERE archived_at IS NULL;

-- Recipe step products indexes
CREATE INDEX idx_recipe_step_products_step ON recipe_step_products (belongs_to_recipe_step) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_products_measurement_unit ON recipe_step_products (measurement_unit) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_products_type ON recipe_step_products (type) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_products_step_index ON recipe_step_products (belongs_to_recipe_step, index) WHERE archived_at IS NULL;

-- Recipe step ingredients indexes
CREATE INDEX idx_recipe_step_ingredients_step ON recipe_step_ingredients (belongs_to_recipe_step) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_ingredients_step_all ON recipe_step_ingredients (belongs_to_recipe_step); -- Non-partial for edge cases
CREATE INDEX idx_recipe_step_ingredients_ingredient ON recipe_step_ingredients (ingredient_id) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_ingredients_measurement_unit ON recipe_step_ingredients (measurement_unit) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_ingredients_product ON recipe_step_ingredients (recipe_step_product_id) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_ingredients_product_recipe ON recipe_step_ingredients (recipe_step_product_recipe_id) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_ingredients_optional ON recipe_step_ingredients (belongs_to_recipe_step, optional) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_ingredients_step_index ON recipe_step_ingredients (belongs_to_recipe_step, index) WHERE archived_at IS NULL;
CREATE UNIQUE INDEX idx_recipe_step_ingredients_step_index_option_index_unique ON recipe_step_ingredients (belongs_to_recipe_step, index, option_index) WHERE archived_at IS NULL;

-- Recipe step instruments indexes
CREATE INDEX idx_recipe_step_instruments_step ON recipe_step_instruments (belongs_to_recipe_step) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_instruments_instrument ON recipe_step_instruments (instrument_id) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_instruments_product ON recipe_step_instruments (recipe_step_product_id) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_instruments_optional ON recipe_step_instruments (belongs_to_recipe_step, optional) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_instruments_step_index ON recipe_step_instruments (belongs_to_recipe_step, index) WHERE archived_at IS NULL;
CREATE UNIQUE INDEX idx_recipe_step_instruments_step_index_option_index_unique ON recipe_step_instruments (belongs_to_recipe_step, index, option_index) WHERE archived_at IS NULL;

-- Recipe step vessels indexes
CREATE INDEX idx_recipe_step_vessels_step ON recipe_step_vessels (belongs_to_recipe_step) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_vessels_vessel ON recipe_step_vessels (valid_vessel_id) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_vessels_product ON recipe_step_vessels (recipe_step_product_id) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_vessels_step_index ON recipe_step_vessels (belongs_to_recipe_step, index) WHERE archived_at IS NULL;
CREATE UNIQUE INDEX idx_recipe_step_vessels_step_index_option_index_unique ON recipe_step_vessels (belongs_to_recipe_step, index, option_index) WHERE archived_at IS NULL;

-- Recipe step completion conditions indexes
CREATE INDEX idx_recipe_step_conditions_step ON recipe_step_completion_conditions (belongs_to_recipe_step) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_step_conditions_state ON recipe_step_completion_conditions (ingredient_state) WHERE archived_at IS NULL;

-- Recipe step completion condition ingredients indexes
CREATE INDEX idx_condition_ingredients_condition ON recipe_step_completion_condition_ingredients (belongs_to_recipe_step_completion_condition) WHERE archived_at IS NULL;
CREATE INDEX idx_condition_ingredients_ingredient ON recipe_step_completion_condition_ingredients (recipe_step_ingredient) WHERE archived_at IS NULL;

-- Recipe media indexes
CREATE INDEX idx_recipe_media_recipe ON recipe_media (belongs_to_recipe) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_media_step ON recipe_media (belongs_to_recipe_step) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_media_recipe_index ON recipe_media (belongs_to_recipe, index) WHERE archived_at IS NULL;

-- Recipe prep tasks indexes
CREATE INDEX idx_recipe_prep_tasks_recipe ON recipe_prep_tasks (belongs_to_recipe) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_prep_tasks_optional ON recipe_prep_tasks (belongs_to_recipe, optional) WHERE archived_at IS NULL;

-- Recipe prep task steps indexes
CREATE INDEX idx_prep_task_steps_task ON recipe_prep_task_steps (belongs_to_recipe_prep_task);
CREATE INDEX idx_prep_task_steps_step ON recipe_prep_task_steps (belongs_to_recipe_step);

-- Recipe ratings indexes
CREATE INDEX idx_recipe_ratings_recipe ON recipe_ratings (recipe_id) WHERE archived_at IS NULL;
CREATE INDEX idx_recipe_ratings_user ON recipe_ratings (by_user) WHERE archived_at IS NULL;

-- Meals indexes
CREATE INDEX idx_meals_archived_at ON meals (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_meals_created_by_user ON meals (created_by_user) WHERE archived_at IS NULL;
CREATE INDEX idx_meals_name ON meals (name) WHERE archived_at IS NULL;
CREATE INDEX idx_meals_eligible_plans ON meals (eligible_for_meal_plans) WHERE archived_at IS NULL AND eligible_for_meal_plans = TRUE;
CREATE INDEX idx_meals_indexing ON meals (last_indexed_at) WHERE archived_at IS NULL;
CREATE INDEX idx_meals_user_created_at ON meals (created_by_user, created_at) WHERE archived_at IS NULL;
CREATE INDEX idx_meals_created_at_id ON meals (created_at, id) WHERE archived_at IS NULL;
CREATE INDEX idx_meals_indexing_needed ON meals (last_indexed_at) WHERE archived_at IS NULL;

-- Meal components indexes
CREATE INDEX idx_meal_components_meal ON meal_components (meal_id) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_components_recipe ON meal_components (recipe_id) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_components_type ON meal_components (meal_component_type) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_components_meal_recipe ON meal_components (meal_id, recipe_id) WHERE archived_at IS NULL;

-- Meal plans indexes
CREATE INDEX idx_meal_plans_archived_at ON meal_plans (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plans_account ON meal_plans (belongs_to_account) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plans_created_by ON meal_plans (created_by_user) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plans_status ON meal_plans (status) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plans_voting_deadline ON meal_plans (voting_deadline) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plans_expired_unresolved ON meal_plans (status, voting_deadline) WHERE archived_at IS NULL AND status = 'awaiting_votes';
CREATE INDEX idx_meal_plans_finalized_no_groceries ON meal_plans (status, grocery_list_initialized) WHERE archived_at IS NULL AND status = 'finalized' AND grocery_list_initialized = FALSE;
CREATE INDEX idx_meal_plans_finalized_no_tasks ON meal_plans (status, tasks_created) WHERE archived_at IS NULL AND status = 'finalized' AND tasks_created = FALSE;
CREATE INDEX idx_meal_plans_account_created_at ON meal_plans (belongs_to_account, created_at) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plans_created_at_id ON meal_plans (created_at, id) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plans_account_status_voting ON meal_plans (belongs_to_account, status, voting_deadline) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plan_finalization ON meal_plans (id, status, belongs_to_account) WHERE archived_at IS NULL AND status = 'finalized';

-- Meal plan events indexes
CREATE INDEX idx_meal_plan_events_plan ON meal_plan_events (belongs_to_meal_plan) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plan_events_meal_name ON meal_plan_events (meal_name) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plan_events_starts_at ON meal_plan_events (starts_at) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plan_events_plan_starts ON meal_plan_events (belongs_to_meal_plan, starts_at) WHERE archived_at IS NULL;

-- Meal plan options indexes
CREATE INDEX idx_meal_plan_options_event ON meal_plan_options (belongs_to_meal_plan_event) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plan_options_meal ON meal_plan_options (meal_id) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plan_options_cook ON meal_plan_options (assigned_cook) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plan_options_dishwasher ON meal_plan_options (assigned_dishwasher) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plan_options_chosen ON meal_plan_options (belongs_to_meal_plan_event, chosen) WHERE archived_at IS NULL AND chosen = TRUE;
CREATE INDEX idx_meal_plan_options_chosen_with_meal ON meal_plan_options (belongs_to_meal_plan_event, chosen, meal_id) WHERE archived_at IS NULL AND chosen = TRUE;

-- Meal plan option votes indexes
CREATE INDEX idx_meal_plan_votes_option ON meal_plan_option_votes (belongs_to_meal_plan_option) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plan_votes_user ON meal_plan_option_votes (by_user) WHERE archived_at IS NULL;
CREATE INDEX idx_meal_plan_votes_option_rank ON meal_plan_option_votes (belongs_to_meal_plan_option, rank) WHERE archived_at IS NULL;

-- Meal plan grocery list items indexes
CREATE INDEX idx_grocery_list_items_plan ON meal_plan_grocery_list_items (belongs_to_meal_plan) WHERE archived_at IS NULL;
CREATE INDEX idx_grocery_list_items_ingredient ON meal_plan_grocery_list_items (valid_ingredient) WHERE archived_at IS NULL;
CREATE INDEX idx_grocery_list_items_measurement_unit ON meal_plan_grocery_list_items (valid_measurement_unit) WHERE archived_at IS NULL;
CREATE INDEX idx_grocery_list_items_purchased_unit ON meal_plan_grocery_list_items (purchased_measurement_unit) WHERE archived_at IS NULL;
CREATE INDEX idx_grocery_list_items_status ON meal_plan_grocery_list_items (status) WHERE archived_at IS NULL;
CREATE INDEX idx_grocery_list_items_plan_status ON meal_plan_grocery_list_items (belongs_to_meal_plan, status) WHERE archived_at IS NULL;
CREATE INDEX idx_grocery_list_choice_group ON meal_plan_grocery_list_items (belongs_to_meal_plan_option, recipe_step_id, ingredient_index) WHERE belongs_to_meal_plan_option IS NOT NULL;

-- Meal plan tasks indexes
CREATE INDEX idx_meal_plan_tasks_option ON meal_plan_tasks (belongs_to_meal_plan_option);
CREATE INDEX idx_meal_plan_tasks_prep_task ON meal_plan_tasks (belongs_to_recipe_prep_task);
CREATE INDEX idx_meal_plan_tasks_assigned_user ON meal_plan_tasks (assigned_to_user);
CREATE INDEX idx_meal_plan_tasks_status ON meal_plan_tasks (status);
CREATE INDEX idx_meal_plan_tasks_user_status ON meal_plan_tasks (assigned_to_user, status);

-- Meal plan recipe option selections indexes
CREATE INDEX idx_option_selections_meal_plan_option ON meal_plan_recipe_option_selections (belongs_to_meal_plan_option);
CREATE INDEX idx_option_selections_recipe_step ON meal_plan_recipe_option_selections (recipe_step_id, ingredient_index, selection_type);

-- Valid prep task configs indexes
CREATE INDEX idx_valid_prep_task_configs_ingredient ON valid_prep_task_configs (valid_ingredient_id) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_prep_task_configs_preparation ON valid_prep_task_configs (valid_preparation_id) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_prep_task_configs_combo ON valid_prep_task_configs (valid_ingredient_id, valid_preparation_id) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_prep_task_configs_container_type ON valid_prep_task_configs (storage_container_type) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_prep_task_configs_temp_range ON valid_prep_task_configs (minimum_storage_temperature_in_celsius, maximum_storage_temperature_in_celsius) WHERE archived_at IS NULL;
CREATE INDEX idx_valid_prep_task_configs_archived_at ON valid_prep_task_configs (archived_at) WHERE archived_at IS NULL;

-- Prevent duplicate meals in meal lists: one meal per list (non-archived).
CREATE UNIQUE INDEX IF NOT EXISTS idx_meal_list_items_meal_list_meal_unique
    ON meal_list_items (belongs_to_meal_list, meal_id)
    WHERE archived_at IS NULL;

-- Prevent duplicate meal options per event: one meal per event (non-archived).
CREATE UNIQUE INDEX IF NOT EXISTS idx_meal_plan_options_event_meal_unique
    ON meal_plan_options (belongs_to_meal_plan_event, meal_id)
    WHERE archived_at IS NULL
    AND belongs_to_meal_plan_event IS NOT NULL;

-- =============================================================================
-- MEDIA FOR PREPARATIONS, INGREDIENTS, AND RECIPE STEPS
-- =============================================================================

CREATE TABLE IF NOT EXISTS preparation_media (
    id TEXT NOT NULL PRIMARY KEY,
    valid_preparation_id TEXT NOT NULL REFERENCES valid_preparations(id) ON DELETE CASCADE,
    for_ingredient_id TEXT REFERENCES valid_ingredients(id) ON DELETE CASCADE,
    uploaded_media_id TEXT NOT NULL REFERENCES uploaded_media(id) ON DELETE CASCADE,
    index INTEGER DEFAULT 0 NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE NULLS NOT DISTINCT (valid_preparation_id, for_ingredient_id, uploaded_media_id, archived_at)
);

CREATE TABLE IF NOT EXISTS ingredient_media (
    id TEXT NOT NULL PRIMARY KEY,
    valid_ingredient_id TEXT NOT NULL REFERENCES valid_ingredients(id) ON DELETE CASCADE,
    uploaded_media_id TEXT NOT NULL REFERENCES uploaded_media(id) ON DELETE CASCADE,
    index INTEGER DEFAULT 0 NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(valid_ingredient_id, uploaded_media_id, archived_at)
);

CREATE TABLE IF NOT EXISTS recipe_step_images (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_recipe_step TEXT NOT NULL REFERENCES recipe_steps(id) ON DELETE CASCADE,
    uploaded_media_id TEXT NOT NULL REFERENCES uploaded_media(id) ON DELETE CASCADE,
    uploaded_by_user TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_preparation_media_preparation ON preparation_media (valid_preparation_id) WHERE archived_at IS NULL;
CREATE INDEX idx_preparation_media_ingredient ON preparation_media (for_ingredient_id) WHERE archived_at IS NULL AND for_ingredient_id IS NOT NULL;

CREATE INDEX idx_ingredient_media_ingredient ON ingredient_media (valid_ingredient_id) WHERE archived_at IS NULL;

CREATE INDEX idx_recipe_step_images_step ON recipe_step_images (belongs_to_recipe_step) WHERE archived_at IS NULL;

-- =============================================================================
-- SEED DATA: mealplanning permissions
-- =============================================================================

INSERT INTO permissions (id, name, description) VALUES
    -- valid instruments
    ('perm_create_valid_instruments', 'create.valid_instruments', 'Create valid instruments'),
    ('perm_read_valid_instruments', 'read.valid_instruments', 'Read valid instruments'),
    ('perm_search_valid_instruments', 'search.valid_instruments', 'Search valid instruments'),
    ('perm_update_valid_instruments', 'update.valid_instruments', 'Update valid instruments'),
    ('perm_archive_valid_instruments', 'archive.valid_instruments', 'Archive valid instruments'),
    -- valid vessels
    ('perm_create_valid_vessels', 'create.valid_vessels', 'Create valid vessels'),
    ('perm_read_valid_vessels', 'read.valid_vessels', 'Read valid vessels'),
    ('perm_search_valid_vessels', 'search.valid_vessels', 'Search valid vessels'),
    ('perm_update_valid_vessels', 'update.valid_vessels', 'Update valid vessels'),
    ('perm_archive_valid_vessels', 'archive.valid_vessels', 'Archive valid vessels'),
    -- valid ingredients
    ('perm_create_valid_ingredients', 'create.valid_ingredients', 'Create valid ingredients'),
    ('perm_read_valid_ingredients', 'read.valid_ingredients', 'Read valid ingredients'),
    ('perm_search_valid_ingredients', 'search.valid_ingredients', 'Search valid ingredients'),
    ('perm_update_valid_ingredients', 'update.valid_ingredients', 'Update valid ingredients'),
    ('perm_archive_valid_ingredients', 'archive.valid_ingredients', 'Archive valid ingredients'),
    -- valid ingredient groups
    ('perm_create_valid_ingredient_groups', 'create.valid_ingredient_groups', 'Create valid ingredient groups'),
    ('perm_read_valid_ingredient_groups', 'read.valid_ingredient_groups', 'Read valid ingredient groups'),
    ('perm_search_valid_ingredient_groups', 'search.valid_ingredient_groups', 'Search valid ingredient groups'),
    ('perm_update_valid_ingredient_groups', 'update.valid_ingredient_groups', 'Update valid ingredient groups'),
    ('perm_archive_valid_ingredient_groups', 'archive.valid_ingredient_groups', 'Archive valid ingredient groups'),
    -- valid preparations
    ('perm_create_valid_preparations', 'create.valid_preparations', 'Create valid preparations'),
    ('perm_read_valid_preparations', 'read.valid_preparations', 'Read valid preparations'),
    ('perm_search_valid_preparations', 'search.valid_preparations', 'Search valid preparations'),
    ('perm_update_valid_preparations', 'update.valid_preparations', 'Update valid preparations'),
    ('perm_archive_valid_preparations', 'archive.valid_preparations', 'Archive valid preparations'),
    -- measurement units
    ('perm_create_measurement_units', 'create.measurement_units', 'Create measurement units'),
    ('perm_read_measurement_units', 'read.measurement_units', 'Read measurement units'),
    ('perm_search_measurement_units', 'search.measurement_units', 'Search measurement units'),
    ('perm_update_measurement_units', 'update.measurement_units', 'Update measurement units'),
    ('perm_archive_measurement_units', 'archive.measurement_units', 'Archive measurement units'),
    -- measurement conversions
    ('perm_create_measurement_conversions', 'create.measurement_conversions', 'Create measurement conversions'),
    ('perm_read_measurement_conversions', 'read.measurement_conversions', 'Read measurement conversions'),
    ('perm_update_measurement_conversions', 'update.measurement_conversions', 'Update measurement conversions'),
    ('perm_archive_measurement_conversions', 'archive.measurement_conversions', 'Archive measurement conversions'),
    -- valid ingredient preparations
    ('perm_create_valid_ingredient_preparations', 'create.valid_ingredient_preparations', 'Create valid ingredient preparations'),
    ('perm_read_valid_ingredient_preparations', 'read.valid_ingredient_preparations', 'Read valid ingredient preparations'),
    ('perm_search_valid_ingredient_preparations', 'search.valid_ingredient_preparations', 'Search valid ingredient preparations'),
    ('perm_update_valid_ingredient_preparations', 'update.valid_ingredient_preparations', 'Update valid ingredient preparations'),
    ('perm_archive_valid_ingredient_preparations', 'archive.valid_ingredient_preparations', 'Archive valid ingredient preparations'),
    -- valid prep task configs
    ('perm_create_valid_prep_task_configs', 'create.valid_prep_task_configs', 'Create valid prep task configs'),
    ('perm_read_valid_prep_task_configs', 'read.valid_prep_task_configs', 'Read valid prep task configs'),
    ('perm_update_valid_prep_task_configs', 'update.valid_prep_task_configs', 'Update valid prep task configs'),
    ('perm_archive_valid_prep_task_configs', 'archive.valid_prep_task_configs', 'Archive valid prep task configs'),
    -- valid ingredient state ingredients
    ('perm_create_valid_ingredient_state_ingredients', 'create.valid_ingredient_state_ingredients', 'Create valid ingredient state ingredients'),
    ('perm_read_valid_ingredient_state_ingredients', 'read.valid_ingredient_state_ingredients', 'Read valid ingredient state ingredients'),
    ('perm_search_valid_ingredient_state_ingredients', 'search.valid_ingredient_state_ingredients', 'Search valid ingredient state ingredients'),
    ('perm_update_valid_ingredient_state_ingredients', 'update.valid_ingredient_state_ingredients', 'Update valid ingredient state ingredients'),
    ('perm_archive_valid_ingredient_state_ingredients', 'archive.valid_ingredient_state_ingredients', 'Archive valid ingredient state ingredients'),
    -- valid preparation instruments
    ('perm_create_valid_preparation_instruments', 'create.valid_preparation_instruments', 'Create valid preparation instruments'),
    ('perm_read_valid_preparation_instruments', 'read.valid_preparation_instruments', 'Read valid preparation instruments'),
    ('perm_search_valid_preparation_instruments', 'search.valid_preparation_instruments', 'Search valid preparation instruments'),
    ('perm_update_valid_preparation_instruments', 'update.valid_preparation_instruments', 'Update valid preparation instruments'),
    ('perm_archive_valid_preparation_instruments', 'archive.valid_preparation_instruments', 'Archive valid preparation instruments'),
    -- valid preparation vessels
    ('perm_create_valid_preparation_vessels', 'create.valid_preparation_vessels', 'Create valid preparation vessels'),
    ('perm_read_valid_preparation_vessels', 'read.valid_preparation_vessels', 'Read valid preparation vessels'),
    ('perm_search_valid_preparation_vessels', 'search.valid_preparation_vessels', 'Search valid preparation vessels'),
    ('perm_update_valid_preparation_vessels', 'update.valid_preparation_vessels', 'Update valid preparation vessels'),
    ('perm_archive_valid_preparation_vessels', 'archive.valid_preparation_vessels', 'Archive valid preparation vessels'),
    -- valid ingredient measurement units
    ('perm_create_valid_ingredient_measurement_units', 'create.valid_ingredient_measurement_units', 'Create valid ingredient measurement units'),
    ('perm_read_valid_ingredient_measurement_units', 'read.valid_ingredient_measurement_units', 'Read valid ingredient measurement units'),
    ('perm_search_valid_ingredient_measurement_units', 'search.valid_ingredient_measurement_units', 'Search valid ingredient measurement units'),
    ('perm_update_valid_ingredient_measurement_units', 'update.valid_ingredient_measurement_units', 'Update valid ingredient measurement units'),
    ('perm_archive_valid_ingredient_measurement_units', 'archive.valid_ingredient_measurement_units', 'Archive valid ingredient measurement units'),
    -- valid ingredient states
    ('perm_create_valid_ingredient_states', 'create.valid_ingredient_states', 'Create valid ingredient states'),
    ('perm_read_valid_ingredient_states', 'read.valid_ingredient_states', 'Read valid ingredient states'),
    ('perm_update_valid_ingredient_states', 'update.valid_ingredient_states', 'Update valid ingredient states'),
    ('perm_archive_valid_ingredient_states', 'archive.valid_ingredient_states', 'Archive valid ingredient states'),
    -- meals
    ('perm_create_meals', 'create.meals', 'Create meals'),
    ('perm_read_meals', 'read.meals', 'Read meals'),
    ('perm_update_meals', 'update.meals', 'Update meals'),
    ('perm_archive_meals', 'archive.meals', 'Archive meals'),
    -- recipes
    ('perm_create_recipes', 'create.recipes', 'Create recipes'),
    ('perm_read_recipes', 'read.recipes', 'Read recipes'),
    ('perm_search_recipes', 'search.recipes', 'Search recipes'),
    ('perm_update_recipes', 'update.recipes', 'Update recipes'),
    ('perm_archive_recipes', 'archive.recipes', 'Archive recipes'),
    ('perm_update_recipe_status', 'update.recipe_status', 'Update recipe status'),
    -- recipe steps
    ('perm_create_recipe_steps', 'create.recipe_steps', 'Create recipe steps'),
    ('perm_read_recipe_steps', 'read.recipe_steps', 'Read recipe steps'),
    ('perm_search_recipe_steps', 'search.recipe_steps', 'Search recipe steps'),
    ('perm_update_recipe_steps', 'update.recipe_steps', 'Update recipe steps'),
    ('perm_archive_recipe_steps', 'archive.recipe_steps', 'Archive recipe steps'),
    -- recipe prep tasks
    ('perm_create_recipe_prep_tasks', 'create.recipe_prep_tasks', 'Create recipe prep tasks'),
    ('perm_read_recipe_prep_tasks', 'read.recipe_prep_tasks', 'Read recipe prep tasks'),
    ('perm_update_recipe_prep_tasks', 'update.recipe_prep_tasks', 'Update recipe prep tasks'),
    ('perm_archive_recipe_prep_tasks', 'archive.recipe_prep_tasks', 'Archive recipe prep tasks'),
    -- recipe step instruments
    ('perm_create_recipe_step_instruments', 'create.recipe_step_instruments', 'Create recipe step instruments'),
    ('perm_read_recipe_step_instruments', 'read.recipe_step_instruments', 'Read recipe step instruments'),
    ('perm_search_recipe_step_instruments', 'search.recipe_step_instruments', 'Search recipe step instruments'),
    ('perm_update_recipe_step_instruments', 'update.recipe_step_instruments', 'Update recipe step instruments'),
    ('perm_archive_recipe_step_instruments', 'archive.recipe_step_instruments', 'Archive recipe step instruments'),
    -- recipe step vessels
    ('perm_create_recipe_step_vessels', 'create.recipe_step_vessels', 'Create recipe step vessels'),
    ('perm_read_recipe_step_vessels', 'read.recipe_step_vessels', 'Read recipe step vessels'),
    ('perm_search_recipe_step_vessels', 'search.recipe_step_vessels', 'Search recipe step vessels'),
    ('perm_update_recipe_step_vessels', 'update.recipe_step_vessels', 'Update recipe step vessels'),
    ('perm_archive_recipe_step_vessels', 'archive.recipe_step_vessels', 'Archive recipe step vessels'),
    -- recipe step ingredients
    ('perm_create_recipe_step_ingredients', 'create.recipe_step_ingredients', 'Create recipe step ingredients'),
    ('perm_read_recipe_step_ingredients', 'read.recipe_step_ingredients', 'Read recipe step ingredients'),
    ('perm_search_recipe_step_ingredients', 'search.recipe_step_ingredients', 'Search recipe step ingredients'),
    ('perm_update_recipe_step_ingredients', 'update.recipe_step_ingredients', 'Update recipe step ingredients'),
    ('perm_archive_recipe_step_ingredients', 'archive.recipe_step_ingredients', 'Archive recipe step ingredients'),
    -- recipe step completion conditions
    ('perm_create_recipe_step_completion_conditions', 'create.recipe_step_completion_conditions', 'Create recipe step completion conditions'),
    ('perm_read_recipe_step_completion_conditions', 'read.recipe_step_completion_conditions', 'Read recipe step completion conditions'),
    ('perm_search_recipe_step_completion_conditions', 'search.recipe_step_completion_conditions', 'Search recipe step completion conditions'),
    ('perm_update_recipe_step_completion_conditions', 'update.recipe_step_completion_conditions', 'Update recipe step completion conditions'),
    ('perm_archive_recipe_step_completion_conditions', 'archive.recipe_step_completion_conditions', 'Archive recipe step completion conditions'),
    -- recipe step products
    ('perm_create_recipe_step_products', 'create.recipe_step_products', 'Create recipe step products'),
    ('perm_read_recipe_step_products', 'read.recipe_step_products', 'Read recipe step products'),
    ('perm_search_recipe_step_products', 'search.recipe_step_products', 'Search recipe step products'),
    ('perm_update_recipe_step_products', 'update.recipe_step_products', 'Update recipe step products'),
    ('perm_archive_recipe_step_products', 'archive.recipe_step_products', 'Archive recipe step products'),
    -- meal plans
    ('perm_create_meal_plans', 'create.meal_plans', 'Create meal plans'),
    ('perm_read_meal_plans', 'read.meal_plans', 'Read meal plans'),
    ('perm_search_meal_plans', 'search.meal_plans', 'Search meal plans'),
    ('perm_update_meal_plans', 'update.meal_plans', 'Update meal plans'),
    ('perm_archive_meal_plans', 'archive.meal_plans', 'Archive meal plans'),
    -- meal plan events
    ('perm_create_meal_plan_events', 'create.meal_plan_events', 'Create meal plan events'),
    ('perm_read_meal_plan_events', 'read.meal_plan_events', 'Read meal plan events'),
    ('perm_update_meal_plan_events', 'update.meal_plan_events', 'Update meal plan events'),
    ('perm_archive_meal_plan_events', 'archive.meal_plan_events', 'Archive meal plan events'),
    -- meal plan options
    ('perm_create_meal_plan_options', 'create.meal_plan_options', 'Create meal plan options'),
    ('perm_read_meal_plan_options', 'read.meal_plan_options', 'Read meal plan options'),
    ('perm_search_meal_plan_options', 'search.meal_plan_options', 'Search meal plan options'),
    ('perm_update_meal_plan_options', 'update.meal_plan_options', 'Update meal plan options'),
    ('perm_archive_meal_plan_options', 'archive.meal_plan_options', 'Archive meal plan options'),
    -- meal plan tasks
    ('perm_create_meal_plan_tasks', 'create.meal_plan_tasks', 'Create meal plan tasks'),
    ('perm_read_meal_plan_tasks', 'read.meal_plan_tasks', 'Read meal plan tasks'),
    ('perm_update_meal_plan_tasks', 'update.meal_plan_tasks', 'Update meal plan tasks'),
    -- meal plan grocery list items
    ('perm_create_meal_plan_grocery_list_items', 'create.meal_plan_grocery_list_items', 'Create meal plan grocery list items'),
    ('perm_read_meal_plan_grocery_list_items', 'read.meal_plan_grocery_list_items', 'Read meal plan grocery list items'),
    ('perm_update_meal_plan_grocery_list_items', 'update.meal_plan_grocery_list_items', 'Update meal plan grocery list items'),
    ('perm_archive_meal_plan_grocery_list_items', 'archive.meal_plan_grocery_list_items', 'Archive meal plan grocery list items'),
    -- meal plan option votes
    ('perm_create_meal_plan_option_votes', 'create.meal_plan_option_votes', 'Create meal plan option votes'),
    ('perm_read_meal_plan_option_votes', 'read.meal_plan_option_votes', 'Read meal plan option votes'),
    ('perm_search_meal_plan_option_votes', 'search.meal_plan_option_votes', 'Search meal plan option votes'),
    ('perm_update_meal_plan_option_votes', 'update.meal_plan_option_votes', 'Update meal plan option votes'),
    ('perm_archive_meal_plan_option_votes', 'archive.meal_plan_option_votes', 'Archive meal plan option votes'),
    -- meal plan recipe option selections
    ('perm_create_meal_plan_recipe_option_selections', 'create.meal_plan_recipe_option_selections', 'Create meal plan recipe option selections'),
    ('perm_read_meal_plan_recipe_option_selections', 'read.meal_plan_recipe_option_selections', 'Read meal plan recipe option selections'),
    ('perm_update_meal_plan_recipe_option_selections', 'update.meal_plan_recipe_option_selections', 'Update meal plan recipe option selections'),
    ('perm_archive_meal_plan_recipe_option_selections', 'archive.meal_plan_recipe_option_selections', 'Archive meal plan recipe option selections'),
    -- user ingredient preferences
    ('perm_create_user_ingredient_preferences', 'create.user_ingredient_preferences', 'Create user ingredient preferences'),
    ('perm_read_user_ingredient_preferences', 'read.user_ingredient_preferences', 'Read user ingredient preferences'),
    ('perm_update_user_ingredient_preferences', 'update.user_ingredient_preferences', 'Update user ingredient preferences'),
    ('perm_archive_user_ingredient_preferences', 'archive.user_ingredient_preferences', 'Archive user ingredient preferences'),
    -- account instrument ownerships
    ('perm_create_account_instrument_ownerships', 'create.account_instrument_ownerships', 'Create account instrument ownerships'),
    ('perm_read_account_instrument_ownerships', 'read.account_instrument_ownerships', 'Read account instrument ownerships'),
    ('perm_update_account_instrument_ownerships', 'update.account_instrument_ownerships', 'Update account instrument ownerships'),
    ('perm_archive_account_instrument_ownerships', 'archive.account_instrument_ownerships', 'Archive account instrument ownerships'),
    -- recipe ratings
    ('perm_create_recipe_ratings', 'create.recipe_ratings', 'Create recipe ratings'),
    ('perm_read_recipe_ratings', 'read.recipe_ratings', 'Read recipe ratings'),
    ('perm_update_recipe_ratings', 'update.recipe_ratings', 'Update recipe ratings'),
    ('perm_archive_recipe_ratings', 'archive.recipe_ratings', 'Archive recipe ratings'),
    -- meal lists
    ('perm_create_meal_lists', 'create.meal_lists', 'Create meal lists'),
    ('perm_read_meal_lists', 'read.meal_lists', 'Read meal lists'),
    ('perm_update_meal_lists', 'update.meal_lists', 'Update meal lists'),
    ('perm_archive_meal_lists', 'archive.meal_lists', 'Archive meal lists'),
    -- recipe lists
    ('perm_create_recipe_lists', 'create.recipe_lists', 'Create recipe lists'),
    ('perm_read_recipe_lists', 'read.recipe_lists', 'Read recipe lists'),
    ('perm_update_recipe_lists', 'update.recipe_lists', 'Update recipe lists'),
    ('perm_archive_recipe_lists', 'archive.recipe_lists', 'Archive recipe lists');

-- =============================================================================
-- SEED DATA: mealplanning role-permission mappings
-- =============================================================================

-- service_admin: mealplanning permissions
INSERT INTO user_role_permissions (id, role_id, permission_id) VALUES
    ('urp_sa_13', 'role_service_admin', 'perm_update_recipe_status'),
    ('urp_sa_15', 'role_service_admin', 'perm_create_meal_plan_tasks'),
    ('urp_sa_16', 'role_service_admin', 'perm_create_meal_plan_grocery_list_items'),
    ('urp_sa_28', 'role_service_admin', 'perm_create_valid_instruments'),
    ('urp_sa_29', 'role_service_admin', 'perm_update_valid_instruments'),
    ('urp_sa_30', 'role_service_admin', 'perm_archive_valid_instruments'),
    ('urp_sa_31', 'role_service_admin', 'perm_create_valid_vessels'),
    ('urp_sa_32', 'role_service_admin', 'perm_update_valid_vessels'),
    ('urp_sa_33', 'role_service_admin', 'perm_archive_valid_vessels'),
    ('urp_sa_34', 'role_service_admin', 'perm_create_valid_ingredients'),
    ('urp_sa_35', 'role_service_admin', 'perm_update_valid_ingredients'),
    ('urp_sa_36', 'role_service_admin', 'perm_archive_valid_ingredients'),
    ('urp_sa_37', 'role_service_admin', 'perm_create_valid_ingredient_groups'),
    ('urp_sa_38', 'role_service_admin', 'perm_update_valid_ingredient_groups'),
    ('urp_sa_39', 'role_service_admin', 'perm_archive_valid_ingredient_groups'),
    ('urp_sa_40', 'role_service_admin', 'perm_create_valid_preparations'),
    ('urp_sa_41', 'role_service_admin', 'perm_update_valid_preparations'),
    ('urp_sa_42', 'role_service_admin', 'perm_archive_valid_preparations'),
    ('urp_sa_43', 'role_service_admin', 'perm_create_measurement_units'),
    ('urp_sa_44', 'role_service_admin', 'perm_update_measurement_units'),
    ('urp_sa_45', 'role_service_admin', 'perm_archive_measurement_units'),
    ('urp_sa_46', 'role_service_admin', 'perm_create_measurement_conversions'),
    ('urp_sa_47', 'role_service_admin', 'perm_update_measurement_conversions'),
    ('urp_sa_48', 'role_service_admin', 'perm_archive_measurement_conversions'),
    ('urp_sa_49', 'role_service_admin', 'perm_create_valid_ingredient_preparations'),
    ('urp_sa_50', 'role_service_admin', 'perm_update_valid_ingredient_preparations'),
    ('urp_sa_51', 'role_service_admin', 'perm_archive_valid_ingredient_preparations'),
    ('urp_sa_52', 'role_service_admin', 'perm_create_valid_prep_task_configs'),
    ('urp_sa_53', 'role_service_admin', 'perm_update_valid_prep_task_configs'),
    ('urp_sa_54', 'role_service_admin', 'perm_archive_valid_prep_task_configs'),
    ('urp_sa_55', 'role_service_admin', 'perm_create_valid_ingredient_state_ingredients'),
    ('urp_sa_56', 'role_service_admin', 'perm_update_valid_ingredient_state_ingredients'),
    ('urp_sa_57', 'role_service_admin', 'perm_archive_valid_ingredient_state_ingredients'),
    ('urp_sa_58', 'role_service_admin', 'perm_create_valid_preparation_instruments'),
    ('urp_sa_59', 'role_service_admin', 'perm_update_valid_preparation_instruments'),
    ('urp_sa_60', 'role_service_admin', 'perm_archive_valid_preparation_instruments'),
    ('urp_sa_61', 'role_service_admin', 'perm_create_valid_preparation_vessels'),
    ('urp_sa_62', 'role_service_admin', 'perm_update_valid_preparation_vessels'),
    ('urp_sa_63', 'role_service_admin', 'perm_archive_valid_preparation_vessels'),
    ('urp_sa_64', 'role_service_admin', 'perm_create_valid_ingredient_measurement_units'),
    ('urp_sa_65', 'role_service_admin', 'perm_update_valid_ingredient_measurement_units'),
    ('urp_sa_66', 'role_service_admin', 'perm_archive_valid_ingredient_measurement_units'),
    ('urp_sa_67', 'role_service_admin', 'perm_create_valid_ingredient_states'),
    ('urp_sa_68', 'role_service_admin', 'perm_update_valid_ingredient_states'),
    ('urp_sa_69', 'role_service_admin', 'perm_archive_valid_ingredient_states');

-- service_data_admin: all mealplanning valid_* permissions
INSERT INTO user_role_permissions (id, role_id, permission_id) VALUES
    ('urp_sda_1', 'role_service_data_admin', 'perm_create_valid_instruments'),
    ('urp_sda_2', 'role_service_data_admin', 'perm_update_valid_instruments'),
    ('urp_sda_3', 'role_service_data_admin', 'perm_archive_valid_instruments'),
    ('urp_sda_4', 'role_service_data_admin', 'perm_create_valid_vessels'),
    ('urp_sda_5', 'role_service_data_admin', 'perm_update_valid_vessels'),
    ('urp_sda_6', 'role_service_data_admin', 'perm_archive_valid_vessels'),
    ('urp_sda_7', 'role_service_data_admin', 'perm_create_valid_ingredients'),
    ('urp_sda_8', 'role_service_data_admin', 'perm_update_valid_ingredients'),
    ('urp_sda_9', 'role_service_data_admin', 'perm_archive_valid_ingredients'),
    ('urp_sda_10', 'role_service_data_admin', 'perm_create_valid_ingredient_groups'),
    ('urp_sda_11', 'role_service_data_admin', 'perm_update_valid_ingredient_groups'),
    ('urp_sda_12', 'role_service_data_admin', 'perm_archive_valid_ingredient_groups'),
    ('urp_sda_13', 'role_service_data_admin', 'perm_create_valid_preparations'),
    ('urp_sda_14', 'role_service_data_admin', 'perm_update_valid_preparations'),
    ('urp_sda_15', 'role_service_data_admin', 'perm_archive_valid_preparations'),
    ('urp_sda_16', 'role_service_data_admin', 'perm_create_measurement_units'),
    ('urp_sda_17', 'role_service_data_admin', 'perm_update_measurement_units'),
    ('urp_sda_18', 'role_service_data_admin', 'perm_archive_measurement_units'),
    ('urp_sda_19', 'role_service_data_admin', 'perm_create_measurement_conversions'),
    ('urp_sda_20', 'role_service_data_admin', 'perm_update_measurement_conversions'),
    ('urp_sda_21', 'role_service_data_admin', 'perm_archive_measurement_conversions'),
    ('urp_sda_22', 'role_service_data_admin', 'perm_create_valid_ingredient_preparations'),
    ('urp_sda_23', 'role_service_data_admin', 'perm_update_valid_ingredient_preparations'),
    ('urp_sda_24', 'role_service_data_admin', 'perm_archive_valid_ingredient_preparations'),
    ('urp_sda_25', 'role_service_data_admin', 'perm_create_valid_prep_task_configs'),
    ('urp_sda_26', 'role_service_data_admin', 'perm_update_valid_prep_task_configs'),
    ('urp_sda_27', 'role_service_data_admin', 'perm_archive_valid_prep_task_configs'),
    ('urp_sda_28', 'role_service_data_admin', 'perm_create_valid_ingredient_state_ingredients'),
    ('urp_sda_29', 'role_service_data_admin', 'perm_update_valid_ingredient_state_ingredients'),
    ('urp_sda_30', 'role_service_data_admin', 'perm_archive_valid_ingredient_state_ingredients'),
    ('urp_sda_31', 'role_service_data_admin', 'perm_create_valid_preparation_instruments'),
    ('urp_sda_32', 'role_service_data_admin', 'perm_update_valid_preparation_instruments'),
    ('urp_sda_33', 'role_service_data_admin', 'perm_archive_valid_preparation_instruments'),
    ('urp_sda_34', 'role_service_data_admin', 'perm_create_valid_preparation_vessels'),
    ('urp_sda_35', 'role_service_data_admin', 'perm_update_valid_preparation_vessels'),
    ('urp_sda_36', 'role_service_data_admin', 'perm_archive_valid_preparation_vessels'),
    ('urp_sda_37', 'role_service_data_admin', 'perm_create_valid_ingredient_measurement_units'),
    ('urp_sda_38', 'role_service_data_admin', 'perm_update_valid_ingredient_measurement_units'),
    ('urp_sda_39', 'role_service_data_admin', 'perm_archive_valid_ingredient_measurement_units'),
    ('urp_sda_40', 'role_service_data_admin', 'perm_create_valid_ingredient_states'),
    ('urp_sda_41', 'role_service_data_admin', 'perm_update_valid_ingredient_states'),
    ('urp_sda_42', 'role_service_data_admin', 'perm_archive_valid_ingredient_states');

-- account_admin: mealplanning permissions
INSERT INTO user_role_permissions (id, role_id, permission_id) VALUES
    ('urp_aa_13', 'role_account_admin', 'perm_create_meal_plans'),
    ('urp_aa_14', 'role_account_admin', 'perm_update_meal_plans'),
    ('urp_aa_15', 'role_account_admin', 'perm_archive_meal_plans'),
    ('urp_aa_16', 'role_account_admin', 'perm_create_meal_plan_events'),
    ('urp_aa_17', 'role_account_admin', 'perm_update_meal_plan_events'),
    ('urp_aa_18', 'role_account_admin', 'perm_archive_meal_plan_events'),
    ('urp_aa_19', 'role_account_admin', 'perm_create_meal_plan_options'),
    ('urp_aa_20', 'role_account_admin', 'perm_update_meal_plan_options'),
    ('urp_aa_21', 'role_account_admin', 'perm_archive_meal_plan_options'),
    ('urp_aa_22', 'role_account_admin', 'perm_create_account_instrument_ownerships'),
    ('urp_aa_23', 'role_account_admin', 'perm_update_account_instrument_ownerships'),
    ('urp_aa_24', 'role_account_admin', 'perm_archive_account_instrument_ownerships'),
    ('urp_aa_31', 'role_account_admin', 'perm_create_meal_lists'),
    ('urp_aa_32', 'role_account_admin', 'perm_read_meal_lists'),
    ('urp_aa_33', 'role_account_admin', 'perm_update_meal_lists'),
    ('urp_aa_34', 'role_account_admin', 'perm_archive_meal_lists'),
    ('urp_aa_35', 'role_account_admin', 'perm_create_recipe_lists'),
    ('urp_aa_36', 'role_account_admin', 'perm_read_recipe_lists'),
    ('urp_aa_37', 'role_account_admin', 'perm_update_recipe_lists'),
    ('urp_aa_38', 'role_account_admin', 'perm_archive_recipe_lists');

-- account_member: mealplanning permissions
INSERT INTO user_role_permissions (id, role_id, permission_id) VALUES
    ('urp_am_12', 'role_account_member', 'perm_create_meals'),
    ('urp_am_13', 'role_account_member', 'perm_read_meals'),
    ('urp_am_14', 'role_account_member', 'perm_update_meals'),
    ('urp_am_15', 'role_account_member', 'perm_archive_meals'),
    ('urp_am_16', 'role_account_member', 'perm_create_recipes'),
    ('urp_am_17', 'role_account_member', 'perm_read_recipes'),
    ('urp_am_18', 'role_account_member', 'perm_search_recipes'),
    ('urp_am_19', 'role_account_member', 'perm_update_recipes'),
    ('urp_am_20', 'role_account_member', 'perm_archive_recipes'),
    ('urp_am_21', 'role_account_member', 'perm_create_recipe_steps'),
    ('urp_am_22', 'role_account_member', 'perm_read_recipe_steps'),
    ('urp_am_23', 'role_account_member', 'perm_search_recipe_steps'),
    ('urp_am_24', 'role_account_member', 'perm_update_recipe_steps'),
    ('urp_am_25', 'role_account_member', 'perm_archive_recipe_steps'),
    ('urp_am_26', 'role_account_member', 'perm_create_recipe_prep_tasks'),
    ('urp_am_27', 'role_account_member', 'perm_read_recipe_prep_tasks'),
    ('urp_am_28', 'role_account_member', 'perm_update_recipe_prep_tasks'),
    ('urp_am_29', 'role_account_member', 'perm_archive_recipe_prep_tasks'),
    ('urp_am_30', 'role_account_member', 'perm_create_recipe_step_instruments'),
    ('urp_am_31', 'role_account_member', 'perm_read_recipe_step_instruments'),
    ('urp_am_32', 'role_account_member', 'perm_search_recipe_step_instruments'),
    ('urp_am_33', 'role_account_member', 'perm_update_recipe_step_instruments'),
    ('urp_am_34', 'role_account_member', 'perm_archive_recipe_step_instruments'),
    ('urp_am_35', 'role_account_member', 'perm_create_recipe_step_vessels'),
    ('urp_am_36', 'role_account_member', 'perm_read_recipe_step_vessels'),
    ('urp_am_37', 'role_account_member', 'perm_search_recipe_step_vessels'),
    ('urp_am_38', 'role_account_member', 'perm_update_recipe_step_vessels'),
    ('urp_am_39', 'role_account_member', 'perm_archive_recipe_step_vessels'),
    ('urp_am_40', 'role_account_member', 'perm_create_recipe_step_ingredients'),
    ('urp_am_41', 'role_account_member', 'perm_read_recipe_step_ingredients'),
    ('urp_am_42', 'role_account_member', 'perm_search_recipe_step_ingredients'),
    ('urp_am_43', 'role_account_member', 'perm_update_recipe_step_ingredients'),
    ('urp_am_44', 'role_account_member', 'perm_archive_recipe_step_ingredients'),
    ('urp_am_45', 'role_account_member', 'perm_create_recipe_step_completion_conditions'),
    ('urp_am_46', 'role_account_member', 'perm_read_recipe_step_completion_conditions'),
    ('urp_am_47', 'role_account_member', 'perm_search_recipe_step_completion_conditions'),
    ('urp_am_48', 'role_account_member', 'perm_update_recipe_step_completion_conditions'),
    ('urp_am_49', 'role_account_member', 'perm_archive_recipe_step_completion_conditions'),
    ('urp_am_50', 'role_account_member', 'perm_create_recipe_step_products'),
    ('urp_am_51', 'role_account_member', 'perm_read_recipe_step_products'),
    ('urp_am_52', 'role_account_member', 'perm_search_recipe_step_products'),
    ('urp_am_53', 'role_account_member', 'perm_update_recipe_step_products'),
    ('urp_am_54', 'role_account_member', 'perm_archive_recipe_step_products'),
    ('urp_am_55', 'role_account_member', 'perm_read_valid_instruments'),
    ('urp_am_56', 'role_account_member', 'perm_search_valid_instruments'),
    ('urp_am_57', 'role_account_member', 'perm_read_valid_vessels'),
    ('urp_am_58', 'role_account_member', 'perm_search_valid_vessels'),
    ('urp_am_59', 'role_account_member', 'perm_read_valid_ingredients'),
    ('urp_am_60', 'role_account_member', 'perm_search_valid_ingredients'),
    ('urp_am_61', 'role_account_member', 'perm_read_valid_ingredient_groups'),
    ('urp_am_62', 'role_account_member', 'perm_search_valid_ingredient_groups'),
    ('urp_am_63', 'role_account_member', 'perm_read_valid_preparations'),
    ('urp_am_64', 'role_account_member', 'perm_search_valid_preparations'),
    ('urp_am_65', 'role_account_member', 'perm_read_measurement_units'),
    ('urp_am_66', 'role_account_member', 'perm_search_measurement_units'),
    ('urp_am_67', 'role_account_member', 'perm_read_measurement_conversions'),
    ('urp_am_68', 'role_account_member', 'perm_read_valid_ingredient_preparations'),
    ('urp_am_69', 'role_account_member', 'perm_search_valid_ingredient_preparations'),
    ('urp_am_70', 'role_account_member', 'perm_read_valid_ingredient_state_ingredients'),
    ('urp_am_71', 'role_account_member', 'perm_search_valid_ingredient_state_ingredients'),
    ('urp_am_72', 'role_account_member', 'perm_read_valid_preparation_instruments'),
    ('urp_am_73', 'role_account_member', 'perm_search_valid_preparation_instruments'),
    ('urp_am_74', 'role_account_member', 'perm_read_valid_preparation_vessels'),
    ('urp_am_75', 'role_account_member', 'perm_search_valid_preparation_vessels'),
    ('urp_am_76', 'role_account_member', 'perm_read_valid_ingredient_measurement_units'),
    ('urp_am_77', 'role_account_member', 'perm_search_valid_ingredient_measurement_units'),
    ('urp_am_78', 'role_account_member', 'perm_read_meal_plans'),
    ('urp_am_79', 'role_account_member', 'perm_search_meal_plans'),
    ('urp_am_80', 'role_account_member', 'perm_read_meal_plan_events'),
    ('urp_am_81', 'role_account_member', 'perm_read_meal_plan_options'),
    ('urp_am_82', 'role_account_member', 'perm_search_meal_plan_options'),
    ('urp_am_83', 'role_account_member', 'perm_read_valid_ingredient_states'),
    ('urp_am_84', 'role_account_member', 'perm_read_meal_plan_grocery_list_items'),
    ('urp_am_85', 'role_account_member', 'perm_update_meal_plan_grocery_list_items'),
    ('urp_am_86', 'role_account_member', 'perm_archive_meal_plan_grocery_list_items'),
    ('urp_am_87', 'role_account_member', 'perm_create_meal_plan_option_votes'),
    ('urp_am_88', 'role_account_member', 'perm_read_meal_plan_option_votes'),
    ('urp_am_89', 'role_account_member', 'perm_search_meal_plan_option_votes'),
    ('urp_am_90', 'role_account_member', 'perm_update_meal_plan_option_votes'),
    ('urp_am_91', 'role_account_member', 'perm_archive_meal_plan_option_votes'),
    ('urp_am_92', 'role_account_member', 'perm_create_meal_plan_recipe_option_selections'),
    ('urp_am_93', 'role_account_member', 'perm_read_meal_plan_recipe_option_selections'),
    ('urp_am_94', 'role_account_member', 'perm_update_meal_plan_recipe_option_selections'),
    ('urp_am_95', 'role_account_member', 'perm_archive_meal_plan_recipe_option_selections'),
    ('urp_am_100', 'role_account_member', 'perm_read_meal_plan_tasks'),
    ('urp_am_101', 'role_account_member', 'perm_update_meal_plan_tasks'),
    ('urp_am_102', 'role_account_member', 'perm_create_user_ingredient_preferences'),
    ('urp_am_103', 'role_account_member', 'perm_read_user_ingredient_preferences'),
    ('urp_am_104', 'role_account_member', 'perm_update_user_ingredient_preferences'),
    ('urp_am_105', 'role_account_member', 'perm_archive_user_ingredient_preferences'),
    ('urp_am_106', 'role_account_member', 'perm_read_account_instrument_ownerships'),
    ('urp_am_107', 'role_account_member', 'perm_create_recipe_ratings'),
    ('urp_am_108', 'role_account_member', 'perm_read_recipe_ratings'),
    ('urp_am_113', 'role_account_member', 'perm_update_recipe_ratings'),
    ('urp_am_114', 'role_account_member', 'perm_archive_recipe_ratings'),
    ('urp_am_125', 'role_account_member', 'perm_read_valid_prep_task_configs');
