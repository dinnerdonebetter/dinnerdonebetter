package mealplangrocerylistitems

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
	serviceName string = "meal_plan_grocery_list_items_service"
)

var _ types.MealPlanGroceryListItemDataService = (*service)(nil)

type (
	// service handles meal plan grocery list items.
	service struct {
		logger                             logging.Logger
		mealPlanGroceryListItemDataManager types.MealPlanGroceryListItemDataManager
		mealPlanIDFetcher                  func(*http.Request) string
		mealPlanEventIDFetcher             func(*http.Request) string
		mealPlanGroceryListItemIDFetcher   func(*http.Request) string
		sessionContextDataFetcher          func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher               messagequeue.Publisher
		encoderDecoder                     encoding.ServerEncoderDecoder
		tracer                             tracing.Tracer
	}
)

// ProvideService builds a new MealPlanGroceryListItem.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	mealPlanGroceryListItemDataManager types.MealPlanGroceryListItemDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.MealPlanGroceryListItemDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up meal plans service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                             logging.EnsureLogger(logger).WithName(serviceName),
		mealPlanIDFetcher:                  routeParamManager.BuildRouteParamStringIDFetcher(mealplansservice.MealPlanIDURIParamKey),
		mealPlanEventIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(mealplaneventsservice.MealPlanEventIDURIParamKey),
		mealPlanGroceryListItemIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(MealPlanGroceryListItemIDURIParamKey),
		sessionContextDataFetcher:          authservice.FetchContextFromRequest,
		mealPlanGroceryListItemDataManager: mealPlanGroceryListItemDataManager,
		dataChangesPublisher:               dataChangesPublisher,
		encoderDecoder:                     encoder,
		tracer:                             tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
