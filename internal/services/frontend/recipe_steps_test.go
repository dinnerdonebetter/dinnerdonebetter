package frontend

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	database "gitlab.com/prixfixe/prixfixe/internal/database"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
	testutils "gitlab.com/prixfixe/prixfixe/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_fetchRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipeID
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStep.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepDataManager.On(
			"GetRecipeStep",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStep.ID,
		).Return(exampleRecipeStep, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/recipe_steps", nil)

		actual, err := s.service.fetchRecipeStep(s.ctx, req)
		assert.Equal(t, exampleRecipeStep, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with fake mode", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.useFakeData = true

		req := httptest.NewRequest(http.MethodGet, "/recipe_steps", nil)

		actual, err := s.service.fetchRecipeStep(s.ctx, req)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with error fetching recipe step", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipeID
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStep.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepDataManager.On(
			"GetRecipeStep",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStep.ID,
		).Return((*types.RecipeStep)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/recipe_steps", nil)

		actual, err := s.service.fetchRecipeStep(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildRecipeStepEditorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipeID
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStep.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepDataManager.On(
			"GetRecipeStep",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStep.ID,
		).Return(exampleRecipeStep, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_steps", nil)

		s.service.buildRecipeStepEditorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipeID
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStep.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepDataManager.On(
			"GetRecipeStep",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStep.ID,
		).Return(exampleRecipeStep, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_steps", nil)

		s.service.buildRecipeStepEditorView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_steps", nil)

		s.service.buildRecipeStepEditorView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching recipe step", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipeID
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStep.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepDataManager.On(
			"GetRecipeStep",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStep.ID,
		).Return((*types.RecipeStep)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_steps", nil)

		s.service.buildRecipeStepEditorView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_fetchRecipeSteps(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepList := fakes.BuildFakeRecipeStepList()

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepDataManager.On(
			"GetRecipeSteps",
			testutils.ContextMatcher,
			exampleRecipeID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeStepList, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/recipe_steps", nil)

		actual, err := s.service.fetchRecipeSteps(s.ctx, req)
		assert.Equal(t, exampleRecipeStepList, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with fake mode", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.useFakeData = true

		req := httptest.NewRequest(http.MethodGet, "/recipe_steps", nil)

		actual, err := s.service.fetchRecipeSteps(s.ctx, req)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepDataManager.On(
			"GetRecipeSteps",
			testutils.ContextMatcher,
			exampleRecipeID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.RecipeStepList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/recipe_steps", nil)

		actual, err := s.service.fetchRecipeSteps(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildRecipeStepsTableView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepList := fakes.BuildFakeRecipeStepList()
		for _, recipeStep := range exampleRecipeStepList.RecipeSteps {
			recipeStep.BelongsToRecipe = exampleRecipeID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepDataManager.On(
			"GetRecipeSteps",
			testutils.ContextMatcher,
			exampleRecipeID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeStepList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_steps", nil)

		s.service.buildRecipeStepsTableView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepList := fakes.BuildFakeRecipeStepList()

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepDataManager.On(
			"GetRecipeSteps",
			testutils.ContextMatcher,
			exampleRecipeID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeStepList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_steps", nil)

		s.service.buildRecipeStepsTableView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_steps", nil)

		s.service.buildRecipeStepsTableView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepDataManager.On(
			"GetRecipeSteps",
			testutils.ContextMatcher,
			exampleRecipeID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.RecipeStepList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_steps", nil)

		s.service.buildRecipeStepsTableView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachRecipeStepUpdateInputToRequest(input *types.RecipeStepUpdateInput) *http.Request {
	form := url.Values{
		recipeStepUpdateInputIndexFormKey:                     {anyToString(input.Index)},
		recipeStepUpdateInputPreparationIDFormKey:             {anyToString(input.PreparationID)},
		recipeStepUpdateInputPrerequisiteStepFormKey:          {anyToString(input.PrerequisiteStep)},
		recipeStepUpdateInputMinEstimatedTimeInSecondsFormKey: {anyToString(input.MinEstimatedTimeInSeconds)},
		recipeStepUpdateInputMaxEstimatedTimeInSecondsFormKey: {anyToString(input.MaxEstimatedTimeInSeconds)},
		recipeStepUpdateInputNotesFormKey:                     {anyToString(input.Notes)},
		recipeStepUpdateInputWhyFormKey:                       {anyToString(input.Why)},
	}

	if input.TemperatureInCelsius != nil {
		form.Set(recipeStepUpdateInputTemperatureInCelsiusFormKey, anyToString(*input.TemperatureInCelsius))
	}

	return httptest.NewRequest(http.MethodPost, "/recipe_steps", strings.NewReader(form.Encode()))
}

func TestService_parseFormEncodedRecipeStepUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipeID
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStep.ID
		}

		expected := fakes.BuildFakeRecipeStepUpdateInputFromRecipeStep(exampleRecipeStep)
		expected.BelongsToRecipe = 0

		req := attachRecipeStepUpdateInputToRequest(expected)

		actual := s.service.parseFormEncodedRecipeStepUpdateInput(s.ctx, req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedRecipeStepUpdateInput(s.ctx, req)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input attached to valid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.RecipeStepUpdateInput{}

		req := attachRecipeStepUpdateInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedRecipeStepUpdateInput(s.ctx, req)
		assert.Nil(t, actual)
	})
}

func TestService_handleRecipeStepUpdateRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipeID
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStep.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepUpdateInputFromRecipeStep(exampleRecipeStep)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepDataManager.On(
			"GetRecipeStep",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStep.ID,
		).Return(exampleRecipeStep, nil)

		mockDB.RecipeStepDataManager.On(
			"UpdateRecipeStep",
			testutils.ContextMatcher,
			exampleRecipeStep,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachRecipeStepUpdateInputToRequest(exampleInput)

		s.service.handleRecipeStepUpdateRequest(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipeID
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStep.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepUpdateInputFromRecipeStep(exampleRecipeStep)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachRecipeStepUpdateInputToRequest(exampleInput)

		s.service.handleRecipeStepUpdateRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.RecipeStepUpdateInput{}

		res := httptest.NewRecorder()
		req := attachRecipeStepUpdateInputToRequest(exampleInput)

		s.service.handleRecipeStepUpdateRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipeID
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStep.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepUpdateInputFromRecipeStep(exampleRecipeStep)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepDataManager.On(
			"GetRecipeStep",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStep.ID,
		).Return((*types.RecipeStep)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachRecipeStepUpdateInputToRequest(exampleInput)

		s.service.handleRecipeStepUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error updating data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipeID
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStep.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepUpdateInputFromRecipeStep(exampleRecipeStep)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepDataManager.On(
			"GetRecipeStep",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStep.ID,
		).Return(exampleRecipeStep, nil)

		mockDB.RecipeStepDataManager.On(
			"UpdateRecipeStep",
			testutils.ContextMatcher,
			exampleRecipeStep,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachRecipeStepUpdateInputToRequest(exampleInput)

		s.service.handleRecipeStepUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_handleRecipeStepArchiveRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipeID
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStep.ID
		}

		exampleRecipeStepList := fakes.BuildFakeRecipeStepList()

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepDataManager.On(
			"ArchiveRecipeStep",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStep.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.RecipeStepDataManager.On(
			"GetRecipeSteps",
			testutils.ContextMatcher,
			exampleRecipeID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeStepList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/recipe_steps", nil)

		s.service.handleRecipeStepArchiveRequest(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/recipe_steps", nil)

		s.service.handleRecipeStepArchiveRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error archiving recipe step", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipeID
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStep.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepDataManager.On(
			"ArchiveRecipeStep",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStep.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/recipe_steps", nil)

		s.service.handleRecipeStepArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error retrieving new list of recipe steps", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipeID
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStep.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepDataManager.On(
			"ArchiveRecipeStep",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStep.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.RecipeStepDataManager.On(
			"GetRecipeSteps",
			testutils.ContextMatcher,
			exampleRecipeID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.RecipeStepList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/recipe_steps", nil)

		s.service.handleRecipeStepArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
