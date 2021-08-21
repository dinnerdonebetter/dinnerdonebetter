package authorization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHouseholdRolePermissionChecker(T *testing.T) {
	T.Parallel()

	T.Run("household user", func(t *testing.T) {
		t.Parallel()

		r := NewHouseholdRolePermissionChecker(HouseholdMemberRole.String())

		assert.False(t, r.CanUpdateHouseholds())
		assert.False(t, r.CanDeleteHouseholds())
		assert.False(t, r.CanAddMemberToHouseholds())
		assert.False(t, r.CanRemoveMemberFromHouseholds())
		assert.False(t, r.CanTransferHouseholdToNewOwner())
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

	T.Run("household admin", func(t *testing.T) {
		t.Parallel()

		r := NewHouseholdRolePermissionChecker(HouseholdAdminRole.String())

		assert.True(t, r.CanUpdateHouseholds())
		assert.True(t, r.CanDeleteHouseholds())
		assert.True(t, r.CanAddMemberToHouseholds())
		assert.True(t, r.CanRemoveMemberFromHouseholds())
		assert.True(t, r.CanTransferHouseholdToNewOwner())
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
