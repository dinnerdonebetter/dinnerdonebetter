package ingredients

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
	// CreateMiddlewareCtxKey is a string alias we can use for referring to ingredient input data in contexts
	CreateMiddlewareCtxKey models.ContextKey = "ingredient_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to ingredient update data in contexts
	UpdateMiddlewareCtxKey models.ContextKey = "ingredient_update_input"

	counterName        metrics.CounterName = "ingredients"
	counterDescription                     = "the number of ingredients managed by the ingredients service"
	topicName          string              = "ingredients"
	serviceName        string              = "ingredients_service"
)

var (
	_ models.IngredientDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list ingredients
	Service struct {
		logger              logging.Logger
		ingredientCounter   metrics.UnitCounter
		ingredientDatabase  models.IngredientDataManager
		userIDFetcher       UserIDFetcher
		ingredientIDFetcher IngredientIDFetcher
		encoderDecoder      encoding.EncoderDecoder
		reporter            newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs
	UserIDFetcher func(*http.Request) uint64

	// IngredientIDFetcher is a function that fetches ingredient IDs
	IngredientIDFetcher func(*http.Request) uint64
)

// ProvideIngredientsService builds a new IngredientsService
func ProvideIngredientsService(
	ctx context.Context,
	logger logging.Logger,
	db models.IngredientDataManager,
	userIDFetcher UserIDFetcher,
	ingredientIDFetcher IngredientIDFetcher,
	encoder encoding.EncoderDecoder,
	ingredientCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	ingredientCounter, err := ingredientCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:              logger.WithName(serviceName),
		ingredientDatabase:  db,
		encoderDecoder:      encoder,
		ingredientCounter:   ingredientCounter,
		userIDFetcher:       userIDFetcher,
		ingredientIDFetcher: ingredientIDFetcher,
		reporter:            reporter,
	}

	ingredientCount, err := svc.ingredientDatabase.GetAllIngredientsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("setting current ingredient count: %w", err)
	}
	svc.ingredientCounter.IncrementBy(ctx, ingredientCount)

	return svc, nil
}
