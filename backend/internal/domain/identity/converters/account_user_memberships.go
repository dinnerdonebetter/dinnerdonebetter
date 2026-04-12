package converters

import (
	types "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/identity"

	"github.com/primandproper/platform/identifiers"
)

// ConvertAccountUserMembershipToAccountUserMembershipDatabaseCreationInput builds a faked AccountUserMembershipCreationRequestInput.
func ConvertAccountUserMembershipToAccountUserMembershipDatabaseCreationInput(account *types.AccountUserMembership) *types.AccountUserMembershipDatabaseCreationInput {
	return &types.AccountUserMembershipDatabaseCreationInput{
		ID:        identifiers.New(),
		Reason:    "",
		UserID:    account.BelongsToUser,
		AccountID: account.ID,
	}
}
