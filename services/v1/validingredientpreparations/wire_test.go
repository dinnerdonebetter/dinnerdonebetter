package validingredientpreparations

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideValidIngredientPreparationDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideValidIngredientPreparationDataManager(database.BuildMockDatabase())
	})
}

func TestProvideValidIngredientPreparationDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideValidIngredientPreparationDataServer(buildTestService())
	})
}
