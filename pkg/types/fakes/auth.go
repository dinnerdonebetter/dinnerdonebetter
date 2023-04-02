package fakes

import (
	"time"

	"github.com/prixfixeco/backend/internal/authorization"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeHouseholdUserMembershipCreationRequestInput builds a faked HouseholdUserMembershipCreationRequestInput.
func BuildFakeHouseholdUserMembershipCreationRequestInput() *types.HouseholdUserMembershipCreationRequestInput {
	return &types.HouseholdUserMembershipCreationRequestInput{
		Reason: fake.Sentence(10),
		UserID: BuildFakeID(),
	}
}

// BuildFakeHouseholdUserMembershipDatabaseCreationInput builds a faked HouseholdUserMembershipCreationRequestInput.
func BuildFakeHouseholdUserMembershipDatabaseCreationInput() *types.HouseholdUserMembershipDatabaseCreationInput {
	input := BuildFakeHouseholdUserMembershipCreationRequestInput()

	return converters.ConvertHouseholdUserMembershipCreationRequestInputToHouseholdUserMembershipDatabaseCreationInput(input)
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

// BuildFakePASETOCreationInput builds a faked PASETOCreationInput.
func BuildFakePASETOCreationInput() *types.PASETOCreationInput {
	return &types.PASETOCreationInput{
		ClientID:    BuildFakeID(),
		RequestTime: time.Now().Unix(),
	}
}
