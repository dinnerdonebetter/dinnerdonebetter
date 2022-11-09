package customerdata

import (
	"context"

	"github.com/prixfixeco/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ Collector = (*MockCollector)(nil)

type (
	// MockCollector represents a service that can collect customer data.
	MockCollector struct {
		mock.Mock
	}
)

// Close is a mock function.
func (m *MockCollector) Close() error {
	return m.Called().Error(0)
}

// AddUser is a mock function.
func (m *MockCollector) AddUser(ctx context.Context, userID string, properties map[string]interface{}) error {
	return m.Called(ctx, userID, properties).Error(0)
}

// EventOccurred is a mock function.
func (m *MockCollector) EventOccurred(ctx context.Context, event types.CustomerEventType, userID string, properties map[string]interface{}) error {
	return m.Called(ctx, event, userID, properties).Error(0)
}
