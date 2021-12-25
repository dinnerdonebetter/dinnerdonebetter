package mealplanoptions

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
		mealPlanOptionIDFetcher   func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		preWritesPublisher        messagequeue.Publisher
		preUpdatesPublisher       messagequeue.Publisher
		preArchivesPublisher      messagequeue.Publisher
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new MealPlanOptionsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	mealPlanOptionDataManager types.MealPlanOptionDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.MealPlanOptionDataService, error) {
	preWritesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreWritesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up meal plan option queue pre-writes publisher: %w", err)
	}

	preUpdatesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreUpdatesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up meal plan option queue pre-updates publisher: %w", err)
	}

	preArchivesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreArchivesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up meal plan option queue pre-archives publisher: %w", err)
	}

	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		mealPlanIDFetcher:         routeParamManager.BuildRouteParamStringIDFetcher(mealplansservice.MealPlanIDURIParamKey),
		mealPlanOptionIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(MealPlanOptionIDURIParamKey),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		mealPlanOptionDataManager: mealPlanOptionDataManager,
		preWritesPublisher:        preWritesPublisher,
		preUpdatesPublisher:       preUpdatesPublisher,
		preArchivesPublisher:      preArchivesPublisher,
		encoderDecoder:            encoder,
		tracer:                    tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
