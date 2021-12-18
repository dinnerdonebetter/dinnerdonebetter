package validpreparations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue/publishers"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/routing"
	"github.com/prixfixeco/api_server/internal/search"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName string = "valid_preparations_service"
)

var _ types.ValidPreparationDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles valid preparations.
	service struct {
		logger                      logging.Logger
		validPreparationDataManager types.ValidPreparationDataManager
		validPreparationIDFetcher   func(*http.Request) string
		sessionContextDataFetcher   func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher        publishers.Publisher
		encoderDecoder              encoding.ServerEncoderDecoder
		tracer                      tracing.Tracer
		search                      SearchIndex
	}
)

// ProvideService builds a new ValidPreparationsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	validPreparationDataManager types.ValidPreparationDataManager,
	encoder encoding.ServerEncoderDecoder,
	searchIndexProvider search.IndexManagerProvider,
	routeParamManager routing.RouteParamManager,
	publisherProvider publishers.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidPreparationDataService, error) {
	searchIndexManager, err := searchIndexProvider.ProvideIndexManager(ctx, logger, "valid_preparations", "name", "description")
	if err != nil {
		return nil, fmt.Errorf("setting up valid preparation search index: %w", err)
	}

	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid preparation queue pre-archives publisher: %w", err)
	}

	svc := &service{
		logger:                      logging.EnsureLogger(logger).WithName(serviceName),
		validPreparationIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidPreparationIDURIParamKey),
		sessionContextDataFetcher:   authservice.FetchContextFromRequest,
		validPreparationDataManager: validPreparationDataManager,
		dataChangesPublisher:        dataChangesPublisher,
		encoderDecoder:              encoder,
		search:                      searchIndexManager,
		tracer:                      tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
