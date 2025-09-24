package notificationsmock

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"

	"github.com/stretchr/testify/mock"
)

var _ notifications.Repository = (*Repository)(nil)

// Repository is a mock implementation of the notifications Repository interface.
type Repository struct {
	mock.Mock
}

// UserNotificationExists is a mock function.
func (m *Repository) UserNotificationExists(ctx context.Context, userID, userNotificationID string) (bool, error) {
	returnValues := m.Called(ctx, userID, userNotificationID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetUserNotification is a mock function.
func (m *Repository) GetUserNotification(ctx context.Context, userID, userNotificationID string) (*notifications.UserNotification, error) {
	returnValues := m.Called(ctx, userID, userNotificationID)
	return returnValues.Get(0).(*notifications.UserNotification), returnValues.Error(1)
}

// GetUserNotifications is a mock function.
func (m *Repository) GetUserNotifications(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[notifications.UserNotification], error) {
	returnValues := m.Called(ctx, userID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[notifications.UserNotification]), returnValues.Error(1)
}

// CreateUserNotification is a mock function.
func (m *Repository) CreateUserNotification(ctx context.Context, input *notifications.UserNotificationDatabaseCreationInput) (*notifications.UserNotification, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*notifications.UserNotification), returnValues.Error(1)
}

// UpdateUserNotification is a mock function.
func (m *Repository) UpdateUserNotification(ctx context.Context, updated *notifications.UserNotification) error {
	return m.Called(ctx, updated).Error(0)
}

// NewRepository creates a new mock repository.
func NewRepository() *Repository {
	return &Repository{}
}
