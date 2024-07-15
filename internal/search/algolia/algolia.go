package algolia

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/search"
	"github.com/dinnerdonebetter/backend/pkg/types"

	algolia "github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

var (
	_ search.Index[types.UserSearchSubset] = (*indexManager[types.UserSearchSubset])(nil)

	ErrNilConfig = errors.New("nil config provided")
)

type (
	indexManager[T search.Searchable] struct {
		logger         logging.Logger
		tracer         tracing.Tracer
		circuitBreaker circuitbreaking.CircuitBreaker
		client         *algolia.Index
		DataType       *T
		timeout        time.Duration
	}
)

func ProvideIndexManager[T search.Searchable](
	_ context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	cfg *Config,
	indexName string,
	circuitBreaker circuitbreaking.CircuitBreaker,
) (search.Index[T], error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	im := &indexManager[T]{
		tracer:         tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(fmt.Sprintf("search_%s", indexName))),
		logger:         logging.EnsureLogger(logger).WithName(indexName),
		client:         algolia.NewClient(cfg.AppID, cfg.APIKey).InitIndex(indexName),
		timeout:        cfg.Timeout,
		circuitBreaker: circuitBreaker,
	}

	return im, nil
}
