package grpc

import (
	"context"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"

	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"google.golang.org/grpc/codes"
)

const (
	o11yName = "identity_service"
)

var _ identitysvc.IdentityServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		identitysvc.UnimplementedIdentityServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		sessionContextDataFetcher func(ctx context.Context) (sessions.ContextData, error)
		identityRepository        identity.Repository
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	identityRepository identity.Repository,
) identitysvc.IdentityServiceServer {
	return &serviceImpl{
		logger:             logging.EnsureLogger(logger).WithName(o11yName),
		tracer:             tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		identityRepository: identityRepository,
	}
}

func (s *serviceImpl) AcceptAccountInvitation(ctx context.Context, request *identitysvc.AcceptAccountInvitationRequest) (*identitysvc.AcceptAccountInvitationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.AcceptAccountInvitationResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) ArchiveAccount(ctx context.Context, request *identitysvc.ArchiveAccountRequest) (*identitysvc.ArchiveAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.ArchiveAccountResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) ArchiveUserMembership(ctx context.Context, request *identitysvc.ArchiveUserMembershipRequest) (*identitysvc.ArchiveUserMembershipResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.ArchiveUserMembershipResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) ArchiveUser(ctx context.Context, request *identitysvc.ArchiveUserRequest) (*identitysvc.ArchiveUserResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.ArchiveUserResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) CancelAccountInvitation(ctx context.Context, request *identitysvc.CancelAccountInvitationRequest) (*identitysvc.CancelAccountInvitationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.CancelAccountInvitationResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) CreateAccount(ctx context.Context, request *identitysvc.CreateAccountRequest) (*identitysvc.CreateAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.CreateAccountResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) CreateAccountInvitation(ctx context.Context, request *identitysvc.CreateAccountInvitationRequest) (*identitysvc.CreateAccountInvitationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.CreateAccountInvitationResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) CreateUser(ctx context.Context, request *identitysvc.CreateUserRequest) (*identitysvc.CreateUserResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.CreateUserResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) GetAccount(ctx context.Context, request *identitysvc.GetAccountRequest) (*identitysvc.GetAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.GetAccountResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) GetAccountInvitation(ctx context.Context, request *identitysvc.GetAccountInvitationRequest) (*identitysvc.GetAccountInvitationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.GetAccountInvitationResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) GetAccountInvitationByID(ctx context.Context, request *identitysvc.GetAccountInvitationByIDRequest) (*identitysvc.GetAccountInvitationByIDResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.GetAccountInvitationByIDResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) GetAccounts(ctx context.Context, request *identitysvc.GetAccountsRequest) (*identitysvc.GetAccountsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.GetAccountsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) GetReceivedAccountInvitations(ctx context.Context, request *identitysvc.GetReceivedAccountInvitationsRequest) (*identitysvc.GetReceivedAccountInvitationsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.GetReceivedAccountInvitationsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) GetSentAccountInvitations(ctx context.Context, request *identitysvc.GetSentAccountInvitationsRequest) (*identitysvc.GetSentAccountInvitationsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.GetSentAccountInvitationsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) GetUser(ctx context.Context, request *identitysvc.GetUserRequest) (*identitysvc.GetUserResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.UserIDKey: request.UserID,
	}, span, s.logger)

	user, err := s.identityRepository.GetUser(ctx, request.UserID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch users from database")
	}

	x := &identitysvc.GetUserResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: ConvertUserToGRPCUser(user),
	}

	return x, nil
}

func (s *serviceImpl) RejectAccountInvitation(ctx context.Context, request *identitysvc.RejectAccountInvitationRequest) (*identitysvc.RejectAccountInvitationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.RejectAccountInvitationResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) GetUsers(ctx context.Context, request *identitysvc.GetUsersRequest) (*identitysvc.GetUsersResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	users, err := s.identityRepository.GetUsers(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch users from database")
	}

	x := &identitysvc.GetUsersResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, user := range users.Data {
		x.Result = append(x.Result, ConvertUserToGRPCUser(user))
	}

	return x, nil
}

func (s *serviceImpl) SearchForUsers(ctx context.Context, request *identitysvc.SearchForUsersRequest) (*identitysvc.SearchForUsersResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.SearchForUsersResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) SetDefaultAccount(ctx context.Context, request *identitysvc.SetDefaultAccountRequest) (*identitysvc.SetDefaultAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.SetDefaultAccountResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) TransferAccountOwnership(ctx context.Context, request *identitysvc.TransferAccountOwnershipRequest) (*identitysvc.TransferAccountOwnershipResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.TransferAccountOwnershipResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) UpdateAccount(ctx context.Context, request *identitysvc.UpdateAccountRequest) (*identitysvc.UpdateAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.UpdateAccountResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) UpdateAccountMemberPermissions(ctx context.Context, request *identitysvc.UpdateAccountMemberPermissionsRequest) (*identitysvc.UpdateAccountMemberPermissionsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.UpdateAccountMemberPermissionsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) UpdateUserDetails(ctx context.Context, request *identitysvc.UpdateUserDetailsRequest) (*identitysvc.UpdateUserDetailsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	input := ConvertGRPCUserDetailsUpdateRequestInputToUserDetailsDatabaseUpdateInput(request.Input)

	if err = s.identityRepository.UpdateUserDetails(ctx, sessionContextData.GetUserID(), input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update user details")
	}

	x := &identitysvc.UpdateUserDetailsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) UpdateUserEmailAddress(ctx context.Context, request *identitysvc.UpdateUserEmailAddressRequest) (*identitysvc.UpdateUserEmailAddressResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	if err = s.identityRepository.UpdateUserEmailAddress(ctx, sessionContextData.GetUserID(), request.NewEmailAddress); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update user email address")
	}

	x := &identitysvc.UpdateUserEmailAddressResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) UpdateUserUsername(ctx context.Context, request *identitysvc.UpdateUserUsernameRequest) (*identitysvc.UpdateUserUsernameResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	if err = s.identityRepository.UpdateUserUsername(ctx, sessionContextData.GetUserID(), request.NewUsername); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update user username")
	}

	x := &identitysvc.UpdateUserUsernameResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) UploadUserAvatar(ctx context.Context, request *identitysvc.UploadUserAvatarRequest) (*identitysvc.UploadUserAvatarResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.UploadUserAvatarResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}
