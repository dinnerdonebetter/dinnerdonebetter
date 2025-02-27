package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/pflag"
)

var (
	checkOnlyFlag = pflag.Bool("check", false, "only check if files match")
	databaseFlag  = pflag.String("database", postgres, "what database to use")
)

const (
	postgres = "postgres"
)

func main() {
	pflag.Parse()

	runErrors := &multierror.Error{}

	databaseToUse := *databaseFlag

	queryOutput := map[string][]*Query{
		"admin.sql":                                        buildAdminQueries(databaseToUse),
		"webhooks.sql":                                     buildWebhooksQueries(databaseToUse),
		"user_notifications.sql":                           buildUserNotificationQueries(databaseToUse),
		"users.sql":                                        buildUsersQueries(databaseToUse),
		"households.sql":                                   buildHouseholdsQueries(databaseToUse),
		"household_user_memberships.sql":                   buildHouseholdUserMembershipsQueries(databaseToUse),
		"webhook_trigger_events.sql":                       buildWebhookTriggerEventsQueries(databaseToUse),
		"password_reset_tokens.sql":                        buildPasswordResetTokensQueries(databaseToUse),
		"oauth2_client_tokens.sql":                         buildOAuth2ClientTokensQueries(databaseToUse),
		"oauth2_clients.sql":                               buildOAuth2ClientsQueries(databaseToUse),
		"service_settings.sql":                             buildServiceSettingQueries(databaseToUse),
		"service_setting_configurations.sql":               buildServiceSettingConfigurationQueries(databaseToUse),
		"household_invitations.sql":                        buildHouseholdInvitationsQueries(databaseToUse),
		"valid_ingredients.sql":                            buildValidIngredientsQueries(databaseToUse),
		"valid_instruments.sql":                            buildValidInstrumentsQueries(databaseToUse),
		"valid_preparations.sql":                           buildValidPreparationsQueries(databaseToUse),
		"valid_measurement_units.sql":                      buildValidMeasurementUnitsQueries(databaseToUse),
		"valid_ingredient_states.sql":                      buildValidIngredientStatesQueries(databaseToUse),
		"valid_vessels.sql":                                buildValidVesselsQueries(databaseToUse),
		"valid_ingredient_groups.sql":                      buildValidIngredientGroupsQueries(databaseToUse),
		"valid_ingredient_preparations.sql":                buildValidIngredientPreparationsQueries(databaseToUse),
		"valid_preparation_vessels.sql":                    buildValidPreparationVesselsQueries(databaseToUse),
		"valid_ingredient_measurement_units.sql":           buildValidIngredientMeasurementUnitsQueries(databaseToUse),
		"valid_measurement_unit_conversions.sql":           buildValidMeasurementUnitConversionsQueries(databaseToUse),
		"valid_ingredient_state_ingredients.sql":           buildValidIngredientStateIngredientsQueries(databaseToUse),
		"valid_preparation_instruments.sql":                buildValidPreparationInstrumentsQueries(databaseToUse),
		"household_instrument_ownerships.sql":              buildHouseholdInstrumentOwnershipQueries(databaseToUse),
		"meal_components.sql":                              buildMealComponentsQueries(databaseToUse),
		"meal_plan_events.sql":                             buildMealPlanEventsQueries(databaseToUse),
		"recipe_media.sql":                                 buildRecipeMediaQueries(databaseToUse),
		"recipe_prep_task_steps.sql":                       buildRecipePrepTaskStepsQueries(databaseToUse),
		"recipe_ratings.sql":                               buildRecipeRatingsQueries(databaseToUse),
		"recipe_step_completion_condition_ingredients.sql": buildRecipeStepCompletionConditionIngredientsQueries(databaseToUse),
		"recipe_prep_tasks.sql":                            buildRecipePrepTasksQueries(databaseToUse),
		"meals.sql":                                        buildMealsQueries(databaseToUse),
		"meal_plans.sql":                                   buildMealPlansQueries(databaseToUse),
		"recipe_step_completion_conditions.sql":            buildRecipeStepCompletionConditionQueries(databaseToUse),
		"meal_plan_option_votes.sql":                       buildMealPlanOptionVotesQueries(databaseToUse),
		"meal_plan_options.sql":                            buildMealPlanOptionsQueries(databaseToUse),
		"meal_plan_tasks.sql":                              buildMealPlanTasksQueries(databaseToUse),
		"recipes.sql":                                      buildRecipesQueries(databaseToUse),
		"recipe_step_ingredients.sql":                      buildRecipeStepIngredientsQueries(databaseToUse),
		"recipe_step_instruments.sql":                      buildRecipeStepInstrumentsQueries(databaseToUse),
		"recipe_step_products.sql":                         buildRecipeStepProductsQueries(databaseToUse),
		"recipe_steps.sql":                                 buildRecipeStepsQueries(databaseToUse),
		"recipe_step_vessels.sql":                          buildRecipeStepVesselsQueries(databaseToUse),
		"user_ingredient_preferences.sql":                  buildUserIngredientPreferencesQueries(databaseToUse),
		"meal_plan_grocery_list_items.sql":                 buildMealPlanGroceryListItemsQueries(databaseToUse),
		"audit_logs.sql":                                   buildAuditLogEntryQueries(databaseToUse),
		"maintenance.sql":                                  buildMaintenanceQueries(databaseToUse),
	}

	checkOnly := *checkOnlyFlag

	for filePath, queries := range queryOutput {
		localFilePath := path.Join("internal", "database", postgres, "sqlc_queries", filePath)
		existingFile, err := os.ReadFile(localFilePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if _, err = os.Create(localFilePath); err != nil {
					log.Fatal(err)
				}
			}
			if err != nil {
				log.Fatal(err)
			}
		}

		var fileContent string
		for i, query := range queries {
			if i != 0 {
				fileContent += "\n"
			}
			fileContent += query.Render()
		}

		fileOutput := ""
		for line := range strings.SplitSeq(strings.TrimSpace(fileContent), "\n") {
			fileOutput += strings.TrimSuffix(line, " ") + "\n"
		}

		if string(existingFile) != fileOutput && checkOnly {
			runErrors = multierror.Append(runErrors, fmt.Errorf("files don't match: %s", filePath))
		}

		if !checkOnly {
			if err = os.WriteFile(localFilePath, []byte(fileOutput), 0o600); err != nil {
				runErrors = multierror.Append(runErrors, err)
			}
		}
	}

	if runErrors.ErrorOrNil() != nil {
		log.Fatal(runErrors)
	}
}
