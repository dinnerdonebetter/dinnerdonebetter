ALTER TABLE households ADD COLUMN "last_payment_provider_sync_occurred_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL;

ALTER TABLE recipes ADD COLUMN "last_validated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL;
ALTER TABLE recipes ADD COLUMN "yields_component_type" component_type NOT NULL DEFAULT 'unspecified';

ALTER TABLE users ADD COLUMN "last_accepted_terms_of_service" TIMESTAMP WITH TIME ZONE DEFAULT NULL;
ALTER TABLE users ADD COLUMN "last_accepted_privacy_policy" TIMESTAMP WITH TIME ZONE DEFAULT NULL;

ALTER TABLE recipe_step_ingredients ADD COLUMN "recipe_step_product_recipe_id" TEXT DEFAULT NULL REFERENCES recipes("id") ON DELETE CASCADE ON UPDATE CASCADE;

CREATE TABLE IF NOT EXISTS recipe_ratings (
    "id" TEXT NOT NULL PRIMARY KEY,
    "recipe_id" TEXT NOT NULL REFERENCES recipes("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "taste"  NUMERIC(14, 2),
    "difficulty"  NUMERIC(14, 2),
    "cleanup"  NUMERIC(14, 2),
    "instructions"  NUMERIC(14, 2),
    "overall"  NUMERIC(14, 2),
    "notes" TEXT NOT NULL DEFAULT '',
    "by_user" TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    UNIQUE("by_user", "recipe_id")
);

CREATE TABLE IF NOT EXISTS household_instrument_ownerships (
    "id" TEXT NOT NULL PRIMARY KEY,
    "notes" TEXT NOT NULL DEFAULT '',
    "quantity" INTEGER NOT NULL DEFAULT 0,
    "valid_instrument_id" TEXT NOT NULL REFERENCES valid_instruments(id) ON DELETE CASCADE,
    "belongs_to_household" TEXT NOT NULL REFERENCES households(id) ON DELETE CASCADE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "last_updated_at" TIMESTAMP WITH TIME ZONE,
    "archived_at" TIMESTAMP WITH TIME ZONE,
    UNIQUE ("valid_instrument_id", "belongs_to_household", "archived_at")
);
