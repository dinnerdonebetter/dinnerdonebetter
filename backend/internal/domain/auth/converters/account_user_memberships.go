package converters

import (
	types "github.com/dinnerdonebetter/backend/internal/domain/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
)

// ConvertAccountUserMembershipToAccountUserMembershipDatabaseCreationInput builds a faked AccountUserMembershipCreationRequestInput.
func ConvertAccountUserMembershipToAccountUserMembershipDatabaseCreationInput(account *types.AccountUserMembership) *types.AccountUserMembershipDatabaseCreationInput {
	return &types.AccountUserMembershipDatabaseCreationInput{
		ID:          identifiers.New(),
		Reason:      "",
		UserID:      account.BelongsToUser,
		AccountID:   account.ID,
		AccountRole: account.AccountRole,
	}
}
