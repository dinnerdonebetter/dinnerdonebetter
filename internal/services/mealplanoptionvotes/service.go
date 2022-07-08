package mealplanoptionvotes

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/database"

	"github.com/prixfixeco/api_server/internal/messagequeue"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/routing"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	mealplanoptionsservice "github.com/prixfixeco/api_server/internal/services/mealplanoptions"
	mealplansservice "github.com/prixfixeco/api_server/internal/services/mealplans"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName string = "meal_plan_option_votes_service"
)

var _ types.MealPlanOptionVoteDataService = (*service)(nil)

type (
	dataManager interface {
		types.MealPlanDataManager
		types.MealPlanOptionDataManager
		types.MealPlanOptionVoteDataManager
	}

	// service handles meal plan option votes.
	service struct {
		logger                      logging.Logger
		dataManager                 dataManager
		mealPlanIDFetcher           func(*http.Request) string
		mealPlanOptionIDFetcher     func(*http.Request) string
		mealPlanOptionVoteIDFetcher func(*http.Request) string
		sessionContextDataFetcher   func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher        messagequeue.Publisher
		encoderDecoder              encoding.ServerEncoderDecoder
		tracer                      tracing.Tracer
	}
)

// ProvideService builds a new MealPlanOptionVotesService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	dataManager database.DataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.MealPlanOptionVoteDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up meal plan option votes service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                      logging.EnsureLogger(logger).WithName(serviceName),
		mealPlanIDFetcher:           routeParamManager.BuildRouteParamStringIDFetcher(mealplansservice.MealPlanIDURIParamKey),
		mealPlanOptionIDFetcher:     routeParamManager.BuildRouteParamStringIDFetcher(mealplanoptionsservice.MealPlanOptionIDURIParamKey),
		mealPlanOptionVoteIDFetcher: routeParamManager.BuildRouteParamStringIDFetcher(MealPlanOptionVoteIDURIParamKey),
		sessionContextDataFetcher:   authservice.FetchContextFromRequest,
		dataManager:                 dataManager,
		dataChangesPublisher:        dataChangesPublisher,
		encoderDecoder:              encoder,
		tracer:                      tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
