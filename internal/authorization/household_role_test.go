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
		assert.Equal(t, r.HasPermission(UpdateHouseholdPermission), r.CanUpdateHouseholds())
		assert.False(t, r.CanDeleteHouseholds())
		assert.Equal(t, r.HasPermission(ArchiveHouseholdPermission), r.CanDeleteHouseholds())
		assert.False(t, r.CanAddMemberToHouseholds())
		assert.Equal(t, r.HasPermission(InviteUserToHouseholdPermission), r.CanAddMemberToHouseholds())
		assert.False(t, r.CanRemoveMemberFromHouseholds())
		assert.Equal(t, r.HasPermission(RemoveMemberHouseholdPermission), r.CanRemoveMemberFromHouseholds())
		assert.False(t, r.CanTransferHouseholdToNewOwner())
		assert.Equal(t, r.HasPermission(TransferHouseholdPermission), r.CanTransferHouseholdToNewOwner())
		assert.False(t, r.CanCreateWebhooks())
		assert.Equal(t, r.HasPermission(CreateWebhooksPermission), r.CanCreateWebhooks())
		assert.False(t, r.CanUpdateWebhooks())
		assert.Equal(t, r.HasPermission(UpdateWebhooksPermission), r.CanUpdateWebhooks())
		assert.False(t, r.CanArchiveWebhooks())
		assert.Equal(t, r.HasPermission(ArchiveWebhooksPermission), r.CanArchiveWebhooks())
		assert.False(t, r.CanCreateAPIClients())
		assert.Equal(t, r.HasPermission(CreateAPIClientsPermission), r.CanCreateAPIClients())
		assert.False(t, r.CanSeeAPIClients())
		assert.Equal(t, r.HasPermission(ReadAPIClientsPermission), r.CanSeeAPIClients())
		assert.False(t, r.CanDeleteAPIClients())
		assert.Equal(t, r.HasPermission(ArchiveAPIClientsPermission), r.CanDeleteAPIClients())
	})

	T.Run("household admin", func(t *testing.T) {
		t.Parallel()

		r := NewHouseholdRolePermissionChecker(HouseholdAdminRole.String())

		assert.True(t, r.CanUpdateHouseholds())
		assert.Equal(t, r.HasPermission(UpdateHouseholdPermission), r.CanUpdateHouseholds())
		assert.True(t, r.CanDeleteHouseholds())
		assert.Equal(t, r.HasPermission(ArchiveHouseholdPermission), r.CanDeleteHouseholds())
		assert.True(t, r.CanAddMemberToHouseholds())
		assert.Equal(t, r.HasPermission(InviteUserToHouseholdPermission), r.CanAddMemberToHouseholds())
		assert.True(t, r.CanRemoveMemberFromHouseholds())
		assert.Equal(t, r.HasPermission(RemoveMemberHouseholdPermission), r.CanRemoveMemberFromHouseholds())
		assert.True(t, r.CanTransferHouseholdToNewOwner())
		assert.Equal(t, r.HasPermission(TransferHouseholdPermission), r.CanTransferHouseholdToNewOwner())
		assert.True(t, r.CanCreateWebhooks())
		assert.Equal(t, r.HasPermission(CreateWebhooksPermission), r.CanCreateWebhooks())
		assert.True(t, r.CanSeeWebhooks())
		assert.Equal(t, r.HasPermission(ReadWebhooksPermission), r.CanSeeWebhooks())
		assert.True(t, r.CanUpdateWebhooks())
		assert.Equal(t, r.HasPermission(UpdateWebhooksPermission), r.CanUpdateWebhooks())
		assert.True(t, r.CanArchiveWebhooks())
		assert.Equal(t, r.HasPermission(ArchiveWebhooksPermission), r.CanArchiveWebhooks())
		assert.True(t, r.CanCreateAPIClients())
		assert.Equal(t, r.HasPermission(CreateAPIClientsPermission), r.CanCreateAPIClients())
		assert.True(t, r.CanSeeAPIClients())
		assert.Equal(t, r.HasPermission(ReadAPIClientsPermission), r.CanSeeAPIClients())
		assert.True(t, r.CanDeleteAPIClients())
		assert.Equal(t, r.HasPermission(ArchiveAPIClientsPermission), r.CanDeleteAPIClients())
	})
}
