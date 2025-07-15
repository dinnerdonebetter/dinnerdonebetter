package grpc

import (
	"context"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/managers"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

type (
	ServiceImpl struct {
		mealplanningsvc.UnimplementedMealPlanningServiceServer
		tracer                   tracing.Tracer
		logger                   logging.Logger
		recipeManager            managers.RecipeManager
		validEnumerationsManager managers.ValidEnumerationsManager
		mealPlanningManager      managers.MealPlanningManager
	}
)

func NewService(
	recipeManager managers.RecipeManager,
	validEnumerationsManager managers.ValidEnumerationsManager,
	mealPlanningManager managers.MealPlanningManager,
) *ServiceImpl {
	return &ServiceImpl{
		recipeManager:            recipeManager,
		validEnumerationsManager: validEnumerationsManager,
		mealPlanningManager:      mealPlanningManager,
	}
}

var (
	errUnimplemented = errors.New("unimplemented")
)

func (s *ServiceImpl) ArchiveUserIngredientPreference(ctx context.Context, request *mealplanningsvc.ArchiveUserIngredientPreferenceRequest) (*mealplanningsvc.ArchiveUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateUserIngredientPreference(ctx context.Context, request *mealplanningsvc.CreateUserIngredientPreferenceRequest) (*mealplanningsvc.CreateUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetUserIngredientPreferences(ctx context.Context, request *mealplanningsvc.GetUserIngredientPreferencesRequest) (*mealplanningsvc.GetUserIngredientPreferencesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateUserIngredientPreference(ctx context.Context, request *mealplanningsvc.UpdateUserIngredientPreferenceRequest) (*mealplanningsvc.UpdateUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}
