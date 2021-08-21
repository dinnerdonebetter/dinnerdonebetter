package fakes

import (
	"time"

	"gitlab.com/prixfixe/prixfixe/internal/authorization"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeSessionContextData builds a faked SessionContextData.
func BuildFakeSessionContextData() *types.SessionContextData {
	fakeHouseholdID := fake.Uint64()

	return &types.SessionContextData{
		HouseholdPermissions: map[uint64]authorization.HouseholdRolePermissionsChecker{
			fakeHouseholdID: authorization.NewHouseholdRolePermissionChecker(authorization.HouseholdAdminRole.String()),
		},
		Requester: types.RequesterInfo{
			Reputation:            types.GoodStandingHouseholdStatus,
			ReputationExplanation: "",
			UserID:                fake.Uint64(),
			ServicePermissions:    authorization.NewServiceRolePermissionChecker(authorization.ServiceUserRole.String()),
		},
		ActiveHouseholdID: fakeHouseholdID,
	}
}

// BuildFakeSessionContextDataForHousehold builds a faked SessionContextData.
func BuildFakeSessionContextDataForHousehold(household *types.Household) *types.SessionContextData {
	fakeHouseholdID := fake.Uint64()

	return &types.SessionContextData{
		HouseholdPermissions: map[uint64]authorization.HouseholdRolePermissionsChecker{
			household.ID: authorization.NewHouseholdRolePermissionChecker(authorization.ServiceUserRole.String()),
		},
		Requester: types.RequesterInfo{
			Reputation:            types.GoodStandingHouseholdStatus,
			ReputationExplanation: "",
			UserID:                fake.Uint64(),
			ServicePermissions:    authorization.NewServiceRolePermissionChecker(authorization.ServiceUserRole.String()),
		},
		ActiveHouseholdID: fakeHouseholdID,
	}
}

// BuildFakeAddUserToHouseholdInput builds a faked AddUserToHouseholdInput.
func BuildFakeAddUserToHouseholdInput() *types.AddUserToHouseholdInput {
	return &types.AddUserToHouseholdInput{
		Reason:         fake.Sentence(10),
		UserID:         fake.Uint64(),
		HouseholdID:    fake.Uint64(),
		HouseholdRoles: []string{authorization.HouseholdMemberRole.String()},
	}
}

// BuildFakeUserPermissionModificationInput builds a faked ModifyUserPermissionsInput.
func BuildFakeUserPermissionModificationInput() *types.ModifyUserPermissionsInput {
	return &types.ModifyUserPermissionsInput{
		Reason:   fake.Sentence(10),
		NewRoles: []string{authorization.HouseholdMemberRole.String()},
	}
}

// BuildFakeTransferHouseholdOwnershipInput builds a faked HouseholdOwnershipTransferInput.
func BuildFakeTransferHouseholdOwnershipInput() *types.HouseholdOwnershipTransferInput {
	return &types.HouseholdOwnershipTransferInput{
		Reason:       fake.Sentence(10),
		CurrentOwner: fake.Uint64(),
		NewOwner:     fake.Uint64(),
	}
}

// BuildFakeChangeActiveHouseholdInput builds a faked ChangeActiveHouseholdInput.
func BuildFakeChangeActiveHouseholdInput() *types.ChangeActiveHouseholdInput {
	return &types.ChangeActiveHouseholdInput{
		HouseholdID: fake.Uint64(),
	}
}

// BuildFakePASETOCreationInput builds a faked PASETOCreationInput.
func BuildFakePASETOCreationInput() *types.PASETOCreationInput {
	return &types.PASETOCreationInput{
		ClientID:    fake.UUID(),
		RequestTime: time.Now().Unix(),
	}
}
