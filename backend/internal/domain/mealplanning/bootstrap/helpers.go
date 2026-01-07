package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
)

// getRecipeIDBySlug returns the recipe ID for the given slug, or nil if not found.
func getRecipeIDBySlug(createdRecipes map[string]*mealplanning.Recipe, slug string) *string {
	if recipe, ok := createdRecipes[slug]; ok && recipe != nil {
		return &recipe.ID
	}
	return nil
}
