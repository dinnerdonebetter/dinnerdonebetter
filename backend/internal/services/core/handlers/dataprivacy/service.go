package dataprivacy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/routing"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads/objectstorage"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "data_privacy_service"
)

var _ types.DataPrivacyService = (*service)(nil)

type (
	// service handles data privacy.
	service struct {
		logger                       logging.Logger
		sessionContextDataFetcher    func(*http.Request) (*sessions.ContextData, error)
		encoderDecoder               encoding.ServerEncoderDecoder
		tracer                       tracing.Tracer
		dataChangesPublisher         messagequeue.Publisher
		reportIDFetcher              func(*http.Request) string
		dataPrivacyDataManager       types.DataPrivacyDataManager
		userDataAggregationPublisher messagequeue.Publisher
		uploader                     uploads.UploadManager
	}
)

// ProvideService builds a new DataPrivacyService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	dataManager types.DataPrivacyDataManager,
	encoder encoding.ServerEncoderDecoder,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	routeParamManager routing.RouteParamManager,
	queueConfig *msgconfig.QueuesConfig,
) (types.DataPrivacyService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	userDataAggregationPublisher, err := publisherProvider.ProvidePublisher(queueConfig.UserDataAggregationTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	uploader, err := objectstorage.NewUploadManager(ctx, logger, tracerProvider, &cfg.Uploads.Storage, routeParamManager)
	if err != nil {
		return nil, fmt.Errorf("initializing %s service upload manager: %w", serviceName, err)
	}

	svc := &service{
		logger:                       logging.EnsureLogger(logger).WithName(serviceName),
		encoderDecoder:               encoder,
		sessionContextDataFetcher:    sessions.FetchContextFromRequest,
		reportIDFetcher:              routeParamManager.BuildRouteParamStringIDFetcher(ReportIDURIParamKey),
		tracer:                       tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		dataPrivacyDataManager:       dataManager,
		dataChangesPublisher:         dataChangesPublisher,
		userDataAggregationPublisher: userDataAggregationPublisher,
		uploader:                     uploader,
	}

	return svc, nil
}
