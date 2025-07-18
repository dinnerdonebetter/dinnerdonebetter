package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "auth_service"
)

var _ authsvc.AuthServiceServer = (*ServiceImpl)(nil)

type (
	ServiceImpl struct {
		authsvc.UnimplementedAuthServiceServer
		tracer             tracing.Tracer
		logger             logging.Logger
		identityRepository identity.Repository
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	identityRepository identity.Repository,
) authsvc.AuthServiceServer {
	return &ServiceImpl{
		logger:             logging.EnsureLogger(logger).WithName(o11yName),
		tracer:             tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		identityRepository: identityRepository,
	}
}

func (s *ServiceImpl) GetAuthStatus(ctx context.Context, request *authsvc.GetAuthStatusRequest) (*authsvc.GetAuthStatusResponse, error) {
	x := &authsvc.GetAuthStatusResponse{}

	return x, nil
}

func (s *ServiceImpl) ExchangeToken(ctx context.Context, request *authsvc.ExchangeTokenRequest) (*authsvc.ExchangeTokenResponse, error) {
	x := &authsvc.ExchangeTokenResponse{}

	return x, nil
}

func (s *ServiceImpl) AdminLoginForToken(ctx context.Context, request *authsvc.AdminLoginForTokenRequest) (*authsvc.AdminLoginForTokenResponse, error) {
	x := &authsvc.AdminLoginForTokenResponse{}

	return x, nil
}

func (s *ServiceImpl) CheckPermissions(ctx context.Context, request *authsvc.CheckPermissionsRequest) (*authsvc.CheckPermissionsResponse, error) {
	x := &authsvc.CheckPermissionsResponse{}

	return x, nil
}

func (s *ServiceImpl) GetActiveAccount(ctx context.Context, request *authsvc.GetActiveAccountRequest) (*authsvc.GetActiveAccountResponse, error) {
	x := &authsvc.GetActiveAccountResponse{}

	return x, nil
}

func (s *ServiceImpl) GetSelf(ctx context.Context, request *authsvc.GetSelfRequest) (*authsvc.GetSelfResponse, error) {
	x := &authsvc.GetSelfResponse{}

	return x, nil
}

func (s *ServiceImpl) LoginForToken(ctx context.Context, request *authsvc.LoginForTokenRequest) (*authsvc.LoginForTokenResponse, error) {
	x := &authsvc.LoginForTokenResponse{}

	return x, nil
}

func (s *ServiceImpl) RedeemPasswordResetToken(ctx context.Context, request *authsvc.RedeemPasswordResetTokenRequest) (*authsvc.RedeemPasswordResetTokenResponse, error) {
	x := &authsvc.RedeemPasswordResetTokenResponse{}

	return x, nil
}

func (s *ServiceImpl) RefreshTOTPSecret(ctx context.Context, request *authsvc.RefreshTOTPSecretRequest) (*authsvc.RefreshTOTPSecretResponse, error) {
	x := &authsvc.RefreshTOTPSecretResponse{}

	return x, nil
}

func (s *ServiceImpl) RequestEmailVerificationEmail(ctx context.Context, request *authsvc.RequestEmailVerificationEmailRequest) (*authsvc.RequestEmailVerificationEmailResponse, error) {
	x := &authsvc.RequestEmailVerificationEmailResponse{}

	return x, nil
}

func (s *ServiceImpl) RequestPasswordResetToken(ctx context.Context, request *authsvc.RequestPasswordResetTokenRequest) (*authsvc.RequestPasswordResetTokenResponse, error) {
	x := &authsvc.RequestPasswordResetTokenResponse{}

	return x, nil
}

func (s *ServiceImpl) RequestUsernameReminder(ctx context.Context, request *authsvc.RequestUsernameReminderRequest) (*authsvc.RequestUsernameReminderResponse, error) {
	x := &authsvc.RequestUsernameReminderResponse{}

	return x, nil
}

func (s *ServiceImpl) VerifyEmailAddress(ctx context.Context, request *authsvc.VerifyEmailAddressRequest) (*authsvc.VerifyEmailAddressResponse, error) {
	x := &authsvc.VerifyEmailAddressResponse{}

	return x, nil
}

func (s *ServiceImpl) VerifyTOTPSecret(ctx context.Context, request *authsvc.VerifyTOTPSecretRequest) (*authsvc.VerifyTOTPSecretResponse, error) {
	x := &authsvc.VerifyTOTPSecretResponse{}

	return x, nil
}

func (s *ServiceImpl) UpdatePassword(ctx context.Context, request *authsvc.UpdatePasswordRequest) (*authsvc.UpdatePasswordResponse, error) {
	x := &authsvc.UpdatePasswordResponse{}

	return x, nil
}
