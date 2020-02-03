package auth_test

import (
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/auth"
)

func TestProvideBcryptHashCost(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		auth.ProvideBcryptHashCost()
	})
}
