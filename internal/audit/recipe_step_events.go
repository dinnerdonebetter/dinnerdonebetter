package audit

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// RecipeStepAssignmentKey is the key we use to indicate that an audit log entry is associated with an item.
	RecipeStepAssignmentKey = "recipe_step_id"
	// RecipeStepCreationEvent is the event type used to indicate an item was created.
	RecipeStepCreationEvent = "recipe_step_created"
	// RecipeStepUpdateEvent is the event type used to indicate an item was updated.
	RecipeStepUpdateEvent = "recipe_step_updated"
	// RecipeStepArchiveEvent is the event type used to indicate an item was archived.
	RecipeStepArchiveEvent = "recipe_step_archived"
)

// BuildRecipeStepCreationEventEntry builds an entry creation input for when a recipe step is created.
func BuildRecipeStepCreationEventEntry(recipeStep *types.RecipeStep, createdByUser uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: RecipeStepCreationEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:      createdByUser,
			RecipeStepAssignmentKey: recipeStep.ID,
			CreationAssignmentKey:   recipeStep,
		},
	}
}

// BuildRecipeStepUpdateEventEntry builds an entry creation input for when a recipe step is updated.
func BuildRecipeStepUpdateEventEntry(changedByUser, recipeStepID uint64, changes []*types.FieldChangeSummary) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: RecipeStepUpdateEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:      changedByUser,
			RecipeStepAssignmentKey: recipeStepID,
			ChangesAssignmentKey:    changes,
		},
	}
}

// BuildRecipeStepArchiveEventEntry builds an entry creation input for when a recipe step is archived.
func BuildRecipeStepArchiveEventEntry(archivedByUser, recipeStepID uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: RecipeStepArchiveEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:      archivedByUser,
			RecipeStepAssignmentKey: recipeStepID,
		},
	}
}
