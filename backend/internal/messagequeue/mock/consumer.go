package mockpublishers

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/messagequeue"

	"github.com/stretchr/testify/mock"
)

var _ messagequeue.ConsumerProvider = (*ConsumerProvider)(nil)

// ConsumerProvider is a mock consumers.ConsumerProvider.
type ConsumerProvider struct {
	mock.Mock
}

// ProvideConsumer implements the interface.
func (m *ConsumerProvider) ProvideConsumer(ctx context.Context, topic string, handlerFunc func(context.Context, []byte) error) (messagequeue.Consumer, error) {
	args := m.Called(ctx, topic, handlerFunc)

	return args.Get(0).(messagequeue.Consumer), args.Error(1)
}

var _ messagequeue.Consumer = (*Consumer)(nil)

// Consumer is a mock consumers.Consumer.
type Consumer struct {
	mock.Mock
}

// Consume implements the interface.
func (m *Consumer) Consume(stopChan chan bool, errors chan error) {
	m.Called(stopChan, errors)
}
