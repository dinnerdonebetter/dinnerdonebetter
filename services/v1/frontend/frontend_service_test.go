package frontend

import (
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/config"

	"gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
)

func TestProvideFrontendService(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ProvideFrontendService(noop.ProvideNoopLogger(), config.FrontendSettings{})
	})
}
