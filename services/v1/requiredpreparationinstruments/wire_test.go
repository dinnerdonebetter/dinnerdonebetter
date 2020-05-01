package requiredpreparationinstruments

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideRequiredPreparationInstrumentDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRequiredPreparationInstrumentDataManager(database.BuildMockDatabase())
	})
}

func TestProvideRequiredPreparationInstrumentDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideRequiredPreparationInstrumentDataServer(buildTestService())
	})
}
