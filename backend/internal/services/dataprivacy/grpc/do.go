package grpc

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authentication/sessions"
	dataprivacymanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/dataprivacy/manager"
	dataprivacysvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/dataprivacy"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v3/uploads"
)

// RegisterDataPrivacyService registers the data privacy gRPC service with the injector.
func RegisterDataPrivacyService(i do.Injector) {
	do.Provide[dataprivacysvc.DataPrivacyServiceServer](i, func(i do.Injector) (dataprivacysvc.DataPrivacyServiceServer, error) {
		return NewDataPrivacyService(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[func(context.Context) (*sessions.ContextData, error)](i),
			do.MustInvoke[dataprivacymanager.DataPrivacyManager](i),
			do.MustInvoke[uploads.UploadManager](i),
		), nil
	})
}
