package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.DataPrivacyDataManager = (*DataPrivacyDataManagerMock)(nil)

// DataPrivacyDataManagerMock is a mocked types.DataPrivacyDataManager for testing.
type DataPrivacyDataManagerMock struct {
	mock.Mock
}

func (m *DataPrivacyDataManagerMock) AggregateUserData(ctx context.Context, userID string) (*types.UserDataCollection, error) {
	returnArgs := m.Called(ctx, userID)

	return returnArgs.Get(0).(*types.UserDataCollection), returnArgs.Error(1)
}

// DeleteUser is a mock function.
func (m *DataPrivacyDataManagerMock) DeleteUser(ctx context.Context, userID string) error {
	return m.Called(ctx, userID).Error(0)
}
