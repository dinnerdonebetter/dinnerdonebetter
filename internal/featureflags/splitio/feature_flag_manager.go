package splitio

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/prixfixeco/backend/internal/featureflags"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"

	"github.com/splitio/go-client/v6/splitio/client"
	"github.com/splitio/go-client/v6/splitio/conf"
)

var (
	ErrMissingHTTPClient = errors.New("missing HTTP client")
	ErrNilConfig         = errors.New("missing config")
	ErrMissingSDKKey     = errors.New("missing SDK key")
)

type (
	Config struct {
		SDKKey      string        `json:"sdkKey" mapstructure:"sdk_key" toml:"sdk_key"`
		InitTimeout time.Duration `json:"initTimeout" mapstructure:"init_timeout" toml:"init_timeout"`
	}

	splitAPIClient interface {
		Treatment(key any, feature string, attributes map[string]any) string
	}

	// featureFlagManager implements the feature flag interface.
	featureFlagManager struct {
		logger      logging.Logger
		tracer      tracing.Tracer
		splitClient splitAPIClient
	}
)

// NewFeatureFlagManager constructs a new featureFlagManager.
func NewFeatureFlagManager(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider, httpClient *http.Client, configMods ...func(*conf.SplitSdkConfig) *conf.SplitSdkConfig) (featureflags.FeatureFlagManager, error) {
	if httpClient == nil {
		return nil, ErrMissingHTTPClient
	}

	if cfg == nil {
		return nil, ErrNilConfig
	}

	if cfg.SDKKey == "" {
		return nil, ErrMissingSDKKey
	}

	splitCfg := conf.Default()
	for _, configMod := range configMods {
		splitCfg = configMod(splitCfg)
	}

	factory, err := client.NewSplitFactory(cfg.SDKKey, splitCfg)
	if err != nil {
		return nil, err
	}

	splitClient := factory.Client()
	if readyErr := splitClient.BlockUntilReady(3); readyErr != nil {
		return nil, readyErr
	}

	ffm := &featureFlagManager{
		logger:      logger,
		tracer:      tracing.NewTracer(tracerProvider.Tracer("splitio_feature_flag_manager")),
		splitClient: splitClient,
	}

	return ffm, nil
}

// CanUseFeature returns whether a user can use a feature or not.
func (f *featureFlagManager) CanUseFeature(ctx context.Context, userID, feature string) (bool, error) {
	_, span := tracing.StartSpan(ctx)
	defer span.End()

	if treatment := f.splitClient.Treatment(userID, feature, nil); treatment == "on" {
		return true, nil
	}

	return false, nil
}
