package sqs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"sync"
	"time"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue/consumers"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

type (
	sqsConsumer struct {
		tracer      tracing.Tracer
		encoder     encoding.ClientEncoder
		logger      logging.Logger
		queueURL    string
		sqsClient   *sqs.SQS
		handlerFunc func(context.Context, []byte) error
		topic       string
	}

	// Config configures a SQS-backed consumer.
	Config struct {
		QueueAddress consumers.MessageQueueAddress `json:"messageQueueAddress" mapstructure:"message_queue_address" toml:"message_queue_address,omitempty"`
	}
)

func provideSQSConsumer(ctx context.Context, logger logging.Logger, sqsClient *sqs.SQS, topic string, handlerFunc func(context.Context, []byte) error) *sqsConsumer {
	return &sqsConsumer{
		topic:       topic,
		handlerFunc: handlerFunc,
		sqsClient:   sqsClient,
		queueURL:    topic,
		logger:      logging.EnsureLogger(logger),
		tracer:      tracing.NewTracer(fmt.Sprintf("%s_consumer", topic)),
		encoder:     encoding.ProvideClientEncoder(logger, encoding.ContentTypeJSON),
	}
}

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
				AttributeNames: []*string{
					aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
				},
				MessageAttributeNames: []*string{
					aws.String(sqs.QueueAttributeNameAll),
				},
				QueueUrl:            aws.String(r.queueURL),
				MaxNumberOfMessages: aws.Int64(1),
			})
			if err != nil {
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
	consumerCache    map[string]consumers.Consumer
	sqsClient        *sqs.SQS
	consumerCacheHat sync.RWMutex
}

// ProvideSQSConsumerProvider returns a ConsumerProvider for a given address.
func ProvideSQSConsumerProvider(logger logging.Logger, queueAddress string) consumers.ConsumerProvider {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sqs.New(sess)

	return &consumerProvider{
		logger:        logging.EnsureLogger(logger),
		sqsClient:     svc,
		consumerCache: map[string]consumers.Consumer{},
	}
}

// ProviderConsumer returns a Consumer for a given topic.
func (p *consumerProvider) ProviderConsumer(ctx context.Context, topic string, handlerFunc func(context.Context, []byte) error) (consumers.Consumer, error) {
	logger := logging.EnsureLogger(p.logger).WithValue("topic", topic)

	p.consumerCacheHat.Lock()
	defer p.consumerCacheHat.Unlock()
	if cachedPub, ok := p.consumerCache[topic]; ok {
		return cachedPub, nil
	}

	c := provideSQSConsumer(ctx, logger, p.sqsClient, topic, handlerFunc)
	p.consumerCache[topic] = c

	return c, nil
}
