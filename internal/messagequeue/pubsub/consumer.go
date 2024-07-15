package pubsub

import (
	"context"
	"sync"

	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"

	"cloud.google.com/go/pubsub"
)

type (
	messageConsumer interface {
		Topic(string) *pubsub.Topic
		CreateSubscription(ctx context.Context, id string, cfg pubsub.SubscriptionConfig) (*pubsub.Subscription, error)
	}

	pubSubConsumer struct {
		logger      logging.Logger
		consumer    messageConsumer
		handlerFunc func(context.Context, []byte) error
		topic       string
	}
)

// buildPubSubConsumer provides a Pub/Sub-backed pubSubConsumer.
func buildPubSubConsumer(
	logger logging.Logger,
	pubsubClient *pubsub.Client,
	topic string,
	handlerFunc func(context.Context, []byte) error,
) messagequeue.Consumer {
	return &pubSubConsumer{
		topic:       topic,
		logger:      logging.EnsureLogger(logger),
		consumer:    pubsubClient,
		handlerFunc: handlerFunc,
	}
}

func (c *pubSubConsumer) Consume(stopChan chan bool, errors chan error) {
	if stopChan == nil {
		stopChan = make(chan bool, 1)
	}

	ctx := context.Background()
	sub, err := c.consumer.CreateSubscription(ctx, c.topic, pubsub.SubscriptionConfig{Topic: c.consumer.Topic(c.topic)})
	if err != nil {
		c.logger.Error(err, "creating subscription")
		errors <- err
		return
	}

	go func() {
		<-stopChan
		if err = sub.Delete(ctx); err != nil {
			c.logger.Error(err, "deleting subscription")
			errors <- err
		}
	}()

	if err = sub.Receive(ctx, func(receivedContext context.Context, m *pubsub.Message) {
		if handleErr := c.handlerFunc(receivedContext, m.Data); handleErr != nil {
			errors <- err
		} else {
			m.Ack()
		}
	}); err != nil {
		c.logger.Error(err, "receiving pub/sub data")
	}
}

type pubsubConsumerProvider struct {
	logger           logging.Logger
	consumerCache    map[string]messagequeue.Consumer
	pubsubClient     *pubsub.Client
	consumerCacheHat sync.RWMutex
}

// ProvidePubSubConsumerProvider returns a ConsumerProvider for a given address.
func ProvidePubSubConsumerProvider(logger logging.Logger, client *pubsub.Client) messagequeue.ConsumerProvider {
	return &pubsubConsumerProvider{
		logger:        logging.EnsureLogger(logger),
		pubsubClient:  client,
		consumerCache: map[string]messagequeue.Consumer{},
	}
}

// Close closes the connection topic.
func (p *pubsubConsumerProvider) Close() {
	if err := p.pubsubClient.Close(); err != nil {
		p.logger.Error(err, "closing pubsub connection")
	}
}

// ProvideConsumer returns a pubSubConsumer for a given topic.
func (p *pubsubConsumerProvider) ProvideConsumer(_ context.Context, topic string, handlerFunc func(context.Context, []byte) error) (messagequeue.Consumer, error) {
	if topic == "" {
		return nil, messagequeue.ErrEmptyTopicName
	}

	logger := logging.EnsureLogger(p.logger.Clone())

	p.consumerCacheHat.Lock()
	defer p.consumerCacheHat.Unlock()
	if cachedPub, ok := p.consumerCache[topic]; ok {
		return cachedPub, nil
	}

	pub := buildPubSubConsumer(logger, p.pubsubClient, topic, handlerFunc)
	p.consumerCache[topic] = pub

	return pub, nil
}
