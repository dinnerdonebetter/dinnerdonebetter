package config

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	routecfg "github.com/dinnerdonebetter/backend/internal/routing"
	"os"
	"runtime/debug"
	"strings"

	analyticsconfig "github.com/dinnerdonebetter/backend/internal/analytics/config"
	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	emailconfig "github.com/dinnerdonebetter/backend/internal/email/config"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	featureflagsconfig "github.com/dinnerdonebetter/backend/internal/featureflags/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
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

	// CeaseOperationEnvVarKey is the env var key used to indicate a function or job should just quit early.
	CeaseOperationEnvVarKey = "CEASE_OPERATION"
	// ServiceEnvironmentEnvVarKey is the env var key we use to refer to the running environment.
	ServiceEnvironmentEnvVarKey = "DINNER_DONE_BETTER_SERVICE_ENVIRONMENT"
	// RunningInGCPEnvVarKey is the env var key we use to indicate we're running in GCP.
	RunningInGCPEnvVarKey = "RUNNING_IN_GCP"
	// FilePathEnvVarKey is the env var key we use to indicate where the config file is located.
	FilePathEnvVarKey = "CONFIGURATION_FILEPATH"
)

type (
	// runMode describes what method of operation the server is under.
	runMode string

	// CloserFunc calls all io.Closers in the service.
	CloserFunc func()

	// InstanceConfig configures an instance of the service. It is composed of all the other setting structs.
	InstanceConfig struct {
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
		"Queues":        cfg.Queues.ValidateWithContext,
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

func ShouldCeaseOperation() bool {
	return strings.TrimSpace(strings.ToLower(os.Getenv(CeaseOperationEnvVarKey))) == "true"
}

func RunningInCloud() bool {
	return os.Getenv(RunningInGCPEnvVarKey) != ""
}

type cloudConfigFetcher func(context.Context) (*InstanceConfig, error)

func FetchForApplication(ctx context.Context, cff cloudConfigFetcher) (*InstanceConfig, error) {
	var cfg *InstanceConfig
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

	if err := cfg.ValidateWithContext(ctx, true); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return cfg, nil
}
