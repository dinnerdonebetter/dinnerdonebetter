package sqs

import (
	"context"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockMessageReceiver struct {
	mock.Mock
}

func (m *mockMessageReceiver) ReceiveMessageWithContext(ctx aws.Context, input *sqs.ReceiveMessageInput, opts ...request.Option) (*sqs.ReceiveMessageOutput, error) {
	retVals := m.Called(ctx, input, opts)
	out := retVals.Get(0)
	if out == nil {
		return nil, retVals.Error(1)
	}
	return out.(*sqs.ReceiveMessageOutput), retVals.Error(1)
}

func (m *mockMessageReceiver) DeleteMessageWithContext(ctx aws.Context, input *sqs.DeleteMessageInput, opts ...request.Option) (*sqs.DeleteMessageOutput, error) {
	retVals := m.Called(ctx, input, opts)
	out := retVals.Get(0)
	if out == nil {
		return nil, retVals.Error(1)
	}
	return out.(*sqs.DeleteMessageOutput), retVals.Error(1)
}

func Test_sqsConsumer_Consume(T *testing.T) {
	T.Parallel()

	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789/test-queue"

	T.Run("successful message handling and deletion", func(t *testing.T) {
		t.Parallel()

		mmr := &mockMessageReceiver{}
		mmr.On(
			"ReceiveMessageWithContext",
			testutils.ContextMatcher,
			mock.MatchedBy(func(in *sqs.ReceiveMessageInput) bool {
				return aws.StringValue(in.QueueUrl) == queueURL &&
					aws.Int64Value(in.MaxNumberOfMessages) == maxNumberOfMessages &&
					aws.Int64Value(in.WaitTimeSeconds) == longPollWaitSeconds
			}),
			mock.Anything,
		).Return(&sqs.ReceiveMessageOutput{
			Messages: []*sqs.Message{
				{
					Body:          aws.String("test-payload"),
					ReceiptHandle: aws.String("receipt-handle-123"),
				},
			},
		}, nil).Once()

		mmr.On(
			"ReceiveMessageWithContext",
			testutils.ContextMatcher,
			mock.Anything,
			mock.Anything,
		).Return(&sqs.ReceiveMessageOutput{Messages: []*sqs.Message{}}, nil)

		mmr.On(
			"DeleteMessageWithContext",
			testutils.ContextMatcher,
			mock.MatchedBy(func(in *sqs.DeleteMessageInput) bool {
				return aws.StringValue(in.QueueUrl) == queueURL &&
					aws.StringValue(in.ReceiptHandle) == "receipt-handle-123"
			}),
			mock.Anything,
		).Return(&sqs.DeleteMessageOutput{}, nil).Once()

		handlerDone := make(chan []byte, 1)
		handler := func(ctx context.Context, body []byte) error {
			handlerDone <- body
			return nil
		}

		consumer := provideSQSConsumer(logging.NewNoopLogger(), mmr, queueURL, handler)
		stopChan := make(chan bool, 1)
		errs := make(chan error, 4)

		go consumer.Consume(stopChan, errs)

		receivedBody := <-handlerDone
		stopChan <- true

		assert.Equal(t, []byte("test-payload"), receivedBody)
		mock.AssertExpectationsForObjects(t, mmr)
	})

	T.Run("handler error does not delete message", func(t *testing.T) {
		t.Parallel()

		anticipatedErr := errors.New("handler failed")
		mmr := &mockMessageReceiver{}
		mmr.On(
			"ReceiveMessageWithContext",
			testutils.ContextMatcher,
			mock.Anything,
			mock.Anything,
		).Return(&sqs.ReceiveMessageOutput{
			Messages: []*sqs.Message{
				{
					Body:          aws.String("fail-payload"),
					ReceiptHandle: aws.String("receipt-handle-456"),
				},
			},
		}, nil).Once()

		mmr.On(
			"ReceiveMessageWithContext",
			testutils.ContextMatcher,
			mock.Anything,
			mock.Anything,
		).Return(&sqs.ReceiveMessageOutput{Messages: []*sqs.Message{}}, nil)

		handler := func(ctx context.Context, body []byte) error {
			return anticipatedErr
		}

		consumer := provideSQSConsumer(logging.NewNoopLogger(), mmr, queueURL, handler)
		stopChan := make(chan bool, 1)
		errs := make(chan error, 4)

		go consumer.Consume(stopChan, errs)

		receivedErr := <-errs
		assert.Error(t, receivedErr)
		assert.Equal(t, anticipatedErr, receivedErr)

		stopChan <- true

		mmr.AssertNotCalled(t, "DeleteMessageWithContext")
		mock.AssertExpectationsForObjects(t, mmr)
	})
}

func TestProvideSQSConsumerProvider(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		cfg := Config{}

		actual := ProvideSQSConsumerProvider(logger, cfg)
		assert.NotNil(t, actual)
	})
}

func Test_consumerProvider_ProvideConsumer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		logger := logging.NewNoopLogger()
		cfg := Config{}

		provider := ProvideSQSConsumerProvider(logger, cfg)
		require.NotNil(t, provider)

		actual, err := provider.ProvideConsumer(ctx, "https://sqs.us-east-1.amazonaws.com/123/test", nil)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	T.Run("with cache hit", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		logger := logging.NewNoopLogger()
		cfg := Config{}
		topic := "https://sqs.us-east-1.amazonaws.com/123/cached-queue"

		provider := ProvideSQSConsumerProvider(logger, cfg)
		require.NotNil(t, provider)

		actual, err := provider.ProvideConsumer(ctx, topic, nil)
		assert.NoError(t, err)
		assert.NotNil(t, actual)

		actual2, err := provider.ProvideConsumer(ctx, topic, nil)
		assert.NoError(t, err)
		assert.NotNil(t, actual2)
		assert.Same(t, actual, actual2)
	})

	T.Run("with empty topic returns error", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		logger := logging.NewNoopLogger()
		cfg := Config{}

		provider := ProvideSQSConsumerProvider(logger, cfg)
		require.NotNil(t, provider)

		actual, err := provider.ProvideConsumer(ctx, "", nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, messagequeue.ErrEmptyTopicName)
	})
}
