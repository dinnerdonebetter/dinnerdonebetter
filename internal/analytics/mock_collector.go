package analytics

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ EventReporter = (*MockEventReporter)(nil)

type (
	// MockEventReporter represents a service that can collect customer data.
	MockEventReporter struct {
		mock.Mock
	}
)

// Close is a mock function.
func (m *MockEventReporter) Close() {
	m.Called()
}

// AddUser is a mock function.
func (m *MockEventReporter) AddUser(ctx context.Context, userID string, properties map[string]any) error {
	return m.Called(ctx, userID, properties).Error(0)
}

// EventOccurred is a mock function.
func (m *MockEventReporter) EventOccurred(ctx context.Context, event types.CustomerEventType, userID string, properties map[string]any) error {
	return m.Called(ctx, event, userID, properties).Error(0)
}
