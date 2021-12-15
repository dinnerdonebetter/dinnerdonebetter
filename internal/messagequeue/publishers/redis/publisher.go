package redis

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue/publishers"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

var (
	// ErrEmptyInputProvided indicates empty input was provided in an unacceptable context.
	ErrEmptyInputProvided = errors.New("empty input provided")
)

type (
	messagePublisher interface {
		Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
	}

	redisPublisher struct {
		tracer    tracing.Tracer
		encoder   encoding.ClientEncoder
		logger    logging.Logger
		publisher messagePublisher
		topic     string
	}
)

func (r *redisPublisher) Publish(ctx context.Context, data interface{}) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	r.logger.Debug("publishing message")

	var b bytes.Buffer
	if err := r.encoder.Encode(ctx, &b, data); err != nil {
		return observability.PrepareError(err, r.logger, span, "encoding topic message")
	}

	return r.publisher.Publish(ctx, r.topic, b.Bytes()).Err()
}

// provideRedisPublisher provides a redis-backed Publisher.
func provideRedisPublisher(logger logging.Logger, tracerProvider trace.TracerProvider, redisClient *redis.Client, topic string) *redisPublisher {
	return &redisPublisher{
		publisher: redisClient,
		topic:     topic,
		encoder:   encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		logger:    logging.EnsureLogger(logger),
		tracer:    tracing.NewTracer(tracerProvider.Tracer(fmt.Sprintf("%s_publisher", topic))),
	}
}

type publisherProvider struct {
	logger            logging.Logger
	publisherCache    map[string]publishers.Publisher
	redisClient       *redis.Client
	tracerProvider    trace.TracerProvider
	publisherCacheHat sync.RWMutex
}

// ProvideRedisPublisherProvider returns a PublisherProvider for a given address.
func ProvideRedisPublisherProvider(logger logging.Logger, tracerProvider trace.TracerProvider, address string) publishers.PublisherProvider {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &publisherProvider{
		logger:         logging.EnsureLogger(logger),
		redisClient:    redisClient,
		publisherCache: map[string]publishers.Publisher{},
		tracerProvider: tracerProvider,
	}
}

// ProviderPublisher returns a Publisher for a given topic.
func (p *publisherProvider) ProviderPublisher(topic string) (publishers.Publisher, error) {
	logger := logging.EnsureLogger(p.logger).WithValue("topic", topic)

	if topic == "" {
		return nil, ErrEmptyInputProvided
	}

	p.publisherCacheHat.Lock()
	defer p.publisherCacheHat.Unlock()
	if cachedPub, ok := p.publisherCache[topic]; ok {
		return cachedPub, nil
	}

	pub := provideRedisPublisher(logger, p.tracerProvider, p.redisClient, topic)
	p.publisherCache[topic] = pub

	return pub, nil
}
