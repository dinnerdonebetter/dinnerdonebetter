package workers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	"github.com/dinnerdonebetter/backend/internal/uploads"
	"github.com/dinnerdonebetter/backend/internal/uploads/objectstorage"
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
		sessionContextDataFetcher    func(*http.Request) (*types.SessionContextData, error)
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
) (types.DataPrivacyService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	userDataAggregationPublisher, err := publisherProvider.ProvidePublisher(cfg.UserDataAggregationTopicName)
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
		sessionContextDataFetcher:    authservice.FetchContextFromRequest,
		reportIDFetcher:              routeParamManager.BuildRouteParamStringIDFetcher(ReportIDURIParamKey),
		tracer:                       tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		dataPrivacyDataManager:       dataManager,
		dataChangesPublisher:         dataChangesPublisher,
		userDataAggregationPublisher: userDataAggregationPublisher,
		uploader:                     uploader,
	}

	return svc, nil
}
