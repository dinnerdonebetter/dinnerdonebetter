-- name: DeleteExpiredOAuth2ClientTokens :execrows
DELETE FROM oauth2_client_tokens WHERE code_expires_at < (NOW() - interval '1 day') AND access_expires_at < (NOW() - interval '1 day') AND refresh_expires_at < (NOW() - interval '1 day');

-- name: DestroyAllData :exec
TRUNCATE meal_plan_grocery_list_items, meal_plans, recipe_step_vessels, service_settings, valid_ingredient_group_members, valid_ingredient_groups, valid_measurement_units, meals, recipe_media, recipe_steps, valid_ingredient_preparations, valid_preparation_vessels, audit_log_entries, household_instrument_ownerships, household_invitations, meal_plan_options, oauth2_clients, recipe_step_ingredients, valid_measurement_unit_conversions, webhooks, households, meal_plan_events, meal_plan_option_votes, recipe_prep_task_steps, recipe_prep_tasks, valid_ingredient_measurement_units, valid_preparation_instruments, recipe_ratings, recipe_step_instruments, users, valid_preparations, household_user_memberships, meal_plan_tasks, oauth2_client_tokens, password_reset_tokens, valid_ingredient_state_ingredients, valid_instruments, webhook_trigger_events, meal_components, valid_vessels, recipe_step_completion_condition_ingredients, recipe_step_completion_conditions, recipe_step_products, recipes, service_setting_configurations, user_ingredient_preferences, user_notifications, valid_ingredient_states, valid_ingredients CASCADE;
