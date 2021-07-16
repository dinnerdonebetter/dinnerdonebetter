package audit

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// ValidPreparationInstrumentAssignmentKey is the key we use to indicate that an audit log entry is associated with an item.
	ValidPreparationInstrumentAssignmentKey = "valid_preparation_instrument_id"
	// ValidPreparationInstrumentCreationEvent is the event type used to indicate an item was created.
	ValidPreparationInstrumentCreationEvent = "valid_preparation_instrument_created"
	// ValidPreparationInstrumentUpdateEvent is the event type used to indicate an item was updated.
	ValidPreparationInstrumentUpdateEvent = "valid_preparation_instrument_updated"
	// ValidPreparationInstrumentArchiveEvent is the event type used to indicate an item was archived.
	ValidPreparationInstrumentArchiveEvent = "valid_preparation_instrument_archived"
)

// BuildValidPreparationInstrumentCreationEventEntry builds an entry creation input for when a valid preparation instrument is created.
func BuildValidPreparationInstrumentCreationEventEntry(validPreparationInstrument *types.ValidPreparationInstrument, createdByUser uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ValidPreparationInstrumentCreationEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:                      createdByUser,
			ValidPreparationInstrumentAssignmentKey: validPreparationInstrument.ID,
			CreationAssignmentKey:                   validPreparationInstrument,
		},
	}
}

// BuildValidPreparationInstrumentUpdateEventEntry builds an entry creation input for when a valid preparation instrument is updated.
func BuildValidPreparationInstrumentUpdateEventEntry(changedByUser, validPreparationInstrumentID uint64, changes []*types.FieldChangeSummary) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ValidPreparationInstrumentUpdateEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:                      changedByUser,
			ValidPreparationInstrumentAssignmentKey: validPreparationInstrumentID,
			ChangesAssignmentKey:                    changes,
		},
	}
}

// BuildValidPreparationInstrumentArchiveEventEntry builds an entry creation input for when a valid preparation instrument is archived.
func BuildValidPreparationInstrumentArchiveEventEntry(archivedByUser, validPreparationInstrumentID uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ValidPreparationInstrumentArchiveEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:                      archivedByUser,
			ValidPreparationInstrumentAssignmentKey: validPreparationInstrumentID,
		},
	}
}
