package sqs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/prixfixeco/api_server/internal/messagequeue"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

type (
	sqsConsumer struct {
		tracer             tracing.Tracer
		encoder            encoding.ClientEncoder
		logger             logging.Logger
		sqsClient          *sqs.SQS
		handlerFunc        func(context.Context, []byte) error
		queueURL           string
		topic              string
		messagesPerReceive uint8
	}
)

// Consume reads messages and applies the handler to their payloads.
// Writes errors to the error chan if it isn't nil.
func (r *sqsConsumer) Consume(stopChan chan bool, errors chan error) {
	pollInterval := time.NewTicker(time.Second)
	if stopChan == nil {
		stopChan = make(chan bool, 1)
	}

	for {
		select {
		case <-pollInterval.C:
			ctx := context.Background()
			msgResult, err := r.sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
				AttributeNames:        []*string{aws.String(sqs.QueueAttributeNameAll)},
				MessageAttributeNames: []*string{aws.String(sqs.QueueAttributeNameAll)},
				QueueUrl:              aws.String(r.queueURL),
				MaxNumberOfMessages:   aws.Int64(int64(r.messagesPerReceive)),
			})
			if err != nil {
				r.logger.Error(err, "receiving SQS message")
				errors <- err
				continue
			}

			for _, msg := range msgResult.Messages {
				if msg.Body != nil {
					if err = r.handlerFunc(ctx, []byte(*msg.Body)); err != nil {
						r.logger.Error(err, "handling SQS message")
						if errors != nil {
							errors <- err
						}
					} else {
						_, deleteErr := r.sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
							QueueUrl:      aws.String(r.queueURL),
							ReceiptHandle: msg.ReceiptHandle,
						})
						if deleteErr != nil {
							r.logger.Error(err, "deleting SQS message")
							errors <- err
						}
					}
				}
			}
		case <-stopChan:
			return
		}
	}
}

type consumerProvider struct {
	logger           logging.Logger
	consumerCache    map[string]messagequeue.Consumer
	sqsClient        *sqs.SQS
	tracerProvider   tracing.TracerProvider
	consumerCacheHat sync.RWMutex
}

var _ messagequeue.ConsumerProvider = (*consumerProvider)(nil)

// ProvideSQSConsumerProvider returns a ConsumerProvider.
func ProvideSQSConsumerProvider(logger logging.Logger, tracerProvider tracing.TracerProvider) messagequeue.ConsumerProvider {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sqs.New(sess)

	return &consumerProvider{
		logger:         logging.EnsureLogger(logger),
		sqsClient:      svc,
		consumerCache:  map[string]messagequeue.Consumer{},
		tracerProvider: tracerProvider,
	}
}

// ProvideConsumer returns a Consumer for a given topic.
func (p *consumerProvider) ProvideConsumer(ctx context.Context, topic string, handlerFunc func(context.Context, []byte) error) (messagequeue.Consumer, error) {
	logger := logging.EnsureLogger(p.logger).Clone().WithValue("topic", topic)

	p.consumerCacheHat.Lock()
	defer p.consumerCacheHat.Unlock()
	if cachedPub, ok := p.consumerCache[topic]; ok {
		return cachedPub, nil
	}

	c := provideSQSConsumer(ctx, logger, p.sqsClient, p.tracerProvider, topic, handlerFunc)
	p.consumerCache[topic] = c

	return c, nil
}

func provideSQSConsumer(_ context.Context, logger logging.Logger, sqsClient *sqs.SQS, tracerProvider tracing.TracerProvider, topic string, handlerFunc func(context.Context, []byte) error) *sqsConsumer {
	return &sqsConsumer{
		topic:              topic,
		handlerFunc:        handlerFunc,
		sqsClient:          sqsClient,
		queueURL:           topic,
		messagesPerReceive: 10, // max value is 10
		logger:             logging.EnsureLogger(logger).Clone().WithName(fmt.Sprintf("%s_consumer", topic)),
		tracer:             tracing.NewTracer(tracerProvider.Tracer(fmt.Sprintf("%s_consumer", topic))),
		encoder:            encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
	}
}
