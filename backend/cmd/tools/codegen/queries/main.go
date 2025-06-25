package main

import (
	"errors"
	"fmt"
	"log"
	"os"
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
		"internal/database/postgres/sqlc_queries/webhooks.sql":                                     buildWebhooksQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/user_notifications.sql":                           buildUserNotificationQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/webhook_trigger_events.sql":                       buildWebhookTriggerEventsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/oauth2_client_tokens.sql":                         buildOAuth2ClientTokensQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/oauth2_clients.sql":                               buildOAuth2ClientsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/valid_ingredients.sql":                            buildValidIngredientsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/valid_instruments.sql":                            buildValidInstrumentsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/valid_preparations.sql":                           buildValidPreparationsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/valid_measurement_units.sql":                      buildValidMeasurementUnitsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/valid_ingredient_states.sql":                      buildValidIngredientStatesQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/valid_vessels.sql":                                buildValidVesselsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/valid_ingredient_groups.sql":                      buildValidIngredientGroupsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/valid_ingredient_preparations.sql":                buildValidIngredientPreparationsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/valid_preparation_vessels.sql":                    buildValidPreparationVesselsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/valid_ingredient_measurement_units.sql":           buildValidIngredientMeasurementUnitsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/valid_measurement_unit_conversions.sql":           buildValidMeasurementUnitConversionsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/valid_ingredient_state_ingredients.sql":           buildValidIngredientStateIngredientsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/valid_preparation_instruments.sql":                buildValidPreparationInstrumentsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/account_instrument_ownerships.sql":                buildAccountInstrumentOwnershipQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/meal_components.sql":                              buildMealComponentsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/meal_plan_events.sql":                             buildMealPlanEventsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/recipe_media.sql":                                 buildRecipeMediaQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/recipe_prep_task_steps.sql":                       buildRecipePrepTaskStepsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/recipe_ratings.sql":                               buildRecipeRatingsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/recipe_step_completion_condition_ingredients.sql": buildRecipeStepCompletionConditionIngredientsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/recipe_prep_tasks.sql":                            buildRecipePrepTasksQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/meals.sql":                                        buildMealsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/meal_plans.sql":                                   buildMealPlansQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/recipe_step_completion_conditions.sql":            buildRecipeStepCompletionConditionQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/meal_plan_option_votes.sql":                       buildMealPlanOptionVotesQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/meal_plan_options.sql":                            buildMealPlanOptionsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/meal_plan_tasks.sql":                              buildMealPlanTasksQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/recipes.sql":                                      buildRecipesQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/recipe_step_ingredients.sql":                      buildRecipeStepIngredientsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/recipe_step_instruments.sql":                      buildRecipeStepInstrumentsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/recipe_step_products.sql":                         buildRecipeStepProductsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/recipe_steps.sql":                                 buildRecipeStepsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/recipe_step_vessels.sql":                          buildRecipeStepVesselsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/user_ingredient_preferences.sql":                  buildUserIngredientPreferencesQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/meal_plan_grocery_list_items.sql":                 buildMealPlanGroceryListItemsQueries(databaseToUse),
		"internal/database/postgres/sqlc_queries/maintenance.sql":                                  buildMaintenanceQueries(databaseToUse),
		// moved files
		"internal/platform/database/postgres/implementations/identity/sqlc_queries/account_invitations.sql":            buildAccountInvitationsQueries(databaseToUse),
		"internal/platform/database/postgres/implementations/identity/sqlc_queries/account_user_memberships.sql":       buildAccountUserMembershipsQueries(databaseToUse),
		"internal/platform/database/postgres/implementations/identity/sqlc_queries/accounts.sql":                       buildAccountsQueries(databaseToUse),
		"internal/platform/database/postgres/implementations/identity/sqlc_queries/audit_logs.sql":                     buildAuditLogEntryQueries(databaseToUse),
		"internal/platform/database/postgres/implementations/identity/sqlc_queries/admin.sql":                          buildAdminQueries(databaseToUse),
		"internal/platform/database/postgres/implementations/identity/sqlc_queries/password_reset_tokens.sql":          buildPasswordResetTokensQueries(databaseToUse),
		"internal/platform/database/postgres/implementations/identity/sqlc_queries/users.sql":                          buildUsersQueries(databaseToUse),
		"internal/platform/database/postgres/implementations/settings/sqlc_queries/service_settings.sql":               buildServiceSettingQueries(databaseToUse),
		"internal/platform/database/postgres/implementations/settings/sqlc_queries/service_setting_configurations.sql": buildServiceSettingConfigurationQueries(databaseToUse),
	}

	checkOnly := *checkOnlyFlag

	for filePath, queries := range queryOutput {
		existingFile, err := os.ReadFile(filePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if _, err = os.Create(filePath); err != nil {
					log.Fatal(fmt.Errorf("creating file: %w", err))
				}
			}
			if err != nil {
				log.Fatal(fmt.Errorf("opening existing file: %w", err))
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
		for _, line := range strings.Split(strings.TrimSpace(fileContent), "\n") {
			fileOutput += strings.TrimSuffix(line, " ") + "\n"
		}

		if string(existingFile) != fileOutput && checkOnly {
			runErrors = multierror.Append(runErrors, fmt.Errorf("files don't match: %s", filePath))
		}

		if !checkOnly {
			if err = os.WriteFile(filePath, []byte(fileOutput), 0o600); err != nil {
				runErrors = multierror.Append(runErrors, err)
			}
		}
	}

	if runErrors.ErrorOrNil() != nil {
		log.Fatal(runErrors)
	}
}
