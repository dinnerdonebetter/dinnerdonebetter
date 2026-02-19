package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	commentsmanager "github.com/dinnerdonebetter/backend/internal/domain/comments/manager"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/managers"
	commentssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/comments"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "comments_service"
)

var _ commentssvc.CommentsServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		commentssvc.UnimplementedCommentsServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		sessionContextDataFetcher func(context.Context) (*sessions.ContextData, error)
		commentsManager           commentsmanager.CommentsDataManager
		mealPlanningManager       managers.MealPlanningManager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	commentsManager commentsmanager.CommentsDataManager,
	mealPlanningManager managers.MealPlanningManager,
) commentssvc.CommentsServiceServer {
	return &serviceImpl{
		logger:                    logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		sessionContextDataFetcher: sessions.FetchContextDataFromContext,
		commentsManager:           commentsManager,
		mealPlanningManager:       mealPlanningManager,
	}
}
