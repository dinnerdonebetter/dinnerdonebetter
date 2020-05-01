package invitations

import (
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
)

func TestProvideInvitationDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideInvitationDataManager(database.BuildMockDatabase())
	})
}

func TestProvideInvitationDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideInvitationDataServer(buildTestService())
	})
}
