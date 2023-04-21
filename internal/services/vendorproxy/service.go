package vendorproxy

import (
	"fmt"
	"net/http"

	"github.com/prixfixeco/backend/internal/analytics"
	"github.com/prixfixeco/backend/internal/encoding"
	"github.com/prixfixeco/backend/internal/featureflags"
	"github.com/prixfixeco/backend/internal/messagequeue"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/routing"
	authservice "github.com/prixfixeco/backend/internal/services/authentication"
	"github.com/prixfixeco/backend/pkg/types"
)

const (
	serviceName string = "vendor_proxy_service"
)

type (
	// Service is the interface for the vendor proxy service.
	Service interface {
		FeatureFlagStatusHandler(res http.ResponseWriter, req *http.Request)
		AnalyticsTrackHandler(res http.ResponseWriter, req *http.Request)
	}

	// service handles vendor proxying.
	service struct {
		logger                    logging.Logger
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher      messagequeue.Publisher
		encoderDecoder            encoding.ServerEncoderDecoder
		featureFlagManager        featureflags.FeatureFlagManager
		featureFlagURLFetcher     func(*http.Request) string
		tracer                    tracing.Tracer
		eventReporter             analytics.EventReporter
	}
)

// ProvideService builds a new ValidInstrumentsService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	featureFlagManager featureflags.FeatureFlagManager,
	eventReporter analytics.EventReporter,
) (Service, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up vendor proxy service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		dataChangesPublisher:      dataChangesPublisher,
		encoderDecoder:            encoder,
		featureFlagURLFetcher:     routeParamManager.BuildRouteParamStringIDFetcher(FeatureFlagURIParamKey),
		tracer:                    tracing.NewTracer(tracerProvider.Tracer(serviceName)),
		featureFlagManager:        featureFlagManager,
		eventReporter:             eventReporter,
	}

	return svc, nil
}
