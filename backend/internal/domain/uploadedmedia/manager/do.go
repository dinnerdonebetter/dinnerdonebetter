package manager

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/uploadedmedia"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

// RegisterUploadedMediaManager registers the uploaded media manager with the injector.
func RegisterUploadedMediaManager(i do.Injector) {
	do.Provide[UploadedMediaManager](i, func(i do.Injector) (UploadedMediaManager, error) {
		return NewUploadedMediaDataManager(
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[uploadedmedia.Repository](i),
		), nil
	})
}
