package audit

import "gitlab.com/prixfixe/prixfixe/pkg/types"

// BuildUserBanEventEntry builds an entry creation input for when a user is banned.
func BuildUserBanEventEntry(banGiver, banRecipient uint64, reason string) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: UserBannedEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey: banGiver,
			UserAssignmentKey:  banRecipient,
			ReasonKey:          reason,
		},
	}
}

// BuildHouseholdTerminationEventEntry builds an entry creation input for when an household is terminated.
func BuildHouseholdTerminationEventEntry(terminator, terminee uint64, reason string) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: HouseholdTerminatedEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey: terminator,
			UserAssignmentKey:  terminee,
			ReasonKey:          reason,
		},
	}
}
