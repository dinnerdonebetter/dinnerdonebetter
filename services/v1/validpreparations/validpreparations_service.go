package validpreparations

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
	// CreateMiddlewareCtxKey is a string alias we can use for referring to valid preparation input data in contexts.
	CreateMiddlewareCtxKey models.ContextKey = "valid_preparation_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to valid preparation update data in contexts.
	UpdateMiddlewareCtxKey models.ContextKey = "valid_preparation_update_input"

	counterName        metrics.CounterName = "validPreparations"
	counterDescription string              = "the number of validPreparations managed by the validPreparations service"
	topicName          string              = "valid_preparations"
	serviceName        string              = "valid_preparations_service"
)

var (
	_ models.ValidPreparationDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list valid preparations
	Service struct {
		logger                      logging.Logger
		validPreparationDataManager models.ValidPreparationDataManager
		validPreparationIDFetcher   ValidPreparationIDFetcher
		validPreparationCounter     metrics.UnitCounter
		encoderDecoder              encoding.EncoderDecoder
		reporter                    newsman.Reporter
	}

	// ValidPreparationIDFetcher is a function that fetches valid preparation IDs.
	ValidPreparationIDFetcher func(*http.Request) uint64
)

// ProvideValidPreparationsService builds a new ValidPreparationsService.
func ProvideValidPreparationsService(
	logger logging.Logger,
	validPreparationDataManager models.ValidPreparationDataManager,
	validPreparationIDFetcher ValidPreparationIDFetcher,
	encoder encoding.EncoderDecoder,
	validPreparationCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	validPreparationCounter, err := validPreparationCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                      logger.WithName(serviceName),
		validPreparationIDFetcher:   validPreparationIDFetcher,
		validPreparationDataManager: validPreparationDataManager,
		encoderDecoder:              encoder,
		validPreparationCounter:     validPreparationCounter,
		reporter:                    reporter,
	}

	return svc, nil
}
