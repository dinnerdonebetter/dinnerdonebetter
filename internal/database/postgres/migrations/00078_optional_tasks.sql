ALTER TABLE recipe_prep_tasks ADD COLUMN "optional" BOOLEAN NOT NULL DEFAULT 'true';
ALTER TABLE meal_plans ALTER COLUMN "created_by_user" SET NOT NULL;
