package preparations

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
	// CreateMiddlewareCtxKey is a string alias we can use for referring to preparation input data in contexts
	CreateMiddlewareCtxKey models.ContextKey = "preparation_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to preparation update data in contexts
	UpdateMiddlewareCtxKey models.ContextKey = "preparation_update_input"

	counterName        metrics.CounterName = "preparations"
	counterDescription                     = "the number of preparations managed by the preparations service"
	topicName          string              = "preparations"
	serviceName        string              = "preparations_service"
)

var (
	_ models.PreparationDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list preparations
	Service struct {
		logger               logging.Logger
		preparationCounter   metrics.UnitCounter
		preparationDatabase  models.PreparationDataManager
		userIDFetcher        UserIDFetcher
		preparationIDFetcher PreparationIDFetcher
		encoderDecoder       encoding.EncoderDecoder
		reporter             newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs
	UserIDFetcher func(*http.Request) uint64

	// PreparationIDFetcher is a function that fetches preparation IDs
	PreparationIDFetcher func(*http.Request) uint64
)

// ProvidePreparationsService builds a new PreparationsService
func ProvidePreparationsService(
	ctx context.Context,
	logger logging.Logger,
	db models.PreparationDataManager,
	userIDFetcher UserIDFetcher,
	preparationIDFetcher PreparationIDFetcher,
	encoder encoding.EncoderDecoder,
	preparationCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	preparationCounter, err := preparationCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:               logger.WithName(serviceName),
		preparationDatabase:  db,
		encoderDecoder:       encoder,
		preparationCounter:   preparationCounter,
		userIDFetcher:        userIDFetcher,
		preparationIDFetcher: preparationIDFetcher,
		reporter:             reporter,
	}

	preparationCount, err := svc.preparationDatabase.GetAllPreparationsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("setting current preparation count: %w", err)
	}
	svc.preparationCounter.IncrementBy(ctx, preparationCount)

	return svc, nil
}
