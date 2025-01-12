package config

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/services/core/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/services/core/authentication"
	"github.com/dinnerdonebetter/backend/internal/services/core/dataprivacy"
	"github.com/dinnerdonebetter/backend/internal/services/core/householdinvitations"
	"github.com/dinnerdonebetter/backend/internal/services/core/households"
	"github.com/dinnerdonebetter/backend/internal/services/core/oauth2clients"
	"github.com/dinnerdonebetter/backend/internal/services/core/servicesettingconfigurations"
	"github.com/dinnerdonebetter/backend/internal/services/core/servicesettings"
	"github.com/dinnerdonebetter/backend/internal/services/core/usernotifications"
	"github.com/dinnerdonebetter/backend/internal/services/core/users"
	"github.com/dinnerdonebetter/backend/internal/services/core/webhooks"
	"github.com/dinnerdonebetter/backend/internal/services/core/workers"
	"github.com/dinnerdonebetter/backend/internal/services/eating/householdinstrumentownerships"
	"github.com/dinnerdonebetter/backend/internal/services/eating/mealplanevents"
	"github.com/dinnerdonebetter/backend/internal/services/eating/mealplangrocerylistitems"
	"github.com/dinnerdonebetter/backend/internal/services/eating/mealplanoptions"
	"github.com/dinnerdonebetter/backend/internal/services/eating/mealplanoptionvotes"
	"github.com/dinnerdonebetter/backend/internal/services/eating/mealplans"
	"github.com/dinnerdonebetter/backend/internal/services/eating/mealplantasks"
	"github.com/dinnerdonebetter/backend/internal/services/eating/meals"
	"github.com/dinnerdonebetter/backend/internal/services/eating/recipepreptasks"
	"github.com/dinnerdonebetter/backend/internal/services/eating/reciperatings"
	"github.com/dinnerdonebetter/backend/internal/services/eating/recipes"
	"github.com/dinnerdonebetter/backend/internal/services/eating/recipestepcompletionconditions"
	"github.com/dinnerdonebetter/backend/internal/services/eating/recipestepingredients"
	"github.com/dinnerdonebetter/backend/internal/services/eating/recipestepinstruments"
	"github.com/dinnerdonebetter/backend/internal/services/eating/recipestepproducts"
	"github.com/dinnerdonebetter/backend/internal/services/eating/recipesteps"
	"github.com/dinnerdonebetter/backend/internal/services/eating/recipestepvessels"
	"github.com/dinnerdonebetter/backend/internal/services/eating/useringredientpreferences"
	"github.com/dinnerdonebetter/backend/internal/services/eating/validingredientgroups"
	"github.com/dinnerdonebetter/backend/internal/services/eating/validingredientmeasurementunits"
	"github.com/dinnerdonebetter/backend/internal/services/eating/validingredientpreparations"
	"github.com/dinnerdonebetter/backend/internal/services/eating/validingredients"
	"github.com/dinnerdonebetter/backend/internal/services/eating/validingredientstateingredients"
	"github.com/dinnerdonebetter/backend/internal/services/eating/validingredientstates"
	"github.com/dinnerdonebetter/backend/internal/services/eating/validinstruments"
	"github.com/dinnerdonebetter/backend/internal/services/eating/validmeasurementunitconversions"
	"github.com/dinnerdonebetter/backend/internal/services/eating/validmeasurementunits"
	"github.com/dinnerdonebetter/backend/internal/services/eating/validpreparationinstruments"
	"github.com/dinnerdonebetter/backend/internal/services/eating/validpreparations"
	"github.com/dinnerdonebetter/backend/internal/services/eating/validpreparationvessels"
	"github.com/dinnerdonebetter/backend/internal/services/eating/validvessels"

	"github.com/hashicorp/go-multierror"
)

type (
	// ServicesConfig collects the various service configurations.
	ServicesConfig struct {
		_ struct{} `json:"-"`

		MealPlanEvents                  mealplanevents.Config                  `envPrefix:"MEAL_PLAN_EVENTS_"                   json:"mealPlanEvents,omitempty"`
		ValidInstrumentMeasurementUnits validingredientmeasurementunits.Config `envPrefix:"VALID_INGREDIENT_MEASUREMENT_UNITS_" json:"validInstrumentMeasurementUnits,omitempty"`
		MealPlanGroceryListItems        mealplangrocerylistitems.Config        `envPrefix:"MEAL_PLAN_GROCERY_LIST_ITEMS_"       json:"mealPlanGroceryListItems,omitempty"`
		ValidIngredientStateIngredients validingredientstateingredients.Config `envPrefix:"VALID_INGREDIENT_STATE_INGREDIENTS_" json:"validIngredientStateIngredients,omitempty"`
		ServiceSettingConfigurations    servicesettingconfigurations.Config    `envPrefix:"SERVICE_SETTING_CONFIGURATIONS_"     json:"serviceSettingConfigurations,omitempty"`
		RecipeRatings                   reciperatings.Config                   `envPrefix:"RECIPE_RATINGS_"                     json:"recipeRatings,omitempty"`
		ValidMeasurementUnitConversions validmeasurementunitconversions.Config `envPrefix:"VALID_MEASUREMENT_UNIT_CONVERSIONS_" json:"validMeasurementUnitConversions,omitempty"`
		ValidIngredientGroups           validingredientgroups.Config           `envPrefix:"VALID_INGREDIENT_GROUPS_"            json:"validIngredientGroups,omitempty"`
		ServiceSettings                 servicesettings.Config                 `envPrefix:"SERVICE_SETTINGS_"                   json:"serviceSettings,omitempty"`
		MealPlanTasks                   mealplantasks.Config                   `envPrefix:"MEAL_PLAN_TASKS_"                    json:"mealPlanTasks,omitempty"`
		RecipeStepProducts              recipestepproducts.Config              `envPrefix:"RECIPE_STEP_PRODUCTS_"               json:"recipeStepProducts,omitempty"`
		RecipeStepIngredients           recipestepingredients.Config           `envPrefix:"RECIPE_STEP_INGREDIENTS_"            json:"recipeStepIngredients,omitempty"`
		RecipePrepTasks                 recipepreptasks.Config                 `envPrefix:"RECIPE_PREP_TASKS_"                  json:"recipePrepTasks,omitempty"`
		RecipeStepCompletionConditions  recipestepcompletionconditions.Config  `envPrefix:"RECIPE_STEP_COMPLETION_CONDITIONS_"  json:"recipeStepCompletionConditions,omitempty"`
		UserIngredientPreferences       useringredientpreferences.Config       `envPrefix:"USER_INGREDIENT_PREFERENCES_"        json:"userIngredientPreferences,omitempty"`
		Households                      households.Config                      `envPrefix:"HOUSEHOLDS_"                         json:"households,omitempty"`
		MealPlans                       mealplans.Config                       `envPrefix:"MEAL_PLANS_"                         json:"mealPlans,omitempty"`
		ValidPreparationInstruments     validpreparationinstruments.Config     `envPrefix:"VALID_PREPARATION_INSTRUMENTS_"      json:"validPreparationInstruments,omitempty"`
		ValidIngredientPreparations     validingredientpreparations.Config     `envPrefix:"VALID_INGREDIENT_PREPARATIONS_"      json:"validIngredientPreparations,omitempty"`
		RecipeStepInstruments           recipestepinstruments.Config           `envPrefix:"RECIPE_STEP_INSTRUMENTS_"            json:"recipeStepInstruments,omitempty"`
		AuditLogEntries                 auditlogentries.Config                 `envPrefix:"AUDIT_LOG_ENTRIES_"                  json:"auditLogEntries,omitempty"`
		HouseholdInstrumentOwnerships   householdinstrumentownerships.Config   `envPrefix:"HOUSEHOLD_INSTRUMENT_OWNERSHIPS_"    json:"householdInstrumentOwnerships,omitempty"`
		MealPlanOptionVotes             mealplanoptionvotes.Config             `envPrefix:"MEAL_PLAN_OPTION_VOTES_"             json:"mealPlanOptionVotes,omitempty"`
		RecipeStepVessels               recipestepvessels.Config               `envPrefix:"RECIPE_STEP_VESSELS_"                json:"recipeStepVessels,omitempty"`
		ValidPreparationVessels         validpreparationvessels.Config         `envPrefix:"VALID_PREPARATION_VESSELS_"          json:"validPreparationVessels,omitempty"`
		MealPlanOptions                 mealplanoptions.Config                 `envPrefix:"MEAL_PLAN_OPTIONS_"                  json:"mealPlanOptions,omitempty"`
		UserNotifications               usernotifications.Config               `envPrefix:"USER_NOTIFICATIONS_"                 json:"userNotifications,omitempty"`
		Workers                         workers.Config                         `envPrefix:"WORKERS_"                            json:"workers,omitempty"`
		RecipeSteps                     recipesteps.Config                     `envPrefix:"RECIPE_STEPS_"                       json:"recipeSteps,omitempty"`
		Users                           users.Config                           `envPrefix:"USERS_"                              json:"users,omitempty"`
		DataPrivacy                     dataprivacy.Config                     `envPrefix:"DATA_PRIVACY_"                       json:"dataPrivacy,omitempty"`
		Recipes                         recipes.Config                         `envPrefix:"RECIPES_"                            json:"recipes,omitempty"`
		Auth                            authentication.Config                  `envPrefix:"AUTH_"                               json:"auth,omitempty"`
		ValidInstruments                validinstruments.Config                `envPrefix:"VALID_INSTRUMENTS_"                  json:"validInstruments,omitempty"`
		Meals                           meals.Config                           `envPrefix:"MEALS_"                              json:"meals,omitempty"`
		ValidMeasurementUnits           validmeasurementunits.Config           `envPrefix:"VALID_MEASUREMENT_UNITS_"            json:"validMeasurementUnits,omitempty"`
		OAuth2Clients                   oauth2clients.Config                   `envPrefix:"OAUTH2_CLIENTS_"                     json:"oauth2Clients,omitempty"`
		ValidVessels                    validvessels.Config                    `envPrefix:"VALID_VESSELS_"                      json:"validVessels,omitempty"`
		HouseholdInvitations            householdinvitations.Config            `envPrefix:"HOUSEHOLD_INVITATIONS_"              json:"householdInvitations,omitempty"`
		ValidPreparations               validpreparations.Config               `envPrefix:"VALID_PREPARATIONS_"                 json:"validPreparations,omitempty"`
		ValidIngredientStates           validingredientstates.Config           `envPrefix:"VALID_INGREDIENT_STATES_"            json:"validIngredientStates,omitempty"`
		ValidIngredients                validingredients.Config                `envPrefix:"VALID_INGREDIENTS_"                  json:"validIngredients,omitempty"`
		Webhooks                        webhooks.Config                        `envPrefix:"WEBHOOKS_"                           json:"webhooks,omitempty"`
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
