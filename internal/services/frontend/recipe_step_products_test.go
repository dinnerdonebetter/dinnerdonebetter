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

func TestService_fetchRecipeStepProduct(T *testing.T) {
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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"GetRecipeStepProduct",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
		).Return(exampleRecipeStepProduct, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		actual, err := s.service.fetchRecipeStepProduct(s.ctx, req)
		assert.Equal(t, exampleRecipeStepProduct, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching recipe step product", func(t *testing.T) {
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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"GetRecipeStepProduct",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
		).Return((*types.RecipeStepProduct)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		actual, err := s.service.fetchRecipeStepProduct(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachRecipeStepProductCreationInputToRequest(input *types.RecipeStepProductCreationInput) *http.Request {
	form := url.Values{
		recipeStepProductCreationInputNameFormKey:          {anyToString(input.Name)},
		recipeStepProductCreationInputQuantityTypeFormKey:  {anyToString(input.QuantityType)},
		recipeStepProductCreationInputQuantityValueFormKey: {anyToString(input.QuantityValue)},
		recipeStepProductCreationInputQuantityNotesFormKey: {anyToString(input.QuantityNotes)},
	}

	return httptest.NewRequest(http.MethodPost, "/recipe_step_products", strings.NewReader(form.Encode()))
}

func TestService_buildRecipeStepProductCreatorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		s.service.buildRecipeStepProductCreatorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		s.service.buildRecipeStepProductCreatorView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		s.service.buildRecipeStepProductCreatorView(false)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		s.service.buildRecipeStepProductCreatorView(true)(res, req)
	})

	T.Run("without base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		s.service.buildRecipeStepProductCreatorView(false)(res, req)
	})
}

func TestService_parseFormEncodedRecipeStepProductCreationInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		expected := fakes.BuildFakeRecipeStepProductCreationInput()
		expected.BelongsToRecipeStep = 0
		req := attachRecipeStepProductCreationInputToRequest(expected)

		actual := s.service.parseFormEncodedRecipeStepProductCreationInput(s.ctx, req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with error extracting form from request", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedRecipeStepProductCreationInput(s.ctx, req)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.RecipeStepProductCreationInput{}
		req := attachRecipeStepProductCreationInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedRecipeStepProductCreationInput(s.ctx, req)
		assert.Nil(t, actual)
	})
}

func TestService_handleRecipeStepProductCreationRequest(T *testing.T) {
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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)

		res := httptest.NewRecorder()
		req := attachRecipeStepProductCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"CreateRecipeStepProduct",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleRecipeStepProduct, nil)
		s.service.dataStore = mockDB

		s.service.handleRecipeStepProductCreationRequest(res, req)

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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachRecipeStepProductCreationInputToRequest(exampleInput)

		s.service.handleRecipeStepProductCreationRequest(res, req)

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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)

		res := httptest.NewRecorder()
		req := attachRecipeStepProductCreationInputToRequest(&types.RecipeStepProductCreationInput{})

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"CreateRecipeStepProduct",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleRecipeStepProduct, nil)
		s.service.dataStore = mockDB

		s.service.handleRecipeStepProductCreationRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error creating recipe step product in database", func(t *testing.T) {
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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)

		res := httptest.NewRecorder()
		req := attachRecipeStepProductCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"CreateRecipeStepProduct",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return((*types.RecipeStepProduct)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		s.service.handleRecipeStepProductCreationRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildRecipeStepProductEditorView(T *testing.T) {
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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"GetRecipeStepProduct",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
		).Return(exampleRecipeStepProduct, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		s.service.buildRecipeStepProductEditorView(true)(res, req)

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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"GetRecipeStepProduct",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
		).Return(exampleRecipeStepProduct, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		s.service.buildRecipeStepProductEditorView(false)(res, req)

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
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		s.service.buildRecipeStepProductEditorView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching recipe step product", func(t *testing.T) {
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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"GetRecipeStepProduct",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
		).Return((*types.RecipeStepProduct)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		s.service.buildRecipeStepProductEditorView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_fetchRecipeStepProducts(T *testing.T) {
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

		exampleRecipeStepProductList := fakes.BuildFakeRecipeStepProductList()

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"GetRecipeStepProducts",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeStepProductList, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		actual, err := s.service.fetchRecipeStepProducts(s.ctx, req)
		assert.Equal(t, exampleRecipeStepProductList, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
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
		mockDB.RecipeStepProductDataManager.On(
			"GetRecipeStepProducts",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.RecipeStepProductList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		actual, err := s.service.fetchRecipeStepProducts(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildRecipeStepProductsTableView(T *testing.T) {
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

		exampleRecipeStepProductList := fakes.BuildFakeRecipeStepProductList()
		for _, recipeStepProduct := range exampleRecipeStepProductList.RecipeStepProducts {
			recipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"GetRecipeStepProducts",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeStepProductList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		s.service.buildRecipeStepProductsTableView(true)(res, req)

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

		exampleRecipeStepProductList := fakes.BuildFakeRecipeStepProductList()

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"GetRecipeStepProducts",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeStepProductList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		s.service.buildRecipeStepProductsTableView(false)(res, req)

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
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		s.service.buildRecipeStepProductsTableView(true)(res, req)

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
		mockDB.RecipeStepProductDataManager.On(
			"GetRecipeStepProducts",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.RecipeStepProductList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/recipe_step_products", nil)

		s.service.buildRecipeStepProductsTableView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachRecipeStepProductUpdateInputToRequest(input *types.RecipeStepProductUpdateInput) *http.Request {
	form := url.Values{
		recipeStepProductUpdateInputNameFormKey:          {anyToString(input.Name)},
		recipeStepProductUpdateInputQuantityTypeFormKey:  {anyToString(input.QuantityType)},
		recipeStepProductUpdateInputQuantityValueFormKey: {anyToString(input.QuantityValue)},
		recipeStepProductUpdateInputQuantityNotesFormKey: {anyToString(input.QuantityNotes)},
	}

	return httptest.NewRequest(http.MethodPost, "/recipe_step_products", strings.NewReader(form.Encode()))
}

func TestService_parseFormEncodedRecipeStepProductUpdateInput(T *testing.T) {
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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		expected := fakes.BuildFakeRecipeStepProductUpdateInputFromRecipeStepProduct(exampleRecipeStepProduct)
		expected.BelongsToRecipeStep = 0

		req := attachRecipeStepProductUpdateInputToRequest(expected)

		actual := s.service.parseFormEncodedRecipeStepProductUpdateInput(s.ctx, req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedRecipeStepProductUpdateInput(s.ctx, req)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input attached to valid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.RecipeStepProductUpdateInput{}

		req := attachRecipeStepProductUpdateInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedRecipeStepProductUpdateInput(s.ctx, req)
		assert.Nil(t, actual)
	})
}

func TestService_handleRecipeStepProductUpdateRequest(T *testing.T) {
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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepProductUpdateInputFromRecipeStepProduct(exampleRecipeStepProduct)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"GetRecipeStepProduct",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
		).Return(exampleRecipeStepProduct, nil)

		mockDB.RecipeStepProductDataManager.On(
			"UpdateRecipeStepProduct",
			testutils.ContextMatcher,
			exampleRecipeStepProduct,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachRecipeStepProductUpdateInputToRequest(exampleInput)

		s.service.handleRecipeStepProductUpdateRequest(res, req)

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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepProductUpdateInputFromRecipeStepProduct(exampleRecipeStepProduct)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachRecipeStepProductUpdateInputToRequest(exampleInput)

		s.service.handleRecipeStepProductUpdateRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.RecipeStepProductUpdateInput{}

		res := httptest.NewRecorder()
		req := attachRecipeStepProductUpdateInputToRequest(exampleInput)

		s.service.handleRecipeStepProductUpdateRequest(res, req)

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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepProductUpdateInputFromRecipeStepProduct(exampleRecipeStepProduct)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"GetRecipeStepProduct",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
		).Return((*types.RecipeStepProduct)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachRecipeStepProductUpdateInputToRequest(exampleInput)

		s.service.handleRecipeStepProductUpdateRequest(res, req)

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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		exampleInput := fakes.BuildFakeRecipeStepProductUpdateInputFromRecipeStepProduct(exampleRecipeStepProduct)

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"GetRecipeStepProduct",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
		).Return(exampleRecipeStepProduct, nil)

		mockDB.RecipeStepProductDataManager.On(
			"UpdateRecipeStepProduct",
			testutils.ContextMatcher,
			exampleRecipeStepProduct,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachRecipeStepProductUpdateInputToRequest(exampleInput)

		s.service.handleRecipeStepProductUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_handleRecipeStepProductArchiveRequest(T *testing.T) {
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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		exampleRecipeStepProductList := fakes.BuildFakeRecipeStepProductList()

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"ArchiveRecipeStepProduct",
			testutils.ContextMatcher,
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.RecipeStepProductDataManager.On(
			"GetRecipeStepProducts",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeStepProductList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/recipe_step_products", nil)

		s.service.handleRecipeStepProductArchiveRequest(res, req)

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
		req := httptest.NewRequest(http.MethodDelete, "/recipe_step_products", nil)

		s.service.handleRecipeStepProductArchiveRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error archiving recipe step product", func(t *testing.T) {
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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"ArchiveRecipeStepProduct",
			testutils.ContextMatcher,
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/recipe_step_products", nil)

		s.service.handleRecipeStepProductArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error retrieving new list of recipe step products", func(t *testing.T) {
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

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStepID
		s.service.recipeStepProductIDFetcher = func(*http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.RecipeStepProductDataManager.On(
			"ArchiveRecipeStepProduct",
			testutils.ContextMatcher,
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.RecipeStepProductDataManager.On(
			"GetRecipeStepProducts",
			testutils.ContextMatcher,
			exampleRecipeID,
			exampleRecipeStepID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.RecipeStepProductList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/recipe_step_products", nil)

		s.service.handleRecipeStepProductArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
