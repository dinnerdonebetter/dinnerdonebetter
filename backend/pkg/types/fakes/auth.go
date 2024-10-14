package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/pkg/types"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeSessionContextData builds a faked SessionContextData.
func BuildFakeSessionContextData() *types.SessionContextData {
	return &types.SessionContextData{
		HouseholdPermissions: map[string]authorization.HouseholdRolePermissionsChecker{},
		Requester: types.RequesterInfo{
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

// BuildFakeJWTResponse builds a faked JWTResponse.
func BuildFakeJWTResponse() *types.JWTResponse {
	return &types.JWTResponse{
		UserID:      BuildFakeID(),
		HouseholdID: BuildFakeID(),
		Token:       fake.UUID(),
	}
}

func BuildFakeUserLoginInput() *types.UserLoginInput {
	return &types.UserLoginInput{
		Username:  BuildFakeUser().Username,
		Password:  buildFakePassword(),
		TOTPToken: buildFakeTOTPToken(),
	}
}
