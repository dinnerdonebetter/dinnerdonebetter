package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	capitalism "gitlab.com/prixfixe/prixfixe/internal/capitalism"
	database "gitlab.com/prixfixe/prixfixe/internal/database"
	dbconfig "gitlab.com/prixfixe/prixfixe/internal/database/config"
	querier "gitlab.com/prixfixe/prixfixe/internal/database/querier"
	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	postgres "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding/postgres"
	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	routing "gitlab.com/prixfixe/prixfixe/internal/routing"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	server "gitlab.com/prixfixe/prixfixe/internal/server"
	auditservice "gitlab.com/prixfixe/prixfixe/internal/services/audit"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	frontendservice "gitlab.com/prixfixe/prixfixe/internal/services/frontend"
	invitationsservice "gitlab.com/prixfixe/prixfixe/internal/services/invitations"
	recipesservice "gitlab.com/prixfixe/prixfixe/internal/services/recipes"
	recipestepingredientsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepingredients"
	recipestepproductsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepproducts"
	recipestepsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipesteps"
	reportsservice "gitlab.com/prixfixe/prixfixe/internal/services/reports"
	validingredientpreparationsservice "gitlab.com/prixfixe/prixfixe/internal/services/validingredientpreparations"
	validingredientsservice "gitlab.com/prixfixe/prixfixe/internal/services/validingredients"
	validinstrumentsservice "gitlab.com/prixfixe/prixfixe/internal/services/validinstruments"
	validpreparationinstrumentsservice "gitlab.com/prixfixe/prixfixe/internal/services/validpreparationinstruments"
	validpreparationsservice "gitlab.com/prixfixe/prixfixe/internal/services/validpreparations"
	webhooksservice "gitlab.com/prixfixe/prixfixe/internal/services/webhooks"
	uploads "gitlab.com/prixfixe/prixfixe/internal/uploads"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// DevelopmentRunMode is the run mode for a development environment.
	DevelopmentRunMode runMode = "development"
	// TestingRunMode is the run mode for a testing environment.
	TestingRunMode runMode = "testing"
	// ProductionRunMode is the run mode for a production environment.
	ProductionRunMode runMode = "production"
	// DefaultRunMode is the default run mode.
	DefaultRunMode = DevelopmentRunMode
	// DefaultStartupDeadline is the default amount of time we allow for server startup.
	DefaultStartupDeadline = time.Minute
)

var (
	errNilDatabaseConnection   = errors.New("nil DB connection provided")
	errNilConfig               = errors.New("nil config provided")
	errInvalidDatabaseProvider = errors.New("invalid database provider")
)

type (
	// runMode describes what method of operation the server is under.
	runMode string

	// ServicesConfigurations collects the various service configurations.
	ServicesConfigurations struct {
		ValidInstruments            validinstrumentsservice.Config            `json:"validInstruments" mapstructure:"valid_instruments" toml:"valid_instruments,omitempty"`
		ValidPreparations           validpreparationsservice.Config           `json:"validPreparations" mapstructure:"valid_preparations" toml:"valid_preparations,omitempty"`
		ValidIngredients            validingredientsservice.Config            `json:"validIngredients" mapstructure:"valid_ingredients" toml:"valid_ingredients,omitempty"`
		Frontend                    frontendservice.Config                    `json:"frontend" mapstructure:"frontend" toml:"frontend,omitempty"`
		ValidPreparationInstruments validpreparationinstrumentsservice.Config `json:"validPreparationInstruments" mapstructure:"valid_preparation_instruments" toml:"valid_preparation_instruments,omitempty"`
		Recipes                     recipesservice.Config                     `json:"recipes" mapstructure:"recipes" toml:"recipes,omitempty"`
		RecipeSteps                 recipestepsservice.Config                 `json:"recipeSteps" mapstructure:"recipe_steps" toml:"recipe_steps,omitempty"`
		ValidIngredientPreparations validingredientpreparationsservice.Config `json:"validIngredientPreparations" mapstructure:"valid_ingredient_preparations" toml:"valid_ingredient_preparations,omitempty"`
		RecipeStepProducts          recipestepproductsservice.Config          `json:"recipeStepProducts" mapstructure:"recipe_step_products" toml:"recipe_step_products,omitempty"`
		Invitations                 invitationsservice.Config                 `json:"invitations" mapstructure:"invitations" toml:"invitations,omitempty"`
		Reports                     reportsservice.Config                     `json:"reports" mapstructure:"reports" toml:"reports,omitempty"`
		RecipeStepIngredients       recipestepingredientsservice.Config       `json:"recipeStepIngredients" mapstructure:"recipe_step_ingredients" toml:"recipe_step_ingredients,omitempty"`
		Auth                        authservice.Config                        `json:"auth" mapstructure:"auth" toml:"auth,omitempty"`
		Webhooks                    webhooksservice.Config                    `json:"webhooks" mapstructure:"webhooks" toml:"webhooks,omitempty"`
		AuditLog                    auditservice.Config                       `json:"auditLog" mapstructure:"audit_log" toml:"audit_log,omitempty"`
	}

	// InstanceConfig configures an instance of the service. It is composed of all the other setting structs.
	InstanceConfig struct {
		Search        search.Config          `json:"search" mapstructure:"search" toml:"search,omitempty"`
		Encoding      encoding.Config        `json:"encoding" mapstructure:"encoding" toml:"encoding,omitempty"`
		Uploads       uploads.Config         `json:"uploads" mapstructure:"uploads" toml:"uploads,omitempty"`
		Observability observability.Config   `json:"observability" mapstructure:"observability" toml:"observability,omitempty"`
		Routing       routing.Config         `json:"routing" mapstructure:"routing" toml:"routing,omitempty"`
		Capitalism    capitalism.Config      `json:"capitalism" mapstructure:"capitalism" toml:"capitalism,omitempty"`
		Meta          MetaSettings           `json:"meta" mapstructure:"meta" toml:"meta,omitempty"`
		Database      dbconfig.Config        `json:"database" mapstructure:"database" toml:"database,omitempty"`
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

	if err := cfg.Capitalism.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Capitalism portion of config: %w", err)
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

	if err := cfg.Services.AuditLog.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating AuditLog portion of config: %w", err)
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

	if err := cfg.Services.ValidPreparations.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating ValidPreparations service portion of config: %w", err)
	}

	if err := cfg.Services.ValidIngredients.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating ValidIngredients service portion of config: %w", err)
	}

	if err := cfg.Services.ValidIngredientPreparations.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating ValidIngredientPreparations service portion of config: %w", err)
	}

	if err := cfg.Services.ValidPreparationInstruments.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating ValidPreparationInstruments service portion of config: %w", err)
	}

	if err := cfg.Services.Recipes.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Recipes service portion of config: %w", err)
	}

	if err := cfg.Services.RecipeSteps.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating RecipeSteps service portion of config: %w", err)
	}

	if err := cfg.Services.RecipeStepIngredients.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating RecipeStepIngredients service portion of config: %w", err)
	}

	if err := cfg.Services.RecipeStepProducts.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating RecipeStepProducts service portion of config: %w", err)
	}

	if err := cfg.Services.Invitations.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Invitations service portion of config: %w", err)
	}

	if err := cfg.Services.Reports.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("error validating Reports service portion of config: %w", err)
	}

	return nil
}

// ProvideDatabaseClient provides a database implementation dependent on the configuration.
// NOTE: you may be tempted to move this to the database/config package. This is a fool's errand.
func ProvideDatabaseClient(ctx context.Context, logger logging.Logger, rawDB *sql.DB, cfg *InstanceConfig) (database.DataManager, error) {
	if rawDB == nil {
		return nil, errNilDatabaseConnection
	}

	if cfg == nil {
		return nil, errNilConfig
	}

	var qb querybuilding.SQLQueryBuilder
	shouldCreateTestUser := cfg.Meta.RunMode != ProductionRunMode

	switch strings.ToLower(strings.TrimSpace(cfg.Database.Provider)) {
	case "postgres":
		qb = postgres.ProvidePostgres(logger)
	default:
		return nil, fmt.Errorf("%w: %q", errInvalidDatabaseProvider, cfg.Database.Provider)
	}

	return querier.ProvideDatabaseClient(ctx, logger, rawDB, &cfg.Database, qb, shouldCreateTestUser)
}
