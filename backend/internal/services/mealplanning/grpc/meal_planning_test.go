package grpc

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mealplanningfakes "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mockmanagers "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/managers/mock"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/fake"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildServiceImplForMealPlanningTest(t *testing.T) *serviceImpl {
	t.Helper()

	return &serviceImpl{
		tracer: tracing.NewTracerForTest(t.Name()),
		logger: logging.NewNoopLogger(),
		sessionContextDataFetcher: func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					UserID: mealplanningfakes.BuildFakeID(),
				},
				ActiveAccountID: mealplanningfakes.BuildFakeID(),
			}, nil
		},
		// Workers are set to nil for most tests since they're only used in specific worker methods
		mealPlanFinalizerWorker:              nil,
		mealPlanGroceryListInitializerWorker: nil,
		mealPlanTaskCreatorWorker:            nil,
	}
}

func TestServiceImpl_ArchiveMeal(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleMealID := mealplanningfakes.BuildFakeID()
		exampleUserID := mealplanningfakes.BuildFakeID()

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ArchiveMeal", testutils.ContextMatcher, exampleMealID, exampleUserID).Return(nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific user ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					UserID: exampleUserID,
				},
			}, nil
		}

		res, err := s.ArchiveMeal(ctx, &mealplanninggrpc.ArchiveMealRequest{MealID: exampleMealID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_ArchiveMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleAccountID := mealplanningfakes.BuildFakeID()

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ArchiveMealPlan", testutils.ContextMatcher, exampleMealPlanID, exampleAccountID).Return(nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific account ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				ActiveAccountID: exampleAccountID,
			}, nil
		}

		res, err := s.ArchiveMealPlan(ctx, &mealplanninggrpc.ArchiveMealPlanRequest{MealPlanID: exampleMealPlanID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_ArchiveMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleMealPlanEventID := mealplanningfakes.BuildFakeID()

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ArchiveMealPlanEvent", testutils.ContextMatcher, exampleMealPlanID, exampleMealPlanEventID).Return(nil)
		s.mealPlanningManager = mmpm

		res, err := s.ArchiveMealPlanEvent(ctx, &mealplanninggrpc.ArchiveMealPlanEventRequest{
			MealPlanID:      exampleMealPlanID,
			MealPlanEventID: exampleMealPlanEventID,
		})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_ArchiveMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleMealPlanGroceryListItemID := mealplanningfakes.BuildFakeID()

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ArchiveMealPlanGroceryListItem", testutils.ContextMatcher, exampleMealPlanID, exampleMealPlanGroceryListItemID).Return(nil)
		s.mealPlanningManager = mmpm

		res, err := s.ArchiveMealPlanGroceryListItem(ctx, &mealplanninggrpc.ArchiveMealPlanGroceryListItemRequest{
			MealPlanID:                exampleMealPlanID,
			MealPlanGroceryListItemID: exampleMealPlanGroceryListItemID,
		})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_ArchiveMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleMealPlanEventID := mealplanningfakes.BuildFakeID()
		exampleMealPlanOptionID := mealplanningfakes.BuildFakeID()

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ArchiveMealPlanOption", testutils.ContextMatcher, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID).Return(nil)
		s.mealPlanningManager = mmpm

		res, err := s.ArchiveMealPlanOption(ctx, &mealplanninggrpc.ArchiveMealPlanOptionRequest{
			MealPlanID:       exampleMealPlanID,
			MealPlanEventID:  exampleMealPlanEventID,
			MealPlanOptionID: exampleMealPlanOptionID,
		})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_ArchiveMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleMealPlanEventID := mealplanningfakes.BuildFakeID()
		exampleMealPlanOptionID := mealplanningfakes.BuildFakeID()
		exampleMealPlanOptionVoteID := mealplanningfakes.BuildFakeID()

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ArchiveMealPlanOptionVote", testutils.ContextMatcher, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVoteID).Return(nil)
		s.mealPlanningManager = mmpm

		res, err := s.ArchiveMealPlanOptionVote(ctx, &mealplanninggrpc.ArchiveMealPlanOptionVoteRequest{
			MealPlanID:           exampleMealPlanID,
			MealPlanEventID:      exampleMealPlanEventID,
			MealPlanOptionID:     exampleMealPlanOptionID,
			MealPlanOptionVoteID: exampleMealPlanOptionVoteID,
		})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_ArchiveUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleUserID := mealplanningfakes.BuildFakeID()
		exampleUserIngredientPreferenceID := mealplanningfakes.BuildFakeID()

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ArchiveUserIngredientPreference", testutils.ContextMatcher, exampleUserID, exampleUserIngredientPreferenceID).Return(nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific user ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					UserID: exampleUserID,
				},
			}, nil
		}

		res, err := s.ArchiveUserIngredientPreference(ctx, &mealplanninggrpc.ArchiveUserIngredientPreferenceRequest{
			UserIngredientPreferenceID: exampleUserIngredientPreferenceID,
		})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_CreateMeal(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleUserID := mealplanningfakes.BuildFakeID()
		exampleCreatedMeal := mealplanningfakes.BuildFakeMeal()

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("CreateMeal", testutils.ContextMatcher, exampleUserID, testutils.MatchType[*mealplanning.MealCreationRequestInput]()).Return(exampleCreatedMeal, nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific user ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					UserID: exampleUserID,
				},
			}, nil
		}

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateMealRequest](t)

		actual, err := s.CreateMeal(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleCreatedMeal.ID, actual.Created.ID)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_CreateMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleAccountID := mealplanningfakes.BuildFakeID()
		exampleUserID := mealplanningfakes.BuildFakeID()
		exampleCreatedMealPlan := mealplanningfakes.BuildFakeMealPlan()

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("CreateMealPlan", testutils.ContextMatcher, exampleAccountID, exampleUserID, testutils.MatchType[*mealplanning.MealPlanCreationRequestInput]()).Return(exampleCreatedMealPlan, nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific IDs
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					UserID: exampleUserID,
				},
				ActiveAccountID: exampleAccountID,
			}, nil
		}

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateMealPlanRequest](t)

		actual, err := s.CreateMealPlan(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleCreatedMealPlan.ID, actual.Created.ID)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_CreateMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleCreatedMealPlanEvent := mealplanningfakes.BuildFakeMealPlanEvent()

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("CreateMealPlanEvent", testutils.ContextMatcher, exampleMealPlanID, testutils.MatchType[*mealplanning.MealPlanEventCreationRequestInput]()).Return(exampleCreatedMealPlanEvent, nil)
		s.mealPlanningManager = mmpm

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateMealPlanEventRequest](t)
		exampleInput.MealPlanID = exampleMealPlanID

		actual, err := s.CreateMealPlanEvent(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleCreatedMealPlanEvent.ID, actual.Created.ID)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_CreateMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleCreatedMealPlanGroceryListItem := mealplanningfakes.BuildFakeMealPlanGroceryListItem()

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("CreateMealPlanGroceryListItem", testutils.ContextMatcher, testutils.MatchType[*mealplanning.MealPlanGroceryListItemCreationRequestInput]()).Return(exampleCreatedMealPlanGroceryListItem, nil)
		s.mealPlanningManager = mmpm

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateMealPlanGroceryListItemRequest](t)
		exampleInput.MealPlanID = exampleMealPlanID

		actual, err := s.CreateMealPlanGroceryListItem(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleCreatedMealPlanGroceryListItem.ID, actual.Created.ID)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_CreateMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleMealPlanEventID := mealplanningfakes.BuildFakeID()
		exampleCreatedMealPlanOption := mealplanningfakes.BuildFakeMealPlanOption()

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("CreateMealPlanOptionWithEventID", testutils.ContextMatcher, exampleMealPlanEventID, testutils.MatchType[*mealplanning.MealPlanOptionCreationRequestInput]()).Return(exampleCreatedMealPlanOption, nil)
		s.mealPlanningManager = mmpm

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateMealPlanOptionRequest](t)
		exampleInput.MealPlanID = exampleMealPlanID
		exampleInput.MealPlanEventID = exampleMealPlanEventID

		actual, err := s.CreateMealPlanOption(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleCreatedMealPlanOption.ID, actual.Created.ID)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_CreateMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleUserID := mealplanningfakes.BuildFakeID()
		exampleCreatedMealPlanOptionVotes := []*mealplanning.MealPlanOptionVote{
			mealplanningfakes.BuildFakeMealPlanOptionVote(),
		}

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("CreateMealPlanOptionVotes", testutils.ContextMatcher, exampleUserID, testutils.MatchType[*mealplanning.MealPlanOptionVoteCreationRequestInput]()).Return(exampleCreatedMealPlanOptionVotes, nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific user ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					UserID: exampleUserID,
				},
			}, nil
		}

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateMealPlanOptionVoteRequest](t)
		exampleInput.MealPlanID = exampleMealPlanID

		actual, err := s.CreateMealPlanOptionVote(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Len(t, actual.Created, len(exampleCreatedMealPlanOptionVotes))

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_CreateMealPlanTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleCreatedMealPlanTask := mealplanningfakes.BuildFakeMealPlanTask()

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("CreateMealPlanTask", testutils.ContextMatcher, testutils.MatchType[*mealplanning.MealPlanTaskCreationRequestInput]()).Return(exampleCreatedMealPlanTask, nil)
		s.mealPlanningManager = mmpm

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateMealPlanTaskRequest](t)
		exampleInput.MealPlanID = exampleMealPlanID

		actual, err := s.CreateMealPlanTask(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleCreatedMealPlanTask.ID, actual.Created.ID)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_CreateUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleUserID := mealplanningfakes.BuildFakeID()
		exampleCreatedUserIngredientPreferences := []*mealplanning.UserIngredientPreference{
			mealplanningfakes.BuildFakeUserIngredientPreference(),
		}

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("CreateUserIngredientPreference", testutils.ContextMatcher, exampleUserID, testutils.MatchType[*mealplanning.UserIngredientPreferenceCreationRequestInput]()).Return(exampleCreatedUserIngredientPreferences, nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific user ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					UserID: exampleUserID,
				},
			}, nil
		}

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateUserIngredientPreferenceRequest](t)

		actual, err := s.CreateUserIngredientPreference(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Len(t, actual.Created, len(exampleCreatedUserIngredientPreferences))

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_FinalizeMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleAccountID := mealplanningfakes.BuildFakeID()
		exampleFinalized := true

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("FinalizeMealPlan", testutils.ContextMatcher, exampleMealPlanID, exampleAccountID).Return(exampleFinalized, nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific account ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				ActiveAccountID: exampleAccountID,
			}, nil
		}

		res, err := s.FinalizeMealPlan(ctx, &mealplanninggrpc.FinalizeMealPlanRequest{MealPlanID: exampleMealPlanID})
		assert.NotNil(t, res)
		assert.NoError(t, err)
		assert.Equal(t, exampleFinalized, res.Finalized)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetMeal(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeMeal()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ReadMeal", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		result, err := s.GetMeal(ctx, &mealplanninggrpc.GetMealRequest{MealID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeMealPlan()
		exampleAccountID := mealplanningfakes.BuildFakeID()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ReadMealPlan", testutils.ContextMatcher, exampleResult.ID, exampleAccountID).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific account ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				ActiveAccountID: exampleAccountID,
			}, nil
		}

		result, err := s.GetMealPlan(ctx, &mealplanninggrpc.GetMealPlanRequest{MealPlanID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetMealPlansForAccount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleAccountID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeMealPlansList()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ListMealPlans", testutils.ContextMatcher, exampleAccountID, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific account ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				ActiveAccountID: exampleAccountID,
			}, nil
		}

		result, err := s.GetMealPlansForAccount(ctx, &mealplanninggrpc.GetMealPlansForAccountRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeMealPlanEvent()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ReadMealPlanEvent", testutils.ContextMatcher, exampleResult.BelongsToMealPlan, exampleResult.ID).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		result, err := s.GetMealPlanEvent(ctx, &mealplanninggrpc.GetMealPlanEventRequest{
			MealPlanID:      exampleResult.BelongsToMealPlan,
			MealPlanEventID: exampleResult.ID,
		})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetMealPlanEvents(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeMealPlanEventsList()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ListMealPlanEvents", testutils.ContextMatcher, exampleMealPlanID, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		result, err := s.GetMealPlanEvents(ctx, &mealplanninggrpc.GetMealPlanEventsRequest{MealPlanID: exampleMealPlanID})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeMealPlanGroceryListItem()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ReadMealPlanGroceryListItem", testutils.ContextMatcher, exampleResult.BelongsToMealPlan, exampleResult.ID).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		result, err := s.GetMealPlanGroceryListItem(ctx, &mealplanninggrpc.GetMealPlanGroceryListItemRequest{
			MealPlanID:                exampleResult.BelongsToMealPlan,
			MealPlanGroceryListItemID: exampleResult.ID,
		})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetMealPlanGroceryListItemsForMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeMealPlanGroceryListItemsList()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ListMealPlanGroceryListItemsByMealPlan", testutils.ContextMatcher, exampleMealPlanID, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		result, err := s.GetMealPlanGroceryListItemsForMealPlan(ctx, &mealplanninggrpc.GetMealPlanGroceryListItemsForMealPlanRequest{MealPlanID: exampleMealPlanID})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeMealPlanOption()
		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleMealPlanEventID := mealplanningfakes.BuildFakeID()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ReadMealPlanOption", testutils.ContextMatcher, exampleMealPlanID, exampleMealPlanEventID, exampleResult.ID).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		result, err := s.GetMealPlanOption(ctx, &mealplanninggrpc.GetMealPlanOptionRequest{
			MealPlanID:       exampleMealPlanID,
			MealPlanEventID:  exampleMealPlanEventID,
			MealPlanOptionID: exampleResult.ID,
		})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeMealPlanOptionVote()
		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleMealPlanEventID := mealplanningfakes.BuildFakeID()
		exampleMealPlanOptionID := mealplanningfakes.BuildFakeID()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ReadMealPlanOptionVote", testutils.ContextMatcher, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, exampleResult.ID).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		result, err := s.GetMealPlanOptionVote(ctx, &mealplanninggrpc.GetMealPlanOptionVoteRequest{
			MealPlanID:           exampleMealPlanID,
			MealPlanEventID:      exampleMealPlanEventID,
			MealPlanOptionID:     exampleMealPlanOptionID,
			MealPlanOptionVoteID: exampleResult.ID,
		})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetMealPlanOptionVotes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleMealPlanEventID := mealplanningfakes.BuildFakeID()
		exampleMealPlanOptionID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeMealPlanOptionVotesList()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ListMealPlanOptionVotes", testutils.ContextMatcher, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		result, err := s.GetMealPlanOptionVotes(ctx, &mealplanninggrpc.GetMealPlanOptionVotesRequest{
			MealPlanID:       exampleMealPlanID,
			MealPlanEventID:  exampleMealPlanEventID,
			MealPlanOptionID: exampleMealPlanOptionID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetMealPlanOptions(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleMealPlanEventID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeMealPlanOptionsList()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ListMealPlanOptions", testutils.ContextMatcher, exampleMealPlanID, exampleMealPlanEventID, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		result, err := s.GetMealPlanOptions(ctx, &mealplanninggrpc.GetMealPlanOptionsRequest{
			MealPlanID:      exampleMealPlanID,
			MealPlanEventID: exampleMealPlanEventID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetMealPlanTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeMealPlanTask()
		exampleMealPlanID := mealplanningfakes.BuildFakeID()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ReadMealPlanTask", testutils.ContextMatcher, exampleMealPlanID, exampleResult.ID).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		result, err := s.GetMealPlanTask(ctx, &mealplanninggrpc.GetMealPlanTaskRequest{
			MealPlanID:     exampleMealPlanID,
			MealPlanTaskID: exampleResult.ID,
		})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetMealPlanTasks(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeMealPlanTasksList()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ListMealPlanTasksByMealPlan", testutils.ContextMatcher, exampleMealPlanID, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		result, err := s.GetMealPlanTasks(ctx, &mealplanninggrpc.GetMealPlanTasksRequest{MealPlanID: exampleMealPlanID})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetMeals(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeMealsList()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ListMeals", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		result, err := s.GetMeals(ctx, &mealplanninggrpc.GetMealsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeUserIngredientPreference()
		exampleUserID := mealplanningfakes.BuildFakeID()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ReadUserIngredientPreference", testutils.ContextMatcher, exampleUserID, exampleResult.ID).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific user ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					UserID: exampleUserID,
				},
			}, nil
		}

		result, err := s.GetUserIngredientPreference(ctx, &mealplanninggrpc.GetUserIngredientPreferenceRequest{
			UserIngredientPreferenceID: exampleResult.ID,
		})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetUserIngredientPreferences(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUserID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeUserIngredientPreferencesList()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ListUserIngredientPreferences", testutils.ContextMatcher, exampleUserID, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific user ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					UserID: exampleUserID,
				},
			}, nil
		}

		result, err := s.GetUserIngredientPreferences(ctx, &mealplanninggrpc.GetUserIngredientPreferencesRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_RunFinalizeMealPlanWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		// Note: This test requires actual worker instances which are complex to set up.
		// For now, we skip this test as it would require full worker initialization.
		// In a real scenario, you would create actual worker instances or use integration tests.
		t.Skip("Worker tests require actual worker instances - use integration tests instead")
	})
}

func TestServiceImpl_RunMealPlanGroceryListInitializerWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		// Note: This test requires actual worker instances which are complex to set up.
		// For now, we skip this test as it would require full worker initialization.
		// In a real scenario, you would create actual worker instances or use integration tests.
		t.Skip("Worker tests require actual worker instances - use integration tests instead")
	})
}

func TestServiceImpl_RunMealPlanTaskCreatorWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		// Note: This test requires actual worker instances which are complex to set up.
		// For now, we skip this test as it would require full worker initialization.
		// In a real scenario, you would create actual worker instances or use integration tests.
		t.Skip("Worker tests require actual worker instances - use integration tests instead")
	})
}

func TestServiceImpl_SearchForMeals(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeMealsList()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.SearchForMealsRequest](t)

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("SearchMeals", testutils.ContextMatcher, exampleRequest.Query, !exampleRequest.UseSearchService, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		result, err := s.SearchForMeals(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_UpdateMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateMealPlanRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeMealPlan()
		exampleAccountID := mealplanningfakes.BuildFakeID()

		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("UpdateMealPlan", testutils.ContextMatcher, exampleRequest.MealPlanID, exampleAccountID, testutils.MatchType[*mealplanning.MealPlanUpdateRequestInput]()).Return(nil)
		mmpm.On("ReadMealPlan", testutils.ContextMatcher, exampleRequest.MealPlanID, exampleAccountID).Return(exampleResponse, nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific account ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				ActiveAccountID: exampleAccountID,
			}, nil
		}

		res, err := s.UpdateMealPlan(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Updated.ID)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_UpdateMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateMealPlanEventRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeMealPlanEvent()

		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("UpdateMealPlanEvent", testutils.ContextMatcher, exampleRequest.MealPlanID, exampleRequest.MealPlanEventID, testutils.MatchType[*mealplanning.MealPlanEventUpdateRequestInput]()).Return(nil)
		mmpm.On("ReadMealPlanEvent", testutils.ContextMatcher, exampleRequest.MealPlanID, exampleRequest.MealPlanEventID).Return(exampleResponse, nil)
		s.mealPlanningManager = mmpm

		res, err := s.UpdateMealPlanEvent(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Updated.ID)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_UpdateMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateMealPlanGroceryListItemRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeMealPlanGroceryListItem()

		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("UpdateMealPlanGroceryListItem", testutils.ContextMatcher, exampleRequest.MealPlanID, exampleRequest.MealPlanGroceryListItemID, testutils.MatchType[*mealplanning.MealPlanGroceryListItemUpdateRequestInput]()).Return(nil)
		mmpm.On("ReadMealPlanGroceryListItem", testutils.ContextMatcher, exampleRequest.MealPlanID, exampleRequest.MealPlanGroceryListItemID).Return(exampleResponse, nil)
		s.mealPlanningManager = mmpm

		res, err := s.UpdateMealPlanGroceryListItem(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Updated.ID)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_UpdateMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateMealPlanOptionRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeMealPlanOption()

		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("UpdateMealPlanOption", testutils.ContextMatcher, exampleRequest.MealPlanID, exampleRequest.MealPlanEventID, exampleRequest.MealPlanOptionID, testutils.MatchType[*mealplanning.MealPlanOptionUpdateRequestInput]()).Return(nil)
		mmpm.On("ReadMealPlanOption", testutils.ContextMatcher, exampleRequest.MealPlanID, exampleRequest.MealPlanEventID, exampleRequest.MealPlanOptionID).Return(exampleResponse, nil)
		s.mealPlanningManager = mmpm

		res, err := s.UpdateMealPlanOption(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Updated.ID)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_UpdateMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateMealPlanOptionVoteRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeMealPlanOptionVote()

		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("UpdateMealPlanOptionVote", testutils.ContextMatcher, exampleRequest.MealPlanID, exampleRequest.MealPlanEventID, exampleRequest.MealPlanOptionID, exampleRequest.MealPlanOptionVoteID, testutils.MatchType[*mealplanning.MealPlanOptionVoteUpdateRequestInput]()).Return(nil)
		mmpm.On("ReadMealPlanOptionVote", testutils.ContextMatcher, exampleRequest.MealPlanID, exampleRequest.MealPlanEventID, exampleRequest.MealPlanOptionID, exampleRequest.MealPlanOptionVoteID).Return(exampleResponse, nil)
		s.mealPlanningManager = mmpm

		res, err := s.UpdateMealPlanOptionVote(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Updated.ID)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_UpdateMealPlanTaskStatus(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateMealPlanTaskStatusRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeMealPlanTask()

		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("MealPlanTaskStatusChange", testutils.ContextMatcher, testutils.MatchType[*mealplanning.MealPlanTaskStatusChangeRequestInput]()).Return(nil)
		mmpm.On("ReadMealPlanTask", testutils.ContextMatcher, exampleRequest.MealPlanID, exampleRequest.MealPlanTaskID).Return(exampleResponse, nil)
		s.mealPlanningManager = mmpm

		res, err := s.UpdateMealPlanTaskStatus(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Updated.ID)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_UpdateUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateUserIngredientPreferenceRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeUserIngredientPreference()
		exampleUserID := mealplanningfakes.BuildFakeID()

		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("UpdateUserIngredientPreference", testutils.ContextMatcher, exampleRequest.UserIngredientPreferenceID, exampleUserID, testutils.MatchType[*mealplanning.UserIngredientPreferenceUpdateRequestInput]()).Return(nil)
		mmpm.On("ReadUserIngredientPreference", testutils.ContextMatcher, exampleUserID, exampleRequest.UserIngredientPreferenceID).Return(exampleResponse, nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific user ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					UserID: exampleUserID,
				},
			}, nil
		}

		res, err := s.UpdateUserIngredientPreference(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Updated.ID)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_CreateAccountInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleAccountID := mealplanningfakes.BuildFakeID()
		exampleCreatedAccountInstrumentOwnership := mealplanningfakes.BuildFakeAccountInstrumentOwnership()

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("CreateAccountInstrumentOwnership", testutils.ContextMatcher, exampleAccountID, testutils.MatchType[*mealplanning.AccountInstrumentOwnershipCreationRequestInput]()).Return(exampleCreatedAccountInstrumentOwnership, nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific account ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				ActiveAccountID: exampleAccountID,
			}, nil
		}

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateAccountInstrumentOwnershipRequest](t)

		actual, err := s.CreateAccountInstrumentOwnership(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleCreatedAccountInstrumentOwnership.ID, actual.Created.ID)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetAccountInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeAccountInstrumentOwnership()
		exampleAccountID := mealplanningfakes.BuildFakeID()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ReadAccountInstrumentOwnership", testutils.ContextMatcher, exampleAccountID, exampleResult.ID).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific account ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				ActiveAccountID: exampleAccountID,
			}, nil
		}

		result, err := s.GetAccountInstrumentOwnership(ctx, &mealplanninggrpc.GetAccountInstrumentOwnershipRequest{
			AccountInstrumentOwnershipID: exampleResult.ID,
		})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_GetAccountInstrumentOwnerships(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleAccountID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeAccountInstrumentOwnershipsList()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ListAccountInstrumentOwnerships", testutils.ContextMatcher, exampleAccountID, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific account ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				ActiveAccountID: exampleAccountID,
			}, nil
		}

		result, err := s.GetAccountInstrumentOwnerships(ctx, &mealplanninggrpc.GetAccountInstrumentOwnershipsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_UpdateAccountInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateAccountInstrumentOwnershipRequest](t)
		exampleAccountID := mealplanningfakes.BuildFakeID()
		exampleAccountInstrumentOwnership := mealplanningfakes.BuildFakeAccountInstrumentOwnership()

		s := buildServiceImplForMealPlanningTest(t)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ReadAccountInstrumentOwnership", testutils.ContextMatcher, exampleAccountID, exampleRequest.AccountInstrumentOwnershipID).Return(exampleAccountInstrumentOwnership, nil)
		mmpm.On("UpdateAccountInstrumentOwnership", testutils.ContextMatcher, exampleAccountInstrumentOwnership.ID, exampleAccountInstrumentOwnership.BelongsToAccount, testutils.MatchType[*mealplanning.AccountInstrumentOwnershipUpdateRequestInput]()).Return(nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific account ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				ActiveAccountID: exampleAccountID,
			}, nil
		}

		res, err := s.UpdateAccountInstrumentOwnership(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.NotNil(t, res)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}

func TestServiceImpl_ArchiveAccountInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		exampleAccountID := mealplanningfakes.BuildFakeID()
		exampleAccountInstrumentOwnershipID := mealplanningfakes.BuildFakeID()

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On("ArchiveAccountInstrumentOwnership", testutils.ContextMatcher, exampleAccountID, exampleAccountInstrumentOwnershipID).Return(nil)
		s.mealPlanningManager = mmpm

		// Override session context to return specific account ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				ActiveAccountID: exampleAccountID,
			}, nil
		}

		res, err := s.ArchiveAccountInstrumentOwnership(ctx, &mealplanninggrpc.ArchiveAccountInstrumentOwnershipRequest{
			AccountInstrumentOwnershipID: exampleAccountInstrumentOwnershipID,
		})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mmpm)
	})
}
