package validingredients

import (
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// CreateMiddlewareCtxKey is a string alias we can use for referring to valid ingredient input data in contexts.
	CreateMiddlewareCtxKey models.ContextKey = "valid_ingredient_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to valid ingredient update data in contexts.
	UpdateMiddlewareCtxKey models.ContextKey = "valid_ingredient_update_input"

	counterName        metrics.CounterName = "validIngredients"
	counterDescription string              = "the number of validIngredients managed by the validIngredients service"
	topicName          string              = "valid_ingredients"
	serviceName        string              = "valid_ingredients_service"
)

var (
	_ models.ValidIngredientDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list valid ingredients
	Service struct {
		logger                     logging.Logger
		validIngredientDataManager models.ValidIngredientDataManager
		validIngredientIDFetcher   ValidIngredientIDFetcher
		validIngredientCounter     metrics.UnitCounter
		encoderDecoder             encoding.EncoderDecoder
		reporter                   newsman.Reporter
	}

	// ValidIngredientIDFetcher is a function that fetches valid ingredient IDs.
	ValidIngredientIDFetcher func(*http.Request) uint64
)

// ProvideValidIngredientsService builds a new ValidIngredientsService.
func ProvideValidIngredientsService(
	logger logging.Logger,
	validIngredientDataManager models.ValidIngredientDataManager,
	validIngredientIDFetcher ValidIngredientIDFetcher,
	encoder encoding.EncoderDecoder,
	validIngredientCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	validIngredientCounter, err := validIngredientCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                     logger.WithName(serviceName),
		validIngredientIDFetcher:   validIngredientIDFetcher,
		validIngredientDataManager: validIngredientDataManager,
		encoderDecoder:             encoder,
		validIngredientCounter:     validIngredientCounter,
		reporter:                   reporter,
	}

	return svc, nil
}
