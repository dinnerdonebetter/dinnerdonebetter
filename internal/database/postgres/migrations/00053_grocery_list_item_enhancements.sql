ALTER TABLE meal_plans ADD COLUMN "grocery_list_initialized" BOOLEAN NOT NULL DEFAULT 'false';
ALTER TABLE meal_plans ADD COLUMN "tasks_created" BOOLEAN NOT NULL DEFAULT 'false';
ALTER TABLE meal_plan_options DROP COLUMN "prep_steps_created";
