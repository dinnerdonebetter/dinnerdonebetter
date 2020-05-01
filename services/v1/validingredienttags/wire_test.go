package validingredienttags

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideValidIngredientTagDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideValidIngredientTagDataManager(database.BuildMockDatabase())
	})
}

func TestProvideValidIngredientTagDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideValidIngredientTagDataServer(buildTestService())
	})
}
