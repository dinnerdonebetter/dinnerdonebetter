ALTER TABLE advanced_prep_steps RENAME TO meal_plan_tasks;

ALTER TYPE prep_step_status ADD VALUE 'ignored';

ALTER TABLE meal_plan_tasks DROP COLUMN "belongs_to_meal_plan_option";
ALTER TABLE meal_plan_tasks DROP COLUMN "satisfies_recipe_step";
ALTER TABLE meal_plan_tasks ADD COLUMN "assigned_to_user" CHAR(27) REFERENCES users("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS meal_plan_task_recipe_steps (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "belongs_to_meal_plan_task" CHAR(27) NOT NULL REFERENCES meal_plan_tasks("id") ON DELETE CASCADE,
    "satisfies_recipe_step" CHAR(27) NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "archived_at" TIMESTAMP WITH TIME ZONE
);
