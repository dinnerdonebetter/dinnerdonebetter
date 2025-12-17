package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
)

func ConvertRecipeListItemCreationRequestInputToRecipeListItemDatabaseCreationInput(x *mealplanning.RecipeListItemCreationRequestInput, recipeListID string) *mealplanning.RecipeListItemDatabaseCreationInput {
	return &mealplanning.RecipeListItemDatabaseCreationInput{
		ID:                  identifiers.New(),
		RecipeID:            x.RecipeID,
		Notes:               x.Notes,
		BelongsToRecipeList: recipeListID,
	}
}
