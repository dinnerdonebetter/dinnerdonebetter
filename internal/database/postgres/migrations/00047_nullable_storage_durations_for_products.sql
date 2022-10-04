ALTER TABLE recipe_step_products ALTER COLUMN maximum_storage_duration_in_seconds DROP DEFAULT;
ALTER TABLE recipe_step_products ALTER COLUMN maximum_storage_duration_in_seconds DROP NOT NULL;
UPDATE recipe_step_products SET maximum_storage_duration_in_seconds = NULL where maximum_storage_duration_in_seconds = 0;
