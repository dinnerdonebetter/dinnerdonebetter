package validingredientpreparations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/messagequeue"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/routing"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName string = "valid_ingredient_preparations_service"
)

var _ types.ValidIngredientPreparationDataService = (*service)(nil)

type (
	// service handles valid ingredient preparations.
	service struct {
		logger                                logging.Logger
		validIngredientPreparationDataManager types.ValidIngredientPreparationDataManager
		validIngredientPreparationIDFetcher   func(*http.Request) string
		sessionContextDataFetcher             func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher                  messagequeue.Publisher
		encoderDecoder                        encoding.ServerEncoderDecoder
		tracer                                tracing.Tracer
	}
)

// ProvideService builds a new ValidIngredientPreparationsService.
func ProvideService(
	_ context.Context,
	logger logging.Logger,
	cfg *Config,
	validIngredientPreparationDataManager types.ValidIngredientPreparationDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidIngredientPreparationDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredient preprarations service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                                logging.EnsureLogger(logger).WithName(serviceName),
		validIngredientPreparationIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientPreparationIDURIParamKey),
		sessionContextDataFetcher:             authservice.FetchContextFromRequest,
		validIngredientPreparationDataManager: validIngredientPreparationDataManager,
		dataChangesPublisher:                  dataChangesPublisher,
		encoderDecoder:                        encoder,
		tracer:                                tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
