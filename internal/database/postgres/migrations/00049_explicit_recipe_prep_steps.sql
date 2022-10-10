DROP TABLE advanced_prep_notifications;

CREATE TYPE storage_container_type AS ENUM ('uncovered', 'covered', 'on a wire rack', 'in an airtight container');

CREATE TABLE IF NOT EXISTS recipe_prep_tasks (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "notes" TEXT NOT NULL DEFAULT '',
    "explicit_storage_instructions" TEXT NOT NULL DEFAULT '',
    "minimum_time_buffer_before_recipe_in_seconds" INTEGER NOT NULL,
    "maximum_time_buffer_before_recipe_in_seconds" INTEGER,
    "storage_type" storage_container_type,
    "minimum_storage_temperature_in_celsius" INTEGER,
    "maximum_storage_temperature_in_celsius" INTEGER,
    "belongs_to_recipe" CHAR(27) NOT NULL REFERENCES recipes("id") ON DELETE CASCADE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE,
    "archived_at" TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS recipe_prep_task_steps (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "satisfies_recipe_step" BOOLEAN NOT NULL DEFAULT 'false',
    "belongs_to_recipe_step" CHAR(27) NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    "belongs_to_recipe_prep_task" CHAR(27) NOT NULL REFERENCES recipe_prep_tasks("id") ON DELETE CASCADE
);

-- ALTER TABLE recipe_step_products DROP COLUMN minimum_storage_temperature_in_celsius;
-- ALTER TABLE recipe_step_products DROP COLUMN maximum_storage_temperature_in_celsius;
-- ALTER TABLE recipe_step_products DROP COLUMN maximum_storage_duration_in_seconds;
-- ALTER TABLE recipe_step_products DROP COLUMN storage_instructions;
