package serverimpl

import (
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/grpc/service"
	"github.com/dinnerdonebetter/backend/internal/lib/analytics"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/tokens"
	"github.com/dinnerdonebetter/backend/internal/lib/featureflags"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/random"
)

const (
	serviceName = "grpc_service"
)

var _ service.EatingServiceServer = (*Server)(nil)

type Server struct {
	service.UnimplementedEatingServiceServer
	authManager          authentication.Manager
	dataChangesPublisher messagequeue.Publisher
	analyticsReporter    analytics.EventReporter
	featureFlagManager   featureflags.FeatureFlagManager
	tracer               tracing.Tracer
	logger               logging.Logger
	config               *config.APIServiceConfig
	dataManager          database.DataManager
	tokenIssuer          tokens.Issuer
	authenticator        authentication.Authenticator
	secretGenerator      random.Generator
}

func NewServer(
	cfg *config.APIServiceConfig,
	authManager authentication.Manager,
	tracerProvider tracing.TracerProvider,
	logger logging.Logger,
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

	s := &Server{
		tracer:               tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		logger:               logging.EnsureLogger(logger).WithName(serviceName),
		dataManager:          dataManager,
		dataChangesPublisher: dataChangesPublisher,
		analyticsReporter:    analyticsReporter,
		featureFlagManager:   featureFlagManager,
		config:               cfg,
		authManager:          authManager,
		tokenIssuer:          tokenIssuer,
		authenticator:        authenticator,
		secretGenerator:      secretGenerator,
	}

	return s, nil
}
