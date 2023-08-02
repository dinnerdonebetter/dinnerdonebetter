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
			Name:      "CreateRecipePrepTaskStep",
			QueryType: sqlcQueryTypeExec,
		},
		"webhooks/create.sql": {
			Name:      "CreateWebhook",
			QueryType: sqlcQueryTypeExec,
		},
		"households/create.sql": {
			Name:      "CreateHousehold",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_completion_condition_ingredients/create.sql": {
			Name:      "CreateRecipeStepCompletionConditionIngredient",
			QueryType: sqlcQueryTypeExec,
		},
		"oauth2_client_tokens/create.sql": {
			Name:      "CreateOAuth2ClientToken",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_tasks/create.sql": {
			Name:      "CreateMealPlanTask",
			QueryType: sqlcQueryTypeExec,
		},
		"users/create.sql": {
			Name:      "CreateUser",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_vessels/create.sql": {
			Name:      "CreateValidVessel",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_completion_conditions/create.sql": {
			Name:      "CreateRecipeStepCompletionCondition",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_states/create.sql": {
			Name:      "CreateValidIngredientState",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_instruments/create.sql": {
			Name:      "CreateValidInstrument",
			QueryType: sqlcQueryTypeExec,
		},
		"household_instrument_ownerships/create.sql": {
			Name:      "CreateHouseholdInstrumentOwnership",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_prep_tasks/create.sql": {
			Name:      "CreateRecipePrepTask",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparation_instruments/create.sql": {
			Name:      "CreateValidPreparationInstrument",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparation_vessels/create.sql": {
			Name:      "CreateValidPreparationVessel",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_vessels/create.sql": {
			Name:      "CreateRecipeStepVessel",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_prep_tasks/update.sql": {
			Name:      "UpdateRecipePrepTask",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_products/update.sql": {
			Name:      "UpdateRecipeStepProduct",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_steps/update.sql": {
			Name:      "UpdateRecipeStep",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_options/update.sql": {
			Name:      "UpdateMealPlanOption",
			QueryType: sqlcQueryTypeExec,
		},
		"household_instrument_ownerships/update.sql": {
			Name:      "UpdateHouseholdInstrumentOwnership",
			QueryType: sqlcQueryTypeExec,
		},
		"recipes/update.sql": {
			Name:      "UpdateRecipe",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_vessels/update.sql": {
			Name:      "UpdateRecipeStepVessel",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparation_vessels/update.sql": {
			Name:      "UpdateValidPreparationVessel",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_instruments/update.sql": {
			Name:      "UpdateRecipeStepInstrument",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_state_ingredients/update.sql": {
			Name:      "UpdateValidIngredientStateIngredient",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_measurement_units/update.sql": {
			Name:      "UpdateValidIngredientMeasurementUnit",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparations/update.sql": {
			Name:      "UpdateValidPreparation",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_ratings/update.sql": {
			Name:      "UpdateRecipeRating",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_instruments/update.sql": {
			Name:      "UpdateValidInstrument",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_grocery_list_items/update.sql": {
			Name:      "UpdateMealPlanGroceryListItem",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_option_votes/update.sql": {
			Name:      "UpdateMealPlanOptionVote",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_states/update.sql": {
			Name:      "UpdateValidIngredientState",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_measurement_units/update.sql": {
			Name:      "UpdateValidMeasurementUnit",
			QueryType: sqlcQueryTypeExec,
		},
		"users/update.sql": {
			Name:      "UpdateUser",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plans/update.sql": {
			Name:      "UpdateMealPlan",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredients/update.sql": {
			Name:      "UpdateValidIngredient",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparation_instruments/update.sql": {
			Name:      "UpdateValidPreparationInstrument",
			QueryType: sqlcQueryTypeExec,
		},
		"households/update.sql": {
			Name:      "UpdateHousehold",
			QueryType: sqlcQueryTypeExec,
		},
		"user_ingredient_preferences/update.sql": {
			Name:      "UpdateUserIngredientPreference",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_measurement_conversions/update.sql": {
			Name:      "UpdateValidMeasurementConversion",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_preparations/update.sql": {
			Name:      "UpdateValidIngredientPreparation",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_vessels/update.sql": {
			Name:      "UpdateValidVessel",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_media/update.sql": {
			Name:      "UpdateRecipeMedia",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_ingredients/update.sql": {
			Name:      "UpdateRecipeStepIngredient",
			QueryType: sqlcQueryTypeExec,
		},
		"service_setting_configurations/update.sql": {
			Name:      "UpdateServiceSettingConfiguration",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_completion_conditions/update.sql": {
			Name:      "UpdateRecipeStepCompletionCondition",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_events/update.sql": {
			Name:      "UpdateMealPlanEvent",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_groups/update.sql": {
			Name:      "UpdateValidIngredientGroup",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_prep_tasks/exists.sql": {
			Name:      "CheckRecipePrepTaskExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_options/exists.sql": {
			Name:      "CheckMealPlanOptionExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"webhooks/exists.sql": {
			Name:      "CheckWebhookExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_measurement_units/exists.sql": {
			Name:      "CheckValidMeasurementUnitExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparation_instruments/exists.sql": {
			Name:      "CheckValidPreparationInstrumentExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_vessels/exists.sql": {
			Name:      "CheckValidVesselExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_ratings/exists.sql": {
			Name:      "CheckRecipeRatingExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"service_settings/exists.sql": {
			Name:      "CheckServiceSettingExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_media/exists.sql": {
			Name:      "CheckRecipeMediaExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_grocery_list_items/exists.sql": {
			Name:      "CheckMealPlanGroceryListItemExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_measurement_units/exists.sql": {
			Name:      "CheckValidIngredientMeasurementUnitExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_measurement_conversions/exists.sql": {
			Name:      "CheckValidMeasurementConversionExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_vessels/exists.sql": {
			Name:      "CheckRecipeStepVesselExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_state_ingredients/exists.sql": {
			Name:      "CheckValidIngredientStateIngredientExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"oauth2_client_tokens/exists.sql": {
			Name:      "CheckOAuth2ClientTokenExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_tasks/exists.sql": {
			Name:      "CheckMealPlanTaskExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_ingredients/exists.sql": {
			Name:      "CheckRecipeStepIngredientExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"service_setting_configurations/exists.sql": {
			Name:      "CheckServiceSettingConfigurationExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_groups/exists.sql": {
			Name:      "CheckValidIngredientGroupExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"user_ingredient_preferences/exists.sql": {
			Name:      "CheckUserIngredientPreferenceExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plans/exists.sql": {
			Name:      "CheckMealPlanExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"household_instrument_ownerships/exists.sql": {
			Name:      "CheckHouseholdInstrumentOwnershipExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_products/exists.sql": {
			Name:      "CheckRecipeStepProductExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_states/exists.sql": {
			Name:      "CheckValidIngredientStateExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_instruments/exists.sql": {
			Name:      "CheckValidInstrumentExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"recipes/exists.sql": {
			Name:      "CheckRecipeExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredients/exists.sql": {
			Name:      "CheckValidIngredientExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_steps/exists.sql": {
			Name:      "CheckRecipeStepExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"meals/exists.sql": {
			Name:      "CheckMealExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_preparations/exists.sql": {
			Name:      "CheckValidIngredientPreparationExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_instruments/exists.sql": {
			Name:      "CheckRecipeStepInstrumentExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_option_votes/exists.sql": {
			Name:      "CheckMealPlanOptionVoteExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparations/exists.sql": {
			Name:      "CheckValidPreparationExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"household_invitations/exists.sql": {
			Name:      "CheckHouseholdInvitationExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_events/exists.sql": {
			Name:      "CheckMealPlanEventExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparation_vessels/exists.sql": {
			Name:      "CheckValidPreparationVesselExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_completion_conditions/exists.sql": {
			Name:      "CheckRecipeStepCompletionConditionExistence",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_ratings/get_many.sql": {
			Name:      "GetRecipeRatings",
			QueryType: sqlcQueryTypeMany,
		},
		"recipe_step_completion_conditions/get_many.sql": {
			Name:      "GetRecipeStepCompletionConditions",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredient_groups/get_many.sql": {
			Name:      "GetValidIngredientGroups",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_instruments/get_many.sql": {
			Name:      "GetValidInstruments",
			QueryType: sqlcQueryTypeMany,
		},
		"webhooks/get_many.sql": {
			Name:      "GetWebhooks",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredients/get_many.sql": {
			Name:      "GetValidIngredients",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_vessels/get_many.sql": {
			Name:      "GetValidVessels",
			QueryType: sqlcQueryTypeMany,
		},
		"household_instrument_ownerships/get_many.sql": {
			Name:      "GetHouseholdInstrumentOwnerships",
			QueryType: sqlcQueryTypeMany,
		},
		"recipes/archive.sql": {
			Name:      "ArchiveRecipe",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_products/archive.sql": {
			Name:      "ArchiveRecipeStepProduct",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_preparations/archive.sql": {
			Name:      "ArchiveValidIngredientPreparation",
			QueryType: sqlcQueryTypeExec,
		},
		"oauth2_clients/archive.sql": {
			Name:      "ArchiveOAuth2Client",
			QueryType: sqlcQueryTypeExec,
		},
		"users/archive.sql": {
			Name:      "ArchiveUser",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_ingredients/archive.sql": {
			Name:      "ArchiveRecipeStepIngredient",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_instruments/archive.sql": {
			Name:      "ArchiveValidInstrument",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_prep_tasks/archive.sql": {
			Name:      "ArchiveRecipePrepTask",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparation_vessels/archive.sql": {
			Name:      "ArchiveValidPreparationVessel",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_vessels/archive.sql": {
			Name:      "ArchiveRecipeStepVessel",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_measurement_units/archive.sql": {
			Name:      "ArchiveValidMeasurementUnit",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_state_ingredients/archive.sql": {
			Name:      "ArchiveValidIngredientStateIngredient",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_vessels/archive.sql": {
			Name:      "ArchiveValidVessel",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_states/archive.sql": {
			Name:      "ArchiveValidIngredientState",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_events/archive.sql": {
			Name:      "ArchiveMealPlanEvent",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_options/archive.sql": {
			Name:      "ArchiveMealPlanOption",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_option_votes/archive.sql": {
			Name:      "ArchiveMealPlanOptionVote",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_instruments/archive.sql": {
			Name:      "ArchiveRecipeStepInstrument",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_measurement_units/archive.sql": {
			Name:      "ArchiveValidIngredientMeasurementUnit",
			QueryType: sqlcQueryTypeExec,
		},
		"meals/archive.sql": {
			Name:      "ArchiveMeal",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparations/archive.sql": {
			Name:      "ArchiveValidPreparation",
			QueryType: sqlcQueryTypeExec,
		},
		"webhooks/archive.sql": {
			Name:      "ArchiveWebhook",
			QueryType: sqlcQueryTypeExec,
		},
		"households/archive.sql": {
			Name:      "ArchiveHousehold",
			QueryType: sqlcQueryTypeExec,
		},
		"household_instrument_ownerships/archive.sql": {
			Name:      "ArchiveHouseholdInstrumentOwnership",
			QueryType: sqlcQueryTypeExec,
		},
		"service_setting_configurations/archive.sql": {
			Name:      "ArchiveServiceSettingConfiguration",
			QueryType: sqlcQueryTypeExec,
		},
		"service_settings/archive.sql": {
			Name:      "ArchiveServiceSetting",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_steps/archive.sql": {
			Name:      "ArchiveRecipeStep",
			QueryType: sqlcQueryTypeExec,
		},
		"user_ingredient_preferences/archive.sql": {
			Name:      "ArchiveUserIngredientPreference",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_grocery_list_items/archive.sql": {
			Name:      "ArchiveMealPlanGroceryListItem",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparation_instruments/archive.sql": {
			Name:      "ArchiveValidPreparationInstrument",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_ratings/archive.sql": {
			Name:      "ArchiveRecipeRating",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plans/archive.sql": {
			Name:      "ArchiveMealPlan",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_completion_conditions/archive.sql": {
			Name:      "ArchiveRecipeStepCompletionCondition",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_measurement_conversions/archive.sql": {
			Name:      "ArchiveValidMeasurementConversion",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_groups/archive.sql": {
			Name:      "ArchiveValidIngredientGroup",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_media/archive.sql": {
			Name:      "ArchiveRecipeMedia",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_state_ingredients/pair_is_valid.sql": {
			Name:      "CheckValidityOfValidIngredientStateIngredientPair",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredients/archive.sql": {
			Name:      "ArchiveValidIngredient",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_groups/get_one.sql": {
			Name:      "GetValidIngredientGroup",
			QueryType: sqlcQueryTypeMany,
		},
		"meal_plans/mark_as_tasks_created.sql": {
			Name:      "MarkMealPlanAsPrepTasksCreated",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plans/get_finalized_without_grocery_list_init.sql": {
			Name:      "GetFinalizedMealPlansWithoutGroceryListInit",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_vessels/get_one.sql": {
			Name:      "GetValidVessel",
			QueryType: sqlcQueryTypeOne,
		},
		"recipes/update_last_indexed_at.sql": {
			Name:      "UpdateRecipeLastIndexedAt",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_states/update_last_indexed_at.sql": {
			Name:      "UpdateValidIngredientStateLastIndexedAt",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_measurement_units/update_last_indexed_at.sql": {
			Name:      "UpdateValidMeasurementUnitLastIndexedAt",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_vessels/update_last_indexed_at.sql": {
			Name:      "UpdateValidVesselLastIndexedAt",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_instruments/update_last_indexed_at.sql": {
			Name:      "UpdateValidInstrumentLastIndexedAt",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredients/update_last_indexed_at.sql": {
			Name:      "UpdateValidIngredientLastIndexedAt",
			QueryType: sqlcQueryTypeExec,
		},
		"meals/update_last_indexed_at.sql": {
			Name:      "UpdateMealLastIndexedAt",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparations/update_last_indexed_at.sql": {
			Name:      "UpdateValidPreparationLastIndexedAt",
			QueryType: sqlcQueryTypeExec,
		},
		"users/update_last_indexed_at.sql": {
			Name:      "UpdateUserLastIndexedAt",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_state_ingredients/get_one.sql": {
			Name:      "GetValidIngredientStateIngredient",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plans/get_one.sql": {
			Name:      "GetMealPlan",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_measurement_units/get_one.sql": {
			Name:      "GetValidMeasurementUnit",
			QueryType: sqlcQueryTypeOne,
		},
		"meals/get_one.sql": {
			Name:      "GetMeal",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_products/get_one.sql": {
			Name:      "GetRecipeStepProduct",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_vessels/get_one.sql": {
			Name:      "GetRecipeStepVessel",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_steps/get_one.sql": {
			Name:      "GetRecipeStep",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_tasks/get_one.sql": {
			Name:      "GetMealPlanTask",
			QueryType: sqlcQueryTypeOne,
		},
		"service_settings/get_one.sql": {
			Name:      "GetServiceSetting",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_completion_conditions/get_one.sql": {
			Name:      "GetRecipeStepCompletionCondition",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparations/get_one.sql": {
			Name:      "GetValidPreparation",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparation_instruments/get_one.sql": {
			Name:      "GetValidPreparationInstrument",
			QueryType: sqlcQueryTypeOne,
		},
		"webhooks/get_one.sql": {
			Name:      "GetWebhook",
			QueryType: sqlcQueryTypeMany,
		},
		"recipe_prep_tasks/get_one.sql": {
			Name:      "GetRecipePrepTask",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparation_vessels/get_one.sql": {
			Name:      "GetValidPreparationVessel",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_measurement_units/get_one.sql": {
			Name:      "GetValidIngredientMeasurementUnit",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_events/get_one.sql": {
			Name:      "GetMealPlanEvent",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_options/get_one.sql": {
			Name:      "GetMealPlanOption",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_measurement_conversions/get_one.sql": {
			Name:      "GetValidMeasurementConversion",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_media/get_one.sql": {
			Name:      "GetRecipeMedia",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_preparations/get_one.sql": {
			Name:      "GetValidIngredientPreparation",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_states/get_one.sql": {
			Name:      "GetValidIngredientState",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_instruments/get_one.sql": {
			Name:      "GetValidInstrument",
			QueryType: sqlcQueryTypeOne,
		},
		"password_reset_tokens/get_one.sql": {
			Name:      "GetPasswordResetToken",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_ratings/get_one.sql": {
			Name:      "GetRecipeRating",
			QueryType: sqlcQueryTypeOne,
		},
		"household_instrument_ownerships/get_one.sql": {
			Name:      "GetHouseholdInstrumentOwnership",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_instruments/get_one.sql": {
			Name:      "GetRecipeStepInstrument",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredients/get_one.sql": {
			Name:      "GetValidIngredient",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_option_votes/get_one.sql": {
			Name:      "GetMealPlanOptionVote",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_ingredients/get_one.sql": {
			Name:      "GetRecipeStepIngredient",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_grocery_list_items/get_one.sql": {
			Name:      "GetMealPlanGroceryListItem",
			QueryType: sqlcQueryTypeOne,
		},
		"users/update_two_factor_secret.sql": {
			Name:      "UpdateUserTwoFactorSecret",
			QueryType: sqlcQueryTypeExec,
		},
		"users/update_avatar_src.sql": {
			Name:      "UpdateUserAvatarSrc",
			QueryType: sqlcQueryTypeExec,
		},
		"users/update_details.sql": {
			Name:      "UpdateUserDetails",
			QueryType: sqlcQueryTypeExec,
		},
		"users/update_email_address.sql": {
			Name:      "UpdateUserEmailAddress",
			QueryType: sqlcQueryTypeExec,
		},
		"users/update_password.sql": {
			Name:      "UpdateUserPassword",
			QueryType: sqlcQueryTypeExec,
		},
		"users/update_username.sql": {
			Name:      "UpdateUserUsername",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_tasks/change_status.sql": {
			Name:      "ChangeMealPlanTaskStatus",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_measurement_units/search.sql": {
			Name:      "SearchForValidMeasurementUnits",
			QueryType: sqlcQueryTypeMany,
		},
		"service_settings/search.sql": {
			Name:      "SearchForServiceSettings",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_vessels/search.sql": {
			Name:      "SearchForValidVessels",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_instruments/search.sql": {
			Name:      "SearchForValidInstruments",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_preparations/search.sql": {
			Name:      "SearchForValidPreparations",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredients/search.sql": {
			Name:      "SearchForValidIngredients",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredient_states/search.sql": {
			Name:      "SearchForValidIngredientStates",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredient_groups/search.sql": {
			Name:      "SearchForValidIngredientGroups",
			QueryType: sqlcQueryTypeMany,
		},
		"users/get_by_email.sql": {
			Name:      "GetUserByEmail",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_measurement_units/get_random.sql": {
			Name:      "GetRandomValidMeasurementUnit",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_vessels/get_random.sql": {
			Name:      "GetRandomValidVessel",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_instruments/get_random.sql": {
			Name:      "GetRandomValidInstrument",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredients/get_random.sql": {
			Name:      "GetRandomValidIngredient",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparations/get_random.sql": {
			Name:      "GetRandomValidPreparation",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_instruments/get_by_id.sql": {
			Name:      "GetValidInstrumentByID",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_measurement_units/get_by_id.sql": {
			Name:      "GetValidMeasurementUnitByID",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredients/get_by_id.sql": {
			Name:      "GetValidIngredientByID",
			QueryType: sqlcQueryTypeOne,
		},
		"recipes/get_by_id.sql": {
			Name:      "GetRecipeByID",
			QueryType: sqlcQueryTypeOne,
		},
		"service_setting_configurations/get_by_id.sql": {
			Name:      "GetServiceSettingConfigurationByID",
			QueryType: sqlcQueryTypeOne,
		},
		"users/get_by_id.sql": {
			Name:      "GetUserByID",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparations/get_by_id.sql": {
			Name:      "GetValidPreparationByID",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_options/finalize.sql": {
			Name:      "FinalizeMealPlanOption",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plans/get_one_past_voting_deadline.sql": {
			Name:      "GetOnePastVotingDeadline",
			QueryType: sqlcQueryTypeOne,
		},
		"oauth2_client_tokens/get_one_by_access.sql": {
			Name:      "GetOAuth2ClientTokenByAccess",
			QueryType: sqlcQueryTypeOne,
		},
		"oauth2_client_tokens/get_one_by_code.sql": {
			Name:      "GetOAuth2ClientTokenByCode",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_steps/get_one_by_recipe_id.sql": {
			Name:      "GetRecipeStepByRecipeID",
			QueryType: sqlcQueryTypeOne,
		},
		"oauth2_client_tokens/get_one_by_refresh.sql": {
			Name:      "GetOAuth2ClientTokenByRefresh",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_options/get_one_by_id.sql": {
			Name:      "GetMealPlanOptionByID",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_step_ingredients/get_for_recipe.sql": {
			Name:      "GetRecipeStepIngredientsForRecipe",
			QueryType: sqlcQueryTypeMany,
		},
		"recipes/get_by_id_and_author_id.sql": {
			Name:      "GetRecipeByIDAndAuthorID",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredient_groups/create_group_member.sql": {
			Name:      "CreateValidIngredientGroupMember",
			QueryType: sqlcQueryTypeExec,
		},
		"oauth2_client_tokens/archive_by_code.sql": {
			Name:      "ArchiveOAuth2ClientTokenByCode",
			QueryType: sqlcQueryTypeExec,
		},
		"household_user_memberships/mark_as_user_default.sql": {
			Name:      "MarkHouseholdUserMembershipAsUserDefault",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_tasks/list_incomplete_by_meal_plan_option.sql": {
			Name:      "ListIncompleteMealPlanTasksByMealPlanOption",
			QueryType: sqlcQueryTypeMany,
		},
		"households/get_by_id_with_memberships.sql": {
			Name:      "GetHouseholdByIDWithMemberships",
			QueryType: sqlcQueryTypeMany,
		},
		"oauth2_client_tokens/archive_by_refresh.sql": {
			Name:      "ArchiveOAuth2ClientTokenByRefresh",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_completion_conditions/get_all_for_recipe.sql": {
			Name:      "GetAllRecipeStepCompletionConditionsForRecipe",
			QueryType: sqlcQueryTypeMany,
		},
		"meal_plan_options/get_for_meal_plan_event.sql": {
			Name:      "GetMealPlanOptionsForMealPlanEvent",
			QueryType: sqlcQueryTypeMany,
		},
		"recipe_step_products/get_for_recipe.sql": {
			Name:      "GetRecipeStepProductsForRecipe",
			QueryType: sqlcQueryTypeMany,
		},
		"household_invitations/get_by_token_and_id.sql": {
			Name:      "GetHouseholdInvitationByTokenAndID",
			QueryType: sqlcQueryTypeOne,
		},
		"service_setting_configurations/get_settings_for_household.sql": {
			Name:      "GetServiceSettingConfigurationsForHousehold",
			QueryType: sqlcQueryTypeMany,
		},
		"service_setting_configurations/get_settings_for_user.sql": {
			Name:      "GetServiceSettingConfigurationsForUser",
			QueryType: sqlcQueryTypeMany,
		},
		"household_invitations/attach_invitations_to_user_id.sql": {
			Name:      "AttachHouseholdInvitationsToUserID",
			QueryType: sqlcQueryTypeExec,
		},
		"service_setting_configurations/get_for_user_by_setting_name.sql": {
			Name:      "GetServiceSettingConfigurationsForUserBySettingName",
			QueryType: sqlcQueryTypeMany,
		},
		"household_user_memberships/get_default_household_id_for_user.sql": {
			Name:      "GetDefaultHouseholdIDForUser",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_preparation_vessels/pair_is_valid.sql": {
			Name:      "ValidPreparationVesselPairIsValid",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_media/for_recipe.sql": {
			Name:      "GetRecipeMediaForRecipe",
			QueryType: sqlcQueryTypeMany,
		},
		"webhooks/get_for_user.sql": {
			Name:      "GetWebhooksForUser",
			QueryType: sqlcQueryTypeMany,
		},
		"meal_plan_events/eligible_for_voting.sql": {
			Name:      "GetMealPlanEventsEligibleForVoting",
			QueryType: sqlcQueryTypeMany,
		},
		"users/mark_email_address_as_verified.sql": {
			Name:      "MarkEmailAddressAsVerified",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_preparation_instruments/pair_is_valid.sql": {
			Name:      "ValidPreparationInstrumentPairIsValid",
			QueryType: sqlcQueryTypeOne,
		},
		"household_user_memberships/remove_user_from_household.sql": {
			Name:      "RemoveUserFromHousehold",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_measurement_units/pair_is_valid.sql": {
			Name:      "ValidIngredientMeasurementUnitPairIsValid",
			QueryType: sqlcQueryTypeOne,
		},
		"admin/set_user_account_status.sql": {
			Name:      "SetUserAccountStatus",
			QueryType: sqlcQueryTypeExec,
		},
		"service_setting_configurations/get_for_household_by_setting_name.sql": {
			Name:      "GetServiceSettingConfigurationsForHouseholdBySettingName",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_measurement_conversions/get_all_from_measurement_unit.sql": {
			Name:      "GetAllValidMeasurementConversionsFromMeasurementUnit",
			QueryType: sqlcQueryTypeMany,
		},
		"household_invitations/get_by_household_and_id.sql": {
			Name:      "GetHouseholdInvitationByHouseholdAndID",
			QueryType: sqlcQueryTypeOne,
		},
		"users/get_by_email_verification_token.sql": {
			Name:      "GetUserByEmailAddressVerificationToken",
			QueryType: sqlcQueryTypeOne,
		},
		"users/mark_two_factor_secret_as_unverified.sql": {
			Name:      "MarkTwoFactorSecretAsUnverified",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plans/get_finalized_for_planning.sql": {
			Name:      "GetFinalizedMealPlansForPlanning",
			QueryType: sqlcQueryTypeMany,
		},
		"user_ingredient_preferences/get_for_user.sql": {
			Name:      "GetUserIngredientPreferencesForUser",
			QueryType: sqlcQueryTypeMany,
		},
		"users/get_admin_by_username.sql": {
			Name:      "GetAdminUserByUsername",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_preparations/get_values_for_ingredient.sql": {
			Name:      "GetValidIngredientPreparationsForIngredient",
			QueryType: sqlcQueryTypeMany,
		},
		"household_user_memberships/create_for_new_user.sql": {
			Name:      "CreateHouseholdUserMembershipForNewUser",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_ingredient_groups/archive_group_member.sql": {
			Name:      "ArchiveValidIngredientGroupMember",
			QueryType: sqlcQueryTypeExec,
		},
		"meal_plan_grocery_list_items/get_all_for_meal_plan.sql": {
			Name:      "GetMealPlanGroceryListItemsForMealPlan",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredient_preparations/get_values_for_preparation.sql": {
			Name:      "GetValidIngredientPreparationsForPreparation",
			QueryType: sqlcQueryTypeMany,
		},
		"meal_plan_events/get_for_meal_plan.sql": {
			Name:      "GetMealPlanEventsForMealPlan",
			QueryType: sqlcQueryTypeMany,
		},
		"household_user_memberships/add_user_to_household.sql": {
			Name:      "AddUserToHousehold",
			QueryType: sqlcQueryTypeExec,
		},
		"oauth2_clients/get_by_client_id.sql": {
			Name:      "GetOAuth2ClientByClientID",
			QueryType: sqlcQueryTypeOne,
		},
		"recipe_prep_tasks/list_all_by_recipe.sql": {
			Name:      "ListAllRecipePrepTasksByRecipe",
			QueryType: sqlcQueryTypeMany,
		},
		"users/exists_with_status.sql": {
			Name:      "UserExistsWithStatus",
			QueryType: sqlcQueryTypeOne,
		},
		"users/accept_terms_of_service_for_user.sql": {
			Name:      "AcceptTermsOfServiceForUser",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_media/for_recipe_step.sql": {
			Name:      "GetRecipeMediaForRecipeStep",
			QueryType: sqlcQueryTypeMany,
		},
		"recipe_step_vessels/get_for_recipe.sql": {
			Name:      "GetRecipeStepVesselsForRecipe",
			QueryType: sqlcQueryTypeMany,
		},
		"users/search_by_username.sql": {
			Name:      "SearchUsersByUsername",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_ingredient_preparations/pair_is_valid.sql": {
			Name:      "ValidIngredientPreparationPairIsValid",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredient_preparations/search_by_preparation_and_ingredient_name.sql": {
			Name:      "SearchValidIngredientPreparationsByPreparationAndIngredientName",
			QueryType: sqlcQueryTypeMany,
		},
		"household_invitations/get_by_email_and_token.sql": {
			Name:      "GetHouseholdInvitationByEmailAndToken",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plan_option_votes/get_for_meal_plan_option.sql": {
			Name:      "GetMealPlanOptionVotesForMealPlanOption",
			QueryType: sqlcQueryTypeMany,
		},
		"household_user_memberships/transfer_membership.sql": {
			Name:      "TransferHouseholdMembership",
			QueryType: sqlcQueryTypeExec,
		},
		"household_user_memberships/transfer_ownership.sql": {
			Name:      "TransferHouseholdOwnership",
			QueryType: sqlcQueryTypeExec,
		},
		"recipe_step_instruments/get_for_recipe.sql": {
			Name:      "GetRecipeStepInstrumentsForRecipe",
			QueryType: sqlcQueryTypeMany,
		},
		"users/mark_two_factor_secret_as_verified.sql": {
			Name:      "MarkTwoFactorSecretAsVerified",
			QueryType: sqlcQueryTypeExec,
		},
		"oauth2_client_tokens/archive_by_access.sql": {
			Name:      "ArchiveOAuth2ClientTokenByAccess",
			QueryType: sqlcQueryTypeExec,
		},
		"recipes/ids_for_meal.sql": {
			Name:      "GetRecipeIDsForMeal",
			QueryType: sqlcQueryTypeMany,
		},
		"meal_plans/mark_as_grocery_list_initialized.sql": {
			Name:      "MarkMealPlanAsGroceryListInitialized",
			QueryType: sqlcQueryTypeExec,
		},
		"household_user_memberships/user_is_member.sql": {
			Name:      "UserIsHouseholdMember",
			QueryType: sqlcQueryTypeOne,
		},
		"meal_plans/get_expired_and_unresolved.sql": {
			Name:      "GetExpiredAndUnresolvedMealPlans",
			QueryType: sqlcQueryTypeMany,
		},
		"password_reset_tokens/redeem.sql": {
			Name:      "RedeemPasswordResetToken",
			QueryType: sqlcQueryTypeExec,
		},
		"households/add_to_household_during_creation.sql": {
			Name:      "AddToHouseholdDuringCreation",
			QueryType: sqlcQueryTypeExec,
		},
		"household_user_memberships/modify_user_permissions.sql": {
			Name:      "ModifyHouseholdUserPermissions",
			QueryType: sqlcQueryTypeExec,
		},
		"users/get_with_verified_two_factor.sql": {
			Name:      "GetUserWithVerifiedTwoFactor",
			QueryType: sqlcQueryTypeOne,
		},
		"users/archive_memberships.sql": {
			Name:      "ArchiveUserMemberships",
			QueryType: sqlcQueryTypeExec,
		},
		"valid_measurement_conversions/get_all_to_measurement_unit.sql": {
			Name:      "GetAllValidMeasurementConversionsToMeasurementUnit",
			QueryType: sqlcQueryTypeMany,
		},
		"meal_plan_tasks/list_all_by_meal_plan.sql": {
			Name:      "ListAllMealPlanTasksByMealPlan",
			QueryType: sqlcQueryTypeMany,
		},
		"valid_measurement_units/search_by_ingredient_id.sql": {
			Name:      "SearchValidMeasurementUnitsByIngredientID",
			QueryType: sqlcQueryTypeMany,
		},
		"meal_plans/finalize.sql": {
			Name:      "FinalizeMealPlan",
			QueryType: sqlcQueryTypeExec,
		},
		"users/get_email_verification_token_by_user_id.sql": {
			Name:      "GetEmailVerificationTokenByUserID",
			QueryType: sqlcQueryTypeOne,
		},
		"household_user_memberships/get_for_user.sql": {
			Name:      "GetHouseholdUserMembershipsForUser",
			QueryType: sqlcQueryTypeMany,
		},
		"household_invitations/set_status.sql": {
			Name:      "SetHouseholdInvitationStatus",
			QueryType: sqlcQueryTypeExec,
		},
		"users/accept_privacy_policy_for_user.sql": {
			Name:      "AcceptPrivacyPolicyForUser",
			QueryType: sqlcQueryTypeExec,
		},
		"oauth2_clients/get_by_database_id.sql": {
			Name:      "GetOAuth2ClientByDatabaseID",
			QueryType: sqlcQueryTypeOne,
		},
		"valid_ingredients/search_by_preparation_and_ingredient_name.sql": {
			Name:      "SearchValidIngredientsByPreparationAndIngredientName",
			QueryType: sqlcQueryTypeMany,
		},
	}
)

func main() {
	queryFolders, err := os.ReadDir(queryFolder)
	if err != nil {
		panic(err)
	}

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

			if err = os.WriteFile(fmt.Sprintf("%s/%s/%s", queryFolder, n, qfn), []byte(query), 0644); err != nil {
				panic(err)
			}
		}
	}
}
