package advancedprepsteps

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
	serviceName string = "advanced_prep_steps_service"
)

var _ types.AdvancedPrepStepDataService = (*service)(nil)

type (
	// service handles advanced prep steps.
	service struct {
		logger                      logging.Logger
		advancedPrepStepDataManager types.AdvancedPrepStepDataManager
		mealPlanIDFetcher           func(*http.Request) string
		mealPlanEventIDFetcher      func(*http.Request) string
		advancedPrepStepIDFetcher   func(*http.Request) string
		sessionContextDataFetcher   func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher        messagequeue.Publisher
		encoderDecoder              encoding.ServerEncoderDecoder
		tracer                      tracing.Tracer
	}
)

// ProvideService builds a new AdvancedPrepStep.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	advancedPrepStepDataManager types.AdvancedPrepStepDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.AdvancedPrepStepDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up meal plans service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                      logging.EnsureLogger(logger).WithName(serviceName),
		mealPlanIDFetcher:           routeParamManager.BuildRouteParamStringIDFetcher(mealplansservice.MealPlanIDURIParamKey),
		mealPlanEventIDFetcher:      routeParamManager.BuildRouteParamStringIDFetcher(mealplaneventsservice.MealPlanEventIDURIParamKey),
		advancedPrepStepIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(AdvancedPrepStepIDURIParamKey),
		sessionContextDataFetcher:   authservice.FetchContextFromRequest,
		advancedPrepStepDataManager: advancedPrepStepDataManager,
		dataChangesPublisher:        dataChangesPublisher,
		encoderDecoder:              encoder,
		tracer:                      tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
