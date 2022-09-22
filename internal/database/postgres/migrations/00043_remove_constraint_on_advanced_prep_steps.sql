ALTER TABLE advanced_prep_steps DROP CONSTRAINT recipe_step_and_meal_plan_option;

ALTER TABLE advanced_prep_steps ADD COLUMN "storage_instructions" TEXT NOT NULL DEFAULT '';
