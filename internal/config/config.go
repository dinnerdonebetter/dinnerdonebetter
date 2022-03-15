package config

import (
	"context"
	"fmt"
	"os"

	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	dbconfig "github.com/prixfixeco/api_server/internal/database/config"
	emailconfig "github.com/prixfixeco/api_server/internal/email/config"
	"github.com/prixfixeco/api_server/internal/encoding"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/routing"
	"github.com/prixfixeco/api_server/internal/server"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	householdinvitationsservice "github.com/prixfixeco/api_server/internal/services/householdinvitations"
	householdsservice "github.com/prixfixeco/api_server/internal/services/households"
	mealplanoptionsservice "github.com/prixfixeco/api_server/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/prixfixeco/api_server/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/prixfixeco/api_server/internal/services/mealplans"
	mealsservice "github.com/prixfixeco/api_server/internal/services/meals"
	recipesservice "github.com/prixfixeco/api_server/internal/services/recipes"
	recipestepingredientsservice "github.com/prixfixeco/api_server/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/prixfixeco/api_server/internal/services/recipestepinstruments"
	recipestepsservice "github.com/prixfixeco/api_server/internal/services/recipesteps"
	usersservice "github.com/prixfixeco/api_server/internal/services/users"
	validingredientpreparationsservice "github.com/prixfixeco/api_server/internal/services/validingredientpreparations"
	validingredientsservice "github.com/prixfixeco/api_server/internal/services/validingredients"
	validinstrumentsservice "github.com/prixfixeco/api_server/internal/services/validinstruments"
	validpreparationsservice "github.com/prixfixeco/api_server/internal/services/validpreparations"
	webhooksservice "github.com/prixfixeco/api_server/internal/services/webhooks"
	websocketsservice "github.com/prixfixeco/api_server/internal/services/websockets"
	"github.com/prixfixeco/api_server/internal/uploads"
)

const (
	// DevelopmentRunMode is the run mode for a development environment.
	DevelopmentRunMode runMode = "development"
	// TestingRunMode is the run mode for a testing environment.
	TestingRunMode runMode = "testing"
	// ProductionRunMode is the run mode for a production environment.
	ProductionRunMode runMode = "production"
)

type (
	// runMode describes what method of operation the server is under.
	runMode string

	// CloserFunc calls all io.Closers in the service.
	CloserFunc func()

	// InstanceConfig configures an instance of the service. It is composed of all the other setting structs.
	InstanceConfig struct {
		_             struct{}
		Observability observability.Config      `json:"observability" mapstructure:"observability" toml:"observability,omitempty"`
		Email         emailconfig.Config        `json:"email" mapstructure:"email" toml:"email,omitempty"`
		CustomerData  customerdataconfig.Config `json:"customerData" mapstructure:"customer_data" toml:"customer_data,omitempty"`
		Encoding      encoding.Config           `json:"encoding" mapstructure:"encoding" toml:"encoding,omitempty"`
		Uploads       uploads.Config            `json:"uploads" mapstructure:"uploads" toml:"uploads,omitempty"`
		Routing       routing.Config            `json:"routing" mapstructure:"routing" toml:"routing,omitempty"`
		Database      dbconfig.Config           `json:"database" mapstructure:"database" toml:"database,omitempty"`
		Meta          MetaSettings              `json:"meta" mapstructure:"meta" toml:"meta,omitempty"`
		Events        msgconfig.Config          `json:"events" mapstructure:"events" toml:"events,omitempty"`
		Server        server.Config             `json:"server" mapstructure:"server" toml:"server,omitempty"`
		Services      ServicesConfigurations    `json:"services" mapstructure:"services" toml:"services,omitempty"`
	}

	// ServicesConfigurations collects the various service configurations.
	ServicesConfigurations struct {
		_                           struct{}
		ValidInstruments            validinstrumentsservice.Config            `json:"validInstruments" mapstructure:"valid_instruments" toml:"valid_instruments,omitempty"`
		ValidIngredients            validingredientsservice.Config            `json:"validIngredients" mapstructure:"valid_ingredients" toml:"valid_ingredients,omitempty"`
		ValidPreparations           validpreparationsservice.Config           `json:"validPreparations" mapstructure:"valid_preparations" toml:"valid_preparations,omitempty"`
		MealPlanOptionVotes         mealplanoptionvotesservice.Config         `json:"mealPlanOptionVotes" mapstructure:"meal_plan_option_votes" toml:"meal_plan_option_votes,omitempty"`
		ValidIngredientPreparations validingredientpreparationsservice.Config `json:"validIngredientPreparations" mapstructure:"valid_ingredient_preparations" toml:"valid_ingredient_preparations,omitempty"`
		Meals                       mealsservice.Config                       `json:"meals" mapstructure:"meals" toml:"meals,omitempty"`
		Recipes                     recipesservice.Config                     `json:"recipes" mapstructure:"recipes" toml:"recipes,omitempty"`
		RecipeSteps                 recipestepsservice.Config                 `json:"recipeSteps" mapstructure:"recipe_steps" toml:"recipe_steps,omitempty"`
		RecipeStepInstruments       recipestepinstrumentsservice.Config       `json:"recipeStepInstruments" mapstructure:"recipe_step_instruments" toml:"recipe_step_instruments,omitempty"`
		RecipeStepIngredients       recipestepingredientsservice.Config       `json:"recipeStepIngredients" mapstructure:"recipe_step_ingredients" toml:"recipe_step_ingredients,omitempty"`
		MealPlans                   mealplansservice.Config                   `json:"mealPlans" mapstructure:"meal_plans" toml:"meal_plans,omitempty"`
		MealPlanOptions             mealplanoptionsservice.Config             `json:"mealPlanOptions" mapstructure:"meal_plan_options" toml:"meal_plan_options,omitempty"`
		Households                  householdsservice.Config                  `json:"households" mapstructure:"households" toml:"households,omitempty"`
		HouseholdInvitations        householdinvitationsservice.Config        `json:"householdInvitations" mapstructure:"household_invitations" toml:"household_invitations,omitempty"`
		Websockets                  websocketsservice.Config                  `json:"websockets" mapstructure:"websockets" toml:"websockets,omitempty"`
		Webhooks                    webhooksservice.Config                    `json:"webhooks" mapstructure:"webhooks" toml:"webhooks,omitempty"`
		Users                       usersservice.Config                       `json:"users" mapstructure:"users" toml:"users,omitempty"`
		Auth                        authservice.Config                        `json:"auth" mapstructure:"auth" toml:"auth,omitempty"`
	}
)

// EncodeToFile renders your config to a file given your favorite encoder.
func (cfg *InstanceConfig) EncodeToFile(path string, marshaller func(v interface{}) ([]byte, error)) error {
	byteSlice, err := marshaller(*cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, byteSlice, 0600)
}

// ValidateWithContext validates a InstanceConfig struct.
func (cfg *InstanceConfig) ValidateWithContext(ctx context.Context, validateServices bool) error {
	if err := cfg.Uploads.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Uploads portion of config: %w", err)
	}

	if err := cfg.Routing.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Routing portion of config: %w", err)
	}

	if err := cfg.Meta.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Meta portion of config: %w", err)
	}

	if err := cfg.Encoding.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Encoding portion of config: %w", err)
	}

	if err := cfg.CustomerData.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating CustomerData portion of config: %w", err)
	}

	if err := cfg.Observability.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Observability portion of config: %w", err)
	}

	if err := cfg.Database.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Database portion of config: %w", err)
	}

	if err := cfg.Server.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Server portion of config: %w", err)
	}

	if err := cfg.Email.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Email portion of config: %w", err)
	}

	if validateServices {
		if err := cfg.Services.Auth.ValidateWithContext(ctx); err != nil {
			return fmt.Errorf("error validating Auth service portion of config: %w", err)
		}

		if err := cfg.Services.Webhooks.ValidateWithContext(ctx); err != nil {
			return fmt.Errorf("error validating Webhooks service portion of config: %w", err)
		}

		if err := cfg.Services.ValidInstruments.ValidateWithContext(ctx); err != nil {
			return fmt.Errorf("error validating ValidInstruments service portion of config: %w", err)
		}

		if err := cfg.Services.ValidIngredients.ValidateWithContext(ctx); err != nil {
			return fmt.Errorf("error validating ValidIngredients service portion of config: %w", err)
		}

		if err := cfg.Services.ValidPreparations.ValidateWithContext(ctx); err != nil {
			return fmt.Errorf("error validating ValidPreparations service portion of config: %w", err)
		}

		if err := cfg.Services.ValidIngredientPreparations.ValidateWithContext(ctx); err != nil {
			return fmt.Errorf("error validating ValidIngredientPreparations service portion of config: %w", err)
		}

		if err := cfg.Services.Recipes.ValidateWithContext(ctx); err != nil {
			return fmt.Errorf("error validating Recipes service portion of config: %w", err)
		}

		if err := cfg.Services.RecipeSteps.ValidateWithContext(ctx); err != nil {
			return fmt.Errorf("error validating RecipeSteps service portion of config: %w", err)
		}

		if err := cfg.Services.RecipeStepInstruments.ValidateWithContext(ctx); err != nil {
			return fmt.Errorf("error validating RecipeStepInstruments service portion of config: %w", err)
		}

		if err := cfg.Services.RecipeStepIngredients.ValidateWithContext(ctx); err != nil {
			return fmt.Errorf("error validating RecipeStepIngredients service portion of config: %w", err)
		}

		if err := cfg.Services.MealPlans.ValidateWithContext(ctx); err != nil {
			return fmt.Errorf("error validating MealPlans service portion of config: %w", err)
		}

		if err := cfg.Services.MealPlanOptions.ValidateWithContext(ctx); err != nil {
			return fmt.Errorf("error validating MealPlanOptions service portion of config: %w", err)
		}

		if err := cfg.Services.MealPlanOptionVotes.ValidateWithContext(ctx); err != nil {
			return fmt.Errorf("error validating MealPlanOptionVotes service portion of config: %w", err)
		}
	}

	return nil
}
