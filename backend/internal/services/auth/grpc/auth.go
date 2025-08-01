package grpc

import (
	"context"

	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
)

func (s *serviceImpl) GetAuthStatus(ctx context.Context, request *authsvc.GetAuthStatusRequest) (*authsvc.GetAuthStatusResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &authsvc.GetAuthStatusResponse{}

	return x, nil
}

func (s *serviceImpl) ExchangeToken(ctx context.Context, request *authsvc.ExchangeTokenRequest) (*authsvc.ExchangeTokenResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &authsvc.ExchangeTokenResponse{}

	return x, nil
}

func (s *serviceImpl) AdminLoginForToken(ctx context.Context, request *authsvc.AdminLoginForTokenRequest) (*authsvc.AdminLoginForTokenResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &authsvc.AdminLoginForTokenResponse{}

	return x, nil
}

func (s *serviceImpl) CheckPermissions(ctx context.Context, request *authsvc.CheckPermissionsRequest) (*authsvc.CheckPermissionsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &authsvc.CheckPermissionsResponse{}

	return x, nil
}

func (s *serviceImpl) GetActiveAccount(ctx context.Context, request *authsvc.GetActiveAccountRequest) (*authsvc.GetActiveAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &authsvc.GetActiveAccountResponse{}

	return x, nil
}

func (s *serviceImpl) GetSelf(ctx context.Context, request *authsvc.GetSelfRequest) (*authsvc.GetSelfResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &authsvc.GetSelfResponse{}

	return x, nil
}

func (s *serviceImpl) LoginForToken(ctx context.Context, request *authsvc.LoginForTokenRequest) (*authsvc.LoginForTokenResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &authsvc.LoginForTokenResponse{}

	return x, nil
}

func (s *serviceImpl) RedeemPasswordResetToken(ctx context.Context, request *authsvc.RedeemPasswordResetTokenRequest) (*authsvc.RedeemPasswordResetTokenResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &authsvc.RedeemPasswordResetTokenResponse{}

	return x, nil
}

func (s *serviceImpl) RefreshTOTPSecret(ctx context.Context, request *authsvc.RefreshTOTPSecretRequest) (*authsvc.RefreshTOTPSecretResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &authsvc.RefreshTOTPSecretResponse{}

	return x, nil
}

func (s *serviceImpl) RequestEmailVerificationEmail(ctx context.Context, request *authsvc.RequestEmailVerificationEmailRequest) (*authsvc.RequestEmailVerificationEmailResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &authsvc.RequestEmailVerificationEmailResponse{}

	return x, nil
}

func (s *serviceImpl) RequestPasswordResetToken(ctx context.Context, request *authsvc.RequestPasswordResetTokenRequest) (*authsvc.RequestPasswordResetTokenResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &authsvc.RequestPasswordResetTokenResponse{}

	return x, nil
}

func (s *serviceImpl) RequestUsernameReminder(ctx context.Context, request *authsvc.RequestUsernameReminderRequest) (*authsvc.RequestUsernameReminderResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &authsvc.RequestUsernameReminderResponse{}

	return x, nil
}

func (s *serviceImpl) VerifyEmailAddress(ctx context.Context, request *authsvc.VerifyEmailAddressRequest) (*authsvc.VerifyEmailAddressResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &authsvc.VerifyEmailAddressResponse{}

	return x, nil
}

func (s *serviceImpl) VerifyTOTPSecret(ctx context.Context, request *authsvc.VerifyTOTPSecretRequest) (*authsvc.VerifyTOTPSecretResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &authsvc.VerifyTOTPSecretResponse{}

	return x, nil
}

func (s *serviceImpl) UpdatePassword(ctx context.Context, request *authsvc.UpdatePasswordRequest) (*authsvc.UpdatePasswordResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &authsvc.UpdatePasswordResponse{}

	return x, nil
}
