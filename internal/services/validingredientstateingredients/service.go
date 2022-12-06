package validingredientstateingredients

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/backend/internal/encoding"
	"github.com/prixfixeco/backend/internal/messagequeue"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/routing"
	authservice "github.com/prixfixeco/backend/internal/services/authentication"
	"github.com/prixfixeco/backend/pkg/types"
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
		validPreparationIDFetcher                 func(*http.Request) string
		sessionContextDataFetcher                 func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher                      messagequeue.Publisher
		encoderDecoder                            encoding.ServerEncoderDecoder
		tracer                                    tracing.Tracer
	}
)

// ProvideService builds a new ValidIngredientStateIngredientsService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	validIngredientStateIngredientDataManager types.ValidIngredientStateIngredientDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidIngredientStateIngredientDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredient preprarations service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                                    logging.EnsureLogger(logger).WithName(serviceName),
		validIngredientStateIngredientIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientStateIngredientIDURIParamKey),
		validPreparationIDFetcher:                 routeParamManager.BuildRouteParamStringIDFetcher(ValidPreparationIDURIParamKey),
		validIngredientIDFetcher:                  routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientIDURIParamKey),
		sessionContextDataFetcher:                 authservice.FetchContextFromRequest,
		validIngredientStateIngredientDataManager: validIngredientStateIngredientDataManager,
		dataChangesPublisher:                      dataChangesPublisher,
		encoderDecoder:                            encoder,
		tracer:                                    tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}