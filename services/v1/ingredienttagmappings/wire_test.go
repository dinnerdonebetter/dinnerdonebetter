package ingredienttagmappings

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideIngredientTagMappingDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideIngredientTagMappingDataManager(database.BuildMockDatabase())
	})
}

func TestProvideIngredientTagMappingDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideIngredientTagMappingDataServer(buildTestService())
	})
}
