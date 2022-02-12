package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"

	"github.com/stretchr/testify/mock"
)

var _ messagePublisher = (*mockMessagePublisher)(nil)

type mockMessagePublisher struct {
	mock.Mock
}

// SendMessageWithContext is a mock function.
func (m *mockMessagePublisher) Publish(ctx context.Context, msg *pubsub.Message) *pubsub.PublishResult {
	retVals := m.Called(ctx, msg)

	return retVals.Get(0).(*pubsub.PublishResult)
}
