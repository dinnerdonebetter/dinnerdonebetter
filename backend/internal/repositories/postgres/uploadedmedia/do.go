package uploadedmedia

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit"
	types "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/uploadedmedia"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v4/database"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

// RegisterUploadedMediaRepository registers the uploaded media repository with the injector.
func RegisterUploadedMediaRepository(i do.Injector) {
	do.Provide[types.Repository](i, func(i do.Injector) (types.Repository, error) {
		return ProvideUploadedMediaRepository(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[audit.Repository](i),
			do.MustInvoke[database.Client](i),
		), nil
	})
}
