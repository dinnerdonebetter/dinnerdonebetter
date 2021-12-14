package mealplanoptionvotes

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue/publishers"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/routing"
	"github.com/prixfixeco/api_server/internal/search"
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
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles meal plan option votes.
	service struct {
		logger                        logging.Logger
		mealPlanOptionVoteDataManager types.MealPlanOptionVoteDataManager
		mealPlanIDFetcher             func(*http.Request) string
		mealPlanOptionIDFetcher       func(*http.Request) string
		mealPlanOptionVoteIDFetcher   func(*http.Request) string
		sessionContextDataFetcher     func(*http.Request) (*types.SessionContextData, error)
		preWritesPublisher            publishers.Publisher
		preUpdatesPublisher           publishers.Publisher
		preArchivesPublisher          publishers.Publisher
		encoderDecoder                encoding.ServerEncoderDecoder
		tracer                        tracing.Tracer
		customerDataCollector         customerdata.Collector
	}
)

// ProvideService builds a new MealPlanOptionVotesService.
func ProvideService(
	_ context.Context,
	logger logging.Logger,
	cfg *Config,
	mealPlanOptionVoteDataManager types.MealPlanOptionVoteDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider publishers.PublisherProvider,
	customerDataCollector customerdata.Collector,
	tracerProvider trace.TracerProvider,
) (types.MealPlanOptionVoteDataService, error) {
	preWritesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreWritesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up meal plan option vote queue pre-writes publisher: %w", err)
	}

	preUpdatesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreUpdatesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up meal plan option vote queue pre-updates publisher: %w", err)
	}

	preArchivesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreArchivesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up meal plan option vote queue pre-archives publisher: %w", err)
	}

	svc := &service{
		logger:                        logging.EnsureLogger(logger).WithName(serviceName),
		mealPlanIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(mealplansservice.MealPlanIDURIParamKey),
		mealPlanOptionIDFetcher:       routeParamManager.BuildRouteParamStringIDFetcher(mealplanoptionsservice.MealPlanOptionIDURIParamKey),
		mealPlanOptionVoteIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(MealPlanOptionVoteIDURIParamKey),
		sessionContextDataFetcher:     authservice.FetchContextFromRequest,
		mealPlanOptionVoteDataManager: mealPlanOptionVoteDataManager,
		preWritesPublisher:            preWritesPublisher,
		preUpdatesPublisher:           preUpdatesPublisher,
		preArchivesPublisher:          preArchivesPublisher,
		encoderDecoder:                encoder,
		tracer:                        tracing.NewTracer(tracerProvider.Tracer(serviceName)),
		customerDataCollector:         customerDataCollector,
	}

	return svc, nil
}
