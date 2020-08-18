package recipeiterations

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

func TestRecipeIterationsService_ListHandler(T *testing.T) {
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

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationList := fakemodels.BuildFakeRecipeIterationList()

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("GetRecipeIterations", mock.Anything, exampleRecipe.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRecipeIterationList, nil)
		s.recipeIterationDataManager = recipeIterationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIterationList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("GetRecipeIterations", mock.Anything, exampleRecipe.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RecipeIterationList)(nil), sql.ErrNoRows)
		s.recipeIterationDataManager = recipeIterationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIterationList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager, ed)
	})

	T.Run("with error fetching recipe iterations from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("GetRecipeIterations", mock.Anything, exampleRecipe.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RecipeIterationList)(nil), errors.New("blah"))
		s.recipeIterationDataManager = recipeIterationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationList := fakemodels.BuildFakeRecipeIterationList()

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("GetRecipeIterations", mock.Anything, exampleRecipe.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRecipeIterationList, nil)
		s.recipeIterationDataManager = recipeIterationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIterationList")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager, ed)
	})
}

func TestRecipeIterationsService_CreateHandler(T *testing.T) {
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

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.recipeIDFetcher = recipeIDFetcher
		s.userIDFetcher = userIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("CreateRecipeIteration", mock.Anything, mock.AnythingOfType("*models.RecipeIterationCreationInput")).Return(exampleRecipeIteration, nil)
		s.recipeIterationDataManager = recipeIterationDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeIterationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIteration")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager, mc, r, ed)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		s.recipeIDFetcher = recipeIDFetcher
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

	T.Run("with error creating recipe iteration", func(t *testing.T) {
		s := buildTestService()

		s.recipeIDFetcher = recipeIDFetcher
		s.userIDFetcher = userIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("CreateRecipeIteration", mock.Anything, mock.AnythingOfType("*models.RecipeIterationCreationInput")).Return((*models.RecipeIteration)(nil), errors.New("blah"))
		s.recipeIterationDataManager = recipeIterationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.recipeIDFetcher = recipeIDFetcher
		s.userIDFetcher = userIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("CreateRecipeIteration", mock.Anything, mock.AnythingOfType("*models.RecipeIterationCreationInput")).Return(exampleRecipeIteration, nil)
		s.recipeIterationDataManager = recipeIterationDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeIterationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIteration")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager, mc, r, ed)
	})
}

func TestRecipeIterationsService_ExistenceHandler(T *testing.T) {
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

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIteration.ID
		}

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("RecipeIterationExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(true, nil)
		s.recipeIterationDataManager = recipeIterationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager)
	})

	T.Run("with no such recipe iteration in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIteration.ID
		}

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("RecipeIterationExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(false, sql.ErrNoRows)
		s.recipeIterationDataManager = recipeIterationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager)
	})

	T.Run("with error fetching recipe iteration from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIteration.ID
		}

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("RecipeIterationExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(false, errors.New("blah"))
		s.recipeIterationDataManager = recipeIterationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager)
	})
}

func TestRecipeIterationsService_ReadHandler(T *testing.T) {
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

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIteration.ID
		}

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("GetRecipeIteration", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(exampleRecipeIteration, nil)
		s.recipeIterationDataManager = recipeIterationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIteration")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager, ed)
	})

	T.Run("with no such recipe iteration in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIteration.ID
		}

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("GetRecipeIteration", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return((*models.RecipeIteration)(nil), sql.ErrNoRows)
		s.recipeIterationDataManager = recipeIterationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager)
	})

	T.Run("with error fetching recipe iteration from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIteration.ID
		}

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("GetRecipeIteration", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return((*models.RecipeIteration)(nil), errors.New("blah"))
		s.recipeIterationDataManager = recipeIterationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIteration.ID
		}

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("GetRecipeIteration", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(exampleRecipeIteration, nil)
		s.recipeIterationDataManager = recipeIterationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIteration")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager, ed)
	})
}

func TestRecipeIterationsService_UpdateHandler(T *testing.T) {
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

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationUpdateInputFromRecipeIteration(exampleRecipeIteration)

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIteration.ID
		}

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("GetRecipeIteration", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(exampleRecipeIteration, nil)
		recipeIterationDataManager.On("UpdateRecipeIteration", mock.Anything, mock.AnythingOfType("*models.RecipeIteration")).Return(nil)
		s.recipeIterationDataManager = recipeIterationDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIteration")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, r, recipeIterationDataManager, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

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

	T.Run("with no rows fetching recipe iteration", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationUpdateInputFromRecipeIteration(exampleRecipeIteration)

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIteration.ID
		}

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("GetRecipeIteration", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return((*models.RecipeIteration)(nil), sql.ErrNoRows)
		s.recipeIterationDataManager = recipeIterationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager)
	})

	T.Run("with error fetching recipe iteration", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationUpdateInputFromRecipeIteration(exampleRecipeIteration)

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIteration.ID
		}

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("GetRecipeIteration", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return((*models.RecipeIteration)(nil), errors.New("blah"))
		s.recipeIterationDataManager = recipeIterationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager)
	})

	T.Run("with error updating recipe iteration", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationUpdateInputFromRecipeIteration(exampleRecipeIteration)

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIteration.ID
		}

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("GetRecipeIteration", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(exampleRecipeIteration, nil)
		recipeIterationDataManager.On("UpdateRecipeIteration", mock.Anything, mock.AnythingOfType("*models.RecipeIteration")).Return(errors.New("blah"))
		s.recipeIterationDataManager = recipeIterationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeIterationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationUpdateInputFromRecipeIteration(exampleRecipeIteration)

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIteration.ID
		}

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("GetRecipeIteration", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(exampleRecipeIteration, nil)
		recipeIterationDataManager.On("UpdateRecipeIteration", mock.Anything, mock.AnythingOfType("*models.RecipeIteration")).Return(nil)
		s.recipeIterationDataManager = recipeIterationDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIteration")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, r, recipeIterationDataManager, ed)
	})
}

func TestRecipeIterationsService_ArchiveHandler(T *testing.T) {
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

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIteration.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("ArchiveRecipeIteration", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(nil)
		s.recipeIterationDataManager = recipeIterationDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.recipeIterationCounter = mc

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationDataManager, mc, r)
	})

	T.Run("with nonexistent recipe", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

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

	T.Run("with no recipe iteration in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIteration.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("ArchiveRecipeIteration", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(sql.ErrNoRows)
		s.recipeIterationDataManager = recipeIterationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIteration.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("ArchiveRecipeIteration", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(errors.New("blah"))
		s.recipeIterationDataManager = recipeIterationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationDataManager)
	})
}
