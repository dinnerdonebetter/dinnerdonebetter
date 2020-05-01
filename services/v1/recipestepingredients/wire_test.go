package recipestepingredients

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideRecipeStepIngredientDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeStepIngredientDataManager(database.BuildMockDatabase())
	})
}

func TestProvideRecipeStepIngredientDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeStepIngredientDataServer(buildTestService())
	})
}
