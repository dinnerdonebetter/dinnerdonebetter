package manager

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/comments"

	"github.com/verygoodsoftwarenotvirus/platform/v5/messagequeue"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v5/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v5/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v5/observability/tracing"

	"github.com/samber/do/v2"
)

// RegisterCommentsDataManager registers the comments data manager with the injector.
func RegisterCommentsDataManager(i do.Injector) {
	do.Provide[CommentsDataManager](i, func(i do.Injector) (CommentsDataManager, error) {
		return NewCommentsDataManager(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[comments.Repository](i),
			do.MustInvoke[*msgconfig.QueuesConfig](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
		)
	})
}
