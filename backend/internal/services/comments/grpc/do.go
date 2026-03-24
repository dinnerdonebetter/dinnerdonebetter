package grpc

import (
	commentsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/comments/manager"
	commentssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/comments"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/tracing"
)

// RegisterCommentsService registers the comments gRPC service with the injector.
func RegisterCommentsService(i do.Injector) {
	do.Provide[CommentsMethodPermissions](i, func(i do.Injector) (CommentsMethodPermissions, error) {
		return ProvideMethodPermissions(), nil
	})

	do.Provide[commentssvc.CommentsServiceServer](i, func(i do.Injector) (commentssvc.CommentsServiceServer, error) {
		return NewService(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[commentsmanager.CommentsDataManager](i),
		), nil
	})
}
