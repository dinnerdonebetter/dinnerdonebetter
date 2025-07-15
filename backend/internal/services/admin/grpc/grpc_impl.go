package grpc

import (
	"context"

	adminsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/admin"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

var _ adminsvc.ServiceAdministrationServiceServer = (*ServiceImpl)(nil)

type (
	ServiceImpl struct {
		adminsvc.UnimplementedServiceAdministrationServiceServer
		tracer tracing.Tracer
		logger logging.Logger
	}
)

func NewService() adminsvc.ServiceAdministrationServiceServer {
	return &ServiceImpl{}
}

func (s *ServiceImpl) AdminUpdateUserStatus(ctx context.Context, request *adminsvc.AdminUpdateUserStatusRequest) (*adminsvc.AdminUpdateUserStatusResponse, error) {
	//TODO implement me
	panic("implement me")
}
