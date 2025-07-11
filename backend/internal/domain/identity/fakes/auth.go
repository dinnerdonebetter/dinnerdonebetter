package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeSessionContextData builds a faked ContextData.
func BuildFakeSessionContextData() *sessions.ContextData {
	return &sessions.ContextData{
		AccountPermissions: map[string]authorization.AccountRolePermissionsChecker{},
		Requester: sessions.RequesterInfo{
			ServicePermissions:       nil,
			AccountStatus:            string(identity.GoodStandingUserAccountStatus),
			AccountStatusExplanation: "fake",
			UserID:                   BuildFakeID(),
			EmailAddress:             fake.Email(),
			Username:                 buildUniqueString(),
		},
		ActiveAccountID: BuildFakeID(),
	}
}

// BuildFakeUserPermissionModificationInput builds a faked ModifyUserPermissionsInput.
func BuildFakeUserPermissionModificationInput() *identity.ModifyUserPermissionsInput {
	return &identity.ModifyUserPermissionsInput{
		Reason:  fake.Sentence(10),
		NewRole: authorization.AccountMemberRole.String(),
	}
}

// BuildFakeTransferAccountOwnershipInput builds a faked AccountOwnershipTransferInput.
func BuildFakeTransferAccountOwnershipInput() *identity.AccountOwnershipTransferInput {
	return &identity.AccountOwnershipTransferInput{
		Reason:       fake.Sentence(10),
		CurrentOwner: fake.UUID(),
		NewOwner:     fake.UUID(),
	}
}

// BuildFakeChangeActiveAccountInput builds a faked ChangeActiveAccountInput.
func BuildFakeChangeActiveAccountInput() *identity.ChangeActiveAccountInput {
	return &identity.ChangeActiveAccountInput{
		AccountID: fake.UUID(),
	}
}

// BuildFakeUserStatusResponse builds a faked UserStatusResponse.
func BuildFakeUserStatusResponse() *identity.UserStatusResponse {
	return &identity.UserStatusResponse{
		UserID:                   BuildFakeID(),
		AccountStatus:            string(identity.GoodStandingUserAccountStatus),
		AccountStatusExplanation: "",
		ActiveAccount:            BuildFakeID(),
		UserIsAuthenticated:      true,
	}
}

// BuildFakeTokenResponse builds a faked TokenResponse.
func BuildFakeTokenResponse() *identity.TokenResponse {
	return &identity.TokenResponse{
		UserID:      BuildFakeID(),
		AccountID:   BuildFakeID(),
		AccessToken: fake.UUID(),
	}
}

func BuildFakeUserLoginInput() *identity.UserLoginInput {
	return &identity.UserLoginInput{
		Username:  BuildFakeUser().Username,
		Password:  buildFakePassword(),
		TOTPToken: buildFakeTOTPToken(),
	}
}
