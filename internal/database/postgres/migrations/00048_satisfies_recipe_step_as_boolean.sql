ALTER TABLE meal_plan_task_recipe_steps RENAME COLUMN satisfies_recipe_step TO applies_to_recipe_step;
ALTER TABLE meal_plan_task_recipe_steps ADD COLUMN "satisfies_recipe_step" BOOLEAN NOT NULL DEFAULT FALSE;
