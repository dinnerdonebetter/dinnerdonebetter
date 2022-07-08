package validingredients

import (
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
	serviceName string = "valid_ingredients_service"
)

var _ types.ValidIngredientDataService = (*service)(nil)

type (
	// service handles valid ingredients.
	service struct {
		logger                     logging.Logger
		validIngredientDataManager types.ValidIngredientDataManager
		validIngredientIDFetcher   func(*http.Request) string
		sessionContextDataFetcher  func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher       messagequeue.Publisher
		encoderDecoder             encoding.ServerEncoderDecoder
		tracer                     tracing.Tracer
	}
)

// ProvideService builds a new ValidIngredientsService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	validIngredientDataManager types.ValidIngredientDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidIngredientDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredients service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                     logging.EnsureLogger(logger).WithName(serviceName),
		validIngredientIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientIDURIParamKey),
		sessionContextDataFetcher:  authservice.FetchContextFromRequest,
		validIngredientDataManager: validIngredientDataManager,
		dataChangesPublisher:       dataChangesPublisher,
		encoderDecoder:             encoder,
		tracer:                     tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
