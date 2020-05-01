package recipeiterations

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideRecipeIterationDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeIterationDataManager(database.BuildMockDatabase())
	})
}

func TestProvideRecipeIterationDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeIterationDataServer(buildTestService())
	})
}
