package converters

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning"

	"github.com/verygoodsoftwarenotvirus/platform/v4/identifiers"
)

func ConvertRecipeListItemCreationRequestInputToRecipeListItemDatabaseCreationInput(x *mealplanning.RecipeListItemCreationRequestInput, recipeListID string) *mealplanning.RecipeListItemDatabaseCreationInput {
	return &mealplanning.RecipeListItemDatabaseCreationInput{
		ID:                  identifiers.New(),
		RecipeID:            x.RecipeID,
		Notes:               x.Notes,
		BelongsToRecipeList: recipeListID,
	}
}
