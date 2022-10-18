package messagequeue

import (
	"context"
)

type (
	// Consumer produces events onto a queue.
	Consumer interface {
		Consume(stopChan chan bool, errors chan error)
	}

	// ConsumerProvider is a function that provides a Consumer for a given topic.
	ConsumerProvider interface {
		ProvideConsumer(ctx context.Context, topic string, handlerFunc func(context.Context, []byte) error) (Consumer, error)
	}
)

type noopConsumer struct{}

// Consume does nothing.
func (n *noopConsumer) Consume(_ chan bool, _ chan error) {}

// NewNoopConsumer is a noop Consumer.
func NewNoopConsumer() Consumer {
	return &noopConsumer{}
}

type noopConsumerProvider struct{}

// ProvideConsumer does nothing.
func (n *noopConsumerProvider) ProvideConsumer(_ context.Context, _ string, _ func(context.Context, []byte) error) (Consumer, error) {
	return NewNoopConsumer(), nil
}

// NewNoopConsumerProvider is a noop ConsumerProvider.
func NewNoopConsumerProvider() ConsumerProvider {
	return &noopConsumerProvider{}
}
