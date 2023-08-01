package main

import (
	"fmt"
	"os"
)

func readFileAndAddAnnotation(filepath string, spec sqlcAnnotation) (string, error) {
	query, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("-- name: %s :%s\n\n%s", spec.Name, spec.QueryType, string(query)), nil
}

const (
	queryFolder = "internal/database/postgres/sqlc_queries"

	sqlcQueryTypeExec = "exec"
	sqlcQueryTypeMany = "many"
	sqlcQueryTypeOne  = "one"
)

type sqlcAnnotation struct {
	Name,
	QueryType string
}

var (
	filenameToAnnotationMap = map[string]sqlcAnnotation{
		"meals/get_needing_indexing.sql": {
			Name:      "GetMealsNeedingIndexing",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_measurement_units/get_needing_indexing.sql": {
			Name:      "GetValidMeasurementUnitsNeedingIndexing",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_preparations/get_needing_indexing.sql": {
			Name:      "GetValidPreparationsNeedingIndexing",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_vessels/get_needing_indexing.sql": {
			Name:      "GetValidVesselsNeedingIndexing",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_instruments/get_needing_indexing.sql": {
			Name:      "GetValidInstrumentsNeedingIndexing",
			QueryType: sqlcQueryTypeMany,
		},
		"users/get_needing_indexing.sql": {
			Name:      "GetUsersNeedingIndexing",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredients/get_needing_indexing.sql": {
			Name:      "GetValidIngredientsNeedingIndexing",
			QueryType: sqlcQueryTypeMany,
		},
		"recipes/get_needing_indexing.sql": {
			Name:      "GetRecipesNeedingIndexing",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredient_states/get_needing_indexing.sql": {
			Name:      "GetValidIngredientStatesNeedingIndexing",
			QueryType: sqlcQueryTypeMany,
		},
		"users/get_by_username.sql": {
			Name:      "GetUserByUsername",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredients/create.sql": {
			Name:      "CreateValidIngredient",
			QueryType: sqlcQueryTypeExec,
		},
		"meals/create.sql": {
			Name:      "CreateMeal",
			QueryType: sqlcQueryTypeExec,
		},
		"password_reset_tokens/create.sql": {
			Name:      "CreatePasswordResetToken",
			QueryType: sqlcQueryTypeExec,
		},
		"oauth2_clients/create.sql": {
			Name:      "CreateOAuth2Client",
			QueryType: sqlcQueryTypeExec,
		},
		"household_invitations/create.sql": {
			Name:      "CreateHouseholdInvitation",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_measurement_conversions/create.sql": {
			Name:      "CreateValidMeasurementConversion",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_media/create.sql": {
			Name:      "CreateRecipeMedia",
			QueryType: sqlcQueryTypeExec,
		},
		"recipes/create.sql": {
			Name:      "CreateRecipe",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_components/create.sql": {
			Name:      "CreateMealComponent",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_measurement_units/create.sql": {
			Name:      "CreateValidIngredientMeasurementUnit",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_measurement_units/create.sql": {
			Name:      "CreateValidMeasurementUnit",
			QueryType: sqlcQueryTypeExec,
		},
		"webhook_trigger_events/create.sql": {
			Name:      "CreateWebhookTriggerEvent",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_steps/create.sql": {
			Name:      "CreateRecipeStep",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_instruments/create.sql": {
			Name:      "CreateRecipeStepInstrument",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_ratings/create.sql": {
			Name:      "CreateRecipeRating",
			QueryType: sqlcQueryTypeExec,
		},
		"service_settings/create.sql": {
			Name:      "CreateServiceSetting",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_grocery_list_items/create.sql": {
			Name:      "CreateMealPlanGroceryListItem",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_events/create.sql": {
			Name:      "CreateMealPlanEvent",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparations/create.sql": {
			Name:      "CreateValidPreparation",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_options/create.sql": {
			Name:      "CreateMealPlanOption",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plans/create.sql": {
			Name:      "CreateMealPlan",
			QueryType: sqlcQueryTypeExec,
		},
		"user_ingredient_preferences/create.sql": {
			Name:      "CreateUserIngredientPreference",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_preparations/create.sql": {
			Name:      "CreateValidIngredientPreparation",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_state_ingredients/create.sql": {
			Name:      "CreateValidIngredientStateIngredient",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_ingredients/create.sql": {
			Name:      "CreateRecipeStepIngredient",
			QueryType: sqlcQueryTypeExec,
		},
		"service_setting_configurations/create.sql": {
			Name:      "CreateServiceSettingConfiguration",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_option_votes/create.sql": {
			Name:      "CreateMealPlanOptionVote",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_groups/create.sql": {
			Name:      "CreateValidIngredientGroup",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_products/create.sql": {
			Name:      "CreateRecipeStepProduct",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_prep_task_steps/create.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"webhooks/create.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"households/create.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_completion_condition_ingredients/create.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"oauth2_client_tokens/create.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_tasks/create.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"users/create.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_vessels/create.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_completion_conditions/create.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_states/create.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_instruments/create.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"household_instrument_ownerships/create.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_prep_tasks/create.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparation_instruments/create.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparation_vessels/create.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_vessels/create.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_prep_tasks/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_products/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_steps/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_options/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"household_instrument_ownerships/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipes/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_vessels/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparation_vessels/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_instruments/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_state_ingredients/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_measurement_units/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparations/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_ratings/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_instruments/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_grocery_list_items/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_option_votes/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_states/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_measurement_units/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"users/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plans/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredients/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparation_instruments/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"households/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"user_ingredient_preferences/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_measurement_conversions/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_preparations/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_vessels/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_media/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_ingredients/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"service_setting_configurations/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_completion_conditions/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_events/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_groups/update.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_prep_tasks/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_options/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"webhooks/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_measurement_units/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparation_instruments/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_vessels/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_ratings/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"service_settings/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_media/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_grocery_list_items/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_measurement_units/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_measurement_conversions/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_vessels/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_state_ingredients/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"oauth2_client_tokens/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_tasks/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_ingredients/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"service_setting_configurations/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_groups/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"user_ingredient_preferences/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plans/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"household_instrument_ownerships/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_products/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_states/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_instruments/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipes/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredients/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_steps/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meals/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_preparations/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_instruments/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_option_votes/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparations/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"household_invitations/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_events/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparation_vessels/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_completion_conditions/exists.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_ratings/get_many.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"recipe_step_completion_conditions/get_many.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredient_groups/get_many.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_instruments/get_many.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"webhooks/get_many.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredients/get_many.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_vessels/get_many.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"household_instrument_ownerships/get_many.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"recipes/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_products/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_preparations/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"oauth2_clients/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"users/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_ingredients/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_instruments/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_prep_tasks/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparation_vessels/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_vessels/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_measurement_units/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_state_ingredients/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_vessels/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_states/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_events/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_options/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_option_votes/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_instruments/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_measurement_units/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meals/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparations/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"webhooks/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"households/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"household_instrument_ownerships/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"service_setting_configurations/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"service_settings/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_steps/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"user_ingredient_preferences/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_grocery_list_items/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparation_instruments/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_ratings/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plans/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_completion_conditions/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_measurement_conversions/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_groups/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_media/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_state_ingredients/pair_is_valid.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredients/archive.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_groups/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"meal_plans/mark_as_steps_created.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plans/get_finalized_without_grocery_list_init.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_vessels/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipes/update_last_indexed_at.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_states/update_last_indexed_at.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_measurement_units/update_last_indexed_at.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_vessels/update_last_indexed_at.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_instruments/update_last_indexed_at.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredients/update_last_indexed_at.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meals/update_last_indexed_at.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparations/update_last_indexed_at.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"users/update_last_indexed_at.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_state_ingredients/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plans/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_measurement_units/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meals/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_products/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_vessels/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_steps/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_tasks/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"service_settings/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_completion_conditions/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparations/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparation_instruments/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"webhooks/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"recipe_prep_tasks/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparation_vessels/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_measurement_units/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_events/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_options/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_measurement_conversions/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_media/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_preparations/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_states/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_instruments/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"password_reset_tokens/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_ratings/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"household_instrument_ownerships/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_instruments/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredients/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_option_votes/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_ingredients/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_grocery_list_items/get_one.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"users/update_two_factor_secret.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"users/update_avatar_src.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"users/update_details.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"users/update_email_address.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"users/update_password.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"users/update_username.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_tasks/change_status.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_measurement_units/search.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"service_settings/search.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_vessels/search.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_instruments/search.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_preparations/search.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredients/search.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredient_states/search.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredient_groups/search.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"users/get_by_email.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_measurement_units/get_random.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_vessels/get_random.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_instruments/get_random.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredients/get_random.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparations/get_random.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_instruments/get_by_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_measurement_units/get_by_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredients/get_by_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipes/get_by_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"service_setting_configurations/get_by_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"users/get_by_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparations/get_by_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_options/finalize.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plans/get_one_past_voting_deadline.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"oauth2_client_tokens/get_one_by_access.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"oauth2_client_tokens/get_one_by_code.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_steps/get_one_by_recipe_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"oauth2_client_tokens/get_one_by_refresh.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_options/get_one_by_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_ingredients/get_for_recipe.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"recipes/get_by_id_and_author_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredient_groups/create_group_member.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"oauth2_client_tokens/archive_by_code.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"household_user_memberships/mark_as_user_default.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_tasks/list_incomplete_by_meal_plan_option.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"households/get_by_id_with_memberships.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"oauth2_client_tokens/archive_by_refresh.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_completion_conditions/get_all_for_recipe.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"meal_plan_options/get_for_meal_plan_event.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"recipe_step_products/get_for_recipe.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"household_invitations/get_by_token_and_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"service_setting_configurations/get_settings_for_household.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"service_setting_configurations/get_settings_for_user.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"household_invitations/attach_invitations_to_user_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"service_setting_configurations/get_for_user_by_setting_name.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"household_user_memberships/get_default_household_id_for_user.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparation_vessels/pair_is_valid.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_media/for_recipe.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"webhooks/get_for_user.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"meal_plan_events/eligible_for_voting.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"users/mark_email_address_as_verified.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparation_instruments/pair_is_valid.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"household_user_memberships/remove_user_from_household.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_measurement_units/pair_is_valid.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"admin/set_user_account_status.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"service_setting_configurations/get_for_household_by_setting_name.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_measurement_conversions/get_all_from_measurement_unit.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"household_invitations/get_by_household_and_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"users/get_by_email_verification_token.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"users/mark_two_factor_secret_as_unverified.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plans/get_finalized_for_planning.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"user_ingredient_preferences/get_for_user.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"users/get_admin_by_username.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_preparations/get_values_for_ingredient.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"household_user_memberships/create_for_new_user.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_groups/archive_group_member.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_grocery_list_items/get_all_for_meal_plan.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredient_preparations/get_values_for_preparation.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"meal_plan_events/get_for_meal_plan.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"household_user_memberships/add_user_to_household.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"oauth2_clients/get_by_client_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_prep_tasks/list_all_by_recipe.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"users/exists_with_status.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"users/accept_terms_of_service_for_user.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_media/for_recipe_step.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"recipe_step_vessels/get_for_recipe.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"users/search_by_username.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredient_preparations/pair_is_valid.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_preparations/search_by_preparation_and_ingredient_name.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"household_invitations/get_by_email_and_token.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_option_votes/get_for_meal_plan_option.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"household_user_memberships/transfer_membership.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"household_user_memberships/transfer_ownership.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_instruments/get_for_recipe.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"users/mark_two_factor_secret_as_verified.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"oauth2_client_tokens/archive_by_access.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"recipes/ids_for_meal.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"meal_plans/mark_as_grocery_list_initialized.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"household_user_memberships/user_is_member.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plans/get_expired_and_unresolved.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"password_reset_tokens/redeem.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"households/add_to_household_during_creation.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"household_user_memberships/modify_user_permissions.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"users/get_with_verified_two_factor.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"users/archive_memberships.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_measurement_conversions/get_all_to_measurement_unit.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"meal_plan_tasks/list_all_by_meal_plan.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_measurement_units/search_by_ingredient_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"meal_plans/finalize.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"users/get_email_verification_token_by_user_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"household_user_memberships/get_for_user.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
		"household_invitations/set_status.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"users/accept_privacy_policy_for_user.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeExec,
		},
		"oauth2_clients/get_by_database_id.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredients/search_by_preparation_and_ingredient_name.sql": {
			Name:      "",
			QueryType: sqlcQueryTypeMany,
		},
	}
)

func main() {
	queryFolders, err := os.ReadDir(queryFolder)
	if err != nil {
		panic(err)
	}

	uniqueQueryNames := map[string]bool{}

	for _, folder := range queryFolders {
		n := folder.Name()

		var queryFiles []os.DirEntry
		queryFiles, err = os.ReadDir(fmt.Sprintf("%s/%s", queryFolder, n))
		if err != nil {
			panic(err)
		}

		for _, queryFile := range queryFiles {
			qfn := queryFile.Name()

			spec := filenameToAnnotationMap[fmt.Sprintf("%s/%s", n, qfn)]

			var query string
			query, err = readFileAndAddAnnotation(fmt.Sprintf("%s/%s/%s", queryFolder, n, qfn), spec)
			if err != nil {
				panic(err)
			}

			fmt.Println(query)
		}
	}
}
