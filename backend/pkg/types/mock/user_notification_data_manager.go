package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.UserNotificationDataManager = (*UserNotificationDataManagerMock)(nil)

// UserNotificationDataManagerMock is a mocked types.UserNotificationDataManager for testing.
type UserNotificationDataManagerMock struct {
	mock.Mock
}

// UserNotificationExists is a mock function.
func (m *UserNotificationDataManagerMock) UserNotificationExists(ctx context.Context, userID, userNotificationID string) (bool, error) {
	returnValues := m.Called(ctx, userID, userNotificationID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetUserNotification is a mock function.
func (m *UserNotificationDataManagerMock) GetUserNotification(ctx context.Context, userID, userNotificationID string) (*types.UserNotification, error) {
	returnValues := m.Called(ctx, userID, userNotificationID)
	return returnValues.Get(0).(*types.UserNotification), returnValues.Error(1)
}

// GetUserNotifications is a mock function.
func (m *UserNotificationDataManagerMock) GetUserNotifications(ctx context.Context, userID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.UserNotification], error) {
	returnValues := m.Called(ctx, userID, filter)
	return returnValues.Get(0).(*types.QueryFilteredResult[types.UserNotification]), returnValues.Error(1)
}

// CreateUserNotification is a mock function.
func (m *UserNotificationDataManagerMock) CreateUserNotification(ctx context.Context, input *types.UserNotificationDatabaseCreationInput) (*types.UserNotification, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.UserNotification), returnValues.Error(1)
}

// UpdateUserNotification is a mock function.
func (m *UserNotificationDataManagerMock) UpdateUserNotification(ctx context.Context, updated *types.UserNotification) error {
	return m.Called(ctx, updated).Error(0)
}
