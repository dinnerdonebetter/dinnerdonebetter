package validinstruments

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
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName string = "valid_instruments_service"
)

var _ types.ValidInstrumentDataService = (*service)(nil)

type (
	// service handles valid instruments.
	service struct {
		logger                     logging.Logger
		validInstrumentDataManager types.ValidInstrumentDataManager
		validInstrumentIDFetcher   func(*http.Request) string
		sessionContextDataFetcher  func(*http.Request) (*types.SessionContextData, error)
		preWritesPublisher         messagequeue.Publisher
		preUpdatesPublisher        messagequeue.Publisher
		preArchivesPublisher       messagequeue.Publisher
		encoderDecoder             encoding.ServerEncoderDecoder
		tracer                     tracing.Tracer
	}
)

// ProvideService builds a new ValidInstrumentsService.
func ProvideService(
	_ context.Context,
	logger logging.Logger,
	cfg *Config,
	validInstrumentDataManager types.ValidInstrumentDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidInstrumentDataService, error) {
	preWritesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreWritesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid instrument queue pre-writes publisher: %w", err)
	}

	preUpdatesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreUpdatesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid instrument queue pre-updates publisher: %w", err)
	}

	preArchivesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreArchivesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid instrument queue pre-archives publisher: %w", err)
	}

	svc := &service{
		logger:                     logging.EnsureLogger(logger).WithName(serviceName),
		validInstrumentIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidInstrumentIDURIParamKey),
		sessionContextDataFetcher:  authservice.FetchContextFromRequest,
		validInstrumentDataManager: validInstrumentDataManager,
		preWritesPublisher:         preWritesPublisher,
		preUpdatesPublisher:        preUpdatesPublisher,
		preArchivesPublisher:       preArchivesPublisher,
		encoderDecoder:             encoder,
		tracer:                     tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
