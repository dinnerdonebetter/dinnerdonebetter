package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "identity_service"
)

var _ identitysvc.IdentityServiceServer = (*ServiceImpl)(nil)

type (
	ServiceImpl struct {
		identitysvc.UnimplementedIdentityServiceServer
		tracer             tracing.Tracer
		logger             logging.Logger
		identityRepository identity.Repository
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	identityRepository identity.Repository,
) identitysvc.IdentityServiceServer {
	return &ServiceImpl{
		logger:             logging.EnsureLogger(logger).WithName(o11yName),
		tracer:             tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		identityRepository: identityRepository,
	}
}

func (s *ServiceImpl) AcceptAccountInvitation(ctx context.Context, request *identitysvc.AcceptAccountInvitationRequest) (*identitysvc.AcceptAccountInvitationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) ArchiveAccount(ctx context.Context, request *identitysvc.ArchiveAccountRequest) (*identitysvc.ArchiveAccountResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) ArchiveUserMembership(ctx context.Context, request *identitysvc.ArchiveUserMembershipRequest) (*identitysvc.ArchiveUserMembershipResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) ArchiveUser(ctx context.Context, request *identitysvc.ArchiveUserRequest) (*identitysvc.ArchiveUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) CancelAccountInvitation(ctx context.Context, request *identitysvc.CancelAccountInvitationRequest) (*identitysvc.CancelAccountInvitationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) CreateAccount(ctx context.Context, request *identitysvc.CreateAccountRequest) (*identitysvc.CreateAccountResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) CreateAccountInvitation(ctx context.Context, request *identitysvc.CreateAccountInvitationRequest) (*identitysvc.CreateAccountInvitationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) CreateUser(ctx context.Context, request *identitysvc.CreateUserRequest) (*identitysvc.CreateUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) GetAccount(ctx context.Context, request *identitysvc.GetAccountRequest) (*identitysvc.GetAccountResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) GetAccountInvitation(ctx context.Context, request *identitysvc.GetAccountInvitationRequest) (*identitysvc.GetAccountInvitationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) GetAccountInvitationByID(ctx context.Context, request *identitysvc.GetAccountInvitationByIDRequest) (*identitysvc.GetAccountInvitationByIDResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) GetAccounts(ctx context.Context, request *identitysvc.GetAccountsRequest) (*identitysvc.GetAccountsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) GetReceivedAccountInvitations(ctx context.Context, request *identitysvc.GetReceivedAccountInvitationsRequest) (*identitysvc.GetReceivedAccountInvitationsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) GetSentAccountInvitations(ctx context.Context, request *identitysvc.GetSentAccountInvitationsRequest) (*identitysvc.GetSentAccountInvitationsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) GetUser(ctx context.Context, request *identitysvc.GetUserRequest) (*identitysvc.GetUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) RejectAccountInvitation(ctx context.Context, request *identitysvc.RejectAccountInvitationRequest) (*identitysvc.RejectAccountInvitationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) GetUsers(ctx context.Context, request *identitysvc.GetUsersRequest) (*identitysvc.GetUsersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) SearchForUsers(ctx context.Context, request *identitysvc.SearchForUsersRequest) (*identitysvc.SearchForUsersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) SetDefaultAccount(ctx context.Context, request *identitysvc.SetDefaultAccountRequest) (*identitysvc.SetDefaultAccountResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) TransferAccountOwnership(ctx context.Context, request *identitysvc.TransferAccountOwnershipRequest) (*identitysvc.TransferAccountOwnershipResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) UpdateAccount(ctx context.Context, request *identitysvc.UpdateAccountRequest) (*identitysvc.UpdateAccountResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) UpdateAccountMemberPermissions(ctx context.Context, request *identitysvc.UpdateAccountMemberPermissionsRequest) (*identitysvc.UpdateAccountMemberPermissionsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) UpdateUserDetails(ctx context.Context, request *identitysvc.UpdateUserDetailsRequest) (*identitysvc.UpdateUserDetailsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) UpdateUserEmailAddress(ctx context.Context, request *identitysvc.UpdateUserEmailAddressRequest) (*identitysvc.UpdateUserEmailAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) UpdateUserUsername(ctx context.Context, request *identitysvc.UpdateUserUsernameRequest) (*identitysvc.UpdateUserUsernameResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) UploadUserAvatar(ctx context.Context, request *identitysvc.UploadUserAvatarRequest) (*identitysvc.UploadUserAvatarResponse, error) {
	//TODO implement me
	panic("implement me")
}
