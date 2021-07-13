package audit

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// ValidInstrumentAssignmentKey is the key we use to indicate that an audit log entry is associated with an item.
	ValidInstrumentAssignmentKey = "valid_instrument_id"
	// ValidInstrumentCreationEvent is the event type used to indicate an item was created.
	ValidInstrumentCreationEvent = "valid_instrument_created"
	// ValidInstrumentUpdateEvent is the event type used to indicate an item was updated.
	ValidInstrumentUpdateEvent = "valid_instrument_updated"
	// ValidInstrumentArchiveEvent is the event type used to indicate an item was archived.
	ValidInstrumentArchiveEvent = "valid_instrument_archived"
)

// BuildValidInstrumentCreationEventEntry builds an entry creation input for when a valid instrument is created.
func BuildValidInstrumentCreationEventEntry(validInstrument *types.ValidInstrument, createdByUser uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ValidInstrumentCreationEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:           createdByUser,
			ValidInstrumentAssignmentKey: validInstrument.ID,
			CreationAssignmentKey:        validInstrument,
		},
	}
}

// BuildValidInstrumentUpdateEventEntry builds an entry creation input for when a valid instrument is updated.
func BuildValidInstrumentUpdateEventEntry(changedByUser, validInstrumentID uint64, changes []*types.FieldChangeSummary) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ValidInstrumentUpdateEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:           changedByUser,
			ValidInstrumentAssignmentKey: validInstrumentID,
			ChangesAssignmentKey:         changes,
		},
	}
}

// BuildValidInstrumentArchiveEventEntry builds an entry creation input for when a valid instrument is archived.
func BuildValidInstrumentArchiveEventEntry(archivedByUser, validInstrumentID uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ValidInstrumentArchiveEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:           archivedByUser,
			ValidInstrumentAssignmentKey: validInstrumentID,
		},
	}
}
