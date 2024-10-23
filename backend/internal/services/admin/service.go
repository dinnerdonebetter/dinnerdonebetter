package admin

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName = "auth_service"
)

type (
	// service handles passwords service-wide.
	service struct {
		logger                    logging.Logger
		userDB                    types.AdminUserDataManager
		encoderDecoder            encoding.ServerEncoderDecoder
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new AuthDataService.
func ProvideService(
	logger logging.Logger,
	userDataManager types.AdminUserDataManager,
	encoder encoding.ServerEncoderDecoder,
	tracerProvider tracing.TracerProvider,
) types.AdminDataService {
	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		encoderDecoder:            encoder,
		userDB:                    userDataManager,
		sessionContextDataFetcher: authentication.FetchContextFromRequest,
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc
}
