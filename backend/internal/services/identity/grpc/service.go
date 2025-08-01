package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identityconverters "github.com/dinnerdonebetter/backend/internal/domain/identity/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/identity/managers"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/identity/grpc/converters"

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
		sessionContextDataFetcher func(ctx context.Context) (*sessions.ContextData, error)
		identityRepository        identity.Repository
		identityDataManager       managers.IdentityDataManager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	sessionContextDataFetcher func(ctx context.Context) (*sessions.ContextData, error),
	identityRepository identity.Repository,
	identityDataManager managers.IdentityDataManager,
) identitysvc.IdentityServiceServer {
	return &serviceImpl{
		logger:                    logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		sessionContextDataFetcher: sessionContextDataFetcher,
		identityRepository:        identityRepository,
		identityDataManager:       identityDataManager,
	}
}

func (s *serviceImpl) buildResponseDetails(ctx context.Context, span tracing.Span) *types.ResponseDetails {
	out := &types.ResponseDetails{}
	if span != nil {
		out.TraceID = span.SpanContext().TraceID().String()
	}

	if sessionContextData, err := s.sessionContextDataFetcher(ctx); err == nil && sessionContextData != nil {
		out.CurrentAccountID = sessionContextData.GetActiveAccountID()
	}

	return out
}

func (s *serviceImpl) AcceptAccountInvitation(ctx context.Context, request *identitysvc.AcceptAccountInvitationRequest) (*identitysvc.AcceptAccountInvitationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.AccountInvitationIDKey: request.AccountInvitationID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	input := converters.ConvertGRPCAccountInvitationUpdateRequestInputToAccountInvitationUpdateRequestInput(request.Input)
	if err = s.identityDataManager.AcceptAccountInvitation(ctx, sessionContextData.GetActiveAccountID(), request.AccountInvitationID, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "accepting account invitation")
	}

	x := &identitysvc.AcceptAccountInvitationResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	return x, nil
}

func (s *serviceImpl) RejectAccountInvitation(ctx context.Context, request *identitysvc.RejectAccountInvitationRequest) (*identitysvc.RejectAccountInvitationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.AccountInvitationIDKey: request.AccountInvitationID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	input := converters.ConvertGRPCAccountInvitationUpdateRequestInputToAccountInvitationUpdateRequestInput(request.Input)
	if err = s.identityDataManager.RejectAccountInvitation(ctx, sessionContextData.GetActiveAccountID(), request.AccountInvitationID, input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to reject account invitation")
	}

	x := &identitysvc.RejectAccountInvitationResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	return x, nil
}

func (s *serviceImpl) CancelAccountInvitation(ctx context.Context, request *identitysvc.CancelAccountInvitationRequest) (*identitysvc.CancelAccountInvitationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.AccountInvitationIDKey: request.AccountInvitationID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	if err = s.identityDataManager.CancelAccountInvitation(ctx, sessionContextData.GetActiveAccountID(), request.AccountInvitationID, request.Input.Note); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Internal, "failed to cancel account invitation")
	}

	x := &identitysvc.CancelAccountInvitationResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	return x, nil
}

func (s *serviceImpl) ArchiveAccount(ctx context.Context, request *identitysvc.ArchiveAccountRequest) (*identitysvc.ArchiveAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.AccountIDKey: request.AccountID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	if err = s.identityDataManager.ArchiveAccount(ctx, request.AccountID, sessionContextData.GetUserID()); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive account")
	}

	x := &identitysvc.ArchiveAccountResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	return x, nil
}

func (s *serviceImpl) ArchiveUserMembership(ctx context.Context, request *identitysvc.ArchiveUserMembershipRequest) (*identitysvc.ArchiveUserMembershipResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	// TODO: validate that the user is authorized to do this?

	if err := s.identityDataManager.ArchiveUserMembership(ctx, request.UserID, request.AccountID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Internal, "failed to archive user membership")
	}

	x := &identitysvc.ArchiveUserMembershipResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	return x, nil
}

func (s *serviceImpl) ArchiveUser(ctx context.Context, request *identitysvc.ArchiveUserRequest) (*identitysvc.ArchiveUserResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	if err := s.identityDataManager.ArchiveUser(ctx, request.UserID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Internal, "failed to archive user")
	}

	x := &identitysvc.ArchiveUserResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	return x, nil
}

func (s *serviceImpl) CreateAccount(ctx context.Context, request *identitysvc.CreateAccountRequest) (*identitysvc.CreateAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	input := converters.ConvertGRPCAccountCreationRequestInputToAccountCreationRequestInput(request.Input)

	created, err := s.identityDataManager.CreateAccount(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Internal, "failed to create account")
	}

	x := &identitysvc.CreateAccountResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
		Created:         converters.ConvertAccountToGRPCAccount(created),
	}

	return x, nil
}

func (s *serviceImpl) CreateAccountInvitation(ctx context.Context, request *identitysvc.CreateAccountInvitationRequest) (*identitysvc.CreateAccountInvitationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	input := converters.ConvertGRPCAccountInvitationCreationRequestInputToAccountInvitationCreationRequestInput(request.Input)
	created, err := s.identityDataManager.CreateAccountInvitation(ctx, request.AccountID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Internal, "failed to create account invitation")
	}

	x := &identitysvc.CreateAccountInvitationResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
		Created:         converters.ConvertAccountInvitationToGRPCAccountInvitation(created),
	}

	return x, nil
}

func (s *serviceImpl) CreateUser(ctx context.Context, request *identitysvc.CreateUserRequest) (*identitysvc.CreateUserResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	input := converters.ConvertGRPCUserRegistrationInputToUserRegistrationInput(request.Input)

	created, err := s.identityDataManager.CreateUser(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Internal, "failed to create user")
	}

	x := &identitysvc.CreateUserResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
		Created:         converters.ConvertUserCreationResponseToGRPCUserCreationResponse(created),
	}

	return x, nil
}

func (s *serviceImpl) GetAccount(ctx context.Context, request *identitysvc.GetAccountRequest) (*identitysvc.GetAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	account, err := s.identityDataManager.GetAccount(ctx, request.AccountID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Internal, "failed to get account")
	}

	x := &identitysvc.GetAccountResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
		Result:          converters.ConvertAccountToGRPCAccount(account),
	}

	return x, nil
}

func (s *serviceImpl) GetAccountInvitation(ctx context.Context, request *identitysvc.GetAccountInvitationRequest) (*identitysvc.GetAccountInvitationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	accountInvitation, err := s.identityDataManager.GetAccountInvitation(ctx, sessionContextData.GetActiveAccountID(), request.AccountInvitationID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Internal, "failed to get account invitation")
	}

	x := &identitysvc.GetAccountInvitationResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
		Result:          converters.ConvertAccountInvitationToGRPCAccountInvitation(accountInvitation),
	}

	return x, nil
}

func (s *serviceImpl) GetAccounts(ctx context.Context, request *identitysvc.GetAccountsRequest) (*identitysvc.GetAccountsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	accounts, _, err := s.identityDataManager.GetAccounts(ctx, sessionContextData.GetUserID(), filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, s.logger, span, codes.Internal, "failed to get accounts")
	}

	x := &identitysvc.GetAccountsResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	for _, account := range accounts {
		x.Result = append(x.Result, converters.ConvertAccountToGRPCAccount(account))
	}

	return x, nil
}

func (s *serviceImpl) GetReceivedAccountInvitations(ctx context.Context, request *identitysvc.GetReceivedAccountInvitationsRequest) (*identitysvc.GetReceivedAccountInvitationsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	logger := filter.AttachToLogger(s.logger.WithSpan(span))

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	invites, _, err := s.identityDataManager.GetReceivedAccountInvitations(ctx, sessionContextData.GetUserID(), filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get received account invitations")
	}

	x := &identitysvc.GetReceivedAccountInvitationsResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	for _, invite := range invites {
		x.Result = append(x.Result, converters.ConvertAccountInvitationToGRPCAccountInvitation(invite))
	}

	return x, nil
}

func (s *serviceImpl) GetSentAccountInvitations(ctx context.Context, request *identitysvc.GetSentAccountInvitationsRequest) (*identitysvc.GetSentAccountInvitationsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	logger := filter.AttachToLogger(s.logger.WithSpan(span))

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	invites, _, err := s.identityDataManager.GetSentAccountInvitations(ctx, sessionContextData.GetUserID(), filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get received account invitations")
	}

	x := &identitysvc.GetSentAccountInvitationsResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	for _, invite := range invites {
		x.Result = append(x.Result, converters.ConvertAccountInvitationToGRPCAccountInvitation(invite))
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
		ResponseDetails: s.buildResponseDetails(ctx, span),
		Result:          converters.ConvertUserToGRPCUser(user),
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
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	for _, user := range users.Data {
		x.Result = append(x.Result, converters.ConvertUserToGRPCUser(user))
	}

	return x, nil
}

func (s *serviceImpl) SearchForUsers(ctx context.Context, request *identitysvc.SearchForUsersRequest) (*identitysvc.SearchForUsersResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.SearchQueryKey: request.Query,
	}, span, s.logger)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	users, _, err := s.identityDataManager.SearchForUsers(ctx, request.Query, request.UseDatabase, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to search for users")
	}

	x := &identitysvc.SearchForUsersResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	for _, user := range users {
		x.Result = append(x.Result, converters.ConvertUserToGRPCUser(user))
	}

	return x, nil
}

func (s *serviceImpl) SetDefaultAccount(ctx context.Context, request *identitysvc.SetDefaultAccountRequest) (*identitysvc.SetDefaultAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.AccountIDKey: request.AccountID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	if err = s.identityDataManager.SetDefaultAccount(ctx, sessionContextData.GetUserID(), request.AccountID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to set default account")
	}

	x := &identitysvc.SetDefaultAccountResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
		Success:         true,
	}

	return x, nil
}

func (s *serviceImpl) TransferAccountOwnership(ctx context.Context, request *identitysvc.TransferAccountOwnershipRequest) (*identitysvc.TransferAccountOwnershipResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.AccountIDKey: request.AccountID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	input := converters.ConvertGRPCAccountOwnershipTransferInputToAccountOwnershipTransferInput(request.Input)

	if err = s.identityDataManager.TransferAccountOwnership(ctx, sessionContextData.GetActiveAccountID(), input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to transfer account ownership")
	}

	x := &identitysvc.TransferAccountOwnershipResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
		Success:         true,
	}

	return x, nil
}

func (s *serviceImpl) UpdateAccount(ctx context.Context, request *identitysvc.UpdateAccountRequest) (*identitysvc.UpdateAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.AccountIDKey: request.AccountID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	input := converters.ConvertGRPCAccountUpdateRequestInputToAccountUpdateRequestInput(request.Input)

	if err = s.identityDataManager.UpdateAccount(ctx, sessionContextData.GetActiveAccountID(), input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update account")
	}

	x := &identitysvc.UpdateAccountResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	return x, nil
}

func (s *serviceImpl) UpdateAccountMemberPermissions(ctx context.Context, request *identitysvc.UpdateAccountMemberPermissionsRequest) (*identitysvc.UpdateAccountMemberPermissionsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.UserIDKey: request.UserID,
	}, span, s.logger)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	input := converters.ConvertGRPCModifyUserPermissionsInputToModifyUserPermissionsInput(request.Input)
	if err = s.identityDataManager.UpdateAccountMemberPermissions(ctx, request.UserID, sessionContextData.GetActiveAccountID(), input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update account member permissions")
	}

	x := &identitysvc.UpdateAccountMemberPermissionsResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
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

	input := converters.ConvertGRPCUserDetailsUpdateRequestInputToUserDetailsDatabaseUpdateInput(request.Input)

	if err = s.identityRepository.UpdateUserDetails(ctx, sessionContextData.GetUserID(), input); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update user details")
	}

	x := &identitysvc.UpdateUserDetailsResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
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
		ResponseDetails: s.buildResponseDetails(ctx, span),
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
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	return x, nil
}

func (s *serviceImpl) UploadUserAvatar(ctx context.Context, request *identitysvc.UploadUserAvatarRequest) (*identitysvc.UploadUserAvatarResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &identitysvc.UploadUserAvatarResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	return x, nil
}

func (s *serviceImpl) AdminUpdateUserStatus(ctx context.Context, request *identitysvc.AdminUpdateUserStatusRequest) (*identitysvc.AdminUpdateUserStatusResponse, error) {
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

	if err = s.identityRepository.UpdateUserAccountStatus(ctx, request.TargetUserID, identityconverters.ConvertGRPCAdminUpdateUserStatusRequestToUserAccountStatusUpdateInput(request)); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating user account status")
	}

	x := &identitysvc.AdminUpdateUserStatusResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	return x, nil
}
