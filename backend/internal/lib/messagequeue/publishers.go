package messagequeue

import (
	"context"
	"errors"
)

type (
	// Publisher produces events onto a queue.
	Publisher interface {
		// Stop halts all publishing.
		Stop()
		// Publish writes a message onto a message queue.
		Publish(ctx context.Context, data any) error
		// PublishAsync writes a message onto a message queue, but logs any encountered errors instead of returning them.
		PublishAsync(ctx context.Context, data any)
	}

	// PublisherProvider is a function that provides a Publisher for a given topic.
	PublisherProvider interface {
		Close()
		ProvidePublisher(topic string) (Publisher, error)
	}
)

var (
	// ErrEmptyTopicName is returned when a topic name is empty.
	ErrEmptyTopicName = errors.New("empty topic name")
)

type noopPublisher struct{}

// Publish does nothing.
func (n *noopPublisher) Publish(context.Context, any) error { return nil }

// PublishAsync does nothing.
func (n *noopPublisher) PublishAsync(context.Context, any) {}

// Stop does nothing.
func (n *noopPublisher) Stop() {}

// NewNoopPublisher is a noop Publisher.
func NewNoopPublisher() Publisher {
	return &noopPublisher{}
}

type noopPublisherProvider struct{}

// ProvidePublisher does nothing.
func (n *noopPublisherProvider) ProvidePublisher(_ string) (Publisher, error) {
	return NewNoopPublisher(), nil
}

// Close does nothing.
func (n *noopPublisherProvider) Close() {}

// NewNoopPublisherProvider is a noop PublisherProvider.
func NewNoopPublisherProvider() PublisherProvider {
	return &noopPublisherProvider{}
}
