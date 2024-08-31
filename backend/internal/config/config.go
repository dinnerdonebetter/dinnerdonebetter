package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/uploads/objectstorage"
	"os"
	"runtime/debug"

	analyticsconfig "github.com/dinnerdonebetter/backend/internal/analytics/config"
	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	emailconfig "github.com/dinnerdonebetter/backend/internal/email/config"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	featureflagsconfig "github.com/dinnerdonebetter/backend/internal/featureflags/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/routing"
	searchcfg "github.com/dinnerdonebetter/backend/internal/search/config"
	"github.com/dinnerdonebetter/backend/internal/server/http"

	"github.com/hashicorp/go-multierror"
)

const (
	// DevelopmentRunMode is the run mode for a development environment.
	DevelopmentRunMode runMode = "development"
	// TestingRunMode is the run mode for a testing environment.
	TestingRunMode runMode = "testing"
	// ProductionRunMode is the run mode for a production environment.
	ProductionRunMode runMode = "production"

	/* #nosec G101 */
	debugCookieHashKey = "HEREISA32CHARSECRETWHICHISMADEUP"
	/* #nosec G101 */
	debugCookieBlockKey = "DIFFERENT32CHARSECRETTHATIMADEUP"
)

type (
	// runMode describes what method of operation the server is under.
	runMode string

	// CloserFunc calls all io.Closers in the service.
	CloserFunc func()

	// InstanceConfig configures an instance of the service. It is composed of all the other setting structs.
	InstanceConfig struct {
		_             struct{}                  `json:"-"`
		Observability observability.Config      `json:"observability" toml:"observability,omitempty"`
		Email         emailconfig.Config        `json:"email"         toml:"email,omitempty"`
		Analytics     analyticsconfig.Config    `json:"analytics"     toml:"analytics,omitempty"`
		Search        searchcfg.Config          `json:"search"        toml:"search,omitempty"`
		FeatureFlags  featureflagsconfig.Config `json:"featureFlags"  toml:"events,omitempty"`
		Encoding      encoding.Config           `json:"encoding"      toml:"encoding,omitempty"`
		Meta          MetaSettings              `json:"meta"          toml:"meta,omitempty"`
		Routing       routing.Config            `json:"routing"       toml:"routing,omitempty"`
		Events        msgconfig.Config          `json:"events"        toml:"events,omitempty"`
		Server        http.Config               `json:"server"        toml:"server,omitempty"`
		Database      dbconfig.Config           `json:"database"      toml:"database,omitempty"`
		Services      ServicesConfig            `json:"services"      toml:"services,omitempty"`
	}
)

// EncodeToFile renders your config to a file given your favorite encoder.
func (cfg *InstanceConfig) EncodeToFile(path string, marshaller func(v any) ([]byte, error)) error {
	if cfg == nil {
		return errors.New("nil config")
	}

	byteSlice, err := marshaller(*cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, byteSlice, 0o600)
}

// ValidateWithContext validates a InstanceConfig struct.
func (cfg *InstanceConfig) ValidateWithContext(ctx context.Context, validateServices bool) error {
	var result *multierror.Error

	validators := map[string]func(context.Context) error{
		"Routing":       cfg.Routing.ValidateWithContext,
		"Meta":          cfg.Meta.ValidateWithContext,
		"Encoding":      cfg.Encoding.ValidateWithContext,
		"Analytics":     cfg.Analytics.ValidateWithContext,
		"Observability": cfg.Observability.ValidateWithContext,
		"Database":      cfg.Database.ValidateWithContext,
		"Server":        cfg.Server.ValidateWithContext,
		"Email":         cfg.Email.ValidateWithContext,
		"FeatureFlags":  cfg.FeatureFlags.ValidateWithContext,
		"Search":        cfg.Search.ValidateWithContext,
	}

	for name, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	if validateServices {
		if err := cfg.Services.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating Services config: %w", err), result)
		}
	}

	return result.ErrorOrNil()
}

func (cfg *InstanceConfig) Commit() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for i := range info.Settings {
			if info.Settings[i].Key == "vcs.revision" {
				return info.Settings[i].Value
			}
		}
	}

	return ""
}

func (cfg *InstanceConfig) Neutralize() {
	if err := os.Setenv("GOOGLE_CLOUD_PROJECT_ID", "something"); err != nil {
		panic(err)
	}

	cfg.Database.RunMigrations = false
	cfg.Database.OAuth2TokenEncryptionKey = "BLAHBLAHBLAHBLAHBLAHBLAHBLAHBLAH"
	cfg.Services.Auth.Cookies.HashKey = debugCookieHashKey
	cfg.Services.Auth.Cookies.BlockKey = debugCookieBlockKey
	cfg.Services.Auth.SSO.Google.ClientID = "blah blah blah blah"
	cfg.Services.Auth.SSO.Google.ClientSecret = "blah blah blah blah"
	cfg.Email.Sendgrid.APIToken = "blah blah blah blah"
	cfg.Analytics.Provider = ""
	cfg.Services.Recipes.Uploads.Storage.GCPConfig = nil
	cfg.Services.Recipes.Uploads.Storage.Provider = objectstorage.FilesystemProvider
	cfg.Services.Recipes.Uploads.Storage.FilesystemConfig = &objectstorage.FilesystemConfig{RootDirectory: "/tmp"}
	cfg.Services.RecipeSteps.Uploads.Storage.GCPConfig = nil
	cfg.Services.RecipeSteps.Uploads.Storage.Provider = objectstorage.FilesystemProvider
	cfg.Services.RecipeSteps.Uploads.Storage.FilesystemConfig = &objectstorage.FilesystemConfig{RootDirectory: "/tmp"}
	cfg.Services.ValidMeasurementUnits.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.ValidInstruments.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.ValidIngredients.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.ValidPreparations.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.ValidIngredientPreparations.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.ValidPreparationInstruments.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.ValidInstrumentMeasurementUnits.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.Recipes.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.RecipeSteps.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.RecipeStepProducts.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.RecipeStepInstruments.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.RecipeStepIngredients.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.Meals.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.MealPlans.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.MealPlanEvents.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.MealPlanOptions.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.MealPlanOptionVotes.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.MealPlanTasks.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.Households.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.HouseholdInvitations.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.Users.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.ValidIngredientGroups.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.Webhooks.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.Auth.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.RecipePrepTasks.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.MealPlanGroceryListItems.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.ValidMeasurementUnitConversions.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.ValidIngredientStates.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.RecipeStepCompletionConditions.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.ValidIngredientStateIngredients.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.RecipeStepVessels.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.ServiceSettings.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.ServiceSettingConfigurations.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.UserIngredientPreferences.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.RecipeRatings.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.HouseholdInstrumentOwnerships.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.ValidVessels.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.ValidPreparationVessels.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.Workers.DataChangesTopicName = "dataChangesTopicName"
	cfg.Services.UserNotifications.DataChangesTopicName = "dataChangesTopicName"
}
