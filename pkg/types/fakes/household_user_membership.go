package fakes

import (
	"github.com/prixfixeco/backend/internal/authorization"
	"github.com/prixfixeco/backend/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeHouseholdUserMembershipWithUser builds a faked HouseholdUserMembershipWithUser.
func BuildFakeHouseholdUserMembershipWithUser() *types.HouseholdUserMembershipWithUser {
	u := BuildFakeUser()
	u.TwoFactorSecret = ""

	return &types.HouseholdUserMembershipWithUser{
		ID:                 BuildFakeID(),
		BelongsToUser:      u,
		BelongsToHousehold: fake.UUID(),
		HouseholdRole:      authorization.HouseholdMemberRole.String(),
		CreatedAt:          BuildFakeTime(),
		ArchivedAt:         nil,
	}
}
