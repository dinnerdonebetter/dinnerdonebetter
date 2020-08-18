package recipestepinstruments

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideRecipeStepInstrumentDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeStepInstrumentDataManager(database.BuildMockDatabase())
	})
}

func TestProvideRecipeStepInstrumentDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRecipeStepInstrumentDataServer(buildTestService())
	})
}
