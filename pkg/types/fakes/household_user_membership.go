package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/backend/internal/authorization"
	"github.com/prixfixeco/backend/pkg/types"
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
		CreatedAt:          fake.Date(),
		ArchivedAt:         nil,
	}
}
