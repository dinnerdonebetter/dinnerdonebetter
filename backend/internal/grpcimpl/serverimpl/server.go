package serverimpl

import (
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/grpc/service"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/tokens"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
)

const (
	serviceName = "grpc_service"
)

var _ service.EatingServiceServer = (*Server)(nil)

type Server struct {
	service.UnimplementedEatingServiceServer
	tracer        tracing.Tracer
	logger        logging.Logger
	config        *config.APIServiceConfig
	dataManager   database.DataManager
	tokenIssuer   tokens.Issuer
	authenticator authentication.Authenticator
}

func NewServer(
	tracerProvider tracing.TracerProvider,
	logger logging.Logger,
	dataManager database.DataManager,
	// ) (service.EatingServiceServer, error) {
) (*Server, error) {
	s := &Server{
		dataManager: dataManager,
		tracer:      tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		logger:      logging.EnsureLogger(logger).WithName(serviceName),
	}

	return s, nil
}
