package validingredientstateingredients

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "valid_ingredient_preparations_service"
)

var _ types.ValidIngredientStateIngredientDataService = (*service)(nil)

type (
	// service handles valid ingredient preparations.
	service struct {
		logger                                    logging.Logger
		validIngredientStateIngredientDataManager types.ValidIngredientStateIngredientDataManager
		validIngredientStateIngredientIDFetcher   func(*http.Request) string
		validIngredientIDFetcher                  func(*http.Request) string
		validIngredientStateIDFetcher             func(*http.Request) string
		sessionContextDataFetcher                 func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher                      messagequeue.Publisher
		encoderDecoder                            encoding.ServerEncoderDecoder
		tracer                                    tracing.Tracer
	}
)

// ProvideService builds a new ValidIngredientStateIngredientsService.
func ProvideService(
	logger logging.Logger,
	validIngredientStateIngredientDataManager types.ValidIngredientStateIngredientDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	queueConfig *msgconfig.QueuesConfig,
) (types.ValidIngredientStateIngredientDataService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	svc := &service{
		logger:                                    logging.EnsureLogger(logger).WithName(serviceName),
		validIngredientStateIngredientIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientStateIngredientIDURIParamKey),
		validIngredientStateIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientStateIDURIParamKey),
		validIngredientIDFetcher:                  routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientIDURIParamKey),
		sessionContextDataFetcher:                 authentication.FetchContextFromRequest,
		validIngredientStateIngredientDataManager: validIngredientStateIngredientDataManager,
		dataChangesPublisher:                      dataChangesPublisher,
		encoderDecoder:                            encoder,
		tracer:                                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
