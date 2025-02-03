package config

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"time"

	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
	analyticscfg "github.com/dinnerdonebetter/backend/internal/lib/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/cookies"
	emailcfg "github.com/dinnerdonebetter/backend/internal/lib/email/config"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	featureflagscfg "github.com/dinnerdonebetter/backend/internal/lib/featureflags/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	routingcfg "github.com/dinnerdonebetter/backend/internal/lib/routing/config"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/lib/server/grpc"
	"github.com/dinnerdonebetter/backend/internal/lib/server/http"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads/objectstorage"

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

	// ConfigurationFilePathEnvVarKey is the env var key we use to indicate where the config file is located.
	ConfigurationFilePathEnvVarKey = "CONFIGURATION_FILEPATH"
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
			SearchDataIndexSchedulerConfig |
			AsyncMessageHandlerConfig |
			AdminWebappConfig
	}

	// APIServiceConfig configures an instance of the service. It is composed of all the other setting structs.
	APIServiceConfig struct {
		_                struct{}               `json:"-"`
		Queues           msgconfig.QueuesConfig `envPrefix:"QUEUES_"        json:"queues"`
		Routing          routingcfg.Config      `envPrefix:"ROUTING_"       json:"routing"`
		Encoding         encoding.Config        `envPrefix:"ENCODING_"      json:"encoding"`
		Events           msgconfig.Config       `envPrefix:"EVENTS_"        json:"events"`
		Observability    observability.Config   `envPrefix:"OBSERVABILITY_" json:"observability"`
		Meta             MetaSettings           `envPrefix:"META_"          json:"meta"`
		Email            emailcfg.Config        `envPrefix:"EMAIL_"         json:"email"`
		Analytics        analyticscfg.Config    `envPrefix:"ANALYTICS_"     json:"analytics"`
		TextSearch       textsearchcfg.Config   `envPrefix:"SEARCH_"        json:"search"`
		FeatureFlags     featureflagscfg.Config `envPrefix:"FEATURE_FLAGS_" json:"featureFlags"`
		HTTPServer       http.Config            `envPrefix:"SERVER_"        json:"server"`
		Database         databasecfg.Config     `envPrefix:"DATABASE_"      json:"database"`
		Services         ServicesConfig         `envPrefix:"SERVICE_"       json:"services"`
		GRPCServer       grpc.Config            `envPrefix:"GRPC_SERVER_"   json:"grpcServer"`
		validateServices bool
	}

	// DBCleanerConfig configures an instance of the database cleaner job.
	DBCleanerConfig struct {
		_ struct{} `json:"-"`

		Observability observability.Config `envPrefix:"OBSERVABILITY_" json:"observability"`
		Database      databasecfg.Config   `envPrefix:"DATABASE_"      json:"database"`
	}

	// EmailProberConfig configures an instance of the email prober job.
	EmailProberConfig struct {
		_ struct{} `json:"-"`

		Email         emailcfg.Config      `envPrefix:"EMAIL_"         json:"email"`
		Observability observability.Config `envPrefix:"OBSERVABILITY_" json:"observability"`
		Database      databasecfg.Config   `envPrefix:"DATABASE_"      json:"database"`
	}

	// MealPlanFinalizerConfig configures an instance of the meal plan finalizer job.
	MealPlanFinalizerConfig struct {
		_ struct{} `json:"-"`

		Queues        msgconfig.QueuesConfig `envPrefix:"QUEUES_"        json:"queues"`
		Events        msgconfig.Config       `envPrefix:"EVENTS_"        json:"events"`
		Observability observability.Config   `envPrefix:"OBSERVABILITY_" json:"observability"`
		Database      databasecfg.Config     `envPrefix:"DATABASE_"      json:"database"`
	}

	// MealPlanGroceryListInitializerConfig configures an instance of the meal plan grocery list initializer job.
	MealPlanGroceryListInitializerConfig struct {
		_ struct{} `json:"-"`

		Queues        msgconfig.QueuesConfig `envPrefix:"QUEUES_"        json:"queues"`
		Analytics     analyticscfg.Config    `envPrefix:"ANALYTICS_"     json:"analytics"`
		Events        msgconfig.Config       `envPrefix:"EVENTS_"        json:"events"`
		Observability observability.Config   `envPrefix:"OBSERVABILITY_" json:"observability"`
		Database      databasecfg.Config     `envPrefix:"DATABASE_"      json:"database"`
	}

	// MealPlanTaskCreatorConfig configures an instance of the meal plan task creator job.
	MealPlanTaskCreatorConfig struct {
		_ struct{} `json:"-"`

		Queues        msgconfig.QueuesConfig `envPrefix:"QUEUES_"        json:"queues"`
		Analytics     analyticscfg.Config    `envPrefix:"ANALYTICS_"     json:"analytics"`
		Events        msgconfig.Config       `envPrefix:"EVENTS_"        json:"events"`
		Observability observability.Config   `envPrefix:"OBSERVABILITY_" json:"observability"`
		Database      databasecfg.Config     `envPrefix:"DATABASE_"      json:"database"`
	}

	// SearchDataIndexSchedulerConfig configures an instance of the search data index scheduler job.
	SearchDataIndexSchedulerConfig struct {
		_ struct{} `json:"-"`

		Queues        msgconfig.QueuesConfig `envPrefix:"QUEUES_"        json:"queues"`
		Events        msgconfig.Config       `envPrefix:"EVENTS_"        json:"events"`
		Observability observability.Config   `envPrefix:"OBSERVABILITY_" json:"observability"`
		Database      databasecfg.Config     `envPrefix:"DATABASE_"      json:"database"`
	}

	// AsyncMessageHandlerConfig configures an instance of the search data index scheduler job.
	AsyncMessageHandlerConfig struct {
		_ struct{} `json:"-"`

		Storage       objectstorage.Config   `envPrefix:"STORAGE_"       json:"storage"`
		Queues        msgconfig.QueuesConfig `envPrefix:"QUEUES_"        json:"queues"`
		Email         emailcfg.Config        `envPrefix:"EMAIL_"         json:"email"`
		Analytics     analyticscfg.Config    `envPrefix:"ANALYTICS_"     json:"analytics"`
		Search        textsearchcfg.Config   `envPrefix:"SEARCH_"        json:"search"`
		Events        msgconfig.Config       `envPrefix:"EVENTS_"        json:"events"`
		Observability observability.Config   `envPrefix:"OBSERVABILITY_" json:"observability"`
		Database      databasecfg.Config     `envPrefix:"DATABASE_"      json:"database"`
	}

	APIServiceOAuth2ConnectionConfig struct {
		_ struct{} `json:"-"`

		APIServerURL          string `env:"API_SERVER_URL"           json:"apiServerURL"`
		OAuth2APIClientID     string `env:"OAUTH2_API_CLIENT_ID"     json:"oauth2APIClientID"`
		OAuth2APIClientSecret string `env:"OAUTH2_API_CLIENT_SECRET" json:"oauth2APIClientSecret"`
	}

	NamedCacheConfig struct {
		_ struct{} `json:"-"`

		CacheCapacity uint64        `env:"CACHE_CAPACITY" json:"cacheCapacity"`
		CacheTTL      time.Duration `env:"CACHE_TTL"      json:"cacheTTL"`
	}

	// AdminWebappConfig configures an instance of the service. It is composed of all the other setting structs.
	AdminWebappConfig struct {
		_ struct{} `json:"-"`

		Cookies              cookies.Config                   `env:"init"                    envPrefix:"COOKIES_"    json:"cookies"`
		APIServiceConnection APIServiceOAuth2ConnectionConfig `envPrefix:"API_SERVICE_"      json:"apiServiceConfig"`
		Routing              routingcfg.Config                `envPrefix:"ROUTING_"          json:"routing"`
		Encoding             encoding.Config                  `envPrefix:"ENCODING_"         json:"encoding"`
		Observability        observability.Config             `envPrefix:"OBSERVABILITY_"    json:"observability"`
		Meta                 MetaSettings                     `envPrefix:"META_"             json:"meta"`
		HTTPServer           http.Config                      `envPrefix:"SERVER_"           json:"server"`
		APIClientCache       NamedCacheConfig                 `envPrefix:"API_CLIENT_CACHE_" json:"apiClientCache"`
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
		"HTTPServer":    cfg.HTTPServer.ValidateWithContext,
		"Email":         cfg.Email.ValidateWithContext,
		"FeatureFlags":  cfg.FeatureFlags.ValidateWithContext,
		"TextSearch":    cfg.TextSearch.ValidateWithContext,
		// no "Events" here, that's a collection of publisher/subscriber configs that can each optionally be setup
	}

	if cfg.validateServices {
		validators["Services"] = cfg.Services.ValidateWithContext
	}

	for name, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}

var _ validation.ValidatableWithContext = (*DBCleanerConfig)(nil)

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

var _ validation.ValidatableWithContext = (*EmailProberConfig)(nil)

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

var _ validation.ValidatableWithContext = (*MealPlanFinalizerConfig)(nil)

// ValidateWithContext validates a MealPlanFinalizerConfig struct.
func (cfg *MealPlanFinalizerConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validators := map[string]func(context.Context) error{
		"Observability": cfg.Observability.ValidateWithContext,
		"Database":      cfg.Database.ValidateWithContext,
		"Queues":        cfg.Queues.ValidateWithContext,
	}

	for name, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}

var _ validation.ValidatableWithContext = (*MealPlanGroceryListInitializerConfig)(nil)

// ValidateWithContext validates a MealPlanGroceryListInitializerConfig struct.
func (cfg *MealPlanGroceryListInitializerConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validators := map[string]func(context.Context) error{
		"Analytics":     cfg.Analytics.ValidateWithContext,
		"Observability": cfg.Observability.ValidateWithContext,
		"Database":      cfg.Database.ValidateWithContext,
		"Queues":        cfg.Queues.ValidateWithContext,
	}

	for name, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}

var _ validation.ValidatableWithContext = (*MealPlanTaskCreatorConfig)(nil)

// ValidateWithContext validates a MealPlanTaskCreatorConfig struct.
func (cfg *MealPlanTaskCreatorConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validators := map[string]func(context.Context) error{
		"Analytics":     cfg.Analytics.ValidateWithContext,
		"Observability": cfg.Observability.ValidateWithContext,
		"Database":      cfg.Database.ValidateWithContext,
		"Queues":        cfg.Queues.ValidateWithContext,
	}

	for name, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}

var _ validation.ValidatableWithContext = (*SearchDataIndexSchedulerConfig)(nil)

// ValidateWithContext validates a SearchDataIndexSchedulerConfig struct.
func (cfg *SearchDataIndexSchedulerConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validators := map[string]func(context.Context) error{
		"Observability": cfg.Observability.ValidateWithContext,
		"Database":      cfg.Database.ValidateWithContext,
		"Queues":        cfg.Queues.ValidateWithContext,
	}

	for name, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}

var _ validation.ValidatableWithContext = (*AsyncMessageHandlerConfig)(nil)

// ValidateWithContext validates a AsyncMessageHandlerConfig struct.
func (cfg *AsyncMessageHandlerConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validators := map[string]func(context.Context) error{
		"Queues":        cfg.Queues.ValidateWithContext,
		"Analytics":     cfg.Analytics.ValidateWithContext,
		"Observability": cfg.Observability.ValidateWithContext,
		"Database":      cfg.Database.ValidateWithContext,
		"Email":         cfg.Email.ValidateWithContext,
		"TextSearch":    cfg.Search.ValidateWithContext,
		"Storage":       cfg.Storage.ValidateWithContext,
	}

	for name, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}

var _ validation.ValidatableWithContext = (*AdminWebappConfig)(nil)

// ValidateWithContext validates a AdminWebappConfig struct.
func (cfg *AdminWebappConfig) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	validators := map[string]func(context.Context) error{
		"Cookies":       cfg.Cookies.ValidateWithContext,
		"Encoding":      cfg.Encoding.ValidateWithContext,
		"Observability": cfg.Observability.ValidateWithContext,
		"Meta":          cfg.Meta.ValidateWithContext,
		"Routing":       cfg.Routing.ValidateWithContext,
		"HTTPServer":    cfg.HTTPServer.ValidateWithContext,
	}

	for name, validator := range validators {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}

func envVarOnSetFunc(tag string, value any, isDefault bool) {
	slog.Info("env var set",
		slog.String("tag", tag),
		slog.String("value", fmt.Sprintf("%+v", value)),
		slog.Bool("isDefault", isDefault),
	)
}

func LoadConfigFromEnvironment[T configurations]() (*T, error) {
	configFilepath := os.Getenv(ConfigurationFilePathEnvVarKey)

	configBytes, err := os.ReadFile(configFilepath)
	if err != nil {
		return nil, fmt.Errorf("reading local config file: %w", err)
	}

	var cfg *T
	if err = json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); err != nil || cfg == nil {
		return nil, fmt.Errorf("decoding config file (%s) contents (%s): %w", configFilepath, string(configBytes), err)
	}

	if err = ApplyEnvironmentVariables(cfg); err != nil {
		return nil, fmt.Errorf("applying environment variables: %w", err)
	}

	return cfg, nil
}
