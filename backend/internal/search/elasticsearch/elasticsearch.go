package elasticsearch

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/search"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var (
	_ search.Index[types.UserSearchSubset] = (*indexManager[types.UserSearchSubset])(nil)
)

type (
	indexManager[T search.Searchable] struct {
		logger                logging.Logger
		tracer                tracing.Tracer
		circuitBreaker        circuitbreaking.CircuitBreaker
		esClient              *elasticsearch.Client
		indexName             string
		indexOperationTimeout time.Duration
	}
)

func provideElasticsearchClient(cfg *Config) (*elasticsearch.Client, error) {
	c, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			cfg.Address,
		},
		Username:      cfg.Username,
		Password:      cfg.Password,
		CACert:        cfg.CACert,
		RetryOnStatus: nil,
		MaxRetries:    10,
		Transport:     nil,
		Logger:        nil,
	})
	if err != nil {
		return nil, fmt.Errorf("initializing search client: %w", err)
	}

	return c, nil
}

func ProvideIndexManager[T search.Searchable](ctx context.Context, logger logging.Logger, tracerProvider tracing.TracerProvider, cfg *Config, indexName string, circuitBreaker circuitbreaking.CircuitBreaker) (search.Index[T], error) {
	c, err := provideElasticsearchClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("initializing search client: %w", err)
	}

	logger = logging.EnsureLogger(logger)

	if ready := elasticsearchIsReadyToInit(ctx, cfg, logger, 10); !ready {
		return nil, fmt.Errorf("initializing search client: %w", err)
	}

	im := &indexManager[T]{
		tracer:                tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(fmt.Sprintf("search_%s", indexName))),
		logger:                logging.EnsureLogger(logger).WithName(indexName),
		esClient:              c,
		indexOperationTimeout: cfg.IndexOperationTimeout,
		indexName:             indexName,
		circuitBreaker:        circuitBreaker,
	}

	if indexErr := im.ensureIndices(ctx); indexErr != nil {
		return nil, indexErr
	}

	return im, nil
}

func elasticsearchIsReadyToInit(
	ctx context.Context,
	cfg *Config,
	l logging.Logger,
	maxAttempts uint8,
) bool {
	attemptCount := 0

	logger := l.WithValues(map[string]any{
		"interval": time.Second.String(),
		"address":  cfg.Address,
	})

	logger.Debug("checking if elasticsearch is ready")

	c, err := provideElasticsearchClient(cfg)
	if err != nil {
		logger.WithValue("attempt_count", attemptCount).Debug("client setup failed, waiting for elasticsearch")
	}

	for {
		var res *esapi.Response
		res, err = (esapi.InfoRequest{}).Do(ctx, c)
		if err != nil && res != nil && !res.IsError() {
			logger.WithValue("attempt_count", attemptCount).Debug("ping failed, waiting for elasticsearch")
			time.Sleep(time.Second)

			attemptCount++
			if attemptCount >= int(maxAttempts) {
				break
			}
		} else {
			return true
		}
	}

	return false
}

func (sm *indexManager[T]) ensureIndices(ctx context.Context) error {
	_, span := sm.tracer.StartSpan(ctx)
	defer span.End()

	if sm.circuitBreaker.CannotProceed() {
		return types.ErrCircuitBroken
	}

	res, err := esapi.IndicesExistsRequest{
		Index:             []string{sm.indexName},
		IgnoreUnavailable: esapi.BoolPtr(false),
		ErrorTrace:        false,
		FilterPath:        nil,
	}.Do(ctx, sm.esClient)
	if err != nil {
		sm.circuitBreaker.Failed()
		return observability.PrepareError(err, span, "checking index existence successfully")
	}

	if res.StatusCode == http.StatusNotFound {
		if _, err = (esapi.IndicesCreateRequest{Index: strings.ToLower(sm.indexName)}).Do(ctx, sm.esClient); err != nil {
			sm.circuitBreaker.Failed()
			return observability.PrepareError(err, span, "checking index existence")
		}
	}

	sm.circuitBreaker.Succeeded()
	return nil
}
