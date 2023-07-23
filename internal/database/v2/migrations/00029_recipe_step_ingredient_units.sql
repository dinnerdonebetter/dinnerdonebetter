-- ALTER TABLE recipe_step_ingredients ADD COLUMN "quantity_type" CHAR(27) REFERENCES valid_measurement_units("id") ON DELETE RESTRICT;

ALTER TABLE recipe_step_ingredients DROP COLUMN "quantity_type";
-- ALTER TABLE recipe_step_products DROP COLUMN "quantity_type";
-- ALTER TABLE recipe_step_products DROP COLUMN "quantity_value";

ALTER TABLE recipe_step_products ADD COLUMN "minimum_quantity_value" DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE recipe_step_products ADD COLUMN "maximum_quantity_value" DOUBLE PRECISION NOT NULL DEFAULT 0; -- NOT NULL ON DELETE RESTRICT;

ALTER TABLE recipe_step_ingredients ADD COLUMN "measurement_unit" CHAR(27) REFERENCES valid_measurement_units("id"); -- NOT NULL ON DELETE RESTRICT;
ALTER TABLE recipe_step_products ADD COLUMN "measurement_unit" CHAR(27) REFERENCES valid_measurement_units("id"); -- NOT NULL ON DELETE RESTRICT;
