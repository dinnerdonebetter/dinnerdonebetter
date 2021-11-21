package customerdata

import (
	"context"

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

// Identify is a mock function.
func (m *MockCollector) Identify(ctx context.Context, userID string, properties map[string]interface{}) error {
	return m.Called(ctx, userID, properties).Error(0)
}

// Track is a mock function.
func (m *MockCollector) Track(ctx context.Context, event, userID string, properties map[string]interface{}) error {
	return m.Called(ctx, event, userID, properties).Error(0)
}
