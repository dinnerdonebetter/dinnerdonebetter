package launchdarkly

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/platform/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/platform/featureflags"
	"github.com/dinnerdonebetter/backend/internal/platform/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	ld "github.com/launchdarkly/go-server-sdk/v6"
	"github.com/launchdarkly/go-server-sdk/v6/ldcomponents"
	ofld "github.com/open-feature/go-sdk-contrib/providers/launchdarkly/pkg"
	"github.com/open-feature/go-sdk/openfeature"
)

const (
	serviceName  = "launchdarkly_feature_flag_manager"
	clientDomain = "launchdarkly_feature_flags"
)

var (
	ErrMissingHTTPClient = errors.New("missing HTTP client")
	ErrNilConfig         = errors.New("missing config")
	ErrMissingSDKKey     = errors.New("missing SDK key")
)

type (
	// featureFlagManager implements the feature flag interface using OpenFeature.
	featureFlagManager struct {
		ldClient       *ld.LDClient
		ofClient       *openfeature.Client
		circuitBreaker circuitbreaking.CircuitBreaker
		logger         logging.Logger
		tracer         tracing.Tracer
	}
)

// NewFeatureFlagManager constructs a new featureFlagManager backed by OpenFeature.
func NewFeatureFlagManager(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider, httpClient *http.Client, circuitBreaker circuitbreaking.CircuitBreaker, configModifiers ...func(ld.Config) ld.Config) (featureflags.FeatureFlagManager, error) {
	if httpClient == nil {
		return nil, ErrMissingHTTPClient
	}

	if cfg == nil {
		return nil, ErrNilConfig
	}

	cfg.CircuitBreakerConfig.EnsureDefaults()

	if cfg.SDKKey == "" {
		return nil, ErrMissingSDKKey
	}

	if cfg.InitTimeout == time.Duration(0) {
		cfg.InitTimeout = 5 * time.Second
	}

	ldConfig := ld.Config{
		HTTP: ldcomponents.HTTPConfiguration().HTTPClientFactory(func() *http.Client { return httpClient }),
	}

	for _, modifier := range configModifiers {
		ldConfig = modifier(ldConfig)
	}

	client, err := ld.MakeCustomClient(
		cfg.SDKKey,
		ldConfig,
		cfg.InitTimeout,
	)
	if err != nil {
		return nil, fmt.Errorf("error initializing LaunchDarkly client: %w", err)
	}

	provider := ofld.NewProvider(client)
	if err = openfeature.SetNamedProviderAndWait(clientDomain, provider); err != nil {
		if closeErr := client.Close(); closeErr != nil {
			logger.Error("error closing OpenFeatureFlag client", closeErr)
		}
		return nil, fmt.Errorf("failed to set OpenFeature provider: %w", err)
	}

	ofClient := openfeature.NewClient(clientDomain)

	ffm := &featureFlagManager{
		logger:         logging.EnsureLogger(logger),
		circuitBreaker: circuitBreaker,
		tracer:         tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		ldClient:       client,
		ofClient:       ofClient,
	}

	return ffm, nil
}

// CanUseFeature returns whether a user can use a feature or not.
func (f *featureFlagManager) CanUseFeature(ctx context.Context, userID, feature string) (bool, error) {
	_, span := f.tracer.StartSpan(ctx)
	defer span.End()

	logger := f.logger.WithValue(keys.UserIDKey, userID).WithValue("feature", feature)

	if !f.circuitBreaker.CanProceed() {
		return false, internalerrors.ErrCircuitBroken
	}

	evalCtx := openfeature.NewEvaluationContext(userID, nil)
	result, err := f.ofClient.BooleanValue(ctx, feature, false, evalCtx)
	if err != nil {
		f.circuitBreaker.Failed()
		return false, observability.PrepareAndLogError(err, logger, span, "checking feature flag variation")
	}

	f.circuitBreaker.Succeeded()
	return result, nil
}

// Close closes the LaunchDarkly client.
func (f *featureFlagManager) Close() error {
	return f.ldClient.Close()
}
