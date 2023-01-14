ALTER TABLE valid_instruments ADD COLUMN "is_vessel" BOOLEAN NOT NULL DEFAULT 'false';
ALTER TABLE valid_instruments ADD COLUMN "is_exclusively_vessel" BOOLEAN NOT NULL DEFAULT 'false';

ALTER TABLE valid_preparations ADD COLUMN "consumes_vessel" BOOLEAN NOT NULL DEFAULT 'false';
ALTER TABLE valid_preparations ADD COLUMN "only_for_vessels" BOOLEAN NOT NULL DEFAULT 'false';
ALTER TABLE valid_preparations ADD COLUMN "minimum_vessel_count" INTEGER NOT NULL DEFAULT 0;
ALTER TABLE valid_preparations ADD COLUMN "maximum_vessel_count" INTEGER NOT NULL DEFAULT 0;

ALTER TYPE recipe_step_product_type ADD VALUE 'vessel';

CREATE TABLE IF NOT EXISTS recipe_step_vessels (
    "id" TEXT NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL DEFAULT '',
    "notes" TEXT NOT NULL DEFAULT '',
    "belongs_to_recipe_step" TEXT NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "recipe_step_product_id" TEXT REFERENCES recipe_step_products("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "valid_instrument_id" TEXT REFERENCES valid_instruments("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "vessel_predicate" TEXT NOT NULL DEFAULT '',
    "minimum_quantity" INTEGER NOT NULL DEFAULT 0,
    "maximum_quantity" INTEGER NOT NULL DEFAULT 0,
    "unavailable_after_step" BOOLEAN NOT NULL DEFAULT 'false',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

ALTER TABLE recipe_step_products ADD COLUMN "index" INTEGER NOT NULL DEFAULT 0;
ALTER TABLE recipe_step_products ADD COLUMN "contained_by_vessel_index" INTEGER NOT NULL DEFAULT 0;
ALTER TABLE recipe_step_ingredients ADD COLUMN "vessel_index" INTEGER NOT NULL DEFAULT 0;
