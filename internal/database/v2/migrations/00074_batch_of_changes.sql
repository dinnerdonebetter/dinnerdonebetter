ALTER TABLE valid_ingredients ADD COLUMN "is_starch" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE valid_ingredients ADD COLUMN "is_protein" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE valid_ingredients ADD COLUMN "is_grain" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE valid_ingredients ADD COLUMN "is_fruit" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE valid_ingredients ADD COLUMN "is_salt" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE valid_ingredients ADD COLUMN "is_fat" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE valid_ingredients ADD COLUMN "is_acid" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE valid_ingredients ADD COLUMN "is_heat" BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE recipes ADD COLUMN "eligible_for_meals" BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE meals ADD COLUMN "eligible_for_meal_plans" BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE recipe_steps ADD COLUMN "start_timer_automatically" BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE recipe_prep_tasks ALTER COLUMN minimum_storage_temperature_in_celsius DROP DEFAULT;
ALTER TABLE recipe_prep_tasks ALTER COLUMN minimum_storage_temperature_in_celsius DROP NOT NULL;

ALTER TABLE recipe_prep_tasks ALTER COLUMN maximum_storage_temperature_in_celsius DROP DEFAULT;
ALTER TABLE recipe_prep_tasks ALTER COLUMN maximum_storage_temperature_in_celsius DROP NOT NULL;

ALTER TABLE recipe_prep_tasks ALTER COLUMN maximum_time_buffer_before_recipe_in_seconds DROP DEFAULT;
ALTER TABLE recipe_prep_tasks ALTER COLUMN maximum_time_buffer_before_recipe_in_seconds DROP NOT NULL;