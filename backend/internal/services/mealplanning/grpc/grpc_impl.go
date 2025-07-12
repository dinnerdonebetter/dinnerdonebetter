package eatinggrpc

import (
	"context"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/managers"
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	eatingsvc "github.com/dinnerdonebetter/backend/internal/grpc/service/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

type (
	ServiceImpl struct {
		eatingsvc.UnimplementedEatingServer
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

func (s *ServiceImpl) PublishArbitraryQueueMessage(ctx context.Context, request *messages.PublishArbitraryQueueMessageRequest) (*messages.PublishArbitraryQueueMessageResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) DestroyAllUserData(ctx context.Context, request *messages.DestroyAllUserDataRequest) (*messages.DestroyAllUserDataResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) FetchUserDataReport(ctx context.Context, request *messages.FetchUserDataReportRequest) (*messages.FetchUserDataReportResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) AdminUpdateUserStatus(ctx context.Context, request *messages.AdminUpdateUserStatusRequest) (*messages.AdminUpdateUserStatusResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) AggregateUserDataReport(ctx context.Context, request *messages.AggregateUserDataReportRequest) (*messages.AggregateUserDataReportResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveUserIngredientPreference(ctx context.Context, request *messages.ArchiveUserIngredientPreferenceRequest) (*messages.ArchiveUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateUserIngredientPreference(ctx context.Context, request *messages.CreateUserIngredientPreferenceRequest) (*messages.CreateUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetUserIngredientPreferences(ctx context.Context, request *messages.GetUserIngredientPreferencesRequest) (*messages.GetUserIngredientPreferencesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateUserIngredientPreference(ctx context.Context, request *messages.UpdateUserIngredientPreferenceRequest) (*messages.UpdateUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}
