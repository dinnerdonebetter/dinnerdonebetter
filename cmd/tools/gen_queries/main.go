package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func main() {
	queryOutput := map[string][]*Query{
		"admin.sql":                          buildAdminQueries(),
		"webhooks.sql":                       buildWebhooksQueries(),
		"users.sql":                          buildUsersQueries(),
		"households.sql":                     buildHouseholdsQueries(),
		"household_user_memberships.sql":     buildHouseholdUserMembershipsQueries(),
		"webhook_trigger_events.sql":         buildWebhookTriggerEventsQueries(),
		"password_reset_tokens.sql":          buildPasswordResetTokensQueries(),
		"oauth2_client_tokens.sql":           buildOAuth2ClientTokensQueries(),
		"oauth2_clients.sql":                 buildOAuth2ClientsQueries(),
		"service_settings.sql":               buildServiceSettingQueries(),
		"service_setting_configurations.sql": buildServiceSettingConfigurationQueries(),
		"household_invitations.sql":          buildHouseholdInvitationsQueries(),
		"valid_ingredients.sql":              buildValidIngredientsQueries(),
		"valid_instruments.sql":              buildValidInstrumentsQueries(),
		"valid_preparations.sql":             buildValidPreparationsQueries(),
		"valid_measurement_units.sql":        buildValidMeasurementUnitsQueries(),
		"valid_ingredient_states.sql":        buildValidIngredientStatesQueries(),
		"valid_vessels.sql":                  buildValidVesselsQueries(),
		"valid_ingredient_groups.sql":        buildValidIngredientGroupsQueries(),
		"valid_preparation_instruments.sql":  buildValidPreparationInstrumentsQueries(),
		"valid_ingredient_preparations.sql":  buildValidIngredientPreparationsQueries(),
		//
		// "valid_ingredient_measurement_units.sql":           buildValidIngredientMeasurementUnitsQueries(),
		// "valid_ingredient_state_ingredients.sql":           buildValidIngredientStateIngredientsQueries(),
		// "valid_measurement_unit_conversions.sql":           buildValidMeasurementUnitConversionsQueries(),
		// "valid_preparation_vessels.sql":                    buildValidPreparationVesselsQueries(),
		// "household_instrument_ownerships.sql":              buildHouseholdInstrumentOwnershipQueries(),
		// "meal_components.sql":                              buildMealComponentsQueries(),
		// "meal_plan_events.sql":                             buildMealPlanEventsQueries(),
		// "meal_plan_grocery_list_items.sql":                 buildMealPlanGroceryListItemsQueries(),
		// "meal_plan_options.sql":                            buildMealPlanOptionsQueries(),
		// "meal_plan_option_votes.sql":                       buildMealPlanOptionVotesQueries(),
		// "meal_plans.sql":                                   buildMealPlansQueries(),
		// "meal_plan_tasks.sql":                              buildMealPlanTasksQueries(),
		// "meals.sql":                                        buildMealsQueries(),
		// "recipe_media.sql":                                 buildRecipeMediaQueries(),
		// "recipe_prep_tasks.sql":                            buildRecipePrepTasksQueries(),
		// "recipe_prep_task_steps.sql":                       buildRecipePrepTaskStepsQueries(),
		// "recipe_ratings.sql":                               buildRecipeRatingsQueries(),
		// "recipes.sql":                                      buildRecipesQueries(),
		// "recipe_step_completion_condition_ingredients.sql": buildRecipeStepCompletionConditionIngredientsQueries(),
		// "recipe_step_completion_conditions.sql":            buildRecipeStepCompletionConditionQueries(),
		// "recipe_step_ingredients.sql":                      buildRecipeStepIngredientsQueries(),
		// "recipe_step_instruments.sql":                      buildRecipeStepInstrumentsQueries(),
		// "recipe_step_products.sql":                         buildRecipeStepProductsQueries(),
		// "recipe_steps.sql":                                 buildRecipeStepsQueries(),
		// "recipe_step_vessels.sql":                          buildRecipeStepVesselsQueries(),
		// "user_ingredient_preferences.sql":                  buildUserIngredientPreferencesQueries(),
	}

	for filePath, queries := range queryOutput {
		existingFile, err := os.ReadFile(path.Join("internal", "database", "postgres", "sqlc_queries", filePath))
		if err != nil {
			panic(err)
		}

		var fileContent string
		for i, query := range queries {
			if i != 0 {
				fileContent += "\n"
			}
			fileContent += query.Render()
		}

		fileOutput := ""
		for _, line := range strings.Split(strings.TrimSpace(fileContent), "\n") {
			fileOutput += strings.TrimSuffix(line, " ") + "\n"
		}

		if string(existingFile) != fileOutput {
			fmt.Printf("files don't match: %s\n", filePath)
		}
	}

	fmt.Println("done")
}
