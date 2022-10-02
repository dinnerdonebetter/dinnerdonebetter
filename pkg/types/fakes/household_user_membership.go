package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/pkg/types"
)

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
