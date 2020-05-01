package validingredientpreparations

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
	// CreateMiddlewareCtxKey is a string alias we can use for referring to valid ingredient preparation input data in contexts.
	CreateMiddlewareCtxKey models.ContextKey = "valid_ingredient_preparation_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to valid ingredient preparation update data in contexts.
	UpdateMiddlewareCtxKey models.ContextKey = "valid_ingredient_preparation_update_input"

	counterName        metrics.CounterName = "validIngredientPreparations"
	counterDescription string              = "the number of validIngredientPreparations managed by the validIngredientPreparations service"
	topicName          string              = "valid_ingredient_preparations"
	serviceName        string              = "valid_ingredient_preparations_service"
)

var (
	_ models.ValidIngredientPreparationDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list valid ingredient preparations
	Service struct {
		logger                                logging.Logger
		validIngredientDataManager            models.ValidIngredientDataManager
		validIngredientPreparationDataManager models.ValidIngredientPreparationDataManager
		validIngredientIDFetcher              ValidIngredientIDFetcher
		validIngredientPreparationIDFetcher   ValidIngredientPreparationIDFetcher
		validIngredientPreparationCounter     metrics.UnitCounter
		encoderDecoder                        encoding.EncoderDecoder
		reporter                              newsman.Reporter
	}

	// ValidIngredientIDFetcher is a function that fetches valid ingredient IDs.
	ValidIngredientIDFetcher func(*http.Request) uint64

	// ValidIngredientPreparationIDFetcher is a function that fetches valid ingredient preparation IDs.
	ValidIngredientPreparationIDFetcher func(*http.Request) uint64
)

// ProvideValidIngredientPreparationsService builds a new ValidIngredientPreparationsService.
func ProvideValidIngredientPreparationsService(
	logger logging.Logger,
	validIngredientDataManager models.ValidIngredientDataManager,
	validIngredientPreparationDataManager models.ValidIngredientPreparationDataManager,
	validIngredientIDFetcher ValidIngredientIDFetcher,
	validIngredientPreparationIDFetcher ValidIngredientPreparationIDFetcher,
	encoder encoding.EncoderDecoder,
	validIngredientPreparationCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	validIngredientPreparationCounter, err := validIngredientPreparationCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                                logger.WithName(serviceName),
		validIngredientIDFetcher:              validIngredientIDFetcher,
		validIngredientPreparationIDFetcher:   validIngredientPreparationIDFetcher,
		validIngredientDataManager:            validIngredientDataManager,
		validIngredientPreparationDataManager: validIngredientPreparationDataManager,
		encoderDecoder:                        encoder,
		validIngredientPreparationCounter:     validIngredientPreparationCounter,
		reporter:                              reporter,
	}

	return svc, nil
}
