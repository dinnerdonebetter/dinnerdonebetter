ALTER TABLE valid_measurement_units ALTER COLUMN volumetric DROP NOT NULL;
ALTER TABLE valid_ingredients ALTER COLUMN is_liquid DROP NOT NULL;
ALTER TABLE meal_plan_options ALTER COLUMN belongs_to_meal_plan_event DROP NOT NULL;
ALTER TABLE recipe_prep_tasks ALTER COLUMN storage_type DROP NOT NULL;
