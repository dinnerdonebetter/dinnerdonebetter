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
	recipemanagement "github.com/dinnerdonebetter/backend/internal/services/eating/recipe_management"
	"github.com/dinnerdonebetter/backend/internal/services/eating/useringredientpreferences"
	validenumerations "github.com/dinnerdonebetter/backend/internal/services/eating/valid_enumerations"

	"github.com/hashicorp/go-multierror"
)

type (
	// ServicesConfig collects the various service configurations.
	ServicesConfig struct {
		_ struct{} `json:"-"`

		AuditLogEntries               auditlogentries.Config               `envPrefix:"AUDIT_LOG_ENTRIES_"               json:"auditLogEntries,omitempty"`
		MealPlans                     mealplans.Config                     `envPrefix:"MEAL_PLANS_"                      json:"mealPlans,omitempty"`
		MealPlanGroceryListItems      mealplangrocerylistitems.Config      `envPrefix:"MEAL_PLAN_GROCERY_LIST_ITEMS_"    json:"mealPlanGroceryListItems,omitempty"`
		ServiceSettingConfigurations  servicesettingconfigurations.Config  `envPrefix:"SERVICE_SETTING_CONFIGURATIONS_"  json:"serviceSettingConfigurations,omitempty"`
		MealPlanTasks                 mealplantasks.Config                 `envPrefix:"MEAL_PLAN_TASKS_"                 json:"mealPlanTasks,omitempty"`
		UserNotifications             usernotifications.Config             `envPrefix:"USER_NOTIFICATIONS_"              json:"userNotifications,omitempty"`
		UserIngredientPreferences     useringredientpreferences.Config     `envPrefix:"USER_INGREDIENT_PREFERENCES_"     json:"userIngredientPreferences,omitempty"`
		Households                    households.Config                    `envPrefix:"HOUSEHOLDS_"                      json:"households,omitempty"`
		MealPlanEvents                mealplanevents.Config                `envPrefix:"MEAL_PLAN_EVENTS_"                json:"mealPlanEvents,omitempty"`
		ServiceSettings               servicesettings.Config               `envPrefix:"SERVICE_SETTINGS_"                json:"serviceSettings,omitempty"`
		HouseholdInstrumentOwnerships householdinstrumentownerships.Config `envPrefix:"HOUSEHOLD_INSTRUMENT_OWNERSHIPS_" json:"householdInstrumentOwnerships,omitempty"`
		MealPlanOptionVotes           mealplanoptionvotes.Config           `envPrefix:"MEAL_PLAN_OPTION_VOTES_"          json:"mealPlanOptionVotes,omitempty"`
		MealPlanOptions               mealplanoptions.Config               `envPrefix:"MEAL_PLAN_OPTIONS_"               json:"mealPlanOptions,omitempty"`
		Workers                       workers.Config                       `envPrefix:"WORKERS_"                         json:"workers,omitempty"`
		Users                         users.Config                         `envPrefix:"USERS_"                           json:"users,omitempty"`
		DataPrivacy                   dataprivacy.Config                   `envPrefix:"DATA_PRIVACY_"                    json:"dataPrivacy,omitempty"`
		Recipes                       recipemanagement.Config              `envPrefix:"RECIPES_"                         json:"recipes,omitempty"`
		Auth                          authentication.Config                `envPrefix:"AUTH_"                            json:"auth,omitempty"`
		OAuth2Clients                 oauth2clients.Config                 `envPrefix:"OAUTH2_CLIENTS_"                  json:"oauth2Clients,omitempty"`
		Meals                         meals.Config                         `envPrefix:"MEALS_"                           json:"meals,omitempty"`
		Webhooks                      webhooks.Config                      `envPrefix:"WEBHOOKS_"                        json:"webhooks,omitempty"`
		HouseholdInvitations          householdinvitations.Config          `envPrefix:"HOUSEHOLD_INVITATIONS_"           json:"householdInvitations,omitempty"`
		ValidEnumerations             validenumerations.Config             `envPrefix:"VALID_ENUMERATIONS_"              json:"validEnumerations,omitempty"`
	}
)

// ValidateWithContext validates a APIServiceConfig struct.
func (cfg *ServicesConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validatorsToRun := map[string]func(context.Context) error{
		"Auth":                         cfg.Auth.ValidateWithContext,
		"Users":                        cfg.Users.ValidateWithContext,
		"Webhooks":                     cfg.Webhooks.ValidateWithContext,
		"ValidEnumerations":            cfg.ValidEnumerations.ValidateWithContext,
		"Recipes":                      cfg.Recipes.ValidateWithContext,
		"MealPlans":                    cfg.MealPlans.ValidateWithContext,
		"MealPlanEvents":               cfg.MealPlanEvents.ValidateWithContext,
		"MealPlanOptions":              cfg.MealPlanOptions.ValidateWithContext,
		"MealPlanOptionVotes":          cfg.MealPlanOptionVotes.ValidateWithContext,
		"MealPlanGroceryListItems":     cfg.MealPlanGroceryListItems.ValidateWithContext,
		"ServiceSettings":              cfg.ServiceSettings.ValidateWithContext,
		"ServiceSettingConfigurations": cfg.ServiceSettingConfigurations.ValidateWithContext,
		"UserIngredientPreferences":    cfg.UserIngredientPreferences.ValidateWithContext,
		"Workers":                      cfg.Workers.ValidateWithContext,
		"UserNotifications":            cfg.UserNotifications.ValidateWithContext,
		"AuditLogEntries":              cfg.AuditLogEntries.ValidateWithContext,
		"DataPrivacy":                  cfg.DataPrivacy.ValidateWithContext,
	}

	for name, validator := range validatorsToRun {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}
