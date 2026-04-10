package comments

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit"
	domaincomments "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/comments"

	"github.com/verygoodsoftwarenotvirus/platform/v5/database"
	"github.com/verygoodsoftwarenotvirus/platform/v5/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v5/observability/tracing"

	"github.com/samber/do/v2"
)

// RegisterCommentsRepository registers the comments repository with the injector.
func RegisterCommentsRepository(i do.Injector) {
	do.Provide[domaincomments.Repository](i, func(i do.Injector) (domaincomments.Repository, error) {
		return ProvideCommentsRepository(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[audit.Repository](i),
			do.MustInvoke[database.Client](i),
		), nil
	})
}
