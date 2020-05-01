package validingredienttags

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
	// CreateMiddlewareCtxKey is a string alias we can use for referring to valid ingredient tag input data in contexts.
	CreateMiddlewareCtxKey models.ContextKey = "valid_ingredient_tag_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to valid ingredient tag update data in contexts.
	UpdateMiddlewareCtxKey models.ContextKey = "valid_ingredient_tag_update_input"

	counterName        metrics.CounterName = "validIngredientTags"
	counterDescription string              = "the number of validIngredientTags managed by the validIngredientTags service"
	topicName          string              = "valid_ingredient_tags"
	serviceName        string              = "valid_ingredient_tags_service"
)

var (
	_ models.ValidIngredientTagDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list valid ingredient tags
	Service struct {
		logger                        logging.Logger
		validIngredientTagDataManager models.ValidIngredientTagDataManager
		validIngredientTagIDFetcher   ValidIngredientTagIDFetcher
		validIngredientTagCounter     metrics.UnitCounter
		encoderDecoder                encoding.EncoderDecoder
		reporter                      newsman.Reporter
	}

	// ValidIngredientTagIDFetcher is a function that fetches valid ingredient tag IDs.
	ValidIngredientTagIDFetcher func(*http.Request) uint64
)

// ProvideValidIngredientTagsService builds a new ValidIngredientTagsService.
func ProvideValidIngredientTagsService(
	logger logging.Logger,
	validIngredientTagDataManager models.ValidIngredientTagDataManager,
	validIngredientTagIDFetcher ValidIngredientTagIDFetcher,
	encoder encoding.EncoderDecoder,
	validIngredientTagCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	validIngredientTagCounter, err := validIngredientTagCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                        logger.WithName(serviceName),
		validIngredientTagIDFetcher:   validIngredientTagIDFetcher,
		validIngredientTagDataManager: validIngredientTagDataManager,
		encoderDecoder:                encoder,
		validIngredientTagCounter:     validIngredientTagCounter,
		reporter:                      reporter,
	}

	return svc, nil
}
