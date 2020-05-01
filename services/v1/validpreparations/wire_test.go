package validpreparations

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideValidPreparationDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideValidPreparationDataManager(database.BuildMockDatabase())
	})
}

func TestProvideValidPreparationDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideValidPreparationDataServer(buildTestService())
	})
}
