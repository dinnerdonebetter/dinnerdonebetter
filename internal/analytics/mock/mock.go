package mock

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

// EventReporter is a mock implementation of analytics.EventReporter.
type EventReporter struct {
	mock.Mock
}

// Close implements the EventReporter interface.
func (m *EventReporter) Close() {
	m.Called()
}

// AddUser implements the EventReporter interface.
func (m *EventReporter) AddUser(ctx context.Context, userID string, properties map[string]any) error {
	return m.Called(ctx, userID, properties).Error(0)
}

// EventOccurred implements the EventReporter interface.
func (m *EventReporter) EventOccurred(ctx context.Context, event types.CustomerEventType, userID string, properties map[string]any) error {
	return m.Called(ctx, event, userID, properties).Error(0)
}
