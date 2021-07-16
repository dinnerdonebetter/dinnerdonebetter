package audit

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// RecipeStepProductAssignmentKey is the key we use to indicate that an audit log entry is associated with an item.
	RecipeStepProductAssignmentKey = "recipe_step_product_id"
	// RecipeStepProductCreationEvent is the event type used to indicate an item was created.
	RecipeStepProductCreationEvent = "recipe_step_product_created"
	// RecipeStepProductUpdateEvent is the event type used to indicate an item was updated.
	RecipeStepProductUpdateEvent = "recipe_step_product_updated"
	// RecipeStepProductArchiveEvent is the event type used to indicate an item was archived.
	RecipeStepProductArchiveEvent = "recipe_step_product_archived"
)

// BuildRecipeStepProductCreationEventEntry builds an entry creation input for when a recipe step product is created.
func BuildRecipeStepProductCreationEventEntry(recipeStepProduct *types.RecipeStepProduct, createdByUser uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: RecipeStepProductCreationEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:             createdByUser,
			RecipeStepProductAssignmentKey: recipeStepProduct.ID,
			CreationAssignmentKey:          recipeStepProduct,
		},
	}
}

// BuildRecipeStepProductUpdateEventEntry builds an entry creation input for when a recipe step product is updated.
func BuildRecipeStepProductUpdateEventEntry(changedByUser, recipeStepProductID uint64, changes []*types.FieldChangeSummary) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: RecipeStepProductUpdateEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:             changedByUser,
			RecipeStepProductAssignmentKey: recipeStepProductID,
			ChangesAssignmentKey:           changes,
		},
	}
}

// BuildRecipeStepProductArchiveEventEntry builds an entry creation input for when a recipe step product is archived.
func BuildRecipeStepProductArchiveEventEntry(archivedByUser, recipeStepProductID uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: RecipeStepProductArchiveEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:             archivedByUser,
			RecipeStepProductAssignmentKey: recipeStepProductID,
		},
	}
}
