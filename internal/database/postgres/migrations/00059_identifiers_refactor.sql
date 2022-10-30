ALTER TABLE users ALTER COLUMN id TYPE TEXT;
ALTER TABLE households ALTER COLUMN id TYPE TEXT;
ALTER TABLE households ALTER COLUMN belongs_to_user TYPE TEXT;
ALTER TABLE household_invitations ALTER COLUMN id TYPE TEXT;
ALTER TABLE household_invitations ALTER COLUMN destination_household TYPE TEXT;
ALTER TABLE household_invitations ALTER COLUMN to_user TYPE TEXT;
ALTER TABLE household_invitations ALTER COLUMN from_user TYPE TEXT;
ALTER TABLE household_user_memberships ALTER COLUMN id TYPE TEXT;
ALTER TABLE household_user_memberships ALTER COLUMN belongs_to_household TYPE TEXT;
ALTER TABLE household_user_memberships ALTER COLUMN belongs_to_user TYPE TEXT;
ALTER TABLE api_clients ALTER COLUMN id TYPE TEXT;
ALTER TABLE api_clients ALTER COLUMN belongs_to_user TYPE TEXT;
ALTER TABLE webhooks ALTER COLUMN id TYPE TEXT;
ALTER TABLE webhooks ALTER COLUMN belongs_to_household TYPE TEXT;
ALTER TABLE webhook_trigger_events ALTER COLUMN id TYPE TEXT;
ALTER TABLE webhook_trigger_events ALTER COLUMN belongs_to_webhook TYPE TEXT;

ALTER TABLE valid_instruments ALTER COLUMN id TYPE TEXT;
ALTER TABLE valid_ingredients ALTER COLUMN id TYPE TEXT;
ALTER TABLE valid_preparations ALTER COLUMN id TYPE TEXT;
ALTER TABLE valid_ingredient_preparations ALTER COLUMN id TYPE TEXT;
ALTER TABLE valid_ingredient_preparations ALTER COLUMN valid_preparation_id TYPE TEXT;
ALTER TABLE valid_ingredient_preparations ALTER COLUMN valid_ingredient_id TYPE TEXT;

ALTER TABLE recipes ALTER COLUMN id TYPE TEXT;
ALTER TABLE recipes ALTER COLUMN created_by_user TYPE TEXT;

ALTER TABLE recipes ALTER COLUMN id TYPE TEXT;

ALTER TABLE recipe_steps ALTER COLUMN id TYPE TEXT;
ALTER TABLE recipe_steps ALTER COLUMN preparation_id TYPE TEXT;
ALTER TABLE recipe_steps ALTER COLUMN belongs_to_recipe TYPE TEXT;

ALTER TABLE recipe_media ALTER COLUMN id TYPE TEXT;
ALTER TABLE recipe_media ALTER COLUMN belongs_to_recipe TYPE TEXT;
ALTER TABLE recipe_media ALTER COLUMN belongs_to_recipe_step TYPE TEXT;

ALTER TABLE recipe_step_ingredients ALTER COLUMN id TYPE TEXT;
ALTER TABLE recipe_step_ingredients ALTER COLUMN belongs_to_recipe_step TYPE TEXT;
ALTER TABLE recipe_step_ingredients ALTER COLUMN recipe_step_product_id TYPE TEXT;
ALTER TABLE recipe_step_ingredients ALTER COLUMN measurement_unit TYPE TEXT;
ALTER TABLE recipe_step_ingredients ALTER COLUMN ingredient_id TYPE TEXT;

ALTER TABLE recipe_step_instruments ALTER COLUMN id TYPE TEXT;
ALTER TABLE recipe_step_instruments ALTER COLUMN belongs_to_recipe_step TYPE TEXT;
ALTER TABLE recipe_step_instruments ALTER COLUMN recipe_step_product_id TYPE TEXT;
ALTER TABLE recipe_step_instruments ALTER COLUMN instrument_id TYPE TEXT;

ALTER TABLE recipe_step_products ALTER COLUMN id TYPE TEXT;
ALTER TABLE recipe_step_products ALTER COLUMN belongs_to_recipe_step TYPE TEXT;
ALTER TABLE recipe_step_products ALTER COLUMN measurement_unit TYPE TEXT;

ALTER TABLE meals ALTER COLUMN id TYPE TEXT;
ALTER TABLE meals ALTER COLUMN created_by_user TYPE TEXT;

ALTER TABLE meal_recipes ALTER COLUMN id TYPE TEXT;
ALTER TABLE meal_recipes ALTER COLUMN meal_id TYPE TEXT;
ALTER TABLE meal_recipes ALTER COLUMN recipe_id TYPE TEXT;

ALTER TABLE meal_plans ALTER COLUMN id TYPE TEXT;
ALTER TABLE meal_plans ALTER COLUMN belongs_to_household TYPE TEXT;

ALTER TABLE meal_plan_options ALTER COLUMN id TYPE TEXT;
ALTER TABLE meal_plan_options ALTER COLUMN meal_id TYPE TEXT;
ALTER TABLE meal_plan_options ALTER COLUMN belongs_to_meal_plan_event TYPE TEXT;
ALTER TABLE meal_plan_options ALTER COLUMN assigned_cook TYPE TEXT;
ALTER TABLE meal_plan_options ALTER COLUMN assigned_dishwasher TYPE TEXT;

ALTER TABLE meal_plan_option_votes ALTER COLUMN id TYPE TEXT;
ALTER TABLE meal_plan_option_votes ALTER COLUMN by_user TYPE TEXT;
ALTER TABLE meal_plan_option_votes ALTER COLUMN belongs_to_meal_plan_option TYPE TEXT;

ALTER TABLE password_reset_tokens ALTER COLUMN id TYPE TEXT;
ALTER TABLE password_reset_tokens ALTER COLUMN belongs_to_user TYPE TEXT;

ALTER TABLE valid_measurement_units ALTER COLUMN id TYPE TEXT;

ALTER TABLE valid_preparation_instruments ALTER COLUMN id TYPE TEXT;
ALTER TABLE valid_preparation_instruments ALTER COLUMN valid_preparation_id TYPE TEXT;
ALTER TABLE valid_preparation_instruments ALTER COLUMN valid_instrument_id TYPE TEXT;

ALTER TABLE valid_ingredient_measurement_units ALTER COLUMN id TYPE TEXT;
ALTER TABLE valid_ingredient_measurement_units ALTER COLUMN valid_ingredient_id TYPE TEXT;
ALTER TABLE valid_ingredient_measurement_units ALTER COLUMN valid_measurement_unit_id TYPE TEXT;

ALTER TABLE meal_plan_events ALTER COLUMN id TYPE TEXT;
ALTER TABLE meal_plan_events ALTER COLUMN belongs_to_meal_plan TYPE TEXT;

ALTER TABLE recipe_prep_tasks ALTER COLUMN id TYPE TEXT;
ALTER TABLE recipe_prep_tasks ALTER COLUMN belongs_to_recipe TYPE TEXT;

ALTER TABLE recipe_prep_task_steps ALTER COLUMN id TYPE TEXT;
ALTER TABLE recipe_prep_task_steps ALTER COLUMN belongs_to_recipe_step TYPE TEXT;
ALTER TABLE recipe_prep_task_steps ALTER COLUMN belongs_to_recipe_prep_task TYPE TEXT;

ALTER TABLE meal_plan_tasks ALTER COLUMN id TYPE TEXT;
ALTER TABLE meal_plan_tasks ALTER COLUMN belongs_to_meal_plan_option TYPE TEXT;
ALTER TABLE meal_plan_tasks ALTER COLUMN belongs_to_recipe_prep_task TYPE TEXT;
ALTER TABLE meal_plan_tasks ALTER COLUMN assigned_to_user TYPE TEXT;

ALTER TABLE meal_plan_grocery_list_items ALTER COLUMN id TYPE TEXT;
ALTER TABLE meal_plan_grocery_list_items ALTER COLUMN belongs_to_meal_plan TYPE TEXT;
ALTER TABLE meal_plan_grocery_list_items ALTER COLUMN valid_ingredient TYPE TEXT;
ALTER TABLE meal_plan_grocery_list_items ALTER COLUMN valid_measurement_unit TYPE TEXT;
ALTER TABLE meal_plan_grocery_list_items ALTER COLUMN purchased_measurement_unit TYPE TEXT;

ALTER TABLE valid_measurement_conversions ALTER COLUMN id TYPE TEXT;
ALTER TABLE valid_measurement_conversions ALTER COLUMN from_unit TYPE TEXT;
ALTER TABLE valid_measurement_conversions ALTER COLUMN to_unit TYPE TEXT;
ALTER TABLE valid_measurement_conversions ALTER COLUMN only_for_ingredient TYPE TEXT;















