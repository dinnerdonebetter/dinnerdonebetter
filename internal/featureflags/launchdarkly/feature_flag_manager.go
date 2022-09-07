package launchdarkly

import (
	"context"
	"fmt"

	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"

	"github.com/prixfixeco/api_server/internal/featureflags"
	"github.com/prixfixeco/api_server/internal/featureflags/config"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

var (
	_ featureflags.FeatureFlagManager = (*FeatureFlagManager)(nil)
	_ launchDarklyClient              = (*ld.LDClient)(nil)
)

type (
	launchDarklyClient interface {
		BoolVariation(key string, user lduser.User, defaultVal bool) (bool, error)
	}

	// FeatureFlagManager implements the feature flag interface.
	FeatureFlagManager struct {
		tracer tracing.Tracer
		logger logging.Logger
		client launchDarklyClient
	}
)

// NewFeatureFlagManager constructs a new FeatureFlagManager.
func NewFeatureFlagManager(cfg config.Config, logger logging.Logger, tracerProvider tracing.TracerProvider) (*FeatureFlagManager, error) {
	client, err := ld.MakeClient(cfg.LaunchDarkly.SDKKey, cfg.LaunchDarkly.InitTimeout)
	if err != nil {
		return nil, fmt.Errorf("error initializing Launch Darkly client: %w", err)
	}

	ffm := &FeatureFlagManager{
		client: client,
		logger: logging.EnsureLogger(logger).WithName("launchdarkly_feature_flag_manager"),
		tracer: tracing.NewTracer(tracerProvider.Tracer("launchdarkly_feature_flag_manager")),
	}

	return ffm, nil
}

// CanUseFeature returns whether a user can use a feature or not.
func (f *FeatureFlagManager) CanUseFeature(ctx context.Context, username, feature string) (bool, error) {
	_, span := f.tracer.StartSpan(ctx)
	defer span.End()

	logger := f.logger.WithValues(map[string]interface{}{
		keys.UsernameKey: username,
		"feature":        feature,
	})

	user := lduser.NewUserBuilder(username).
		Name(username).
		Build()

	showFeature, err := f.client.BoolVariation(feature, user, false)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "fetching variation")
	}

	return showFeature, nil
}
