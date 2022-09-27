ALTER TABLE advanced_prep_steps RENAME TO meal_plan_tasks;

ALTER TABLE meal_plan_tasks ADD COLUMN "assigned_to_user" CHAR(27) REFERENCES users("id") ON DELETE CASCADE;
