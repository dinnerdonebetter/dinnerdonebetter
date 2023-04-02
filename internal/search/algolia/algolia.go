package algolia

import (
	"context"
	"fmt"
	"time"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/search"
	"github.com/prixfixeco/backend/pkg/types"

	algolia "github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

var _ search.IndexManager[types.ValidIngredient] = (*IndexManager[types.ValidIngredient])(nil)

type (
	IndexManager[T search.Searchable] struct {
		logger   logging.Logger
		tracer   tracing.Tracer
		client   *algolia.Index
		DataType *T
		timeout  time.Duration
	}
)

func ProvideIndexManager[T search.Searchable](
	_ context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	cfg *Config,
	indexName string,
) (search.IndexManager[T], error) {
	im := &IndexManager[T]{
		tracer:  tracing.NewTracer(tracerProvider.Tracer(fmt.Sprintf("search_%s", indexName))),
		logger:  logging.EnsureLogger(logger).WithName(indexName),
		client:  algolia.NewClient(cfg.AppID, cfg.APIKey).InitIndex(indexName),
		timeout: cfg.Timeout,
	}

	return im, nil
}
