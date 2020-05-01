package recipesteppreparations

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideRecipeStepPreparationDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeStepPreparationDataManager(database.BuildMockDatabase())
	})
}

func TestProvideRecipeStepPreparationDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeStepPreparationDataServer(buildTestService())
	})
}
