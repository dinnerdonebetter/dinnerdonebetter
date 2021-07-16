package audit

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// ValidPreparationAssignmentKey is the key we use to indicate that an audit log entry is associated with an item.
	ValidPreparationAssignmentKey = "valid_preparation_id"
	// ValidPreparationCreationEvent is the event type used to indicate an item was created.
	ValidPreparationCreationEvent = "valid_preparation_created"
	// ValidPreparationUpdateEvent is the event type used to indicate an item was updated.
	ValidPreparationUpdateEvent = "valid_preparation_updated"
	// ValidPreparationArchiveEvent is the event type used to indicate an item was archived.
	ValidPreparationArchiveEvent = "valid_preparation_archived"
)

// BuildValidPreparationCreationEventEntry builds an entry creation input for when a valid preparation is created.
func BuildValidPreparationCreationEventEntry(validPreparation *types.ValidPreparation, createdByUser uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ValidPreparationCreationEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:            createdByUser,
			ValidPreparationAssignmentKey: validPreparation.ID,
			CreationAssignmentKey:         validPreparation,
		},
	}
}

// BuildValidPreparationUpdateEventEntry builds an entry creation input for when a valid preparation is updated.
func BuildValidPreparationUpdateEventEntry(changedByUser, validPreparationID uint64, changes []*types.FieldChangeSummary) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ValidPreparationUpdateEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:            changedByUser,
			ValidPreparationAssignmentKey: validPreparationID,
			ChangesAssignmentKey:          changes,
		},
	}
}

// BuildValidPreparationArchiveEventEntry builds an entry creation input for when a valid preparation is archived.
func BuildValidPreparationArchiveEventEntry(archivedByUser, validPreparationID uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ValidPreparationArchiveEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:            archivedByUser,
			ValidPreparationAssignmentKey: validPreparationID,
		},
	}
}
