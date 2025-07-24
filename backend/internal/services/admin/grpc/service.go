package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	adminsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/admin"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/admin/grpc/converters"

	"google.golang.org/grpc/codes"
)

const (
	o11yName = "admin_service"
)

var _ adminsvc.ServiceAdministrationServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		adminsvc.UnimplementedServiceAdministrationServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		sessionContextDataFetcher func(context.Context) (sessions.ContextData, error)
		identityRepository        identity.AdminUserDataManager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	identityRepository identity.AdminUserDataManager,
) adminsvc.ServiceAdministrationServiceServer {
	return &serviceImpl{
		logger:             logging.EnsureLogger(logger).WithName(o11yName),
		tracer:             tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		identityRepository: identityRepository,
	}
}

func (s *serviceImpl) AdminUpdateUserStatus(ctx context.Context, request *adminsvc.AdminUpdateUserStatusRequest) (*adminsvc.AdminUpdateUserStatusResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}

	if !sessionContextData.Requester.ServicePermissions.CanUpdateUserAccountStatuses() {
		return nil, observability.PrepareAndLogGRPCStatus(nil, logger, span, codes.Unauthenticated, "user account status update requester does not have permission")
	}

	input := converters.ConvertGRPCAdminUpdateUserStatusRequestToUserAccountStatusUpdateInput(request)
	if err = s.identityRepository.UpdateUserAccountStatus(ctx, request.TargetUserID, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating user account status")
	}

	x := &adminsvc.AdminUpdateUserStatusResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}
