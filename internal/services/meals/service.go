package meals

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/routing"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName string = "meals_service"
)

var _ types.MealDataService = (*service)(nil)

type (
	// service handles meals.
	service struct {
		logger                    logging.Logger
		mealDataManager           types.MealDataManager
		mealIDFetcher             func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher      messagequeue.Publisher
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new MealsService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	mealDataManager types.MealDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.MealDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up meals service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		mealIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(MealIDURIParamKey),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		mealDataManager:           mealDataManager,
		dataChangesPublisher:      dataChangesPublisher,
		encoderDecoder:            encoder,
		tracer:                    tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
