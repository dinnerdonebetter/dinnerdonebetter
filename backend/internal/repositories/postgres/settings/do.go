package settings

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v4/database"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

// RegisterSettingsRepository registers the settings repository with the injector.
func RegisterSettingsRepository(i do.Injector) {
	do.Provide[*Repository](i, func(i do.Injector) (*Repository, error) {
		return ProvideSettingsRepository(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[audit.Repository](i),
			do.MustInvoke[database.Client](i),
		), nil
	})
}
