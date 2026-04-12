package uploadedmedia

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit"
	types "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/uploadedmedia"

	"github.com/primandproper/platform/database"
	"github.com/primandproper/platform/observability/logging"
	"github.com/primandproper/platform/observability/tracing"

	"github.com/samber/do/v2"
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
