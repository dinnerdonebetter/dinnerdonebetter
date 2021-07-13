package audit

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// ValidIngredientPreparationAssignmentKey is the key we use to indicate that an audit log entry is associated with an item.
	ValidIngredientPreparationAssignmentKey = "valid_ingredient_preparation_id"
	// ValidIngredientPreparationCreationEvent is the event type used to indicate an item was created.
	ValidIngredientPreparationCreationEvent = "valid_ingredient_preparation_created"
	// ValidIngredientPreparationUpdateEvent is the event type used to indicate an item was updated.
	ValidIngredientPreparationUpdateEvent = "valid_ingredient_preparation_updated"
	// ValidIngredientPreparationArchiveEvent is the event type used to indicate an item was archived.
	ValidIngredientPreparationArchiveEvent = "valid_ingredient_preparation_archived"
)

// BuildValidIngredientPreparationCreationEventEntry builds an entry creation input for when a valid ingredient preparation is created.
func BuildValidIngredientPreparationCreationEventEntry(validIngredientPreparation *types.ValidIngredientPreparation, createdByUser uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ValidIngredientPreparationCreationEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:                      createdByUser,
			ValidIngredientPreparationAssignmentKey: validIngredientPreparation.ID,
			CreationAssignmentKey:                   validIngredientPreparation,
		},
	}
}

// BuildValidIngredientPreparationUpdateEventEntry builds an entry creation input for when a valid ingredient preparation is updated.
func BuildValidIngredientPreparationUpdateEventEntry(changedByUser, validIngredientPreparationID uint64, changes []*types.FieldChangeSummary) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ValidIngredientPreparationUpdateEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:                      changedByUser,
			ValidIngredientPreparationAssignmentKey: validIngredientPreparationID,
			ChangesAssignmentKey:                    changes,
		},
	}
}

// BuildValidIngredientPreparationArchiveEventEntry builds an entry creation input for when a valid ingredient preparation is archived.
func BuildValidIngredientPreparationArchiveEventEntry(archivedByUser, validIngredientPreparationID uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: ValidIngredientPreparationArchiveEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:                      archivedByUser,
			ValidIngredientPreparationAssignmentKey: validIngredientPreparationID,
		},
	}
}
