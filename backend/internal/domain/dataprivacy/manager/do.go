package manager

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/dataprivacy"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

// RegisterDataPrivacyManager registers the data privacy manager with the injector.
func RegisterDataPrivacyManager(i do.Injector) {
	do.Provide[DataPrivacyManager](i, func(i do.Injector) (DataPrivacyManager, error) {
		return NewDataPrivacyManager(
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[dataprivacy.Repository](i),
		), nil
	})
}
