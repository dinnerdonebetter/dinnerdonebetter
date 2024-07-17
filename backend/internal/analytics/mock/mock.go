package analyticsmock

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/analytics"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ analytics.EventReporter = (*EventReporter)(nil)

type (
	// EventReporter represents a service that can collect customer data.
	EventReporter struct {
		mock.Mock
	}
)

// Close implements the EventReporter interface.
func (m *EventReporter) Close() {
	m.Called()
}

// AddUser implements the EventReporter interface.
func (m *EventReporter) AddUser(ctx context.Context, userID string, properties map[string]any) error {
	return m.Called(ctx, userID, properties).Error(0)
}

// EventOccurred implements the EventReporter interface.
func (m *EventReporter) EventOccurred(ctx context.Context, event types.ServiceEventType, userID string, properties map[string]any) error {
	return m.Called(ctx, event, userID, properties).Error(0)
}
