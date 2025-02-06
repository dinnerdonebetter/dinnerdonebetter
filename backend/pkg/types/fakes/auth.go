package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/tokens"
	"github.com/dinnerdonebetter/backend/pkg/types"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeSessionContextData builds a faked ContextData.
func BuildFakeSessionContextData() *sessions.ContextData {
	return &sessions.ContextData{
		HouseholdPermissions: map[string]authorization.HouseholdRolePermissionsChecker{},
		Requester: sessions.RequesterInfo{
			ServicePermissions:       nil,
			AccountStatus:            string(types.GoodStandingUserAccountStatus),
			AccountStatusExplanation: "fake",
			UserID:                   BuildFakeID(),
			EmailAddress:             fake.Email(),
			Username:                 buildUniqueString(),
		},
		ActiveHouseholdID: BuildFakeID(),
	}
}

// BuildFakeUserPermissionModificationInput builds a faked ModifyUserPermissionsInput.
func BuildFakeUserPermissionModificationInput() *types.ModifyUserPermissionsInput {
	return &types.ModifyUserPermissionsInput{
		Reason:  fake.Sentence(10),
		NewRole: authorization.HouseholdMemberRole.String(),
	}
}

// BuildFakeTransferHouseholdOwnershipInput builds a faked HouseholdOwnershipTransferInput.
func BuildFakeTransferHouseholdOwnershipInput() *types.HouseholdOwnershipTransferInput {
	return &types.HouseholdOwnershipTransferInput{
		Reason:       fake.Sentence(10),
		CurrentOwner: fake.UUID(),
		NewOwner:     fake.UUID(),
	}
}

// BuildFakeChangeActiveHouseholdInput builds a faked ChangeActiveHouseholdInput.
func BuildFakeChangeActiveHouseholdInput() *types.ChangeActiveHouseholdInput {
	return &types.ChangeActiveHouseholdInput{
		HouseholdID: fake.UUID(),
	}
}

// BuildFakeUserStatusResponse builds a faked UserStatusResponse.
func BuildFakeUserStatusResponse() *types.UserStatusResponse {
	return &types.UserStatusResponse{
		UserID:                   BuildFakeID(),
		AccountStatus:            string(types.GoodStandingUserAccountStatus),
		AccountStatusExplanation: "",
		ActiveHousehold:          BuildFakeID(),
		UserIsAuthenticated:      true,
	}
}

// BuildFakeTokenResponse builds a faked TokenResponse.
func BuildFakeTokenResponse() *tokens.TokenResponse {
	return &tokens.TokenResponse{
		UserID:      BuildFakeID(),
		HouseholdID: BuildFakeID(),
		AccessToken: fake.UUID(),
	}
}

func BuildFakeUserLoginInput() *types.UserLoginInput {
	return &types.UserLoginInput{
		Username:  BuildFakeUser().Username,
		Password:  buildFakePassword(),
		TOTPToken: buildFakeTOTPToken(),
	}
}
