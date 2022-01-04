package mealplanoptionvotes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/messagequeue"

	"github.com/prixfixeco/api_server/internal/customerdata"
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
	// service handles meal plan option votes.
	service struct {
		logger                        logging.Logger
		mealPlanOptionVoteDataManager types.MealPlanOptionVoteDataManager
		mealPlanIDFetcher             func(*http.Request) string
		mealPlanOptionIDFetcher       func(*http.Request) string
		mealPlanOptionVoteIDFetcher   func(*http.Request) string
		sessionContextDataFetcher     func(*http.Request) (*types.SessionContextData, error)
		preWritesPublisher            messagequeue.Publisher
		preUpdatesPublisher           messagequeue.Publisher
		preArchivesPublisher          messagequeue.Publisher
		dataChangesPublisher          messagequeue.Publisher
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
	publisherProvider messagequeue.PublisherProvider,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
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

	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe step product queue data changes publisher: %w", err)
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
		dataChangesPublisher:          dataChangesPublisher,
		encoderDecoder:                encoder,
		tracer:                        tracing.NewTracer(tracerProvider.Tracer(serviceName)),
		customerDataCollector:         customerDataCollector,
	}

	return svc, nil
}
