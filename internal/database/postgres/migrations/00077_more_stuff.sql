ALTER TABLE recipe_prep_tasks ADD COLUMN "name" TEXT NOT NULL DEFAULT '';
ALTER TABLE recipe_prep_tasks ADD COLUMN "description" TEXT NOT NULL DEFAULT '';
ALTER TABLE meal_plans ADD COLUMN "created_by_user" TEXT REFERENCES users("id");
ALTER TABLE users ADD COLUMN "first_name" TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN "last_name" TEXT NOT NULL DEFAULT '';
ALTER TABLE household_invitations ADD COLUMN "addressed_to_name" TEXT NOT NULL DEFAULT '';
