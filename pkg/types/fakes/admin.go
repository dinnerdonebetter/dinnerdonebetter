package fakes

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeUserReputationUpdateInput builds a faked UserReputationUpdateInput.
func BuildFakeUserReputationUpdateInput() *types.UserReputationUpdateInput {
	return &types.UserReputationUpdateInput{
		TargetUserID:  uint64(fake.Uint32()),
		NewReputation: types.GoodStandingHouseholdStatus,
		Reason:        fake.Sentence(10),
	}
}
