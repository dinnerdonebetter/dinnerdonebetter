package audit

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// ValidIngredientAssignmentKey is the key we use to indicate that an audit log entry is associated with an item.
	ValidIngredientAssignmentKey = "valid_ingredient_id"
	// ValidIngredientCreationEvent is the event type used to indicate an item was created.
	ValidIngredientCreationEvent = "valid_ingredient_created"
	// ValidIngredientUpdateEvent is the event type used to indicate an item was updated.
	ValidIngredientUpdateEvent = "valid_ingredient_updated"
	// ValidIngredientArchiveEvent is the event type used to indicate an item was archived.
	ValidIngredientArchiveEvent = "valid_ingredient_archived"
)

// BuildValidIngredientCreationEventEntry builds an entry creation input for when a valid ingredient is created.
func BuildValidIngredientCreationEventEntry(validIngredient *types.ValidIngredient, createdByUser uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ValidIngredientCreationEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:           createdByUser,
			ValidIngredientAssignmentKey: validIngredient.ID,
			CreationAssignmentKey:        validIngredient,
		},
	}
}

// BuildValidIngredientUpdateEventEntry builds an entry creation input for when a valid ingredient is updated.
func BuildValidIngredientUpdateEventEntry(changedByUser, validIngredientID uint64, changes []*types.FieldChangeSummary) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ValidIngredientUpdateEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:           changedByUser,
			ValidIngredientAssignmentKey: validIngredientID,
			ChangesAssignmentKey:         changes,
		},
	}
}

// BuildValidIngredientArchiveEventEntry builds an entry creation input for when a valid ingredient is archived.
func BuildValidIngredientArchiveEventEntry(archivedByUser, validIngredientID uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ValidIngredientArchiveEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:           archivedByUser,
			ValidIngredientAssignmentKey: validIngredientID,
		},
	}
}
