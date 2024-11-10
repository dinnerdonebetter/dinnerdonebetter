package admin

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName = "auth_service"
)

type (
	// service handles miscellaneous admin tasks.
	service struct {
		logger                    logging.Logger
		userDB                    types.AdminUserDataManager
		encoderDecoder            encoding.ServerEncoderDecoder
		publisherProvider         messagequeue.PublisherProvider
		tracer                    tracing.Tracer
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		queuesConfig              config.QueueSettings
	}
)

// ProvideService builds a new AdminDataService.
func ProvideService(
	logger logging.Logger,
	userDataManager types.AdminUserDataManager,
	encoder encoding.ServerEncoderDecoder,
	tracerProvider tracing.TracerProvider,
	queuesConfig config.QueueSettings,
	publisherProvider messagequeue.PublisherProvider,
) types.AdminDataService {
	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		encoderDecoder:            encoder,
		userDB:                    userDataManager,
		queuesConfig:              queuesConfig,
		publisherProvider:         publisherProvider,
		sessionContextDataFetcher: authentication.FetchContextFromRequest,
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc
}
