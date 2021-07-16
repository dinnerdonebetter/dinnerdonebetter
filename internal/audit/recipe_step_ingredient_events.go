package audit

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// RecipeStepIngredientAssignmentKey is the key we use to indicate that an audit log entry is associated with an item.
	RecipeStepIngredientAssignmentKey = "recipe_step_ingredient_id"
	// RecipeStepIngredientCreationEvent is the event type used to indicate an item was created.
	RecipeStepIngredientCreationEvent = "recipe_step_ingredient_created"
	// RecipeStepIngredientUpdateEvent is the event type used to indicate an item was updated.
	RecipeStepIngredientUpdateEvent = "recipe_step_ingredient_updated"
	// RecipeStepIngredientArchiveEvent is the event type used to indicate an item was archived.
	RecipeStepIngredientArchiveEvent = "recipe_step_ingredient_archived"
)

// BuildRecipeStepIngredientCreationEventEntry builds an entry creation input for when a recipe step ingredient is created.
func BuildRecipeStepIngredientCreationEventEntry(recipeStepIngredient *types.RecipeStepIngredient, createdByUser uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: RecipeStepIngredientCreationEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:                createdByUser,
			RecipeStepIngredientAssignmentKey: recipeStepIngredient.ID,
			CreationAssignmentKey:             recipeStepIngredient,
		},
	}
}

// BuildRecipeStepIngredientUpdateEventEntry builds an entry creation input for when a recipe step ingredient is updated.
func BuildRecipeStepIngredientUpdateEventEntry(changedByUser, recipeStepIngredientID uint64, changes []*types.FieldChangeSummary) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: RecipeStepIngredientUpdateEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:                changedByUser,
			RecipeStepIngredientAssignmentKey: recipeStepIngredientID,
			ChangesAssignmentKey:              changes,
		},
	}
}

// BuildRecipeStepIngredientArchiveEventEntry builds an entry creation input for when a recipe step ingredient is archived.
func BuildRecipeStepIngredientArchiveEventEntry(archivedByUser, recipeStepIngredientID uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: RecipeStepIngredientArchiveEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:                archivedByUser,
			RecipeStepIngredientAssignmentKey: recipeStepIngredientID,
		},
	}
}
