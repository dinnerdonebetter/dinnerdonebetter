package config

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	analyticsconfig "github.com/dinnerdonebetter/backend/internal/analytics/config"
	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	emailconfig "github.com/dinnerdonebetter/backend/internal/email/config"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	featureflagsconfig "github.com/dinnerdonebetter/backend/internal/featureflags/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	routecfg "github.com/dinnerdonebetter/backend/internal/routing"
	searchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"
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
		_ struct{} `json:"-"`

		Observability observability.Config      `json:"observability" toml:"observability,omitempty"`
		Queues        QueuesConfig              `json:"queues"        toml:"queues,omitempty"`
		Email         emailconfig.Config        `json:"email"         toml:"email,omitempty"`
		Analytics     analyticsconfig.Config    `json:"analytics"     toml:"analytics,omitempty"`
		Search        searchcfg.Config          `json:"search"        toml:"search,omitempty"`
		FeatureFlags  featureflagsconfig.Config `json:"featureFlags"  toml:"events,omitempty"`
		Encoding      encoding.Config           `json:"encoding"      toml:"encoding,omitempty"`
		Meta          MetaSettings              `json:"meta"          toml:"meta,omitempty"`
		Routing       routecfg.Config           `json:"routing"       toml:"routing,omitempty"`
		Events        msgconfig.Config          `json:"events"        toml:"events,omitempty"`
		Server        http.Config               `json:"server"        toml:"server,omitempty"`
		Database      dbconfig.Config           `json:"database"      toml:"database,omitempty"`
		Services      ServicesConfig            `json:"services"      toml:"services,omitempty"`
	}

	// DBCleanerConfig configures an instance of the database cleaner job.
	DBCleanerConfig struct {
		_ struct{} `json:"-"`

		Observability observability.Config `json:"observability" toml:"observability,omitempty"`
		Database      dbconfig.Config      `json:"database"      toml:"database,omitempty"`
	}

	// EmailProberConfig configures an instance of the email prober job.
	EmailProberConfig struct {
		_ struct{} `json:"-"`

		Observability observability.Config `json:"observability" toml:"observability,omitempty"`
		Email         emailconfig.Config   `json:"email"         toml:"email,omitempty"`
		Database      dbconfig.Config      `json:"database"      toml:"database,omitempty"`
	}

	// MealPlanFinalizerConfig configures an instance of the meal plan finalizer job.
	MealPlanFinalizerConfig struct {
		_ struct{} `json:"-"`

		Observability observability.Config `json:"observability" toml:"observability,omitempty"`
		Events        msgconfig.Config     `json:"events"        toml:"events,omitempty"`
		Database      dbconfig.Config      `json:"database"      toml:"database,omitempty"`
	}

	// MealPlanGroceryListInitializerConfig configures an instance of the meal plan grocery list initializer job.
	MealPlanGroceryListInitializerConfig struct {
		_ struct{} `json:"-"`

		Observability observability.Config   `json:"observability" toml:"observability,omitempty"`
		Analytics     analyticsconfig.Config `json:"analytics"     toml:"analytics,omitempty"`
		Events        msgconfig.Config       `json:"events"        toml:"events,omitempty"`
		Database      dbconfig.Config        `json:"database"      toml:"database,omitempty"`
	}

	// MealPlanTaskCreatorConfig configures an instance of the meal plan task creator job.
	MealPlanTaskCreatorConfig struct {
		_ struct{} `json:"-"`

		Observability observability.Config   `json:"observability" toml:"observability,omitempty"`
		Analytics     analyticsconfig.Config `json:"analytics"     toml:"analytics,omitempty"`
		Events        msgconfig.Config       `json:"events"        toml:"events,omitempty"`
		Database      dbconfig.Config        `json:"database"      toml:"database,omitempty"`
	}

	// SearchDataIndexSchedulerConfig configures an instance of the search data index scheduler job.
	SearchDataIndexSchedulerConfig struct {
		_ struct{} `json:"-"`

		Observability observability.Config `json:"observability" toml:"observability,omitempty"`
		Events        msgconfig.Config     `json:"events"        toml:"events,omitempty"`
		Database      dbconfig.Config      `json:"database"      toml:"database,omitempty"`
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

// ValidateWithContext validates a APIServiceConfig struct.
func (cfg *APIServiceConfig) ValidateWithContext(ctx context.Context, validateServices bool) error {
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

	if validateServices {
		if err := cfg.Services.ValidateWithContext(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating Services config: %w", err), result)
		}
	}

	return result.ErrorOrNil()
}

// ValidateWithContext validates a APIServiceConfig struct.
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

func FetchForApplication[T configurations](ctx context.Context, cff genericCloudConfigFetcher[T]) (*T, error) {
	var cfg *T
	if RunningInCloud() {
		c, err := cff(ctx)
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

	return cfg, nil
}
