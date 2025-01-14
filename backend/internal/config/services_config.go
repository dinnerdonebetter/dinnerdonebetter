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
	mealplanning "github.com/dinnerdonebetter/backend/internal/services/eating/meal_planning"
	recipemanagement "github.com/dinnerdonebetter/backend/internal/services/eating/recipe_management"
	validenumerations "github.com/dinnerdonebetter/backend/internal/services/eating/valid_enumerations"

	"github.com/hashicorp/go-multierror"
)

type (
	// ServicesConfig collects the various service configurations.
	ServicesConfig struct {
		_ struct{} `json:"-"`

		AuditLogEntries              auditlogentries.Config              `envPrefix:"AUDIT_LOG_ENTRIES_"              json:"auditLogEntries,omitempty"`
		ServiceSettingConfigurations servicesettingconfigurations.Config `envPrefix:"SERVICE_SETTING_CONFIGURATIONS_" json:"serviceSettingConfigurations,omitempty"`
		UserNotifications            usernotifications.Config            `envPrefix:"USER_NOTIFICATIONS_"             json:"userNotifications,omitempty"`
		Households                   households.Config                   `envPrefix:"HOUSEHOLDS_"                     json:"households,omitempty"`
		ServiceSettings              servicesettings.Config              `envPrefix:"SERVICE_SETTINGS_"               json:"serviceSettings,omitempty"`
		Workers                      workers.Config                      `envPrefix:"WORKERS_"                        json:"workers,omitempty"`
		Users                        users.Config                        `envPrefix:"USERS_"                          json:"users,omitempty"`
		DataPrivacy                  dataprivacy.Config                  `envPrefix:"DATA_PRIVACY_"                   json:"dataPrivacy,omitempty"`
		Recipes                      recipemanagement.Config             `envPrefix:"RECIPES_"                        json:"recipes,omitempty"`
		Auth                         authentication.Config               `envPrefix:"AUTH_"                           json:"auth,omitempty"`
		OAuth2Clients                oauth2clients.Config                `envPrefix:"OAUTH2_CLIENTS_"                 json:"oauth2Clients,omitempty"`
		MealPlanning                 mealplanning.Config                 `envPrefix:"MEAL_PLANNING_"                  json:"meals,omitempty"`
		Webhooks                     webhooks.Config                     `envPrefix:"WEBHOOKS_"                       json:"webhooks,omitempty"`
		HouseholdInvitations         householdinvitations.Config         `envPrefix:"HOUSEHOLD_INVITATIONS_"          json:"householdInvitations,omitempty"`
		ValidEnumerations            validenumerations.Config            `envPrefix:"VALID_ENUMERATIONS_"             json:"validEnumerations,omitempty"`
	}
)

// ValidateWithContext validates a APIServiceConfig struct.
func (cfg *ServicesConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validatorsToRun := map[string]func(context.Context) error{
		"AuditLogEntries":              cfg.AuditLogEntries.ValidateWithContext,
		"ServiceSettingConfigurations": cfg.ServiceSettingConfigurations.ValidateWithContext,
		"UserNotifications":            cfg.UserNotifications.ValidateWithContext,
		"Households":                   cfg.Households.ValidateWithContext,
		"ServiceSettings":              cfg.ServiceSettings.ValidateWithContext,
		"Workers":                      cfg.Workers.ValidateWithContext,
		"Users":                        cfg.Users.ValidateWithContext,
		"DataPrivacy":                  cfg.DataPrivacy.ValidateWithContext,
		"Recipes":                      cfg.Recipes.ValidateWithContext,
		"Auth":                         cfg.Auth.ValidateWithContext,
		"OAuth2Clients":                cfg.OAuth2Clients.ValidateWithContext,
		"MealPlanning":                 cfg.MealPlanning.ValidateWithContext,
		"Webhooks":                     cfg.Webhooks.ValidateWithContext,
		"HouseholdInvitations":         cfg.HouseholdInvitations.ValidateWithContext,
		"ValidEnumerations":            cfg.ValidEnumerations.ValidateWithContext,
	}

	for name, validator := range validatorsToRun {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}
