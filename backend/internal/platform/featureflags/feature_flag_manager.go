package featureflags

import (
	"context"
)

type (
	// FeatureFlagManager manages feature flags.
	FeatureFlagManager interface {
		CanUseFeature(ctx context.Context, userID, feature string) (bool, error)
		Close() error
	}
)

func NewNoopFeatureFlagManager() FeatureFlagManager {
	return &NoopFeatureFlagManager{}
}

// NoopFeatureFlagManager is a no-op FeatureFlagManager.
type NoopFeatureFlagManager struct{}

// CanUseFeature implements the FeatureFlagManager interface.
func (*NoopFeatureFlagManager) CanUseFeature(context.Context, string, string) (bool, error) {
	return false, nil
}

// Close implements the FeatureFlagManager interface.
func (*NoopFeatureFlagManager) Close() error {
	return nil
}
