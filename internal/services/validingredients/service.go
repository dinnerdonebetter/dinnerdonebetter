package validingredients

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
	serviceName string = "valid_ingredients_service"
)

var _ types.ValidIngredientDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles valid ingredients.
	service struct {
		logger                     logging.Logger
		validIngredientDataManager types.ValidIngredientDataManager
		validIngredientIDFetcher   func(*http.Request) string
		sessionContextDataFetcher  func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher       publishers.Publisher
		encoderDecoder             encoding.ServerEncoderDecoder
		tracer                     tracing.Tracer
		search                     SearchIndex
	}
)

// ProvideService builds a new ValidIngredientsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	validIngredientDataManager types.ValidIngredientDataManager,
	encoder encoding.ServerEncoderDecoder,
	searchIndexProvider search.IndexManagerProvider,
	routeParamManager routing.RouteParamManager,
	publisherProvider publishers.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidIngredientDataService, error) {
	searchIndexManager, err := searchIndexProvider.ProvideIndexManager(ctx, logger, "valid_ingredients", "name", "variant", "description", "warning")
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredient search index: %w", err)
	}

	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredient queue pre-archives publisher: %w", err)
	}

	svc := &service{
		logger:                     logging.EnsureLogger(logger).WithName(serviceName),
		validIngredientIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientIDURIParamKey),
		sessionContextDataFetcher:  authservice.FetchContextFromRequest,
		validIngredientDataManager: validIngredientDataManager,
		dataChangesPublisher:       dataChangesPublisher,
		encoderDecoder:             encoder,
		search:                     searchIndexManager,
		tracer:                     tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
