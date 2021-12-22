package redis

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/prixfixeco/api_server/internal/messagequeue"

	"github.com/go-redis/redis/v8"

	"github.com/prixfixeco/api_server/internal/encoding"
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

	var b bytes.Buffer
	if err := r.encoder.Encode(ctx, &b, data); err != nil {
		return observability.PrepareError(err, r.logger, span, "encoding topic message")
	}

	return r.publisher.Publish(ctx, r.topic, b.Bytes()).Err()
}

// provideRedisPublisher provides a redis-backed Publisher.
func provideRedisPublisher(logger logging.Logger, tracerProvider tracing.TracerProvider, redisClient messagePublisher, topic string) *redisPublisher {
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
	publisherCache    map[string]messagequeue.Publisher
	redisClient       messagePublisher
	tracerProvider    tracing.TracerProvider
	publisherCacheHat sync.RWMutex
}

// ProvideRedisPublisherProvider returns a PublisherProvider for a given address.
func ProvideRedisPublisherProvider(logger logging.Logger, tracerProvider tracing.TracerProvider, cfg Config) messagequeue.PublisherProvider {
	logger.WithValue("queue_addresses", cfg.QueueAddresses).
		WithValue("username", cfg.Username).
		WithValue("password", cfg.Password).Info("setting up redis publisher")

	var redisClient messagePublisher
	if len(cfg.QueueAddresses) > 1 {
		redisClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        cfg.QueueAddresses,
			Username:     cfg.Username,
			Password:     cfg.Password,
			DialTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		})
	} else if len(cfg.QueueAddresses) == 1 {
		redisClient = redis.NewClient(&redis.Options{
			Addr:         cfg.QueueAddresses[0],
			Username:     cfg.Username,
			Password:     cfg.Password,
			DialTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		})
	}

	return &publisherProvider{
		logger:         logging.EnsureLogger(logger),
		redisClient:    redisClient,
		publisherCache: map[string]messagequeue.Publisher{},
		tracerProvider: tracerProvider,
	}
}

// ProviderPublisher returns a Publisher for a given topic.
func (p *publisherProvider) ProviderPublisher(topic string) (messagequeue.Publisher, error) {
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
