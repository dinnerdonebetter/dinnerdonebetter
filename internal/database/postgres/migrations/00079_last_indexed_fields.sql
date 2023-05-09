ALTER TABLE recipes ADD COLUMN "last_indexed_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL;
ALTER TABLE meals ADD COLUMN "last_indexed_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL;
ALTER TABLE valid_ingredients ADD COLUMN "last_indexed_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL;
ALTER TABLE valid_instruments ADD COLUMN "last_indexed_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL;
ALTER TABLE valid_measurement_units ADD COLUMN "last_indexed_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL;
ALTER TABLE valid_preparations ADD COLUMN "last_indexed_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL;

ALTER TABLE recipe_steps ADD COLUMN "time_scale_factor" NUMERIC(10, 5) NOT NULL DEFAULT 1.0;
ALTER TABLE recipe_step_products ADD COLUMN "quantity_scale_factor" NUMERIC(10, 5) NOT NULL DEFAULT 1.0;
-- ALTER TABLE recipe_step_ingredients ADD COLUMN "quantity_scale_factor" NUMERIC(10, 5) NOT NULL DEFAULT 1.0;
-- ALTER TABLE recipe_step_ingredients DROP COLUMN IF EXISTS "requires_defrost";
-- ALTER TABLE recipe_step_vessels ADD COLUMN "quantity_scale_factor" NUMERIC(10, 5) NOT NULL DEFAULT 1.0;
-- ALTER TABLE recipe_step_instruments ADD COLUMN "quantity_scale_factor" NUMERIC(10, 5) NOT NULL DEFAULT 1.0;
