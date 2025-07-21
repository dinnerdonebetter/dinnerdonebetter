package grpc

import (
	"context"

	adminsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/admin"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "admin_service"
)

var _ adminsvc.ServiceAdministrationServiceServer = (*ServiceImpl)(nil)

type (
	ServiceImpl struct {
		adminsvc.UnimplementedServiceAdministrationServiceServer
		tracer tracing.Tracer
		logger logging.Logger
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
) adminsvc.ServiceAdministrationServiceServer {
	return &ServiceImpl{
		logger: logging.EnsureLogger(logger).WithName(o11yName),
		tracer: tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
	}
}

func (s *ServiceImpl) AdminUpdateUserStatus(ctx context.Context, request *adminsvc.AdminUpdateUserStatusRequest) (*adminsvc.AdminUpdateUserStatusResponse, error) {
	//TODO implement me
	panic("implement me")
}
