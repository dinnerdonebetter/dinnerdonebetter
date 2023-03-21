package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type FeatureFlagManager struct {
	mock.Mock
}

func (m *FeatureFlagManager) CanUseFeature(ctx context.Context, username, feature string) (bool, error) {
	returnValues := m.Called(ctx, username, feature)
	return returnValues.Bool(0), returnValues.Error(1)
}
