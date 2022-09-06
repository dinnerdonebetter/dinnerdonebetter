package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeHouseholdUserMembership builds a faked HouseholdUserMembership.
func BuildFakeHouseholdUserMembership() *types.HouseholdUserMembership {
	return &types.HouseholdUserMembership{
		ID:                 ksuid.New().String(),
		BelongsToUser:      fake.UUID(),
		BelongsToHousehold: fake.UUID(),
		HouseholdRoles:     []string{authorization.HouseholdMemberRole.String()},
		CreatedAt:          fake.Date(),
		ArchivedAt:         nil,
	}
}

// BuildFakeHouseholdUserMembershipWithUser builds a faked HouseholdUserMembershipWithUser.
func BuildFakeHouseholdUserMembershipWithUser() *types.HouseholdUserMembershipWithUser {
	u := BuildFakeUser()
	u.TwoFactorSecret = ""

	return &types.HouseholdUserMembershipWithUser{
		ID:                 ksuid.New().String(),
		BelongsToUser:      u,
		BelongsToHousehold: fake.UUID(),
		HouseholdRoles:     []string{authorization.HouseholdMemberRole.String()},
		CreatedAt:          fake.Date(),
		ArchivedAt:         nil,
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
