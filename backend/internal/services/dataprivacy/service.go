package workers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "data_privacy_service"
)

var _ types.DataPrivacyService = (*service)(nil)

type (
	// service handles data privacy.
	service struct {
		logger                    logging.Logger
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
		dataChangesPublisher      messagequeue.Publisher
		userDataManager           types.UserDataManager
	}
)

// ProvideService builds a new ValidVesselsService.
func ProvideService(
	_ context.Context,
	logger logging.Logger,
	cfg *Config,
	dataManager types.UserDataManager,
	encoder encoding.ServerEncoderDecoder,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.DataPrivacyService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		encoderDecoder:            encoder,
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		userDataManager:           dataManager,
		dataChangesPublisher:      dataChangesPublisher,
	}

	return svc, nil
}
