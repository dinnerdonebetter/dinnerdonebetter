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

//func Test_pubsubPublisher_Publish(T *testing.T) {
//	T.SkipNow()
//
//	T.Run("standard", func(t *testing.T) {
//		t.Parallel()
//
//		ctx := context.Background()
//		logger := zerolog.NewZerologLogger()
//
//		provider := ProvidePubSubPublisherProvider(logger, trace.NewNoopTracerProvider(), t.Name())
//		require.NotNil(t, provider)
//
//		a, err := provider.ProviderPublisher(t.Name())
//		assert.NotNil(t, a)
//		assert.NoError(t, err)
//
//		actual, ok := a.(*publisher)
//		require.True(t, ok)
//
//		inputData := &struct {
//			Name string `json:"name"`
//		}{
//			Name: t.Name(),
//		}
//
//		mmp := &mockMessagePublisher{}
//		mmp.On(
//			"Publish",
//			testutils.ContextMatcher,
//			mock.MatchedBy(func(*pubsub.Message) bool { return true }),
//		).Return(&sqs.SendMessageOutput{}, nil)
//
//		actual.publisher = mmp
//
//		err = actual.Publish(ctx, inputData)
//		assert.NoError(t, err)
//
//		mock.AssertExpectationsForObjects(t, mmp)
//	})
//
//	T.Run("with error encoding value", func(t *testing.T) {
//		t.Parallel()
//
//		ctx := context.Background()
//		logger := zerolog.NewZerologLogger()
//
//		provider := ProvidePubSubPublisherProvider(ctx, logger, trace.NewNoopTracerProvider(), t.Name())
//		require.NotNil(t, provider)
//
//		a, err := provider.ProviderPublisher(t.Name())
//		assert.NotNil(t, a)
//		assert.NoError(t, err)
//
//		actual, ok := a.(*publisher)
//		require.True(t, ok)
//
//		inputData := &struct {
//			Name json.Number `json:"name"`
//		}{
//			Name: json.Number(t.Name()),
//		}
//
//		err = actual.Publish(ctx, inputData)
//		assert.Error(t, err)
//	})
//}
//
//func TestProvideSQSPublisherProvider(T *testing.T) {
//	T.Parallel()
//
//	T.Run("standard", func(t *testing.T) {
//		t.Parallel()
//
//		ctx := context.Background()
//		logger := zerolog.NewZerologLogger()
//
//		actual := ProvidePubSubPublisherProvider(ctx, logger, trace.NewNoopTracerProvider(), t.Name())
//		assert.NotNil(t, actual)
//	})
//}
//
//func Test_publisherProvider_ProviderPublisher(T *testing.T) {
//	T.Parallel()
//
//	T.Run("standard", func(t *testing.T) {
//		t.Parallel()
//
//		ctx := context.Background()
//		logger := zerolog.NewZerologLogger()
//
//		provider := ProvidePubSubPublisherProvider(ctx, logger, trace.NewNoopTracerProvider(), t.Name())
//		require.NotNil(t, provider)
//
//		actual, err := provider.ProviderPublisher(t.Name())
//		assert.NotNil(t, actual)
//		assert.NoError(t, err)
//	})
//
//	T.Run("with cache hit", func(t *testing.T) {
//		t.Parallel()
//
//		ctx := context.Background()
//		logger := zerolog.NewZerologLogger()
//
//		provider := ProvidePubSubPublisherProvider(ctx, logger, trace.NewNoopTracerProvider(), t.Name())
//		require.NotNil(t, provider)
//
//		actual, err := provider.ProviderPublisher(t.Name())
//		assert.NotNil(t, actual)
//		assert.NoError(t, err)
//
//		actual, err = provider.ProviderPublisher(t.Name())
//		assert.NotNil(t, actual)
//		assert.NoError(t, err)
//	})
//}
