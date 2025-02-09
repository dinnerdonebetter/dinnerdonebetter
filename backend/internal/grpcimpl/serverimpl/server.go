package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/grpc/service"
	"github.com/dinnerdonebetter/backend/internal/lib/analytics"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/tokens"
	"github.com/dinnerdonebetter/backend/internal/lib/featureflags"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/random"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/eating/indexing"
)

const (
	serviceName = "grpc_service"
)

var _ service.EatingServiceServer = (*Server)(nil)

type Server struct {
	service.UnimplementedEatingServiceServer
	authManager                authentication.Manager
	dataChangesPublisher       messagequeue.Publisher
	analyticsReporter          analytics.EventReporter
	featureFlagManager         featureflags.FeatureFlagManager
	tracer                     tracing.Tracer
	logger                     logging.Logger
	config                     *config.APIServiceConfig
	dataManager                database.DataManager
	tokenIssuer                tokens.Issuer
	validIngredientSearchIndex textsearch.Index[eatingindexing.ValidIngredientSearchSubset]
	authenticator              authentication.Authenticator
	secretGenerator            random.Generator
}

func NewServer(
	ctx context.Context,
	cfg *config.APIServiceConfig,
	authManager authentication.Manager,
	tracerProvider tracing.TracerProvider,
	logger logging.Logger,
	metricsProvider metrics.Provider,
	dataManager database.DataManager,
	publisherProvider messagequeue.PublisherProvider,
	analyticsReporter analytics.EventReporter,
	featureFlagManager featureflags.FeatureFlagManager,
	tokenIssuer tokens.Issuer,
	authenticator authentication.Authenticator,
	secretGenerator random.Generator,
) (*Server, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.Queues.DataChangesTopicName)
	if err != nil {
		return nil, err
	}

	validIngredientSearchManager, err := textsearchcfg.ProvideIndex[eatingindexing.ValidIngredientSearchSubset](ctx, logger, tracerProvider, metricsProvider, &cfg.TextSearch, eatingindexing.IndexTypeValidIngredients)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "initializing valid ingredient index manager")
	}

	s := &Server{
		tracer:                     tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		logger:                     logging.EnsureLogger(logger).WithName(serviceName),
		dataManager:                dataManager,
		dataChangesPublisher:       dataChangesPublisher,
		analyticsReporter:          analyticsReporter,
		featureFlagManager:         featureFlagManager,
		config:                     cfg,
		authManager:                authManager,
		validIngredientSearchIndex: validIngredientSearchManager,
		tokenIssuer:                tokenIssuer,
		authenticator:              authenticator,
		secretGenerator:            secretGenerator,
	}

	return s, nil
}
