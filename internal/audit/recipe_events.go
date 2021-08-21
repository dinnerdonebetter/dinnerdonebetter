package audit

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// RecipeAssignmentKey is the key we use to indicate that an audit log entry is associated with an item.
	RecipeAssignmentKey = "recipe_id"
	// RecipeCreationEvent is the event type used to indicate an item was created.
	RecipeCreationEvent = "recipe_created"
	// RecipeUpdateEvent is the event type used to indicate an item was updated.
	RecipeUpdateEvent = "recipe_updated"
	// RecipeArchiveEvent is the event type used to indicate an item was archived.
	RecipeArchiveEvent = "recipe_archived"
)

// BuildRecipeCreationEventEntry builds an entry creation input for when a recipe is created.
func BuildRecipeCreationEventEntry(recipe *types.Recipe, createdByUser uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: RecipeCreationEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:     createdByUser,
			RecipeAssignmentKey:    recipe.ID,
			CreationAssignmentKey:  recipe,
			HouseholdAssignmentKey: recipe.BelongsToHousehold,
		},
	}
}

// BuildRecipeUpdateEventEntry builds an entry creation input for when a recipe is updated.
func BuildRecipeUpdateEventEntry(changedByUser, recipeID, householdID uint64, changes []*types.FieldChangeSummary) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: RecipeUpdateEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:     changedByUser,
			HouseholdAssignmentKey: householdID,
			RecipeAssignmentKey:    recipeID,
			ChangesAssignmentKey:   changes,
		},
	}
}

// BuildRecipeArchiveEventEntry builds an entry creation input for when a recipe is archived.
func BuildRecipeArchiveEventEntry(archivedByUser, householdID, recipeID uint64) *types.AuditLogEntryCreationInput {
	return &types.AuditLogEntryCreationInput{
		EventType: RecipeArchiveEvent,
		Context: map[string]interface{}{
			ActorAssignmentKey:     archivedByUser,
			HouseholdAssignmentKey: householdID,
			RecipeAssignmentKey:    recipeID,
		},
	}
}
