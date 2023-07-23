CREATE TYPE grocery_list_item_status AS ENUM ('unknown', 'already owned', 'needs', 'unavailable', 'acquired');

CREATE TABLE IF NOT EXISTS meal_plan_grocery_list_items (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "belongs_to_meal_plan_option" CHAR(27) NOT NULL REFERENCES meal_plan_options("id") ON DELETE CASCADE,
    "valid_ingredient" CHAR(27) NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    "valid_measurement_unit" CHAR(27) NOT NULL REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    "minimum_quantity_needed" INTEGER NOT NULL,
    "maximum_quantity_needed" INTEGER NOT NULL,
    "quantity_purchased" INTEGER,
    "purchased_measurement_unit" CHAR(27) REFERENCES valid_measurement_units("id") ON DELETE CASCADE,
    "purchased_upc" TEXT,
    "purchase_price" INTEGER,
    "status_explanation" TEXT NOT NULL DEFAULT '',
    "status" grocery_list_item_status NOT NULL DEFAULT 'unknown',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "completed_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
);