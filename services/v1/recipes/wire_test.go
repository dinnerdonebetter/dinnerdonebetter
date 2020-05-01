package recipes

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideRecipeDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeDataManager(database.BuildMockDatabase())
	})
}

func TestProvideRecipeDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeDataServer(buildTestService())
	})
}
