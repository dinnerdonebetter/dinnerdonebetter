package featureflags

import "context"

type (
	// FeatureFlagManager manages feature flags.
	FeatureFlagManager interface {
		CanUseFeature(ctx context.Context, username, feature string) (bool, error)
	}
)
