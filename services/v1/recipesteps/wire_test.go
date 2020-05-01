package recipesteps

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideRecipeStepDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeStepDataManager(database.BuildMockDatabase())
	})
}

func TestProvideRecipeStepDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeStepDataServer(buildTestService())
	})
}
