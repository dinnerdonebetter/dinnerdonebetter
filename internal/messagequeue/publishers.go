package messagequeue

import (
	"context"
)

type (
	// Publisher produces events onto a queue.
	Publisher interface {
		Publish(ctx context.Context, data interface{}) error
	}

	// PublisherProvider is a function that provides a Publisher for a given topic.
	PublisherProvider interface {
		ProviderPublisher(topic string) (Publisher, error)
	}
)

type noopPublisher struct{}

// Publish does nothing.
func (n *noopPublisher) Publish(_ context.Context, _ interface{}) error { return nil }

// NewNoopPublisher is a noop Publisher.
func NewNoopPublisher() Publisher {
	return &noopPublisher{}
}

type noopPublisherProvider struct{}

// ProviderPublisher does nothing.
func (n *noopPublisherProvider) ProviderPublisher(_ string) (Publisher, error) {
	return NewNoopPublisher(), nil
}

// NewNoopPublisherProvider is a noop PublisherProvider.
func NewNoopPublisherProvider() PublisherProvider {
	return &noopPublisherProvider{}
}
