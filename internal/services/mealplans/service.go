package mealplans

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	publishers "gitlab.com/prixfixe/prixfixe/internal/messagequeue/publishers"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	routing "gitlab.com/prixfixe/prixfixe/internal/routing"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	serviceName string = "meal_plans_service"
)

var _ types.MealPlanDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles meal plans.
	service struct {
		logger                    logging.Logger
		mealPlanDataManager       types.MealPlanDataManager
		mealPlanIDFetcher         func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		preWritesPublisher        publishers.Publisher
		preUpdatesPublisher       publishers.Publisher
		preArchivesPublisher      publishers.Publisher
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new MealPlansService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	mealPlanDataManager types.MealPlanDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider publishers.PublisherProvider,
) (types.MealPlanDataService, error) {
	preWritesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreWritesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up meal plan queue pre-writes publisher: %w", err)
	}

	preUpdatesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreUpdatesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up meal plan queue pre-updates publisher: %w", err)
	}

	preArchivesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreArchivesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up meal plan queue pre-archives publisher: %w", err)
	}

	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		mealPlanIDFetcher:         routeParamManager.BuildRouteParamStringIDFetcher(MealPlanIDURIParamKey),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		mealPlanDataManager:       mealPlanDataManager,
		preWritesPublisher:        preWritesPublisher,
		preUpdatesPublisher:       preUpdatesPublisher,
		preArchivesPublisher:      preArchivesPublisher,
		encoderDecoder:            encoder,
		tracer:                    tracing.NewTracer(serviceName),
	}

	return svc, nil
}
