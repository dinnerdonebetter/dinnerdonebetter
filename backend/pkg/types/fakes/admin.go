package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/pkg/types"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeUserAccountStatusUpdateInput builds a faked UserAccountStatusUpdateInput.
func BuildFakeUserAccountStatusUpdateInput() *types.UserAccountStatusUpdateInput {
	return &types.UserAccountStatusUpdateInput{
		TargetUserID: BuildFakeID(),
		NewStatus:    string(types.GoodStandingUserAccountStatus),
		Reason:       fake.Sentence(10),
	}
}

func BuildFakeModifyUserPermissionsInput() *types.ModifyUserPermissionsInput {
	return &types.ModifyUserPermissionsInput{
		Reason:  fake.Sentence(10),
		NewRole: authorization.HouseholdMemberRoleName,
	}
}
