package grpc

import (
	"context"

	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/services/identity/grpc/converters"

	"google.golang.org/grpc/codes"
)

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

func (s *serviceImpl) GetReceivedAccountInvitations(ctx context.Context, request *identitysvc.GetReceivedAccountInvitationsRequest) (*identitysvc.GetReceivedAccountInvitationsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	logger := filter.AttachToLogger(s.logger.WithSpan(span))

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to get session context data")
	}

	invites, err := s.identityDataManager.GetReceivedAccountInvitations(ctx, sessionContextData.GetUserID(), filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get received account invitations")
	}

	x := &identitysvc.GetReceivedAccountInvitationsResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	for _, invite := range invites.Data {
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

	invites, err := s.identityDataManager.GetSentAccountInvitations(ctx, sessionContextData.GetUserID(), filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to get received account invitations")
	}

	x := &identitysvc.GetSentAccountInvitationsResponse{
		ResponseDetails: s.buildResponseDetails(ctx, span),
	}

	for _, invite := range invites.Data {
		x.Result = append(x.Result, converters.ConvertAccountInvitationToGRPCAccountInvitation(invite))
	}

	return x, nil
}
