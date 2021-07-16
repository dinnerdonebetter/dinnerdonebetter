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

func TestService_fetchRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"GetRecipeStepIngredient",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return(exampleRecipeStepIngredient, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		actual, err := s.service.fetchRecipeStepIngredient(s.ctx, req)
		assert.Equal(t, exampleRecipeStepIngredient, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with fake mode", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.useFakeData = true

		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		actual, err := s.service.fetchRecipeStepIngredient(s.ctx, req)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with error fetching recipe step ingredient", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"GetRecipeStepIngredient",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return((*types.RecipeStepIngredient)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		actual, err := s.service.fetchRecipeStepIngredient(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachRecipeStepIngredientCreationInputToRequest(input *types.RecipeStepIngredientCreationInput) *http.Request {
	form := url.Values{
		recipeStepIngredientCreationInputNameFormKey:                {anyToString(input.Name)},
		recipeStepIngredientCreationInputQuantityTypeFormKey:        {anyToString(input.QuantityType)},
		recipeStepIngredientCreationInputQuantityValueFormKey:       {anyToString(input.QuantityValue)},
		recipeStepIngredientCreationInputQuantityNotesFormKey:       {anyToString(input.QuantityNotes)},
		recipeStepIngredientCreationInputProductOfRecipeStepFormKey: {anyToString(input.ProductOfRecipeStep)},
		recipeStepIngredientCreationInputIngredientNotesFormKey:     {anyToString(input.IngredientNotes)},
	}

	if input.IngredientID != nil {
		form.Set(recipeStepIngredientCreationInputIngredientIDFormKey, anyToString(*input.IngredientID))
	}

	return httptest.NewRequest(http.MethodPost, "/recipe_step_ingredients", strings.NewReader(form.Encode()))
}

func TestService_buildRecipeStepIngredientCreatorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		s.service.buildRecipeStepIngredientCreatorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		s.service.buildRecipeStepIngredientCreatorView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		s.service.buildRecipeStepIngredientCreatorView(false)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		s.service.buildRecipeStepIngredientCreatorView(true)(res, req)
	})

	T.Run("without base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		s.service.buildRecipeStepIngredientCreatorView(false)(res, req)
	})
}

func TestService_parseFormEncodedRecipeStepIngredientCreationInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		expected := fakes.BuildFakeRecipeStepIngredientCreationInput()
		expected.BelongsToRecipeStep = 0
		req := attachRecipeStepIngredientCreationInputToRequest(expected)

		actual := s.service.parseFormEncodedRecipeStepIngredientCreationInput(s.ctx, req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with error extracting form from request", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedRecipeStepIngredientCreationInput(s.ctx, req)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.RecipeStepIngredientCreationInput{}
		req := attachRecipeStepIngredientCreationInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedRecipeStepIngredientCreationInput(s.ctx, req)
		assert.Nil(t, actual)
	})
}

func TestService_handleRecipeStepIngredientCreationRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		res := httptest.NewRecorder()
		req := attachRecipeStepIngredientCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"CreateRecipeStepIngredient",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleRecipeStepIngredient, nil)
		s.service.dataStore = mockDB

		s.service.handleRecipeStepIngredientCreationRequest(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.NotEmpty(t, res.Header().Get(htmxRedirectionHeader))

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachRecipeStepIngredientCreationInputToRequest(exampleInput)

		s.service.handleRecipeStepIngredientCreationRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		res := httptest.NewRecorder()
		req := attachRecipeStepIngredientCreationInputToRequest(&types.RecipeStepIngredientCreationInput{})

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"CreateRecipeStepIngredient",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleRecipeStepIngredient, nil)
		s.service.dataStore = mockDB

		s.service.handleRecipeStepIngredientCreationRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error creating recipe step ingredient in database", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		res := httptest.NewRecorder()
		req := attachRecipeStepIngredientCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"CreateRecipeStepIngredient",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return((*types.RecipeStepIngredient)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		s.service.handleRecipeStepIngredientCreationRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildRecipeStepIngredientEditorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"GetRecipeStepIngredient",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return(exampleRecipeStepIngredient, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		s.service.buildRecipeStepIngredientEditorView(true)(res, req)

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

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"GetRecipeStepIngredient",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return(exampleRecipeStepIngredient, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		s.service.buildRecipeStepIngredientEditorView(false)(res, req)

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
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		s.service.buildRecipeStepIngredientEditorView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching recipe step ingredient", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"GetRecipeStepIngredient",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return((*types.RecipeStepIngredient)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		s.service.buildRecipeStepIngredientEditorView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_fetchRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"GetRecipeStepIngredients",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeStepIngredientList, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		actual, err := s.service.fetchRecipeStepIngredients(s.ctx, req)
		assert.Equal(t, exampleRecipeStepIngredientList, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with fake mode", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.useFakeData = true

		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		actual, err := s.service.fetchRecipeStepIngredients(s.ctx, req)
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

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"GetRecipeStepIngredients",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.RecipeStepIngredientList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		actual, err := s.service.fetchRecipeStepIngredients(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildRecipeStepIngredientsTableView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()
		for _, recipeStepIngredient := range exampleRecipeStepIngredientList.RecipeStepIngredients {
			recipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"GetRecipeStepIngredients",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeStepIngredientList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		s.service.buildRecipeStepIngredientsTableView(true)(res, req)

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

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"GetRecipeStepIngredients",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeStepIngredientList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		s.service.buildRecipeStepIngredientsTableView(false)(res, req)

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
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		s.service.buildRecipeStepIngredientsTableView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"GetRecipeStepIngredients",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.RecipeStepIngredientList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_ingredients", nil)

		s.service.buildRecipeStepIngredientsTableView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachRecipeStepIngredientUpdateInputToRequest(input *types.RecipeStepIngredientUpdateInput) *http.Request {
	form := url.Values{
		recipeStepIngredientUpdateInputNameFormKey:                {anyToString(input.Name)},
		recipeStepIngredientUpdateInputQuantityTypeFormKey:        {anyToString(input.QuantityType)},
		recipeStepIngredientUpdateInputQuantityValueFormKey:       {anyToString(input.QuantityValue)},
		recipeStepIngredientUpdateInputQuantityNotesFormKey:       {anyToString(input.QuantityNotes)},
		recipeStepIngredientUpdateInputProductOfRecipeStepFormKey: {anyToString(input.ProductOfRecipeStep)},
		recipeStepIngredientUpdateInputIngredientNotesFormKey:     {anyToString(input.IngredientNotes)},
	}

	if input.IngredientID != nil {
		form.Set(recipeStepIngredientUpdateInputIngredientIDFormKey, anyToString(*input.IngredientID))
	}

	return httptest.NewRequest(http.MethodPost, "/recipe_step_ingredients", strings.NewReader(form.Encode()))
}

func TestService_parseFormEncodedRecipeStepIngredientUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		expected := fakes.BuildFakeRecipeStepIngredientUpdateInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
		expected.BelongsToRecipeStep = 0

		req := attachRecipeStepIngredientUpdateInputToRequest(expected)

		actual := s.service.parseFormEncodedRecipeStepIngredientUpdateInput(s.ctx, req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedRecipeStepIngredientUpdateInput(s.ctx, req)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input attached to valid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.RecipeStepIngredientUpdateInput{}

		req := attachRecipeStepIngredientUpdateInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedRecipeStepIngredientUpdateInput(s.ctx, req)
		assert.Nil(t, actual)
	})
}

func TestService_handleRecipeStepIngredientUpdateRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepIngredientUpdateInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"GetRecipeStepIngredient",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return(exampleRecipeStepIngredient, nil)

		mockDB.RecipeStepIngredientDataManager.On(
			"UpdateRecipeStepIngredient",
			testutils.ContextMatcher,
			exampleRecipeStepIngredient,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachRecipeStepIngredientUpdateInputToRequest(exampleInput)

		s.service.handleRecipeStepIngredientUpdateRequest(res, req)

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

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepIngredientUpdateInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachRecipeStepIngredientUpdateInputToRequest(exampleInput)

		s.service.handleRecipeStepIngredientUpdateRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.RecipeStepIngredientUpdateInput{}

		res := httptest.NewRecorder()
		req := attachRecipeStepIngredientUpdateInputToRequest(exampleInput)

		s.service.handleRecipeStepIngredientUpdateRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepIngredientUpdateInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"GetRecipeStepIngredient",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return((*types.RecipeStepIngredient)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachRecipeStepIngredientUpdateInputToRequest(exampleInput)

		s.service.handleRecipeStepIngredientUpdateRequest(res, req)

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

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepIngredientUpdateInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"GetRecipeStepIngredient",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		).Return(exampleRecipeStepIngredient, nil)

		mockDB.RecipeStepIngredientDataManager.On(
			"UpdateRecipeStepIngredient",
			testutils.ContextMatcher,
			exampleRecipeStepIngredient,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachRecipeStepIngredientUpdateInputToRequest(exampleInput)

		s.service.handleRecipeStepIngredientUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_handleRecipeStepIngredientArchiveRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"ArchiveRecipeStepIngredient",
			testutils.ContextMatcher,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.RecipeStepIngredientDataManager.On(
			"GetRecipeStepIngredients",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeStepIngredientList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/recipe_step_ingredients", nil)

		s.service.handleRecipeStepIngredientArchiveRequest(res, req)

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
		req := httptest.NewRequest(http.MethodDelete, "/recipe_step_ingredients", nil)

		s.service.handleRecipeStepIngredientArchiveRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error archiving recipe step ingredient", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"ArchiveRecipeStepIngredient",
			testutils.ContextMatcher,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/recipe_step_ingredients", nil)

		s.service.handleRecipeStepIngredientArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error retrieving new list of recipe step ingredients", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeID := fakes.BuildFakeID()
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeID
		}

		exampleRecipeStepID := fakes.BuildFakeID()
		s.service.recipeStepIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepID
		}

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepIngredientDataManager.On(
			"ArchiveRecipeStepIngredient",
			testutils.ContextMatcher,
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.RecipeStepIngredientDataManager.On(
			"GetRecipeStepIngredients",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.RecipeStepIngredientList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/recipe_step_ingredients", nil)

		s.service.handleRecipeStepIngredientArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
