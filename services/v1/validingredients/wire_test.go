package validingredients

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideValidIngredientDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideValidIngredientDataManager(database.BuildMockDatabase())
	})
}

func TestProvideValidIngredientDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideValidIngredientDataServer(buildTestService())
	})
}
