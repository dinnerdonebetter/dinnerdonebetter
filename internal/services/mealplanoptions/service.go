package mealplanoptions

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/routing"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	mealplaneventsservice "github.com/prixfixeco/api_server/internal/services/mealplanevents"
	mealplansservice "github.com/prixfixeco/api_server/internal/services/mealplans"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName string = "meal_plan_options_service"
)

var _ types.MealPlanOptionDataService = (*service)(nil)

type (
	// service handles meal plan options.
	service struct {
		logger                    logging.Logger
		mealPlanOptionDataManager types.MealPlanOptionDataManager
		mealPlanIDFetcher         func(*http.Request) string
		mealPlanEventIDFetcher    func(*http.Request) string
		mealPlanOptionIDFetcher   func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher      messagequeue.Publisher
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new MealPlanOptionsService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	mealPlanOptionDataManager types.MealPlanOptionDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.MealPlanOptionDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up meal plan options service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		mealPlanIDFetcher:         routeParamManager.BuildRouteParamStringIDFetcher(mealplansservice.MealPlanIDURIParamKey),
		mealPlanEventIDFetcher:    routeParamManager.BuildRouteParamStringIDFetcher(mealplaneventsservice.MealPlanEventIDURIParamKey),
		mealPlanOptionIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(MealPlanOptionIDURIParamKey),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		mealPlanOptionDataManager: mealPlanOptionDataManager,
		dataChangesPublisher:      dataChangesPublisher,
		encoderDecoder:            encoder,
		tracer:                    tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
