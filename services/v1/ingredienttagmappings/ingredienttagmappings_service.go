package ingredienttagmappings

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
	// CreateMiddlewareCtxKey is a string alias we can use for referring to ingredient tag mapping input data in contexts.
	CreateMiddlewareCtxKey models.ContextKey = "ingredient_tag_mapping_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to ingredient tag mapping update data in contexts.
	UpdateMiddlewareCtxKey models.ContextKey = "ingredient_tag_mapping_update_input"

	counterName        metrics.CounterName = "ingredientTagMappings"
	counterDescription string              = "the number of ingredientTagMappings managed by the ingredientTagMappings service"
	topicName          string              = "ingredient_tag_mappings"
	serviceName        string              = "ingredient_tag_mappings_service"
)

var (
	_ models.IngredientTagMappingDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list ingredient tag mappings
	Service struct {
		logger                          logging.Logger
		validIngredientDataManager      models.ValidIngredientDataManager
		ingredientTagMappingDataManager models.IngredientTagMappingDataManager
		validIngredientIDFetcher        ValidIngredientIDFetcher
		ingredientTagMappingIDFetcher   IngredientTagMappingIDFetcher
		ingredientTagMappingCounter     metrics.UnitCounter
		encoderDecoder                  encoding.EncoderDecoder
		reporter                        newsman.Reporter
	}

	// ValidIngredientIDFetcher is a function that fetches valid ingredient IDs.
	ValidIngredientIDFetcher func(*http.Request) uint64

	// IngredientTagMappingIDFetcher is a function that fetches ingredient tag mapping IDs.
	IngredientTagMappingIDFetcher func(*http.Request) uint64
)

// ProvideIngredientTagMappingsService builds a new IngredientTagMappingsService.
func ProvideIngredientTagMappingsService(
	logger logging.Logger,
	validIngredientDataManager models.ValidIngredientDataManager,
	ingredientTagMappingDataManager models.IngredientTagMappingDataManager,
	validIngredientIDFetcher ValidIngredientIDFetcher,
	ingredientTagMappingIDFetcher IngredientTagMappingIDFetcher,
	encoder encoding.EncoderDecoder,
	ingredientTagMappingCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	ingredientTagMappingCounter, err := ingredientTagMappingCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                          logger.WithName(serviceName),
		validIngredientIDFetcher:        validIngredientIDFetcher,
		ingredientTagMappingIDFetcher:   ingredientTagMappingIDFetcher,
		validIngredientDataManager:      validIngredientDataManager,
		ingredientTagMappingDataManager: ingredientTagMappingDataManager,
		encoderDecoder:                  encoder,
		ingredientTagMappingCounter:     ingredientTagMappingCounter,
		reporter:                        reporter,
	}

	return svc, nil
}
