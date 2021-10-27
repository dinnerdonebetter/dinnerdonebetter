package validingredientpreparations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/encoding"
	publishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	routing "github.com/prixfixeco/api_server/internal/routing"
	"github.com/prixfixeco/api_server/internal/search"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName string = "valid_ingredient_preparations_service"
)

var _ types.ValidIngredientPreparationDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles valid ingredient preparations.
	service struct {
		logger                                logging.Logger
		validIngredientPreparationDataManager types.ValidIngredientPreparationDataManager
		validIngredientPreparationIDFetcher   func(*http.Request) string
		sessionContextDataFetcher             func(*http.Request) (*types.SessionContextData, error)
		preWritesPublisher                    publishers.Publisher
		preUpdatesPublisher                   publishers.Publisher
		preArchivesPublisher                  publishers.Publisher
		encoderDecoder                        encoding.ServerEncoderDecoder
		tracer                                tracing.Tracer
	}
)

// ProvideService builds a new ValidIngredientPreparationsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	validIngredientPreparationDataManager types.ValidIngredientPreparationDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider publishers.PublisherProvider,
) (types.ValidIngredientPreparationDataService, error) {
	preWritesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreWritesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredient preparation queue pre-writes publisher: %w", err)
	}

	preUpdatesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreUpdatesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredient preparation queue pre-updates publisher: %w", err)
	}

	preArchivesPublisher, err := publisherProvider.ProviderPublisher(cfg.PreArchivesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredient preparation queue pre-archives publisher: %w", err)
	}

	svc := &service{
		logger:                                logging.EnsureLogger(logger).WithName(serviceName),
		validIngredientPreparationIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientPreparationIDURIParamKey),
		sessionContextDataFetcher:             authservice.FetchContextFromRequest,
		validIngredientPreparationDataManager: validIngredientPreparationDataManager,
		preWritesPublisher:                    preWritesPublisher,
		preUpdatesPublisher:                   preUpdatesPublisher,
		preArchivesPublisher:                  preArchivesPublisher,
		encoderDecoder:                        encoder,
		tracer:                                tracing.NewTracer(serviceName),
	}

	return svc, nil
}
