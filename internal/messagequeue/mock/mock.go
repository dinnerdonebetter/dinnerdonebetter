package mockpublishers

import (
	"context"

	"github.com/prixfixeco/backend/internal/messagequeue"

	"github.com/stretchr/testify/mock"
)

var _ messagequeue.Publisher = (*Publisher)(nil)

// Publisher implements our interface.
type Publisher struct {
	mock.Mock
}

// Publish implements our interface.
func (m *Publisher) Publish(ctx context.Context, data any) error {
	return m.Called(ctx, data).Error(0)
}

// Stop implements our interface.
func (m *Publisher) Stop() {
	m.Called()
}

// ProducerProvider implements our interface.
type ProducerProvider struct {
	mock.Mock
}

// ProviderPublisher implements our interface.
func (m *ProducerProvider) ProviderPublisher(topic string) (messagequeue.Publisher, error) {
	args := m.Called(topic)

	return args.Get(0).(messagequeue.Publisher), args.Error(1)
}
