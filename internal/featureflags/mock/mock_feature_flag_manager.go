package mock

import (
	"context"

	"github.com/prixfixeco/backend/internal/featureflags"
	"github.com/prixfixeco/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ featureflags.FeatureFlagManager = (*FeatureFlagManager)(nil)

type FeatureFlagManager struct {
	mock.Mock
}

func (m *FeatureFlagManager) CanUseFeature(ctx context.Context, username, feature string) (bool, error) {
	returnValues := m.Called(ctx, username, feature)
	return returnValues.Bool(0), returnValues.Error(1)
}

func (m *FeatureFlagManager) Identify(ctx context.Context, user *types.User) error {
	return m.Called(ctx, user).Error(0)
}
