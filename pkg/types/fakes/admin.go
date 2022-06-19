package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeUserAccountStatusUpdateInput builds a faked UserAccountStatusUpdateInput.
func BuildFakeUserAccountStatusUpdateInput() *types.UserAccountStatusUpdateInput {
	return &types.UserAccountStatusUpdateInput{
		TargetUserID: ksuid.New().String(),
		NewStatus:    types.GoodStandingUserAccountStatus,
		Reason:       fake.Sentence(10),
	}
}
