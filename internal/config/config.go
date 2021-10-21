package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	database "gitlab.com/prixfixe/prixfixe/internal/database"
	dbconfig "gitlab.com/prixfixe/prixfixe/internal/database/config"
	"gitlab.com/prixfixe/prixfixe/internal/database/queriers/postgres"
	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	msgconfig "gitlab.com/prixfixe/prixfixe/internal/messagequeue/config"
	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	routing "gitlab.com/prixfixe/prixfixe/internal/routing"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	server "gitlab.com/prixfixe/prixfixe/internal/server"
	accountsservice "gitlab.com/prixfixe/prixfixe/internal/services/accounts"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	frontendservice "gitlab.com/prixfixe/prixfixe/internal/services/frontend"
	mealplanoptionsservice "gitlab.com/prixfixe/prixfixe/internal/services/mealplanoptions"
	mealplanoptionvotesservice "gitlab.com/prixfixe/prixfixe/internal/services/mealplanoptionvotes"
	mealplansservice "gitlab.com/prixfixe/prixfixe/internal/services/mealplans"
	recipesservice "gitlab.com/prixfixe/prixfixe/internal/services/recipes"
	recipestepingredientsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepingredients"
	recipestepinstrumentsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepinstruments"
	recipestepproductsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepproducts"
	recipestepsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipesteps"
	validingredientpreparationsservice "gitlab.com/prixfixe/prixfixe/internal/services/validingredientpreparations"
	validingredientsservice "gitlab.com/prixfixe/prixfixe/internal/services/validingredients"
	validinstrumentsservice "gitlab.com/prixfixe/prixfixe/internal/services/validinstruments"
	validpreparationsservice "gitlab.com/prixfixe/prixfixe/internal/services/validpreparations"
	webhooksservice "gitlab.com/prixfixe/prixfixe/internal/services/webhooks"
	websocketsservice "gitlab.com/prixfixe/prixfixe/internal/services/websockets"
	uploads "gitlab.com/prixfixe/prixfixe/internal/uploads"
)

const (
	// DevelopmentRunMode is the run mode for a development environment.
	DevelopmentRunMode runMode = "development"
	// TestingRunMode is the run mode for a testing environment.
	TestingRunMode runMode = "testing"
	// ProductionRunMode is the run mode for a production environment.
	ProductionRunMode runMode = "production"
)

var (
	errNilConfig               = errors.New("nil config provided")
	errInvalidDatabaseProvider = errors.New("invalid database provider")
)

type (
	// runMode describes what method of operation the server is under.
	runMode string

	// ServicesConfigurations collects the various service configurations.
	ServicesConfigurations struct {
		_                           struct{}
		ValidInstruments            validinstrumentsservice.Config            `json:"validInstruments" mapstructure:"valid_instruments" toml:"valid_instruments,omitempty"`
		ValidIngredients            validingredientsservice.Config            `json:"validIngredients" mapstructure:"valid_ingredients" toml:"valid_ingredients,omitempty"`
		ValidPreparations           validpreparationsservice.Config           `json:"validPreparations" mapstructure:"valid_preparations" toml:"valid_preparations,omitempty"`
		MealPlanOptionVotes         mealplanoptionvotesservice.Config         `json:"mealPlanOptionVotes" mapstructure:"meal_plan_option_votes" toml:"meal_plan_option_votes,omitempty"`
		ValidIngredientPreparations validingredientpreparationsservice.Config `json:"validIngredientPreparations" mapstructure:"valid_ingredient_preparations" toml:"valid_ingredient_preparations,omitempty"`
		Recipes                     recipesservice.Config                     `json:"recipes" mapstructure:"recipes" toml:"recipes,omitempty"`
		RecipeSteps                 recipestepsservice.Config                 `json:"recipeSteps" mapstructure:"recipe_steps" toml:"recipe_steps,omitempty"`
		RecipeStepInstruments       recipestepinstrumentsservice.Config       `json:"recipeStepInstruments" mapstructure:"recipe_step_instruments" toml:"recipe_step_instruments,omitempty"`
		RecipeStepIngredients       recipestepingredientsservice.Config       `json:"recipeStepIngredients" mapstructure:"recipe_step_ingredients" toml:"recipe_step_ingredients,omitempty"`
		MealPlans                   mealplansservice.Config                   `json:"mealPlans" mapstructure:"meal_plans" toml:"meal_plans,omitempty"`
		MealPlanOptions             mealplanoptionsservice.Config             `json:"mealPlanOptions" mapstructure:"meal_plan_options" toml:"meal_plan_options,omitempty"`
		RecipeStepProducts          recipestepproductsservice.Config          `json:"recipeStepProducts" mapstructure:"recipe_step_products" toml:"recipe_step_products,omitempty"`
		Accounts                    accountsservice.Config                    `json:"accounts" mapstructure:"accounts" toml:"accounts,omitempty"`
		Websockets                  websocketsservice.Config                  `json:"websockets" mapstructure:"websockets" toml:"websockets,omitempty"`
		Webhooks                    webhooksservice.Config                    `json:"webhooks" mapstructure:"webhooks" toml:"webhooks,omitempty"`
		Frontend                    frontendservice.Config                    `json:"frontend" mapstructure:"frontend" toml:"frontend,omitempty"`
		Auth                        authservice.Config                        `json:"auth" mapstructure:"auth" toml:"auth,omitempty"`
	}

	// InstanceConfig configures an instance of the service. It is composed of all the other setting structs.
	InstanceConfig struct {
		_             struct{}
		Events        msgconfig.Config       `json:"events" mapstructure:"events" toml:"events,omitempty"`
		Search        search.Config          `json:"search" mapstructure:"search" toml:"search,omitempty"`
		Encoding      encoding.Config        `json:"encoding" mapstructure:"encoding" toml:"encoding,omitempty"`
		Uploads       uploads.Config         `json:"uploads" mapstructure:"uploads" toml:"uploads,omitempty"`
		Observability observability.Config   `json:"observability" mapstructure:"observability" toml:"observability,omitempty"`
		Routing       routing.Config         `json:"routing" mapstructure:"routing" toml:"routing,omitempty"`
		Database      dbconfig.Config        `json:"database" mapstructure:"database" toml:"database,omitempty"`
		Meta          MetaSettings           `json:"meta" mapstructure:"meta" toml:"meta,omitempty"`
		Services      ServicesConfigurations `json:"services" mapstructure:"services" toml:"services,omitempty"`
		Server        server.Config          `json:"server" mapstructure:"server" toml:"server,omitempty"`
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

var _ validation.ValidatableWithContext = (*InstanceConfig)(nil)

// ValidateWithContext validates a InstanceConfig struct.
func (cfg *InstanceConfig) ValidateWithContext(ctx context.Context) error {
	if err := cfg.Search.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Search portion of config: %w", err)
	}

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

	if err := cfg.Encoding.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Encoding portion of config: %w", err)
	}

	if err := cfg.Observability.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Observability portion of config: %w", err)
	}

	if err := cfg.Database.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Database portion of config: %w", err)
	}

	if err := cfg.Server.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating HTTPServer portion of config: %w", err)
	}

	if err := cfg.Services.Auth.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Auth service portion of config: %w", err)
	}

	if err := cfg.Services.Frontend.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Frontend service portion of config: %w", err)
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

	if err := cfg.Services.RecipeStepProducts.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating RecipeStepProducts service portion of config: %w", err)
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

	return nil
}

// ProvideDatabaseClient provides a database implementation dependent on the configuration.
// NOTE: you may be tempted to move this to the database/config package. This is a fool's errand.
func ProvideDatabaseClient(ctx context.Context, logger logging.Logger, cfg *InstanceConfig) (database.DataManager, error) {
	if cfg == nil {
		return nil, errNilConfig
	}

	shouldCreateTestUser := cfg.Meta.RunMode != ProductionRunMode

	switch strings.ToLower(strings.TrimSpace(cfg.Database.Provider)) {
	case dbconfig.PostgresProvider:
		return postgres.ProvideDatabaseClient(ctx, logger, &cfg.Database, shouldCreateTestUser)
	default:
		return nil, fmt.Errorf("%w: %q", errInvalidDatabaseProvider, cfg.Database.Provider)
	}
}
