package recipeiterationsteps

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideRecipeIterationStepDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeIterationStepDataManager(database.BuildMockDatabase())
	})
}

func TestProvideRecipeIterationStepDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeIterationStepDataServer(buildTestService())
	})
}
