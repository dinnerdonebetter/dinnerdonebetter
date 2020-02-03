package iterationmedias

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// CreateMiddlewareCtxKey is a string alias we can use for referring to iteration media input data in contexts
	CreateMiddlewareCtxKey models.ContextKey = "iteration_media_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to iteration media update data in contexts
	UpdateMiddlewareCtxKey models.ContextKey = "iteration_media_update_input"

	counterName        metrics.CounterName = "iterationMedias"
	counterDescription                     = "the number of iterationMedias managed by the iterationMedias service"
	topicName          string              = "iteration_medias"
	serviceName        string              = "iteration_medias_service"
)

var (
	_ models.IterationMediaDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list iteration medias
	Service struct {
		logger                  logging.Logger
		iterationMediaCounter   metrics.UnitCounter
		iterationMediaDatabase  models.IterationMediaDataManager
		userIDFetcher           UserIDFetcher
		iterationMediaIDFetcher IterationMediaIDFetcher
		encoderDecoder          encoding.EncoderDecoder
		reporter                newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs
	UserIDFetcher func(*http.Request) uint64

	// IterationMediaIDFetcher is a function that fetches iteration media IDs
	IterationMediaIDFetcher func(*http.Request) uint64
)

// ProvideIterationMediasService builds a new IterationMediasService
func ProvideIterationMediasService(
	ctx context.Context,
	logger logging.Logger,
	db models.IterationMediaDataManager,
	userIDFetcher UserIDFetcher,
	iterationMediaIDFetcher IterationMediaIDFetcher,
	encoder encoding.EncoderDecoder,
	iterationMediaCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	iterationMediaCounter, err := iterationMediaCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                  logger.WithName(serviceName),
		iterationMediaDatabase:  db,
		encoderDecoder:          encoder,
		iterationMediaCounter:   iterationMediaCounter,
		userIDFetcher:           userIDFetcher,
		iterationMediaIDFetcher: iterationMediaIDFetcher,
		reporter:                reporter,
	}

	iterationMediaCount, err := svc.iterationMediaDatabase.GetAllIterationMediasCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("setting current iteration media count: %w", err)
	}
	svc.iterationMediaCounter.IncrementBy(ctx, iterationMediaCount)

	return svc, nil
}
