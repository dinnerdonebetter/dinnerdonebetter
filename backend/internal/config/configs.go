package config

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	analyticscfg "github.com/dinnerdonebetter/backend/internal/analytics/config"
	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
	emailcfg "github.com/dinnerdonebetter/backend/internal/email/config"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	featureflagscfg "github.com/dinnerdonebetter/backend/internal/featureflags/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/routing"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/server/http"

	"github.com/caarlos0/env/v11"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-multierror"
)

const (
	// DevelopmentRunMode is the run mode for a development environment.
	DevelopmentRunMode runMode = "development"
	// TestingRunMode is the run mode for a testing environment.
	TestingRunMode runMode = "testing"
	// ProductionRunMode is the run mode for a production environment.
	ProductionRunMode runMode = "production"

	EnvVarPrefix = "DINNER_DONE_BETTER_"

	// FilePathEnvVarKey is the env var key we use to indicate where the config file is located.
	FilePathEnvVarKey = "CONFIGURATION_FILEPATH"
)

type (
	// runMode describes what method of operation the server is under.
	runMode string

	// CloserFunc calls all io.Closers in the service.
	CloserFunc func()

	configurations interface {
		APIServiceConfig |
			DBCleanerConfig |
			EmailProberConfig |
			MealPlanFinalizerConfig |
			MealPlanGroceryListInitializerConfig |
			MealPlanTaskCreatorConfig |
			SearchDataIndexSchedulerConfig
	}

	genericCloudConfigFetcher[T configurations] func(context.Context) (*T, error)

	// APIServiceConfig configures an instance of the service. It is composed of all the other setting structs.
	APIServiceConfig struct {
		_                struct{}               `json:"-"`
		Queues           msgconfig.QueuesConfig `envPrefix:"QUEUES_"        json:"queues"`
		Email            emailcfg.Config        `envPrefix:"EMAIL_"         json:"email"`
		Analytics        analyticscfg.Config    `envPrefix:"ANALYTICS_"     json:"analytics"`
		Search           textsearchcfg.Config   `envPrefix:"SEARCH_"        json:"search"`
		FeatureFlags     featureflagscfg.Config `envPrefix:"FEATURE_FLAGS_" json:"featureFlags"`
		Encoding         encoding.Config        `envPrefix:"ENCODING_"      json:"encoding"`
		Events           msgconfig.Config       `envPrefix:"EVENTS_"        json:"events"`
		Observability    observability.Config   `envPrefix:"OBSERVABILITY_" json:"observability"`
		Meta             MetaSettings           `envPrefix:"META_"          json:"meta"`
		Routing          routing.Config         `envPrefix:"ROUTING_"       json:"routing"`
		Server           http.Config            `envPrefix:"SERVER_"        json:"server"`
		Database         databasecfg.Config     `envPrefix:"DATABASE_"      json:"database"`
		Services         ServicesConfig         `envPrefix:"SERVICE_"       json:"services"`
		validateServices bool                   `json:"-"`
	}

	// DBCleanerConfig configures an instance of the database cleaner job.
	DBCleanerConfig struct {
		_ struct{} `json:"-"`

		Observability observability.Config `envPrefix:"OBSERVABILITY_" json:"observability"`
		Database      databasecfg.Config   `envPrefix:"DATABASE_"      json:"database"`
	}

	// EmailProberConfig configures an instance of the email prober job.
	EmailProberConfig struct {
		_             struct{}             `json:"-"`
		Email         emailcfg.Config      `envPrefix:"EMAIL_"         json:"email"`
		Observability observability.Config `envPrefix:"OBSERVABILITY_" json:"observability"`
		Database      databasecfg.Config   `envPrefix:"DATABASE_"      json:"database"`
	}

	// MealPlanFinalizerConfig configures an instance of the meal plan finalizer job.
	MealPlanFinalizerConfig struct {
		_             struct{}             `json:"-"`
		Observability observability.Config `envPrefix:"OBSERVABILITY_" json:"observability"`
		Events        msgconfig.Config     `envPrefix:"EVENTS_"        json:"events"`
		Database      databasecfg.Config   `envPrefix:"DATABASE_"      json:"database"`
	}

	// MealPlanGroceryListInitializerConfig configures an instance of the meal plan grocery list initializer job.
	MealPlanGroceryListInitializerConfig struct {
		_             struct{}             `json:"-"`
		Analytics     analyticscfg.Config  `envPrefix:"ANALYTICS_"     json:"analytics"`
		Observability observability.Config `envPrefix:"OBSERVABILITY_" json:"observability"`
		Events        msgconfig.Config     `envPrefix:"EVENTS_"        json:"events"`
		Database      databasecfg.Config   `envPrefix:"DATABASE_"      json:"database"`
	}

	// MealPlanTaskCreatorConfig configures an instance of the meal plan task creator job.
	MealPlanTaskCreatorConfig struct {
		_             struct{}             `json:"-"`
		Analytics     analyticscfg.Config  `envPrefix:"ANALYTICS_"     json:"analytics"`
		Observability observability.Config `envPrefix:"OBSERVABILITY_" json:"observability"`
		Events        msgconfig.Config     `envPrefix:"EVENTS_"        json:"events"`
		Database      databasecfg.Config   `envPrefix:"DATABASE_"      json:"database"`
	}

	// SearchDataIndexSchedulerConfig configures an instance of the search data index scheduler job.
	SearchDataIndexSchedulerConfig struct {
		_             struct{}             `json:"-"`
		Observability observability.Config `envPrefix:"OBSERVABILITY_" json:"observability"`
		Events        msgconfig.Config     `envPrefix:"EVENTS_"        json:"events"`
		Database      databasecfg.Config   `envPrefix:"DATABASE_"      json:"database"`
	}
)

// EncodeToFile renders your config to a file given your favorite encoder.
func (cfg *APIServiceConfig) EncodeToFile(path string, marshaller func(v any) ([]byte, error)) error {
	if cfg == nil {
		return errors.New("nil config")
	}

	byteSlice, err := marshaller(*cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, byteSlice, 0o600)
}

func (cfg *APIServiceConfig) Commit() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for i := range info.Settings {
			if info.Settings[i].Key == "vcs.revision" {
				return info.Settings[i].Value
			}
		}
	}

	return ""
}

var _ validation.ValidatableWithContext = (*APIServiceConfig)(nil)

// ValidateWithContext validates a APIServiceConfig struct.
func (cfg *APIServiceConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validators := map[string]func(context.Context) error{
		"Routing":       cfg.Routing.ValidateWithContext,
		"Meta":          cfg.Meta.ValidateWithContext,
		"Queues":        cfg.Queues.ValidateWithContext,
		"Encoding":      cfg.Encoding.ValidateWithContext,
		"Analytics":     cfg.Analytics.ValidateWithContext,
		"Observability": cfg.Observability.ValidateWithContext,
		"Database":      cfg.Database.ValidateWithContext,
		"Server":        cfg.Server.ValidateWithContext,
		"Email":         cfg.Email.ValidateWithContext,
		"FeatureFlags":  cfg.FeatureFlags.ValidateWithContext,
		"Search":        cfg.Search.ValidateWithContext,
		// no "Events" here, that's a collection of publisher/subscriber configs that can each optionally be setup
	}

	for name, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	if cfg.validateServices {
		if err := cfg.Services.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating Services config: %w", err), result)
		}
	}

	return result.ErrorOrNil()
}

// ValidateWithContext validates a DBCleanerConfig struct.
func (cfg *DBCleanerConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validators := map[string]func(context.Context) error{
		"Observability": cfg.Observability.ValidateWithContext,
		"Database":      cfg.Database.ValidateWithContext,
	}

	for name, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}

// ValidateWithContext validates a EmailProberConfig struct.
func (cfg *EmailProberConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validators := map[string]func(context.Context) error{
		"Observability": cfg.Observability.ValidateWithContext,
		"Database":      cfg.Database.ValidateWithContext,
		"Email":         cfg.Email.ValidateWithContext,
	}

	for name, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}

// ValidateWithContext validates a MealPlanFinalizerConfig struct.
func (cfg *MealPlanFinalizerConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validators := map[string]func(context.Context) error{
		"Observability": cfg.Observability.ValidateWithContext,
		"Database":      cfg.Database.ValidateWithContext,
	}

	for name, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}

// ValidateWithContext validates a MealPlanGroceryListInitializerConfig struct.
func (cfg *MealPlanGroceryListInitializerConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validators := map[string]func(context.Context) error{
		"Analytics":     cfg.Analytics.ValidateWithContext,
		"Observability": cfg.Observability.ValidateWithContext,
		"Database":      cfg.Database.ValidateWithContext,
	}

	for name, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}

// ValidateWithContext validates a MealPlanTaskCreatorConfig struct.
func (cfg *MealPlanTaskCreatorConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validators := map[string]func(context.Context) error{
		"Analytics":     cfg.Analytics.ValidateWithContext,
		"Observability": cfg.Observability.ValidateWithContext,
		"Database":      cfg.Database.ValidateWithContext,
	}

	for name, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}

// ValidateWithContext validates a SearchDataIndexSchedulerConfig struct.
func (cfg *SearchDataIndexSchedulerConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validators := map[string]func(context.Context) error{
		"Observability": cfg.Observability.ValidateWithContext,
		"Database":      cfg.Database.ValidateWithContext,
	}

	for name, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}

func FetchForApplication[T configurations](ctx context.Context, f genericCloudConfigFetcher[T]) (*T, error) {
	var cfg *T
	if RunningInTheCloud() {
		c, err := f(ctx)
		if err != nil {
			return nil, fmt.Errorf("fetching config from GCP: %w", err)
		}

		cfg = c
	} else if configFilepath := os.Getenv(FilePathEnvVarKey); configFilepath != "" {
		configBytes, err := os.ReadFile(configFilepath)
		if err != nil {
			return nil, fmt.Errorf("reading local config file: %w", err)
		}

		if err = json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); err != nil || cfg == nil {
			return nil, fmt.Errorf("decoding config file contents: %w", err)
		}
	} else {
		return nil, errors.New("not running in the cloud, and no config filepath provided")
	}

	if err := env.ParseWithOptions(cfg, env.Options{Prefix: EnvVarPrefix}); err != nil {
		return nil, err
	}

	return cfg, nil
}
