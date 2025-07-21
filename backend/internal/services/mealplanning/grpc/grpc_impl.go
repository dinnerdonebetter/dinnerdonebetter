package grpc

import (
	"context"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/managers"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "mealplanning_service"
)

type (
	ServiceImpl struct {
		mealplanningsvc.UnimplementedMealPlanningServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		sessionContextDataFetcher func(context.Context) (sessions.ContextData, error)
		recipeManager             managers.RecipeManager
		validEnumerationsManager  managers.ValidEnumerationsManager
		mealPlanningManager       managers.MealPlanningManager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	recipeManager managers.RecipeManager,
	validEnumerationsManager managers.ValidEnumerationsManager,
	mealPlanningManager managers.MealPlanningManager,
) *ServiceImpl {
	return &ServiceImpl{
		logger:                   logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                   tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		recipeManager:            recipeManager,
		validEnumerationsManager: validEnumerationsManager,
		mealPlanningManager:      mealPlanningManager,
	}
}

var (
	errUnimplemented = errors.New("unimplemented")
)
