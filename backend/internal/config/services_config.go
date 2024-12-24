package config

import (
	"context"
	"fmt"

	auditlogentriesservice "github.com/dinnerdonebetter/backend/internal/services/auditlogentries"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	dataprivacyservice "github.com/dinnerdonebetter/backend/internal/services/dataprivacy"
	householdinstrumentownershipsservice "github.com/dinnerdonebetter/backend/internal/services/householdinstrumentownerships"
	householdinvitationsservice "github.com/dinnerdonebetter/backend/internal/services/householdinvitations"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/households"
	mealplaneventsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanevents"
	"github.com/dinnerdonebetter/backend/internal/services/mealplangrocerylistitems"
	mealplanoptionsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/mealplans"
	"github.com/dinnerdonebetter/backend/internal/services/mealplantasks"
	mealsservice "github.com/dinnerdonebetter/backend/internal/services/meals"
	oauth2clientsservice "github.com/dinnerdonebetter/backend/internal/services/oauth2clients"
	"github.com/dinnerdonebetter/backend/internal/services/recipepreptasks"
	reciperatingsservice "github.com/dinnerdonebetter/backend/internal/services/reciperatings"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/recipes"
	recipestepcompletionconditionsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepcompletionconditions"
	recipestepingredientsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepproducts"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/recipesteps"
	recipestepvesselsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepvessels"
	servicesettingconfigurationsservice "github.com/dinnerdonebetter/backend/internal/services/servicesettingconfigurations"
	servicesettingsservice "github.com/dinnerdonebetter/backend/internal/services/servicesettings"
	useringredientpreferencesservice "github.com/dinnerdonebetter/backend/internal/services/useringredientpreferences"
	usernotificationsservice "github.com/dinnerdonebetter/backend/internal/services/usernotifications"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/users"
	validingredientgroupsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientgroups"
	validingredientmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientmeasurementunits"
	validingredientpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientpreparations"
	validingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredients"
	"github.com/dinnerdonebetter/backend/internal/services/validingredientstateingredients"
	"github.com/dinnerdonebetter/backend/internal/services/validingredientstates"
	validinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validinstruments"
	"github.com/dinnerdonebetter/backend/internal/services/validmeasurementunitconversions"
	validmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validmeasurementunits"
	validpreparationinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparationinstruments"
	validpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparations"
	"github.com/dinnerdonebetter/backend/internal/services/validpreparationvessels"
	validvesselsservice "github.com/dinnerdonebetter/backend/internal/services/validvessels"
	webhooksservice "github.com/dinnerdonebetter/backend/internal/services/webhooks"
	workersservice "github.com/dinnerdonebetter/backend/internal/services/workers"

	"github.com/hashicorp/go-multierror"
)

type (
	// ServicesConfig collects the various service configurations.
	ServicesConfig struct {
		_ struct{} `json:"-"`

		AuditLogEntries                 auditlogentriesservice.Config                 `envPrefix:"AUDIT_LOG_ENTRIES_"                  json:"auditLogEntries,omitempty"`
		MealPlanGroceryListItems        mealplangrocerylistitems.Config               `envPrefix:"MEAL_PLAN_GROCERY_LIST_ITEMS_"       json:"mealPlanGroceryListItems,omitempty"`
		ValidInstrumentMeasurementUnits validingredientmeasurementunitsservice.Config `envPrefix:"VALID_INGREDIENT_MEASUREMENT_UNITS_" json:"validInstrumentMeasurementUnits,omitempty"`
		ServiceSettingConfigurations    servicesettingconfigurationsservice.Config    `envPrefix:"SERVICE_SETTING_CONFIGURATIONS_"     json:"serviceSettingConfigurations,omitempty"`
		RecipeRatings                   reciperatingsservice.Config                   `envPrefix:"RECIPE_RATINGS_"                     json:"recipeRatings,omitempty"`
		ValidMeasurementUnitConversions validmeasurementunitconversions.Config        `envPrefix:"VALID_MEASUREMENT_UNIT_CONVERSIONS_" json:"validMeasurementUnitConversions,omitempty"`
		ValidIngredientGroups           validingredientgroupsservice.Config           `envPrefix:"VALID_INGREDIENT_GROUPS_"            json:"validIngredientGroups,omitempty"`
		ServiceSettings                 servicesettingsservice.Config                 `envPrefix:"SERVICE_SETTINGS_"                   json:"serviceSettings,omitempty"`
		MealPlanTasks                   mealplantasks.Config                          `envPrefix:"MEAL_PLAN_TASKS_"                    json:"mealPlanTasks,omitempty"`
		RecipeStepInstruments           recipestepinstrumentsservice.Config           `envPrefix:"RECIPE_STEP_INSTRUMENTS_"            json:"recipeStepInstruments,omitempty"`
		RecipeStepIngredients           recipestepingredientsservice.Config           `envPrefix:"RECIPE_STEP_INGREDIENTS_"            json:"recipeStepIngredients,omitempty"`
		HouseholdInstrumentOwnerships   householdinstrumentownershipsservice.Config   `envPrefix:"HOUSEHOLD_INSTRUMENT_OWNERSHIPS_"    json:"householdInstrumentOwnerships,omitempty"`
		RecipePrepTasks                 recipepreptasks.Config                        `envPrefix:"RECIPE_PREP_TASKS_"                  json:"recipePrepTasks,omitempty"`
		RecipeStepCompletionConditions  recipestepcompletionconditionsservice.Config  `envPrefix:"RECIPE_STEP_COMPLETION_CONDITIONS_"  json:"recipeStepCompletionConditions,omitempty"`
		UserIngredientPreferences       useringredientpreferencesservice.Config       `envPrefix:"USER_INGREDIENT_PREFERENCES_"        json:"userIngredientPreferences,omitempty"`
		Households                      householdsservice.Config                      `envPrefix:"HOUSEHOLDS_"                         json:"households,omitempty"`
		MealPlans                       mealplansservice.Config                       `envPrefix:"MEAL_PLANS_"                         json:"mealPlans,omitempty"`
		ValidPreparationInstruments     validpreparationinstrumentsservice.Config     `envPrefix:"VALID_PREPARATION_INSTRUMENTS_"      json:"validPreparationInstruments,omitempty"`
		ValidIngredientPreparations     validingredientpreparationsservice.Config     `envPrefix:"VALID_INGREDIENT_PREPARATIONS_"      json:"validIngredientPreparations,omitempty"`
		RecipeStepProducts              recipestepproductsservice.Config              `envPrefix:"RECIPE_STEP_PRODUCTS_"               json:"recipeStepProducts,omitempty"`
		ValidIngredientStateIngredients validingredientstateingredients.Config        `envPrefix:"VALID_INGREDIENT_STATE_INGREDIENTS_" json:"validIngredientStateIngredients,omitempty"`
		MealPlanEvents                  mealplaneventsservice.Config                  `envPrefix:"MEAL_PLAN_EVENTS_"                   json:"mealPlanEvents,omitempty"`
		MealPlanOptionVotes             mealplanoptionvotesservice.Config             `envPrefix:"MEAL_PLAN_OPTION_VOTES_"             json:"mealPlanOptionVotes,omitempty"`
		RecipeStepVessels               recipestepvesselsservice.Config               `envPrefix:"RECIPE_STEP_VESSELS_"                json:"recipeStepVessels,omitempty"`
		ValidPreparationVessels         validpreparationvessels.Config                `envPrefix:"VALID_PREPARATION_VESSELS_"          json:"validPreparationVessels,omitempty"`
		Workers                         workersservice.Config                         `envPrefix:"WORKERS_"                            json:"workers,omitempty"`
		UserNotifications               usernotificationsservice.Config               `envPrefix:"USER_NOTIFICATIONS_"                 json:"userNotifications,omitempty"`
		MealPlanOptions                 mealplanoptionsservice.Config                 `envPrefix:"MEAL_PLAN_OPTIONS_"                  json:"mealPlanOptions,omitempty"`
		DataPrivacy                     dataprivacyservice.Config                     `envPrefix:"DATA_PRIVACY_"                       json:"dataPrivacy,omitempty"`
		RecipeSteps                     recipestepsservice.Config                     `envPrefix:"RECIPE_STEPS_"                       json:"recipeSteps,omitempty"`
		Users                           usersservice.Config                           `envPrefix:"USERS_"                              json:"users,omitempty"`
		ValidInstruments                validinstrumentsservice.Config                `envPrefix:"VALID_INSTRUMENTS_"                  json:"validInstruments,omitempty"`
		ValidMeasurementUnits           validmeasurementunitsservice.Config           `envPrefix:"VALID_MEASUREMENT_UNITS_"            json:"validMeasurementUnits,omitempty"`
		OAuth2Clients                   oauth2clientsservice.Config                   `envPrefix:"OAUTH2_CLIENTS_"                     json:"oauth2Clients,omitempty"`
		Webhooks                        webhooksservice.Config                        `envPrefix:"WEBHOOKS_"                           json:"webhooks,omitempty"`
		ValidIngredients                validingredientsservice.Config                `envPrefix:"VALID_INGREDIENTS_"                  json:"validIngredients,omitempty"`
		Meals                           mealsservice.Config                           `envPrefix:"MEALS_"                              json:"meals,omitempty"`
		ValidVessels                    validvesselsservice.Config                    `envPrefix:"VALID_VESSELS_"                      json:"validVessels,omitempty"`
		HouseholdInvitations            householdinvitationsservice.Config            `envPrefix:"HOUSEHOLD_INVITATIONS_"              json:"householdInvitations,omitempty"`
		ValidPreparations               validpreparationsservice.Config               `envPrefix:"VALID_PREPARATIONS_"                 json:"validPreparations,omitempty"`
		ValidIngredientStates           validingredientstates.Config                  `envPrefix:"VALID_INGREDIENT_STATES_"            json:"validIngredientStates,omitempty"`
		Recipes                         recipesservice.Config                         `envPrefix:"RECIPES_"                            json:"recipes,omitempty"`
		Auth                            authservice.Config                            `envPrefix:"AUTH_"                               json:"auth,omitempty"`
	}
)

// ValidateWithContext validates a APIServiceConfig struct.
func (cfg *ServicesConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validatorsToRun := map[string]func(context.Context) error{
		"Auth":                            cfg.Auth.ValidateWithContext,
		"Users":                           cfg.Users.ValidateWithContext,
		"Webhooks":                        cfg.Webhooks.ValidateWithContext,
		"ValidInstruments":                cfg.ValidInstruments.ValidateWithContext,
		"ValidIngredients":                cfg.ValidIngredients.ValidateWithContext,
		"ValidIngredientGroups":           cfg.ValidIngredientGroups.ValidateWithContext,
		"ValidPreparations":               cfg.ValidPreparations.ValidateWithContext,
		"ValidMeasurementUnits":           cfg.ValidMeasurementUnits.ValidateWithContext,
		"ValidIngredientPreparations":     cfg.ValidIngredientPreparations.ValidateWithContext,
		"ValidIngredientStateIngredients": cfg.ValidIngredientStateIngredients.ValidateWithContext,
		"ValidPreparationInstruments":     cfg.ValidPreparationInstruments.ValidateWithContext,
		"ValidInstrumentMeasurementUnits": cfg.ValidInstrumentMeasurementUnits.ValidateWithContext,
		"Recipes":                         cfg.Recipes.ValidateWithContext,
		"RecipeSteps":                     cfg.RecipeSteps.ValidateWithContext,
		"RecipeStepInstruments":           cfg.RecipeStepInstruments.ValidateWithContext,
		"RecipeStepVessels":               cfg.RecipeStepVessels.ValidateWithContext,
		"RecipeStepIngredients":           cfg.RecipeStepIngredients.ValidateWithContext,
		"RecipeStepCompletionConditions":  cfg.RecipeStepCompletionConditions.ValidateWithContext,
		"MealPlans":                       cfg.MealPlans.ValidateWithContext,
		"MealPlanEvents":                  cfg.MealPlanEvents.ValidateWithContext,
		"MealPlanOptions":                 cfg.MealPlanOptions.ValidateWithContext,
		"MealPlanOptionVotes":             cfg.MealPlanOptionVotes.ValidateWithContext,
		"RecipePrepTasks":                 cfg.RecipePrepTasks.ValidateWithContext,
		"MealPlanGroceryListItems":        cfg.MealPlanGroceryListItems.ValidateWithContext,
		"ValidMeasurementUnitConversions": cfg.ValidMeasurementUnitConversions.ValidateWithContext,
		"ValidIngredientStates":           cfg.ValidIngredientStates.ValidateWithContext,
		"ServiceSettings":                 cfg.ServiceSettings.ValidateWithContext,
		"ServiceSettingConfigurations":    cfg.ServiceSettingConfigurations.ValidateWithContext,
		"UserIngredientPreferences":       cfg.UserIngredientPreferences.ValidateWithContext,
		"ValidVessels":                    cfg.ValidVessels.ValidateWithContext,
		"ValidPreparationVessels":         cfg.ValidPreparationVessels.ValidateWithContext,
		"Workers":                         cfg.Workers.ValidateWithContext,
		"UserNotifications":               cfg.UserNotifications.ValidateWithContext,
		"AuditLogEntries":                 cfg.AuditLogEntries.ValidateWithContext,
		"DataPrivacy":                     cfg.DataPrivacy.ValidateWithContext,
	}

	for name, validator := range validatorsToRun {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}
