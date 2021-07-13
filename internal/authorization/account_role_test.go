package authorization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccountRolePermissionChecker(T *testing.T) {
	T.Parallel()

	T.Run("account user", func(t *testing.T) {
		t.Parallel()

		r := NewAccountRolePermissionChecker(AccountMemberRole.String())

		assert.False(t, r.CanUpdateAccounts())
		assert.False(t, r.CanDeleteAccounts())
		assert.False(t, r.CanAddMemberToAccounts())
		assert.False(t, r.CanRemoveMemberFromAccounts())
		assert.False(t, r.CanTransferAccountToNewOwner())
		assert.False(t, r.CanCreateWebhooks())
		assert.False(t, r.CanSeeWebhooks())
		assert.False(t, r.CanUpdateWebhooks())
		assert.False(t, r.CanArchiveWebhooks())
		assert.False(t, r.CanCreateAPIClients())
		assert.False(t, r.CanSeeAPIClients())
		assert.False(t, r.CanDeleteAPIClients())
		assert.False(t, r.CanSeeAuditLogEntriesForWebhooks())
		assert.False(t, r.CanSeeAuditLogEntriesForValidInstruments())
		assert.False(t, r.CanSeeAuditLogEntriesForValidPreparations())
		assert.False(t, r.CanSeeAuditLogEntriesForValidIngredients())
		assert.False(t, r.CanSeeAuditLogEntriesForValidIngredientPreparations())
		assert.False(t, r.CanSeeAuditLogEntriesForValidPreparationInstruments())
		assert.False(t, r.CanSeeAuditLogEntriesForRecipes())
		assert.False(t, r.CanSeeAuditLogEntriesForRecipeSteps())
		assert.False(t, r.CanSeeAuditLogEntriesForRecipeStepIngredients())
		assert.False(t, r.CanSeeAuditLogEntriesForRecipeStepProducts())
		assert.False(t, r.CanSeeAuditLogEntriesForInvitations())
		assert.False(t, r.CanSeeAuditLogEntriesForReports())
	})

	T.Run("account admin", func(t *testing.T) {
		t.Parallel()

		r := NewAccountRolePermissionChecker(AccountAdminRole.String())

		assert.True(t, r.CanUpdateAccounts())
		assert.True(t, r.CanDeleteAccounts())
		assert.True(t, r.CanAddMemberToAccounts())
		assert.True(t, r.CanRemoveMemberFromAccounts())
		assert.True(t, r.CanTransferAccountToNewOwner())
		assert.True(t, r.CanCreateWebhooks())
		assert.True(t, r.CanSeeWebhooks())
		assert.True(t, r.CanUpdateWebhooks())
		assert.True(t, r.CanArchiveWebhooks())
		assert.True(t, r.CanCreateAPIClients())
		assert.True(t, r.CanSeeAPIClients())
		assert.True(t, r.CanDeleteAPIClients())
		assert.True(t, r.CanSeeAuditLogEntriesForWebhooks())
		assert.True(t, r.CanSeeAuditLogEntriesForValidInstruments())
		assert.True(t, r.CanSeeAuditLogEntriesForValidPreparations())
		assert.True(t, r.CanSeeAuditLogEntriesForValidIngredients())
		assert.True(t, r.CanSeeAuditLogEntriesForValidIngredientPreparations())
		assert.True(t, r.CanSeeAuditLogEntriesForValidPreparationInstruments())
		assert.True(t, r.CanSeeAuditLogEntriesForRecipes())
		assert.True(t, r.CanSeeAuditLogEntriesForRecipeSteps())
		assert.True(t, r.CanSeeAuditLogEntriesForRecipeStepIngredients())
		assert.True(t, r.CanSeeAuditLogEntriesForRecipeStepProducts())
		assert.True(t, r.CanSeeAuditLogEntriesForInvitations())
		assert.True(t, r.CanSeeAuditLogEntriesForReports())
	})
}
