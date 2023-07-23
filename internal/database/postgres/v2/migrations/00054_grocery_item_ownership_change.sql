ALTER TABLE meal_plan_grocery_list_items DROP COLUMN "belongs_to_meal_plan_option";
ALTER TABLE meal_plan_grocery_list_items ADD COLUMN  "belongs_to_meal_plan" CHAR(27) NOT NULL REFERENCES meal_plans("id") ON DELETE CASCADE;
