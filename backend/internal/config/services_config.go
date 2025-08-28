package config

import (
	"context"
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/services/auth/handlers/authentication"

	dataprivacycfg "github.com/dinnerdonebetter/backend/internal/services/dataprivacy/config"
	identitycfg "github.com/dinnerdonebetter/backend/internal/services/identity/config"
	mealplanningcfg "github.com/dinnerdonebetter/backend/internal/services/mealplanning/config"
	oauthcfg "github.com/dinnerdonebetter/backend/internal/services/oauth/config"

	"github.com/hashicorp/go-multierror"
)

type (
	// ServicesConfig collects the various service configurations.
	ServicesConfig struct {
		_ struct{} `json:"-"`

		Users         identitycfg.Config     `envPrefix:"USERS_"                          json:"users,omitempty"`
		DataPrivacy   dataprivacycfg.Config  `envPrefix:"DATA_PRIVACY_"                   json:"dataPrivacy,omitempty"`
		MealPlanning  mealplanningcfg.Config `envPrefix:"MEAL_PLANNING_"                  json:"mealPlanning,omitempty"`
		Auth          authentication.Config  `envPrefix:"AUTH_"                           json:"auth,omitempty"`
		OAuth2Clients oauthcfg.Config        `envPrefix:"OAUTH2_CLIENTS_"                 json:"oauth2Clients,omitempty"`
	}
)

// ValidateWithContext validates a APIServiceConfig struct.
func (cfg *ServicesConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validatorsToRun := map[string]func(context.Context) error{
		"Users":         cfg.Users.ValidateWithContext,
		"DataPrivacy":   cfg.DataPrivacy.ValidateWithContext,
		"MealPlanning":  cfg.MealPlanning.ValidateWithContext,
		"OAuth2Clients": cfg.OAuth2Clients.ValidateWithContext,
	}

	for name, validator := range validatorsToRun {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}
