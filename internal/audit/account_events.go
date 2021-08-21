package audit

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// HouseholdAssignmentKey is the key we use to indicate that an audit log entry is associated with an household.
	HouseholdAssignmentKey = "household_id"
	// HouseholdCreationEvent events indicate a user created an household.
	HouseholdCreationEvent = "household_created"
	// HouseholdUpdateEvent events indicate a user updated an household.
	HouseholdUpdateEvent = "household_updated"
	// HouseholdArchiveEvent events indicate a user deleted an household.
	HouseholdArchiveEvent = "household_archived"
)

// BuildHouseholdCreationEventEntry builds an entry creation input for when an household is created.
func BuildHouseholdCreationEventEntry(household *types.Household, createdByUser uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: HouseholdCreationEvent,
		Context: map[string]interface{}{
			HouseholdAssignmentKey: household.ID,
			UserAssignmentKey:      household.BelongsToUser,
			ActorAssignmentKey:     createdByUser,
			CreationAssignmentKey:  household,
		},
	}
}

// BuildHouseholdUpdateEventEntry builds an entry creation input for when an household is updated.
func BuildHouseholdUpdateEventEntry(userID, householdID, changedByUser uint64, changes []*types.FieldChangeSummary) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: HouseholdUpdateEvent,
		Context: map[string]interface{}{
			UserAssignmentKey:      userID,
			HouseholdAssignmentKey: householdID,
			ChangesAssignmentKey:   changes,
			ActorAssignmentKey:     changedByUser,
		},
	}
}

// BuildHouseholdArchiveEventEntry builds an entry creation input for when an household is archived.
func BuildHouseholdArchiveEventEntry(userID, householdID, archivedByUser uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: HouseholdArchiveEvent,
		Context: map[string]interface{}{
			UserAssignmentKey:      userID,
			HouseholdAssignmentKey: householdID,
			ActorAssignmentKey:     archivedByUser,
		},
	}
}
