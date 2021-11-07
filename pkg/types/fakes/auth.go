package fakes

import (
	"time"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeSessionContextData builds a faked SessionContextData.
func BuildFakeSessionContextData() *types.SessionContextData {
	fakeHouseholdID := fake.UUID()

	return &types.SessionContextData{
		HouseholdPermissions: map[string]authorization.HouseholdRolePermissionsChecker{
			fakeHouseholdID: authorization.NewHouseholdRolePermissionChecker(authorization.HouseholdAdminRole.String()),
		},
		Requester: types.RequesterInfo{
			Reputation:            types.GoodStandingHouseholdStatus,
			ReputationExplanation: "",
			UserID:                ksuid.New().String(),
			ServicePermissions:    authorization.NewServiceRolePermissionChecker(authorization.ServiceUserRole.String()),
		},
		ActiveHouseholdID: fakeHouseholdID,
	}
}

// BuildFakeSessionContextDataForHousehold builds a faked SessionContextData.
func BuildFakeSessionContextDataForHousehold(household *types.Household) *types.SessionContextData {
	fakeHouseholdID := fake.UUID()

	return &types.SessionContextData{
		HouseholdPermissions: map[string]authorization.HouseholdRolePermissionsChecker{
			household.ID: authorization.NewHouseholdRolePermissionChecker(authorization.ServiceUserRole.String()),
		},
		Requester: types.RequesterInfo{
			Reputation:            types.GoodStandingHouseholdStatus,
			ReputationExplanation: "",
			UserID:                ksuid.New().String(),
			ServicePermissions:    authorization.NewServiceRolePermissionChecker(authorization.ServiceUserRole.String()),
		},
		ActiveHouseholdID: fakeHouseholdID,
	}
}

// BuildFakeHouseholdUserMembershipCreationRequestInput builds a faked HouseholdUserMembershipCreationRequestInput.
func BuildFakeHouseholdUserMembershipCreationRequestInput() *types.HouseholdUserMembershipCreationRequestInput {
	return &types.HouseholdUserMembershipCreationRequestInput{
		Reason:         fake.Sentence(10),
		UserID:         ksuid.New().String(),
		HouseholdID:    ksuid.New().String(),
		HouseholdRoles: []string{authorization.HouseholdMemberRole.String()},
	}
}

// BuildFakeHouseholdUserMembershipDatabaseCreationInput builds a faked HouseholdUserMembershipCreationRequestInput.
func BuildFakeHouseholdUserMembershipDatabaseCreationInput() *types.HouseholdUserMembershipDatabaseCreationInput {
	input := BuildFakeHouseholdUserMembershipCreationRequestInput()

	return BuildFakeHouseholdUserMembershipDatabaseCreationInputFromHouseholdUserMembershipCreationRequestInput(input)
}

// BuildFakeHouseholdUserMembershipDatabaseCreationInputFromHouseholdUserMembershipCreationRequestInput builds a faked HouseholdUserMembershipCreationRequestInput.
func BuildFakeHouseholdUserMembershipDatabaseCreationInputFromHouseholdUserMembershipCreationRequestInput(input *types.HouseholdUserMembershipCreationRequestInput) *types.HouseholdUserMembershipDatabaseCreationInput {
	return &types.HouseholdUserMembershipDatabaseCreationInput{
		ID:             input.ID,
		Reason:         input.Reason,
		UserID:         input.UserID,
		HouseholdID:    input.HouseholdID,
		HouseholdRoles: input.HouseholdRoles,
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

// BuildFakePASETOCreationInput builds a faked PASETOCreationInput.
func BuildFakePASETOCreationInput() *types.PASETOCreationInput {
	return &types.PASETOCreationInput{
		ClientID:    ksuid.New().String(),
		RequestTime: time.Now().Unix(),
	}
}
