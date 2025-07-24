package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	oauth2server "github.com/go-oauth2/oauth2/v4"
)

const (
	o11yName = "auth_service"
)

var _ authsvc.AuthServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		authsvc.UnimplementedAuthServiceServer
		tracer              tracing.Tracer
		logger              logging.Logger
		identityRepository  identity.Repository
		oauth2ClientManager oauth2server.Manager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	identityRepository identity.Repository,
) authsvc.AuthServiceServer {
	return &serviceImpl{
		logger:             logging.EnsureLogger(logger).WithName(o11yName),
		tracer:             tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		identityRepository: identityRepository,
	}
}

func (s *serviceImpl) GetAuthStatus(ctx context.Context, request *authsvc.GetAuthStatusRequest) (*authsvc.GetAuthStatusResponse, error) {
	x := &authsvc.GetAuthStatusResponse{}

	return x, nil
}

func (s *serviceImpl) ExchangeToken(ctx context.Context, request *authsvc.ExchangeTokenRequest) (*authsvc.ExchangeTokenResponse, error) {
	x := &authsvc.ExchangeTokenResponse{}

	return x, nil
}

func (s *serviceImpl) AdminLoginForToken(ctx context.Context, request *authsvc.AdminLoginForTokenRequest) (*authsvc.AdminLoginForTokenResponse, error) {
	x := &authsvc.AdminLoginForTokenResponse{}

	return x, nil
}

func (s *serviceImpl) CheckPermissions(ctx context.Context, request *authsvc.CheckPermissionsRequest) (*authsvc.CheckPermissionsResponse, error) {
	x := &authsvc.CheckPermissionsResponse{}

	return x, nil
}

func (s *serviceImpl) GetActiveAccount(ctx context.Context, request *authsvc.GetActiveAccountRequest) (*authsvc.GetActiveAccountResponse, error) {
	x := &authsvc.GetActiveAccountResponse{}

	return x, nil
}

func (s *serviceImpl) GetSelf(ctx context.Context, request *authsvc.GetSelfRequest) (*authsvc.GetSelfResponse, error) {
	x := &authsvc.GetSelfResponse{}

	return x, nil
}

func (s *serviceImpl) LoginForToken(ctx context.Context, request *authsvc.LoginForTokenRequest) (*authsvc.LoginForTokenResponse, error) {
	x := &authsvc.LoginForTokenResponse{}

	return x, nil
}

func (s *serviceImpl) RedeemPasswordResetToken(ctx context.Context, request *authsvc.RedeemPasswordResetTokenRequest) (*authsvc.RedeemPasswordResetTokenResponse, error) {
	x := &authsvc.RedeemPasswordResetTokenResponse{}

	return x, nil
}

func (s *serviceImpl) RefreshTOTPSecret(ctx context.Context, request *authsvc.RefreshTOTPSecretRequest) (*authsvc.RefreshTOTPSecretResponse, error) {
	x := &authsvc.RefreshTOTPSecretResponse{}

	return x, nil
}

func (s *serviceImpl) RequestEmailVerificationEmail(ctx context.Context, request *authsvc.RequestEmailVerificationEmailRequest) (*authsvc.RequestEmailVerificationEmailResponse, error) {
	x := &authsvc.RequestEmailVerificationEmailResponse{}

	return x, nil
}

func (s *serviceImpl) RequestPasswordResetToken(ctx context.Context, request *authsvc.RequestPasswordResetTokenRequest) (*authsvc.RequestPasswordResetTokenResponse, error) {
	x := &authsvc.RequestPasswordResetTokenResponse{}

	return x, nil
}

func (s *serviceImpl) RequestUsernameReminder(ctx context.Context, request *authsvc.RequestUsernameReminderRequest) (*authsvc.RequestUsernameReminderResponse, error) {
	x := &authsvc.RequestUsernameReminderResponse{}

	return x, nil
}

func (s *serviceImpl) VerifyEmailAddress(ctx context.Context, request *authsvc.VerifyEmailAddressRequest) (*authsvc.VerifyEmailAddressResponse, error) {
	x := &authsvc.VerifyEmailAddressResponse{}

	return x, nil
}

func (s *serviceImpl) VerifyTOTPSecret(ctx context.Context, request *authsvc.VerifyTOTPSecretRequest) (*authsvc.VerifyTOTPSecretResponse, error) {
	x := &authsvc.VerifyTOTPSecretResponse{}

	return x, nil
}

func (s *serviceImpl) UpdatePassword(ctx context.Context, request *authsvc.UpdatePasswordRequest) (*authsvc.UpdatePasswordResponse, error) {
	x := &authsvc.UpdatePasswordResponse{}

	return x, nil
}
