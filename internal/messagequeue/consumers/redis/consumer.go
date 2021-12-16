package redis

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue/consumers"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

var (
	// ErrEmptyInputProvided indicates empty input was provided in an unacceptable context.
	ErrEmptyInputProvided = errors.New("empty input provided")
)

type (
	channelProvider interface {
		Channel(...redis.ChannelOption) <-chan *redis.Message
	}

	redisConsumer struct {
		tracer       tracing.Tracer
		encoder      encoding.ClientEncoder
		logger       logging.Logger
		redisClient  *redis.Client
		handlerFunc  func(context.Context, []byte) error
		subscription channelProvider
		topic        string
	}

	// Config configures a Redis-backed consumer.
	Config struct {
		QueueAddress consumers.MessageQueueAddress `json:"messageQueueAddress" mapstructure:"message_queue_address" toml:"message_queue_address,omitempty"`
	}
)

func provideRedisConsumer(ctx context.Context, logger logging.Logger, redisClient *redis.Client, tracerProvider tracing.TracerProvider, topic string, handlerFunc func(context.Context, []byte) error) *redisConsumer {
	subscription := redisClient.Subscribe(ctx, topic)

	return &redisConsumer{
		topic:        topic,
		handlerFunc:  handlerFunc,
		redisClient:  redisClient,
		subscription: subscription,
		logger:       logging.EnsureLogger(logger),
		tracer:       tracing.NewTracer(tracerProvider.Tracer(fmt.Sprintf("%s_consumer", topic))),
		encoder:      encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
	}
}

// Consume reads messages and applies the handler to their payloads.
// Writes errors to the error chan if it isn't nil.
func (r *redisConsumer) Consume(stopChan chan bool, errs chan error) {
	if stopChan == nil {
		stopChan = make(chan bool, 1)
	}
	subChan := r.subscription.Channel()

	for {
		select {
		case msg := <-subChan:
			ctx := context.Background()
			if err := r.handlerFunc(ctx, []byte(msg.Payload)); err != nil {
				r.logger.Error(err, "handling message")
				if errs != nil {
					errs <- err
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
	redisClient      *redis.Client
	tracerProvider   tracing.TracerProvider
	consumerCacheHat sync.RWMutex
}

// ProvideRedisConsumerProvider returns a ConsumerProvider for a given address.
func ProvideRedisConsumerProvider(logger logging.Logger, tracerProvider tracing.TracerProvider, queueAddress string) consumers.ConsumerProvider {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     queueAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &consumerProvider{
		logger:         logging.EnsureLogger(logger),
		redisClient:    redisClient,
		tracerProvider: tracerProvider,
		consumerCache:  map[string]consumers.Consumer{},
	}
}

// ProvideConsumer returns a Consumer for a given topic.
func (p *consumerProvider) ProvideConsumer(ctx context.Context, topic string, handlerFunc func(context.Context, []byte) error) (consumers.Consumer, error) {
	logger := logging.EnsureLogger(p.logger).WithValue("topic", topic)

	if topic == "" {
		return nil, ErrEmptyInputProvided
	}

	p.consumerCacheHat.Lock()
	defer p.consumerCacheHat.Unlock()
	if cachedPub, ok := p.consumerCache[topic]; ok {
		return cachedPub, nil
	}

	c := provideRedisConsumer(ctx, logger, p.redisClient, p.tracerProvider, topic, handlerFunc)
	p.consumerCache[topic] = c

	return c, nil
}
