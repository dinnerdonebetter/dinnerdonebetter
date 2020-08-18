package recipestepproducts

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockencoding "gitlab.com/prixfixe/prixfixe/internal/v1/encoding/mock"
	mockmetrics "gitlab.com/prixfixe/prixfixe/internal/v1/metrics/mock"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"
	mockmodels "gitlab.com/prixfixe/prixfixe/models/v1/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mocknewsman "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
)

func TestRecipeStepProductsService_ListHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	exampleRecipe := fakemodels.BuildFakeRecipe()
	exampleRecipe.BelongsToUser = exampleUser.ID
	recipeIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipe.ID
	}

	exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
	exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
	recipeStepIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipeStep.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProductList := fakemodels.BuildFakeRecipeStepProductList()

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("GetRecipeStepProducts", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRecipeStepProductList, nil)
		s.recipeStepProductDataManager = recipeStepProductDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepProductList")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("GetRecipeStepProducts", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RecipeStepProductList)(nil), sql.ErrNoRows)
		s.recipeStepProductDataManager = recipeStepProductDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepProductList")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager, ed)
	})

	T.Run("with error fetching recipe step products from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("GetRecipeStepProducts", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RecipeStepProductList)(nil), errors.New("blah"))
		s.recipeStepProductDataManager = recipeStepProductDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProductList := fakemodels.BuildFakeRecipeStepProductList()

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("GetRecipeStepProducts", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRecipeStepProductList, nil)
		s.recipeStepProductDataManager = recipeStepProductDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepProductList")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager, ed)
	})
}

func TestRecipeStepProductsService_CreateHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	exampleRecipe := fakemodels.BuildFakeRecipe()
	exampleRecipe.BelongsToUser = exampleUser.ID
	recipeIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipe.ID
	}

	exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
	exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
	recipeStepIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipeStep.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher
		s.userIDFetcher = userIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("CreateRecipeStepProduct", mock.Anything, mock.AnythingOfType("*models.RecipeStepProductCreationInput")).Return(exampleRecipeStepProduct, nil)
		s.recipeStepProductDataManager = recipeStepProductDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepProductCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepProduct")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), createMiddlewareCtxKey, exampleInput))

		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager, mc, r, ed)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher
		s.userIDFetcher = userIDFetcher

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error creating recipe step product", func(t *testing.T) {
		s := buildTestService()

		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher
		s.userIDFetcher = userIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("CreateRecipeStepProduct", mock.Anything, mock.AnythingOfType("*models.RecipeStepProductCreationInput")).Return((*models.RecipeStepProduct)(nil), errors.New("blah"))
		s.recipeStepProductDataManager = recipeStepProductDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), createMiddlewareCtxKey, exampleInput))

		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher
		s.userIDFetcher = userIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("CreateRecipeStepProduct", mock.Anything, mock.AnythingOfType("*models.RecipeStepProductCreationInput")).Return(exampleRecipeStepProduct, nil)
		s.recipeStepProductDataManager = recipeStepProductDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepProductCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepProduct")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), createMiddlewareCtxKey, exampleInput))

		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager, mc, r, ed)
	})
}

func TestRecipeStepProductsService_ExistenceHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	exampleRecipe := fakemodels.BuildFakeRecipe()
	exampleRecipe.BelongsToUser = exampleUser.ID
	recipeIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipe.ID
	}

	exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
	exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
	recipeStepIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipeStep.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepProductIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("RecipeStepProductExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return(true, nil)
		s.recipeStepProductDataManager = recipeStepProductDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager)
	})

	T.Run("with no such recipe step product in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepProductIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("RecipeStepProductExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return(false, sql.ErrNoRows)
		s.recipeStepProductDataManager = recipeStepProductDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager)
	})

	T.Run("with error fetching recipe step product from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepProductIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("RecipeStepProductExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return(false, errors.New("blah"))
		s.recipeStepProductDataManager = recipeStepProductDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager)
	})
}

func TestRecipeStepProductsService_ReadHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	exampleRecipe := fakemodels.BuildFakeRecipe()
	exampleRecipe.BelongsToUser = exampleUser.ID
	recipeIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipe.ID
	}

	exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
	exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
	recipeStepIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipeStep.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepProductIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("GetRecipeStepProduct", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return(exampleRecipeStepProduct, nil)
		s.recipeStepProductDataManager = recipeStepProductDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepProduct")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager, ed)
	})

	T.Run("with no such recipe step product in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepProductIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("GetRecipeStepProduct", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return((*models.RecipeStepProduct)(nil), sql.ErrNoRows)
		s.recipeStepProductDataManager = recipeStepProductDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager)
	})

	T.Run("with error fetching recipe step product from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepProductIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("GetRecipeStepProduct", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return((*models.RecipeStepProduct)(nil), errors.New("blah"))
		s.recipeStepProductDataManager = recipeStepProductDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepProductIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("GetRecipeStepProduct", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return(exampleRecipeStepProduct, nil)
		s.recipeStepProductDataManager = recipeStepProductDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepProduct")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager, ed)
	})
}

func TestRecipeStepProductsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	exampleRecipe := fakemodels.BuildFakeRecipe()
	exampleRecipe.BelongsToUser = exampleUser.ID
	recipeIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipe.ID
	}

	exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
	exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
	recipeStepIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipeStep.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepProductUpdateInputFromRecipeStepProduct(exampleRecipeStepProduct)

		s.recipeStepProductIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("GetRecipeStepProduct", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return(exampleRecipeStepProduct, nil)
		recipeStepProductDataManager.On("UpdateRecipeStepProduct", mock.Anything, mock.AnythingOfType("*models.RecipeStepProduct")).Return(nil)
		s.recipeStepProductDataManager = recipeStepProductDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepProduct")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, r, recipeStepProductDataManager, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.UpdateHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with no rows fetching recipe step product", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepProductUpdateInputFromRecipeStepProduct(exampleRecipeStepProduct)

		s.recipeStepProductIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("GetRecipeStepProduct", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return((*models.RecipeStepProduct)(nil), sql.ErrNoRows)
		s.recipeStepProductDataManager = recipeStepProductDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager)
	})

	T.Run("with error fetching recipe step product", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepProductUpdateInputFromRecipeStepProduct(exampleRecipeStepProduct)

		s.recipeStepProductIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("GetRecipeStepProduct", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return((*models.RecipeStepProduct)(nil), errors.New("blah"))
		s.recipeStepProductDataManager = recipeStepProductDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager)
	})

	T.Run("with error updating recipe step product", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepProductUpdateInputFromRecipeStepProduct(exampleRecipeStepProduct)

		s.recipeStepProductIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("GetRecipeStepProduct", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return(exampleRecipeStepProduct, nil)
		recipeStepProductDataManager.On("UpdateRecipeStepProduct", mock.Anything, mock.AnythingOfType("*models.RecipeStepProduct")).Return(errors.New("blah"))
		s.recipeStepProductDataManager = recipeStepProductDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepProductDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepProductUpdateInputFromRecipeStepProduct(exampleRecipeStepProduct)

		s.recipeStepProductIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("GetRecipeStepProduct", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return(exampleRecipeStepProduct, nil)
		recipeStepProductDataManager.On("UpdateRecipeStepProduct", mock.Anything, mock.AnythingOfType("*models.RecipeStepProduct")).Return(nil)
		s.recipeStepProductDataManager = recipeStepProductDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepProduct")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, r, recipeStepProductDataManager, ed)
	})
}

func TestRecipeStepProductsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	exampleRecipe := fakemodels.BuildFakeRecipe()
	exampleRecipe.BelongsToUser = exampleUser.ID
	recipeIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipe.ID
	}

	exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
	exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
	recipeStepIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipeStep.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepProductIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("ArchiveRecipeStepProduct", mock.Anything, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return(nil)
		s.recipeStepProductDataManager = recipeStepProductDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.recipeStepProductCounter = mc

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler(res, req)

		assert.Equal(t, http.StatusNoContent, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepProductDataManager, mc, r)
	})

	T.Run("with nonexistent recipe", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(false, nil)
		s.recipeDataManager = recipeDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("with error checking recipe existence", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, errors.New("blah"))
		s.recipeDataManager = recipeDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("with nonexistent recipe step", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(false, nil)
		s.recipeStepDataManager = recipeStepDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager)
	})

	T.Run("with error checking recipe step existence", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, errors.New("blah"))
		s.recipeStepDataManager = recipeStepDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager)
	})

	T.Run("with no recipe step product in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepProductIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("ArchiveRecipeStepProduct", mock.Anything, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return(sql.ErrNoRows)
		s.recipeStepProductDataManager = recipeStepProductDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepProductDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepProductIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepProduct.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepProductDataManager := &mockmodels.RecipeStepProductDataManager{}
		recipeStepProductDataManager.On("ArchiveRecipeStepProduct", mock.Anything, exampleRecipeStep.ID, exampleRecipeStepProduct.ID).Return(errors.New("blah"))
		s.recipeStepProductDataManager = recipeStepProductDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepProductDataManager)
	})
}
