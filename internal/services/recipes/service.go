package recipes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/prixfixeco/backend/internal/encoding"
	"github.com/prixfixeco/backend/internal/features/recipeanalysis"
	"github.com/prixfixeco/backend/internal/messagequeue"
	"github.com/prixfixeco/backend/internal/objectstorage"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/routing"
	authservice "github.com/prixfixeco/backend/internal/services/authentication"
	"github.com/prixfixeco/backend/internal/uploads"
	"github.com/prixfixeco/backend/internal/uploads/images"
	"github.com/prixfixeco/backend/pkg/types"
)

const (
	serviceName string = "recipes_service"
)

var _ types.RecipeDataService = (*service)(nil)

type (
	// service handles recipes.
	service struct {
		tracer                    tracing.Tracer
		encoderDecoder            encoding.ServerEncoderDecoder
		recipeDataManager         types.RecipeDataManager
		recipeMediaDataManager    types.RecipeMediaDataManager
		recipeAnalyzer            recipeanalysis.RecipeAnalyzer
		imageUploadProcessor      images.MediaUploadProcessor
		logger                    logging.Logger
		dataChangesPublisher      messagequeue.Publisher
		uploadManager             uploads.UploadManager
		timeFunc                  func() time.Time
		recipeIDFetcher           func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		cfg                       Config
	}
)

func defaultTimeFunc() time.Time {
	return time.Now()
}

var errInvalidConfig = errors.New("config cannot be nil")

// ProvideService builds a new RecipesService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	recipeDataManager types.RecipeDataManager,
	recipeMediaDataManager types.RecipeMediaDataManager,
	recipeGrapher recipeanalysis.RecipeAnalyzer,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	imageUploadProcessor images.MediaUploadProcessor,
	tracerProvider tracing.TracerProvider,
) (types.RecipeDataService, error) {
	if cfg == nil {
		return nil, errInvalidConfig
	}

	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe service data changes publisher: %w", err)
	}

	uploader, err := objectstorage.NewUploadManager(ctx, logger, tracerProvider, &cfg.Uploads.Storage, routeParamManager)
	if err != nil {
		return nil, fmt.Errorf("initializing recipe service upload manager: %w", err)
	}

	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:           routeParamManager.BuildRouteParamStringIDFetcher(RecipeIDURIParamKey),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		recipeDataManager:         recipeDataManager,
		cfg:                       *cfg,
		recipeMediaDataManager:    recipeMediaDataManager,
		dataChangesPublisher:      dataChangesPublisher,
		encoderDecoder:            encoder,
		timeFunc:                  defaultTimeFunc,
		recipeAnalyzer:            recipeGrapher,
		uploadManager:             uploader,
		imageUploadProcessor:      imageUploadProcessor,
		tracer:                    tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
