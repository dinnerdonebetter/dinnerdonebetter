package recipetags

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideRecipeTagDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeTagDataManager(database.BuildMockDatabase())
	})
}

func TestProvideRecipeTagDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeTagDataServer(buildTestService())
	})
}
