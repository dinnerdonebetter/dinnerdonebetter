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

func buildServiceImplForRecipesTest(t *testing.T) *serviceImpl {
	t.Helper()

	return &serviceImpl{
		tracer: tracing.NewTracerForTest(t.Name()),
		logger: logging.NewNoopLogger(),
		sessionContextDataFetcher: func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					UserID: mealplanningfakes.BuildFakeID(),
				},
			}, nil
		},
	}
}

func TestServiceImpl_ArchiveRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleUserID := mealplanningfakes.BuildFakeID()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ArchiveRecipe", testutils.ContextMatcher, exampleRecipeID, exampleUserID).Return(nil)
		s.recipeManager = mrm

		// Override session context to return specific user ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					UserID: exampleUserID,
				},
			}, nil
		}

		res, err := s.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeID: exampleRecipeID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_ArchiveRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipePrepTaskID := mealplanningfakes.BuildFakeID()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ArchiveRecipePrepTask", testutils.ContextMatcher, exampleRecipeID, exampleRecipePrepTaskID).Return(nil)
		s.recipeManager = mrm

		res, err := s.ArchiveRecipePrepTask(ctx, &mealplanninggrpc.ArchiveRecipePrepTaskRequest{
			RecipeID:         exampleRecipeID,
			RecipePrepTaskID: exampleRecipePrepTaskID,
		})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_ArchiveRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeRatingID := mealplanningfakes.BuildFakeID()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ArchiveRecipeRating", testutils.ContextMatcher, exampleRecipeID, exampleRecipeRatingID).Return(nil)
		s.recipeManager = mrm

		res, err := s.ArchiveRecipeRating(ctx, &mealplanninggrpc.ArchiveRecipeRatingRequest{
			RecipeID:       exampleRecipeID,
			RecipeRatingID: exampleRecipeRatingID,
		})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_ArchiveRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepID := mealplanningfakes.BuildFakeID()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ArchiveRecipeStep", testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID).Return(nil)
		s.recipeManager = mrm

		res, err := s.ArchiveRecipeStep(ctx, &mealplanninggrpc.ArchiveRecipeStepRequest{
			RecipeID:     exampleRecipeID,
			RecipeStepID: exampleRecipeStepID,
		})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_ArchiveRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepCompletionConditionID := mealplanningfakes.BuildFakeID()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ArchiveRecipeStepCompletionCondition", testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCompletionConditionID).Return(nil)
		s.recipeManager = mrm

		res, err := s.ArchiveRecipeStepCompletionCondition(ctx, &mealplanninggrpc.ArchiveRecipeStepCompletionConditionRequest{
			RecipeID:                        exampleRecipeID,
			RecipeStepID:                    exampleRecipeStepID,
			RecipeStepCompletionConditionID: exampleRecipeStepCompletionConditionID,
		})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_ArchiveRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepIngredientID := mealplanningfakes.BuildFakeID()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ArchiveRecipeStepIngredient", testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredientID).Return(nil)
		s.recipeManager = mrm

		res, err := s.ArchiveRecipeStepIngredient(ctx, &mealplanninggrpc.ArchiveRecipeStepIngredientRequest{
			RecipeID:               exampleRecipeID,
			RecipeStepID:           exampleRecipeStepID,
			RecipeStepIngredientID: exampleRecipeStepIngredientID,
		})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_ArchiveRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepInstrumentID := mealplanningfakes.BuildFakeID()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ArchiveRecipeStepInstrument", testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepInstrumentID).Return(nil)
		s.recipeManager = mrm

		res, err := s.ArchiveRecipeStepInstrument(ctx, &mealplanninggrpc.ArchiveRecipeStepInstrumentRequest{
			RecipeID:               exampleRecipeID,
			RecipeStepID:           exampleRecipeStepID,
			RecipeStepInstrumentID: exampleRecipeStepInstrumentID,
		})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_ArchiveRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepProductID := mealplanningfakes.BuildFakeID()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ArchiveRecipeStepProduct", testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProductID).Return(nil)
		s.recipeManager = mrm

		res, err := s.ArchiveRecipeStepProduct(ctx, &mealplanninggrpc.ArchiveRecipeStepProductRequest{
			RecipeID:            exampleRecipeID,
			RecipeStepID:        exampleRecipeStepID,
			RecipeStepProductID: exampleRecipeStepProductID,
		})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_ArchiveRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepVesselID := mealplanningfakes.BuildFakeID()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ArchiveRecipeStepVessel", testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepVesselID).Return(nil)
		s.recipeManager = mrm

		res, err := s.ArchiveRecipeStepVessel(ctx, &mealplanninggrpc.ArchiveRecipeStepVesselRequest{
			RecipeID:           exampleRecipeID,
			RecipeStepID:       exampleRecipeStepID,
			RecipeStepVesselID: exampleRecipeStepVesselID,
		})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_CloneRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleUserID := mealplanningfakes.BuildFakeID()
		exampleClonedRecipe := mealplanningfakes.BuildFakeRecipe()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("CloneRecipe", testutils.ContextMatcher, exampleRecipeID, exampleUserID).Return(exampleClonedRecipe, nil)
		s.recipeManager = mrm

		// Override session context to return specific user ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					UserID: exampleUserID,
				},
			}, nil
		}

		res, err := s.CloneRecipe(ctx, &mealplanninggrpc.CloneRecipeRequest{RecipeID: exampleRecipeID})
		assert.NotNil(t, res)
		assert.NoError(t, err)
		assert.Equal(t, exampleClonedRecipe.ID, res.Cloned.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_CreateRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleUserID := mealplanningfakes.BuildFakeID()
		exampleCreatedRecipe := mealplanningfakes.BuildFakeRecipe()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("CreateRecipe", testutils.ContextMatcher, exampleUserID, testutils.MatchType[*mealplanning.RecipeCreationRequestInput]()).Return(exampleCreatedRecipe, nil)
		s.recipeManager = mrm

		// Override session context to return specific user ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					UserID: exampleUserID,
				},
			}, nil
		}

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateRecipeRequest](t)

		actual, err := s.CreateRecipe(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleCreatedRecipe.ID, actual.Created.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_CreateRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleCreatedRecipePrepTask := mealplanningfakes.BuildFakeRecipePrepTask()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("CreateRecipePrepTask", testutils.ContextMatcher, exampleRecipeID, testutils.MatchType[*mealplanning.RecipePrepTaskCreationRequestInput]()).Return(exampleCreatedRecipePrepTask, nil)
		s.recipeManager = mrm

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateRecipePrepTaskRequest](t)
		exampleInput.RecipeID = exampleRecipeID

		actual, err := s.CreateRecipePrepTask(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleCreatedRecipePrepTask.ID, actual.Created.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_CreateRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleUserID := mealplanningfakes.BuildFakeID()
		exampleCreatedRecipeRating := mealplanningfakes.BuildFakeRecipeRating()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("CreateRecipeRating", testutils.ContextMatcher, exampleRecipeID, testutils.MatchType[*mealplanning.RecipeRatingCreationRequestInput]()).Return(exampleCreatedRecipeRating, nil)
		s.recipeManager = mrm

		// Override session context to return specific user ID
		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					UserID: exampleUserID,
				},
			}, nil
		}

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateRecipeRatingRequest](t)
		exampleInput.RecipeID = exampleRecipeID

		actual, err := s.CreateRecipeRating(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleCreatedRecipeRating.ID, actual.Created.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_CreateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleCreatedRecipeStep := mealplanningfakes.BuildFakeRecipeStep()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("CreateRecipeStep", testutils.ContextMatcher, exampleRecipeID, testutils.MatchType[*mealplanning.RecipeStepCreationRequestInput]()).Return(exampleCreatedRecipeStep, nil)
		s.recipeManager = mrm

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateRecipeStepRequest](t)
		exampleInput.RecipeID = exampleRecipeID

		actual, err := s.CreateRecipeStep(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleCreatedRecipeStep.ID, actual.Created.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_CreateRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepID := mealplanningfakes.BuildFakeID()
		exampleCreatedRecipeStepCompletionCondition := mealplanningfakes.BuildFakeRecipeStepCompletionCondition()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("CreateRecipeStepCompletionCondition", testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.MatchType[*mealplanning.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput]()).Return(exampleCreatedRecipeStepCompletionCondition, nil)
		s.recipeManager = mrm

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateRecipeStepCompletionConditionRequest](t)
		exampleInput.RecipeID = exampleRecipeID
		exampleInput.RecipeStepID = exampleRecipeStepID

		actual, err := s.CreateRecipeStepCompletionCondition(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleCreatedRecipeStepCompletionCondition.ID, actual.Created.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_CreateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepID := mealplanningfakes.BuildFakeID()
		exampleCreatedRecipeStepIngredient := mealplanningfakes.BuildFakeRecipeStepIngredient()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("CreateRecipeStepIngredient", testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.MatchType[*mealplanning.RecipeStepIngredientCreationRequestInput]()).Return(exampleCreatedRecipeStepIngredient, nil)
		s.recipeManager = mrm

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateRecipeStepIngredientRequest](t)
		exampleInput.RecipeID = exampleRecipeID
		exampleInput.RecipeStepID = exampleRecipeStepID

		actual, err := s.CreateRecipeStepIngredient(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleCreatedRecipeStepIngredient.ID, actual.Created.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_CreateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepID := mealplanningfakes.BuildFakeID()
		exampleCreatedRecipeStepInstrument := mealplanningfakes.BuildFakeRecipeStepInstrument()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("CreateRecipeStepInstrument", testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.MatchType[*mealplanning.RecipeStepInstrumentCreationRequestInput]()).Return(exampleCreatedRecipeStepInstrument, nil)
		s.recipeManager = mrm

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateRecipeStepInstrumentRequest](t)
		exampleInput.RecipeID = exampleRecipeID
		exampleInput.RecipeStepID = exampleRecipeStepID

		actual, err := s.CreateRecipeStepInstrument(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleCreatedRecipeStepInstrument.ID, actual.Created.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_CreateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepID := mealplanningfakes.BuildFakeID()
		exampleCreatedRecipeStepProduct := mealplanningfakes.BuildFakeRecipeStepProduct()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("CreateRecipeStepProduct", testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.MatchType[*mealplanning.RecipeStepProductCreationRequestInput]()).Return(exampleCreatedRecipeStepProduct, nil)
		s.recipeManager = mrm

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateRecipeStepProductRequest](t)
		exampleInput.RecipeID = exampleRecipeID
		exampleInput.RecipeStepID = exampleRecipeStepID

		actual, err := s.CreateRecipeStepProduct(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleCreatedRecipeStepProduct.ID, actual.Created.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_CreateRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepID := mealplanningfakes.BuildFakeID()
		exampleCreatedRecipeStepVessel := mealplanningfakes.BuildFakeRecipeStepVessel()

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("CreateRecipeStepVessel", testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.MatchType[*mealplanning.RecipeStepVesselCreationRequestInput]()).Return(exampleCreatedRecipeStepVessel, nil)
		s.recipeManager = mrm

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateRecipeStepVesselRequest](t)
		exampleInput.RecipeID = exampleRecipeID
		exampleInput.RecipeStepID = exampleRecipeStepID

		actual, err := s.CreateRecipeStepVessel(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleCreatedRecipeStepVessel.ID, actual.Created.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetMermaidDiagramForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleMermaidDiagram := "graph TD\nA[Recipe]"

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("RecipeMermaid", testutils.ContextMatcher, exampleRecipeID).Return(exampleMermaidDiagram, nil)
		s.recipeManager = mrm

		result, err := s.GetMermaidDiagramForRecipe(ctx, &mealplanninggrpc.GetMermaidDiagramForRecipeRequest{RecipeID: exampleRecipeID})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, exampleMermaidDiagram, result.Response)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeRecipe()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ReadRecipe", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipe(ctx, &mealplanninggrpc.GetRecipeRequest{RecipeID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_EstimateRecipePrepTasks(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleEstimatedPrepSteps := []*mealplanning.MealPlanTaskDatabaseCreationEstimate{
			{
				CreationExplanation: "test explanation",
			},
		}

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("RecipeEstimatedPrepSteps", testutils.ContextMatcher, exampleRecipeID).Return(exampleEstimatedPrepSteps, nil)
		s.recipeManager = mrm

		result, err := s.EstimateRecipePrepTasks(ctx, &mealplanninggrpc.EstimateRecipePrepTasksRequest{RecipeID: exampleRecipeID})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleEstimatedPrepSteps))

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeRecipePrepTask()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ReadRecipePrepTask", testutils.ContextMatcher, exampleResult.BelongsToRecipe, exampleResult.ID).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipePrepTask(ctx, &mealplanninggrpc.GetRecipePrepTaskRequest{
			RecipeID:         exampleResult.BelongsToRecipe,
			RecipePrepTaskID: exampleResult.ID,
		})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipePrepTasks(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeRecipePrepTasksList()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ListRecipePrepTask", testutils.ContextMatcher, exampleRecipeID, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipePrepTasks(ctx, &mealplanninggrpc.GetRecipePrepTasksRequest{RecipeID: exampleRecipeID})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeRecipeRating()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ReadRecipeRating", testutils.ContextMatcher, exampleResult.RecipeID, exampleResult.ID).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipeRating(ctx, &mealplanninggrpc.GetRecipeRatingRequest{
			RecipeID:       exampleResult.RecipeID,
			RecipeRatingID: exampleResult.ID,
		})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipeRatingsForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeRecipeRatingsList()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ListRecipeRatings", testutils.ContextMatcher, exampleRecipeID, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipeRatingsForRecipe(ctx, &mealplanninggrpc.GetRecipeRatingsForRecipeRequest{RecipeID: exampleRecipeID})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeRecipeStep()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ReadRecipeStep", testutils.ContextMatcher, exampleResult.BelongsToRecipe, exampleResult.ID).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipeStep(ctx, &mealplanninggrpc.GetRecipeStepRequest{
			RecipeID:     exampleResult.BelongsToRecipe,
			RecipeStepID: exampleResult.ID,
		})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeRecipeStepCompletionCondition()
		exampleRecipeID := mealplanningfakes.BuildFakeID()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ReadRecipeStepCompletionCondition", testutils.ContextMatcher, exampleRecipeID, exampleResult.BelongsToRecipeStep, exampleResult.ID).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipeStepCompletionCondition(ctx, &mealplanninggrpc.GetRecipeStepCompletionConditionRequest{
			RecipeID:                        exampleRecipeID,
			RecipeStepID:                    exampleResult.BelongsToRecipeStep,
			RecipeStepCompletionConditionID: exampleResult.ID,
		})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipeStepCompletionConditions(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeRecipeStepCompletionConditionsList()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ListRecipeStepCompletionConditions", testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipeStepCompletionConditions(ctx, &mealplanninggrpc.GetRecipeStepCompletionConditionsRequest{
			RecipeID:     exampleRecipeID,
			RecipeStepID: exampleRecipeStepID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeRecipeStepIngredient()
		exampleRecipeID := mealplanningfakes.BuildFakeID()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ReadRecipeStepIngredient", testutils.ContextMatcher, exampleRecipeID, exampleResult.BelongsToRecipeStep, exampleResult.ID).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipeStepIngredient(ctx, &mealplanninggrpc.GetRecipeStepIngredientRequest{
			RecipeID:               exampleRecipeID,
			RecipeStepID:           exampleResult.BelongsToRecipeStep,
			RecipeStepIngredientID: exampleResult.ID,
		})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeRecipeStepIngredientsList()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ListRecipeStepIngredients", testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipeStepIngredients(ctx, &mealplanninggrpc.GetRecipeStepIngredientsRequest{
			RecipeID:     exampleRecipeID,
			RecipeStepID: exampleRecipeStepID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeRecipeStepInstrument()
		exampleRecipeID := mealplanningfakes.BuildFakeID()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ReadRecipeStepInstrument", testutils.ContextMatcher, exampleRecipeID, exampleResult.BelongsToRecipeStep, exampleResult.ID).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipeStepInstrument(ctx, &mealplanninggrpc.GetRecipeStepInstrumentRequest{
			RecipeID:               exampleRecipeID,
			RecipeStepID:           exampleResult.BelongsToRecipeStep,
			RecipeStepInstrumentID: exampleResult.ID,
		})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipeStepInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeRecipeStepInstrumentsList()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ListRecipeStepInstruments", testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipeStepInstruments(ctx, &mealplanninggrpc.GetRecipeStepInstrumentsRequest{
			RecipeID:     exampleRecipeID,
			RecipeStepID: exampleRecipeStepID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeRecipeStepProduct()
		exampleRecipeID := mealplanningfakes.BuildFakeID()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ReadRecipeStepProduct", testutils.ContextMatcher, exampleRecipeID, exampleResult.BelongsToRecipeStep, exampleResult.ID).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipeStepProduct(ctx, &mealplanninggrpc.GetRecipeStepProductRequest{
			RecipeID:            exampleRecipeID,
			RecipeStepID:        exampleResult.BelongsToRecipeStep,
			RecipeStepProductID: exampleResult.ID,
		})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipeStepProducts(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeRecipeStepProductsList()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ListRecipeStepProducts", testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipeStepProducts(ctx, &mealplanninggrpc.GetRecipeStepProductsRequest{
			RecipeID:     exampleRecipeID,
			RecipeStepID: exampleRecipeStepID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeRecipeStepVessel()
		exampleRecipeID := mealplanningfakes.BuildFakeID()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ReadRecipeStepVessel", testutils.ContextMatcher, exampleRecipeID, exampleResult.BelongsToRecipeStep, exampleResult.ID).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipeStepVessel(ctx, &mealplanninggrpc.GetRecipeStepVesselRequest{
			RecipeID:           exampleRecipeID,
			RecipeStepID:       exampleResult.BelongsToRecipeStep,
			RecipeStepVesselID: exampleResult.ID,
		})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipeStepVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleRecipeStepID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeRecipeStepVesselsList()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ListRecipeStepVessels", testutils.ContextMatcher, exampleRecipeID, exampleRecipeStepID, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipeStepVessels(ctx, &mealplanninggrpc.GetRecipeStepVesselsRequest{
			RecipeID:     exampleRecipeID,
			RecipeStepID: exampleRecipeStepID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipeSteps(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeRecipeStepsList()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ListRecipeSteps", testutils.ContextMatcher, exampleRecipeID, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipeSteps(ctx, &mealplanninggrpc.GetRecipeStepsRequest{RecipeID: exampleRecipeID})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_GetRecipes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeRecipesList()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("ListRecipes", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.GetRecipes(ctx, &mealplanninggrpc.GetRecipesRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_SearchForRecipes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeRecipesList()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.SearchForRecipesRequest](t)

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("SearchRecipes", testutils.ContextMatcher, exampleRequest.Query, exampleRequest.UseSearchService, testutils.QueryFilterMatcher).Return(exampleResult, nil)
		s.recipeManager = mrm

		result, err := s.SearchForRecipes(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_UpdateRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateRecipeRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeRecipe()

		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("UpdateRecipe", testutils.ContextMatcher, exampleRequest.RecipeID, testutils.MatchType[*mealplanning.RecipeUpdateRequestInput]()).Return(nil)
		mrm.On("ReadRecipe", testutils.ContextMatcher, exampleRequest.RecipeID).Return(exampleResponse, nil)
		s.recipeManager = mrm

		res, err := s.UpdateRecipe(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Updated.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_UpdateRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateRecipePrepTaskRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeRecipePrepTask()

		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("UpdateRecipePrepTask", testutils.ContextMatcher, exampleRequest.RecipeID, exampleRequest.RecipePrepTaskID, testutils.MatchType[*mealplanning.RecipePrepTaskUpdateRequestInput]()).Return(nil)
		mrm.On("ReadRecipePrepTask", testutils.ContextMatcher, exampleRequest.RecipeID, exampleRequest.RecipePrepTaskID).Return(exampleResponse, nil)
		s.recipeManager = mrm

		res, err := s.UpdateRecipePrepTask(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Updated.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_UpdateRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateRecipeRatingRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeRecipeRating()

		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("UpdateRecipeRating", testutils.ContextMatcher, exampleRequest.RecipeID, exampleRequest.RecipeRatingID, testutils.MatchType[*mealplanning.RecipeRatingUpdateRequestInput]()).Return(nil)
		mrm.On("ReadRecipeRating", testutils.ContextMatcher, exampleRequest.RecipeID, exampleRequest.RecipeRatingID).Return(exampleResponse, nil)
		s.recipeManager = mrm

		res, err := s.UpdateRecipeRating(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Updated.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_UpdateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateRecipeStepRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeRecipeStep()

		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("UpdateRecipeStep", testutils.ContextMatcher, exampleRequest.RecipeID, exampleRequest.RecipeStepID, testutils.MatchType[*mealplanning.RecipeStepUpdateRequestInput]()).Return(nil)
		mrm.On("ReadRecipeStep", testutils.ContextMatcher, exampleRequest.RecipeID, exampleRequest.RecipeStepID).Return(exampleResponse, nil)
		s.recipeManager = mrm

		res, err := s.UpdateRecipeStep(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Updated.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_UpdateRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateRecipeStepCompletionConditionRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeRecipeStepCompletionCondition()

		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("UpdateRecipeStepCompletionCondition", testutils.ContextMatcher, exampleRequest.RecipeID, exampleRequest.RecipeStepID, exampleRequest.RecipeStepCompletionConditionID, testutils.MatchType[*mealplanning.RecipeStepCompletionConditionUpdateRequestInput]()).Return(nil)
		mrm.On("ReadRecipeStepCompletionCondition", testutils.ContextMatcher, exampleRequest.RecipeID, exampleRequest.RecipeStepID, exampleRequest.RecipeStepCompletionConditionID).Return(exampleResponse, nil)
		s.recipeManager = mrm

		res, err := s.UpdateRecipeStepCompletionCondition(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Updated.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_UpdateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateRecipeStepIngredientRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeRecipeStepIngredient()

		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("UpdateRecipeStepIngredient", testutils.ContextMatcher, exampleRequest.RecipeID, exampleRequest.RecipeStepID, exampleRequest.RecipeStepIngredientID, testutils.MatchType[*mealplanning.RecipeStepIngredientUpdateRequestInput]()).Return(nil)
		mrm.On("ReadRecipeStepIngredient", testutils.ContextMatcher, exampleRequest.RecipeID, exampleRequest.RecipeStepID, exampleRequest.RecipeStepIngredientID).Return(exampleResponse, nil)
		s.recipeManager = mrm

		res, err := s.UpdateRecipeStepIngredient(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Updated.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_UpdateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateRecipeStepInstrumentRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeRecipeStepInstrument()

		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("UpdateRecipeStepInstrument", testutils.ContextMatcher, exampleRequest.RecipeID, exampleRequest.RecipeStepID, exampleRequest.RecipeStepInstrumentID, testutils.MatchType[*mealplanning.RecipeStepInstrumentUpdateRequestInput]()).Return(nil)
		mrm.On("ReadRecipeStepInstrument", testutils.ContextMatcher, exampleRequest.RecipeID, exampleRequest.RecipeStepID, exampleRequest.RecipeStepInstrumentID).Return(exampleResponse, nil)
		s.recipeManager = mrm

		res, err := s.UpdateRecipeStepInstrument(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Updated.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_UpdateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateRecipeStepProductRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeRecipeStepProduct()

		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("UpdateRecipeStepProduct", testutils.ContextMatcher, exampleRequest.RecipeID, exampleRequest.RecipeStepID, exampleRequest.RecipeStepProductID, testutils.MatchType[*mealplanning.RecipeStepProductUpdateRequestInput]()).Return(nil)
		mrm.On("ReadRecipeStepProduct", testutils.ContextMatcher, exampleRequest.RecipeID, exampleRequest.RecipeStepID, exampleRequest.RecipeStepProductID).Return(exampleResponse, nil)
		s.recipeManager = mrm

		res, err := s.UpdateRecipeStepProduct(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Updated.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}

func TestServiceImpl_UpdateRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateRecipeStepVesselRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeRecipeStepVessel()

		s := buildServiceImplForRecipesTest(t)

		mrm := &mockmanagers.MockRecipeManager{}
		mrm.On("UpdateRecipeStepVessel", testutils.ContextMatcher, exampleRequest.RecipeID, exampleRequest.RecipeStepID, exampleRequest.RecipeStepVesselID, testutils.MatchType[*mealplanning.RecipeStepVesselUpdateRequestInput]()).Return(nil)
		mrm.On("ReadRecipeStepVessel", testutils.ContextMatcher, exampleRequest.RecipeID, exampleRequest.RecipeStepID, exampleRequest.RecipeStepVesselID).Return(exampleResponse, nil)
		s.recipeManager = mrm

		res, err := s.UpdateRecipeStepVessel(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Updated.ID)

		mock.AssertExpectationsForObjects(t, mrm)
	})
}
