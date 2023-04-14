package featureflags

import "context"

type (
	// FeatureFlagManager manages feature flags.
	FeatureFlagManager interface {
		CanUseFeature(ctx context.Context, userID, feature string) (bool, error)
	}
)

func NewNoopFeatureFlagManager() FeatureFlagManager {
	return &NoopFeatureFlagManager{}
}

// NoopFeatureFlagManager is a no-op FeatureFlagManager.
type NoopFeatureFlagManager struct{}

func (*NoopFeatureFlagManager) CanUseFeature(_ context.Context, _, _ string) (bool, error) {
	return false, nil
}
