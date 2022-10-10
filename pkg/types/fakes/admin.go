package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeUserAccountStatusUpdateInput builds a faked UserAccountStatusUpdateInput.
func BuildFakeUserAccountStatusUpdateInput() *types.UserAccountStatusUpdateInput {
	return &types.UserAccountStatusUpdateInput{
		TargetUserID: BuildFakeID(),
		NewStatus:    types.GoodStandingUserAccountStatus,
		Reason:       fake.Sentence(10),
	}
}
