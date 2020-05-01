package validinstruments

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideValidInstrumentDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideValidInstrumentDataManager(database.BuildMockDatabase())
	})
}

func TestProvideValidInstrumentDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideValidInstrumentDataServer(buildTestService())
	})
}
