package grpc

import (
	"context"

	identityconverters "github.com/dinnerdonebetter/backend/internal/domain/identity/converters"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	errorsgrpc "github.com/dinnerdonebetter/backend/internal/platform/errors/grpc"

	"google.golang.org/grpc/codes"
)

func (s *serviceImpl) AdminUpdateUserStatus(ctx context.Context, request *identitysvc.AdminUpdateUserStatusRequest) (*identitysvc.AdminUpdateUserStatusResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, errorsgrpc.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}

	if !sessionContextData.Requester.ServicePermissions.CanUpdateUserAccountStatuses() {
		return nil, errorsgrpc.PrepareAndLogGRPCStatus(nil, logger, span, codes.Unauthenticated, "user account status update requester does not have permission")
	}

	if err = s.identityDataManager.AdminUpdateUserStatus(ctx, identityconverters.ConvertGRPCAdminUpdateUserStatusRequestToUserAccountStatusUpdateInput(request)); err != nil {
		return nil, errorsgrpc.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating user account status")
	}

	x := &identitysvc.AdminUpdateUserStatusResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	return x, nil
}
