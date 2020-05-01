package iterationmedias

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideIterationMediaDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideIterationMediaDataManager(database.BuildMockDatabase())
	})
}

func TestProvideIterationMediaDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideIterationMediaDataServer(buildTestService())
	})
}
