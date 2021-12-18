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
