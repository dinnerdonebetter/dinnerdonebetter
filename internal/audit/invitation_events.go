package audit

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// InvitationAssignmentKey is the key we use to indicate that an audit log entry is associated with an item.
	InvitationAssignmentKey = "invitation_id"
	// InvitationCreationEvent is the event type used to indicate an item was created.
	InvitationCreationEvent = "invitation_created"
	// InvitationUpdateEvent is the event type used to indicate an item was updated.
	InvitationUpdateEvent = "invitation_updated"
	// InvitationArchiveEvent is the event type used to indicate an item was archived.
	InvitationArchiveEvent = "invitation_archived"
)

// BuildInvitationCreationEventEntry builds an entry creation input for when an invitation is created.
func BuildInvitationCreationEventEntry(invitation *types.Invitation, createdByUser uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: InvitationCreationEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:      createdByUser,
			InvitationAssignmentKey: invitation.ID,
			CreationAssignmentKey:   invitation,
			HouseholdAssignmentKey:  invitation.BelongsToHousehold,
		},
	}
}

// BuildInvitationUpdateEventEntry builds an entry creation input for when an invitation is updated.
func BuildInvitationUpdateEventEntry(changedByUser, invitationID, householdID uint64, changes []*types.FieldChangeSummary) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: InvitationUpdateEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:      changedByUser,
			HouseholdAssignmentKey:  householdID,
			InvitationAssignmentKey: invitationID,
			ChangesAssignmentKey:    changes,
		},
	}
}

// BuildInvitationArchiveEventEntry builds an entry creation input for when an invitation is archived.
func BuildInvitationArchiveEventEntry(archivedByUser, householdID, invitationID uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: InvitationArchiveEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:      archivedByUser,
			HouseholdAssignmentKey:  householdID,
			InvitationAssignmentKey: invitationID,
		},
	}
}
