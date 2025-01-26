package admin

import (
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessioncontext"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName = "admin_service"
)

type (
	// service handles miscellaneous admin tasks.
	service struct {
		logger                    logging.Logger
		userDB                    types.AdminUserDataManager
		encoderDecoder            encoding.ServerEncoderDecoder
		publisherProvider         messagequeue.PublisherProvider
		tracer                    tracing.Tracer
		sessionContextDataFetcher func(*http.Request) (*sessioncontext.SessionContextData, error)
		queuesConfig              msgconfig.QueuesConfig
	}
)

// ProvideService builds a new AdminDataService.
func ProvideService(
	logger logging.Logger,
	userDataManager types.AdminUserDataManager,
	encoder encoding.ServerEncoderDecoder,
	tracerProvider tracing.TracerProvider,
	queuesConfig *msgconfig.QueuesConfig,
	publisherProvider messagequeue.PublisherProvider,
) (types.AdminDataService, error) {
	if queuesConfig == nil {
		return nil, errors.New("nil queues config provided")
	}

	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		encoderDecoder:            encoder,
		userDB:                    userDataManager,
		queuesConfig:              *queuesConfig,
		publisherProvider:         publisherProvider,
		sessionContextDataFetcher: sessioncontext.FetchContextFromRequest,
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
