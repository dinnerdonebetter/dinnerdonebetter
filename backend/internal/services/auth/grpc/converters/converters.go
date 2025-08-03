package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication"
)

func ConvertGRPCUpdatePasswordRequestToPasswordUpdateInput(request *authsvc.UpdatePasswordRequest) *identity.PasswordUpdateInput {
	return &identity.PasswordUpdateInput{
		NewPassword:     request.NewPassword,
		CurrentPassword: request.CurrentPassword,
		TOTPToken:       request.TOTPToken,
	}
}

func ConvertGRPCVerifyTOTPSecretRequestToTOTPSecretVerificationInput(request *authsvc.VerifyTOTPSecretRequest) *identity.TOTPSecretVerificationInput {
	return &identity.TOTPSecretVerificationInput{
		TOTPToken: request.TOTPToken,
		UserID:    request.UserID,
	}
}

func ConvertGRPCVerifyEmailAddressRequestToEmailAddressVerificationRequestInput(request *authsvc.VerifyEmailAddressRequest) *identity.EmailAddressVerificationRequestInput {
	return &identity.EmailAddressVerificationRequestInput{
		Token: request.Token,
	}
}

func ConvertGRPCRequestUsernameReminderRequestToUsernameReminderRequestInput(request *authsvc.RequestUsernameReminderRequest) *identity.UsernameReminderRequestInput {
	return &identity.UsernameReminderRequestInput{
		EmailAddress: request.EmailAddress,
	}
}

func ConvertGRPCRequestPasswordResetTokenRequestToPasswordResetTokenCreationRequestInput(request *authsvc.RequestPasswordResetTokenRequest) *identity.PasswordResetTokenCreationRequestInput {
	return &identity.PasswordResetTokenCreationRequestInput{EmailAddress: request.EmailAddress}
}

func ConvertGRPCRefreshTOTPSecretRequestToTOTPSecretRefreshInput(request *authsvc.RefreshTOTPSecretRequest) *identity.TOTPSecretRefreshInput {
	return &identity.TOTPSecretRefreshInput{
		CurrentPassword: request.CurrentPassword,
		TOTPToken:       request.TOTPToken,
	}
}

func ConvertTOTPSecretRefreshResponseToGRPCTOTPSecretRefreshResponse(input *identity.TOTPSecretRefreshResponse) *authsvc.TOTPSecretRefreshResponse {
	return &authsvc.TOTPSecretRefreshResponse{
		TwoFactorQRCode: input.TwoFactorQRCode,
		TwoFactorSecret: input.TwoFactorSecret,
	}
}

func ConvertGRPCRedeemPasswordResetTokenRequestToPasswordResetTokenRedemptionRequestInput(request *authsvc.RedeemPasswordResetTokenRequest) *identity.PasswordResetTokenRedemptionRequestInput {
	return &identity.PasswordResetTokenRedemptionRequestInput{
		Token:       request.Token,
		NewPassword: request.NewPassword,
	}
}

func ConvertGRPCCheckPermissionsRequestToUserPermissionsRequestInput(request *authsvc.UserPermissionsRequestInput) *identity.UserPermissionsRequestInput {
	return &identity.UserPermissionsRequestInput{
		Permissions: request.Permissions,
	}
}

func ConvertGRPCAdminLoginForTokenRequestToUserLoginInput(request *authsvc.UserLoginInput) *authentication.UserLoginInput {
	return &authentication.UserLoginInput{
		Username:  request.Username,
		Password:  request.Password,
		TOTPToken: request.TOTPToken,
	}
}

func ConvertTokenResponseToGRPCTokenResponse(input *identity.TokenResponse) *authsvc.TokenResponse {
	return &authsvc.TokenResponse{
		UserID:       input.UserID,
		AccountID:    input.AccountID,
		AccessToken:  input.AccessToken,
		RefreshToken: input.RefreshToken,
		ExpiresUTC:   grpcconverters.ConvertTimeToPBTimestamp(input.ExpiresUTC),
	}
}
