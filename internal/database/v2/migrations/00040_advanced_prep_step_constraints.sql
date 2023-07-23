CREATE TYPE prep_step_status AS ENUM ('unfinished', 'postponed', 'canceled', 'finished');

ALTER TABLE advanced_prep_steps RENAME COLUMN completed_at TO settled_at;
ALTER TABLE advanced_prep_steps ADD COLUMN "status" prep_step_status NOT NULL DEFAULT 'unfinished';
ALTER TABLE advanced_prep_steps ADD COLUMN "creation_explanation" TEXT NOT NULL DEFAULT '';
ALTER TABLE advanced_prep_steps ADD COLUMN "status_explanation" TEXT NOT NULL DEFAULT '';
ALTER TABLE advanced_prep_steps ADD CONSTRAINT recipe_step_and_meal_plan_option unique(belongs_to_meal_plan_option, satisfies_recipe_step, status);

ALTER TABLE valid_ingredients ALTER COLUMN minimum_ideal_storage_temperature_in_celsius DROP DEFAULT;
ALTER TABLE valid_ingredients ALTER COLUMN minimum_ideal_storage_temperature_in_celsius DROP NOT NULL;
ALTER TABLE valid_ingredients ALTER COLUMN maximum_ideal_storage_temperature_in_celsius DROP DEFAULT;
ALTER TABLE valid_ingredients ALTER COLUMN maximum_ideal_storage_temperature_in_celsius DROP NOT NULL;
UPDATE valid_ingredients SET minimum_ideal_storage_temperature_in_celsius = NULL where minimum_ideal_storage_temperature_in_celsius = 0;
UPDATE valid_ingredients SET maximum_ideal_storage_temperature_in_celsius = NULL where maximum_ideal_storage_temperature_in_celsius = 0;
