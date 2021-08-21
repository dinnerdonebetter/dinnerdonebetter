package fakes

import (
	"gitlab.com/prixfixe/prixfixe/internal/authorization"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeHouseholdUserMembership builds a faked HouseholdUserMembership.
func BuildFakeHouseholdUserMembership() *types.HouseholdUserMembership {
	return &types.HouseholdUserMembership{
		ID:                 uint64(fake.Uint32()),
		BelongsToUser:      fake.Uint64(),
		BelongsToHousehold: fake.Uint64(),
		HouseholdRoles:     []string{authorization.HouseholdMemberRole.String()},
		CreatedOn:          0,
		ArchivedOn:         nil,
	}
}

// BuildFakeHouseholdUserMembershipList builds a faked HouseholdUserMembershipList.
func BuildFakeHouseholdUserMembershipList() *types.HouseholdUserMembershipList {
	var examples []*types.HouseholdUserMembership
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeHouseholdUserMembership())
	}

	return &types.HouseholdUserMembershipList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		HouseholdUserMemberships: examples,
	}
}

// BuildFakeHouseholdUserMembershipUpdateInputFromHouseholdUserMembership builds a faked HouseholdUserMembershipUpdateInput from an household user membership.
func BuildFakeHouseholdUserMembershipUpdateInputFromHouseholdUserMembership(householdUserMembership *types.HouseholdUserMembership) *types.HouseholdUserMembershipUpdateInput {
	return &types.HouseholdUserMembershipUpdateInput{
		BelongsToUser:      householdUserMembership.BelongsToUser,
		BelongsToHousehold: householdUserMembership.BelongsToHousehold,
	}
}

// BuildFakeHouseholdUserMembershipCreationInput builds a faked HouseholdUserMembershipCreationInput.
func BuildFakeHouseholdUserMembershipCreationInput() *types.HouseholdUserMembershipCreationInput {
	return BuildFakeHouseholdUserMembershipCreationInputFromHouseholdUserMembership(BuildFakeHouseholdUserMembership())
}

// BuildFakeHouseholdUserMembershipCreationInputFromHouseholdUserMembership builds a faked HouseholdUserMembershipCreationInput from an household user membership.
func BuildFakeHouseholdUserMembershipCreationInputFromHouseholdUserMembership(householdUserMembership *types.HouseholdUserMembership) *types.HouseholdUserMembershipCreationInput {
	return &types.HouseholdUserMembershipCreationInput{
		BelongsToUser:      householdUserMembership.BelongsToUser,
		BelongsToHousehold: householdUserMembership.BelongsToHousehold,
	}
}
