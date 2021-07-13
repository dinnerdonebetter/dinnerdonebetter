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

func TestService_fetchRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"GetRecipe",
			testutils.ContextMatcher,
			exampleRecipe.ID,
		).Return(exampleRecipe, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		actual, err := s.service.fetchRecipe(s.ctx, req)
		assert.Equal(t, exampleRecipe, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with fake mode", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.useFakeData = true

		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		actual, err := s.service.fetchRecipe(s.ctx, req)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with error fetching recipe", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"GetRecipe",
			testutils.ContextMatcher,
			exampleRecipe.ID,
		).Return((*types.Recipe)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		actual, err := s.service.fetchRecipe(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachRecipeCreationInputToRequest(input *types.RecipeCreationInput) *http.Request {
	form := url.Values{
		recipeCreationInputNameFormKey:        {anyToString(input.Name)},
		recipeCreationInputSourceFormKey:      {anyToString(input.Source)},
		recipeCreationInputDescriptionFormKey: {anyToString(input.Description)},
	}

	if input.InspiredByRecipeID != nil {
		form.Set(recipeCreationInputInspiredByRecipeIDFormKey, anyToString(*input.InspiredByRecipeID))
	}

	return httptest.NewRequest(http.MethodPost, "/recipes", strings.NewReader(form.Encode()))
}

func TestService_buildRecipeCreatorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		s.service.buildRecipeCreatorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		s.service.buildRecipeCreatorView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		s.service.buildRecipeCreatorView(false)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		s.service.buildRecipeCreatorView(true)(res, req)
	})

	T.Run("without base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		s.service.buildRecipeCreatorView(false)(res, req)
	})
}

func TestService_parseFormEncodedRecipeCreationInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		expected := fakes.BuildFakeRecipeCreationInput()
		expected.BelongsToAccount = s.exampleAccount.ID
		req := attachRecipeCreationInputToRequest(expected)

		actual := s.service.parseFormEncodedRecipeCreationInput(s.ctx, req, s.sessionCtxData)
		assert.Equal(t, expected, actual)
	})

	T.Run("with error extracting form from request", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedRecipeCreationInput(s.ctx, req, s.sessionCtxData)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.RecipeCreationInput{}
		req := attachRecipeCreationInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedRecipeCreationInput(s.ctx, req, s.sessionCtxData)
		assert.Nil(t, actual)
	})
}

func TestService_handleRecipeCreationRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		exampleInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
		exampleInput.BelongsToAccount = s.sessionCtxData.ActiveAccountID

		res := httptest.NewRecorder()
		req := attachRecipeCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"CreateRecipe",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleRecipe, nil)
		s.service.dataStore = mockDB

		s.service.handleRecipeCreationRequest(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.NotEmpty(t, res.Header().Get(htmxRedirectionHeader))

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		exampleInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachRecipeCreationInputToRequest(exampleInput)

		s.service.handleRecipeCreationRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		exampleInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
		exampleInput.BelongsToAccount = s.sessionCtxData.ActiveAccountID

		res := httptest.NewRecorder()
		req := attachRecipeCreationInputToRequest(&types.RecipeCreationInput{})

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"CreateRecipe",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleRecipe, nil)
		s.service.dataStore = mockDB

		s.service.handleRecipeCreationRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error creating recipe in database", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		exampleInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
		exampleInput.BelongsToAccount = s.sessionCtxData.ActiveAccountID

		res := httptest.NewRecorder()
		req := attachRecipeCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"CreateRecipe",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return((*types.Recipe)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		s.service.handleRecipeCreationRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildRecipeEditorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"GetRecipe",
			testutils.ContextMatcher,
			exampleRecipe.ID,
		).Return(exampleRecipe, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		s.service.buildRecipeEditorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"GetRecipe",
			testutils.ContextMatcher,
			exampleRecipe.ID,
		).Return(exampleRecipe, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		s.service.buildRecipeEditorView(false)(res, req)

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
		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		s.service.buildRecipeEditorView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching recipe", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"GetRecipe",
			testutils.ContextMatcher,
			exampleRecipe.ID,
		).Return((*types.Recipe)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		s.service.buildRecipeEditorView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_fetchRecipes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeList := fakes.BuildFakeRecipeList()

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"GetRecipes",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeList, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		actual, err := s.service.fetchRecipes(s.ctx, req)
		assert.Equal(t, exampleRecipeList, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with fake mode", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.useFakeData = true

		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		actual, err := s.service.fetchRecipes(s.ctx, req)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"GetRecipes",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.RecipeList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		actual, err := s.service.fetchRecipes(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildRecipesTableView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeList := fakes.BuildFakeRecipeList()
		for _, recipe := range exampleRecipeList.Recipes {
			recipe.BelongsToAccount = s.exampleAccount.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"GetRecipes",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		s.service.buildRecipesTableView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipeList := fakes.BuildFakeRecipeList()

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"GetRecipes",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		s.service.buildRecipesTableView(false)(res, req)

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
		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		s.service.buildRecipesTableView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"GetRecipes",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.RecipeList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipes", nil)

		s.service.buildRecipesTableView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachRecipeUpdateInputToRequest(input *types.RecipeUpdateInput) *http.Request {
	form := url.Values{
		recipeUpdateInputNameFormKey:        {anyToString(input.Name)},
		recipeUpdateInputSourceFormKey:      {anyToString(input.Source)},
		recipeUpdateInputDescriptionFormKey: {anyToString(input.Description)},
	}

	if input.InspiredByRecipeID != nil {
		form.Set(recipeUpdateInputInspiredByRecipeIDFormKey, anyToString(*input.InspiredByRecipeID))
	}

	return httptest.NewRequest(http.MethodPost, "/recipes", strings.NewReader(form.Encode()))
}

func TestService_parseFormEncodedRecipeUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		expected := fakes.BuildFakeRecipeUpdateInputFromRecipe(exampleRecipe)

		req := attachRecipeUpdateInputToRequest(expected)

		actual := s.service.parseFormEncodedRecipeUpdateInput(s.ctx, req, s.sessionCtxData)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedRecipeUpdateInput(s.ctx, req, s.sessionCtxData)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input attached to valid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.RecipeUpdateInput{}

		req := attachRecipeUpdateInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedRecipeUpdateInput(s.ctx, req, s.sessionCtxData)
		assert.Nil(t, actual)
	})
}

func TestService_handleRecipeUpdateRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		exampleInput := fakes.BuildFakeRecipeUpdateInputFromRecipe(exampleRecipe)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"GetRecipe",
			testutils.ContextMatcher,
			exampleRecipe.ID,
		).Return(exampleRecipe, nil)

		mockDB.RecipeDataManager.On(
			"UpdateRecipe",
			testutils.ContextMatcher,
			exampleRecipe,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachRecipeUpdateInputToRequest(exampleInput)

		s.service.handleRecipeUpdateRequest(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		exampleInput := fakes.BuildFakeRecipeUpdateInputFromRecipe(exampleRecipe)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachRecipeUpdateInputToRequest(exampleInput)

		s.service.handleRecipeUpdateRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.RecipeUpdateInput{}

		res := httptest.NewRecorder()
		req := attachRecipeUpdateInputToRequest(exampleInput)

		s.service.handleRecipeUpdateRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		exampleInput := fakes.BuildFakeRecipeUpdateInputFromRecipe(exampleRecipe)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"GetRecipe",
			testutils.ContextMatcher,
			exampleRecipe.ID,
		).Return((*types.Recipe)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachRecipeUpdateInputToRequest(exampleInput)

		s.service.handleRecipeUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error updating data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		exampleInput := fakes.BuildFakeRecipeUpdateInputFromRecipe(exampleRecipe)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"GetRecipe",
			testutils.ContextMatcher,
			exampleRecipe.ID,
		).Return(exampleRecipe, nil)

		mockDB.RecipeDataManager.On(
			"UpdateRecipe",
			testutils.ContextMatcher,
			exampleRecipe,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachRecipeUpdateInputToRequest(exampleInput)

		s.service.handleRecipeUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_handleRecipeArchiveRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		exampleRecipeList := fakes.BuildFakeRecipeList()

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"ArchiveRecipe",
			testutils.ContextMatcher,
			exampleRecipe.ID,
			s.sessionCtxData.ActiveAccountID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.RecipeDataManager.On(
			"GetRecipes",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/recipes", nil)

		s.service.handleRecipeArchiveRequest(res, req)

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
		req := httptest.NewRequest(http.MethodDelete, "/recipes", nil)

		s.service.handleRecipeArchiveRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error archiving recipe", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"ArchiveRecipe",
			testutils.ContextMatcher,
			exampleRecipe.ID,
			s.sessionCtxData.ActiveAccountID,
			s.sessionCtxData.Requester.UserID,
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/recipes", nil)

		s.service.handleRecipeArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error retrieving new list of recipes", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.BelongsToAccount = s.exampleAccount.ID
		s.service.recipeIDFetcher = func(*http.Request) uint64 {
			return exampleRecipe.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeDataManager.On(
			"ArchiveRecipe",
			testutils.ContextMatcher,
			exampleRecipe.ID,
			s.sessionCtxData.ActiveAccountID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.RecipeDataManager.On(
			"GetRecipes",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.RecipeList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/recipes", nil)

		s.service.handleRecipeArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
