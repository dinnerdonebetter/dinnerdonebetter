package recipestepevents

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideRecipeStepEventDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeStepEventDataManager(database.BuildMockDatabase())
	})
}

func TestProvideRecipeStepEventDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeStepEventDataServer(buildTestService())
	})
}
