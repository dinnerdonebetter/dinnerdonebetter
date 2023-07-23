ALTER TABLE meal_components ADD COLUMN "recipe_scale" NUMERIC(14, 2) NOT NULL DEFAULT 1.0;
ALTER TABLE meal_plan_options ADD COLUMN "meal_scale" NUMERIC(14, 2) NOT NULL DEFAULT 1.0;
ALTER TABLE meals ADD COLUMN "min_estimated_portions" NUMERIC(14, 2) NOT NULL DEFAULT 1.0;
ALTER TABLE meals ADD COLUMN "max_estimated_portions" NUMERIC(14, 2);
