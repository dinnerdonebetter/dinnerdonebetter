package recipestepproducts

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideRecipeStepProductDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeStepProductDataManager(database.BuildMockDatabase())
	})
}

func TestProvideRecipeStepProductDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeStepProductDataServer(buildTestService())
	})
}
