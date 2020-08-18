package recipestepinstruments

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

func TestRecipeStepInstrumentsService_ListHandler(T *testing.T) {
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

		exampleRecipeStepInstrumentList := fakemodels.BuildFakeRecipeStepInstrumentList()

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("GetRecipeStepInstruments", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRecipeStepInstrumentList, nil)
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepInstrumentList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("GetRecipeStepInstruments", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RecipeStepInstrumentList)(nil), sql.ErrNoRows)
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepInstrumentList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager, ed)
	})

	T.Run("with error fetching recipe step instruments from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("GetRecipeStepInstruments", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RecipeStepInstrumentList)(nil), errors.New("blah"))
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepInstrumentList := fakemodels.BuildFakeRecipeStepInstrumentList()

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("GetRecipeStepInstruments", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRecipeStepInstrumentList, nil)
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepInstrumentList")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager, ed)
	})
}

func TestRecipeStepInstrumentsService_CreateHandler(T *testing.T) {
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

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("CreateRecipeStepInstrument", mock.Anything, mock.AnythingOfType("*models.RecipeStepInstrumentCreationInput")).Return(exampleRecipeStepInstrument, nil)
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepInstrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepInstrument")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager, mc, r, ed)
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

	T.Run("with error creating recipe step instrument", func(t *testing.T) {
		s := buildTestService()

		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher
		s.userIDFetcher = userIDFetcher

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("CreateRecipeStepInstrument", mock.Anything, mock.AnythingOfType("*models.RecipeStepInstrumentCreationInput")).Return((*models.RecipeStepInstrument)(nil), errors.New("blah"))
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher
		s.userIDFetcher = userIDFetcher

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("CreateRecipeStepInstrument", mock.Anything, mock.AnythingOfType("*models.RecipeStepInstrumentCreationInput")).Return(exampleRecipeStepInstrument, nil)
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepInstrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepInstrument")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager, mc, r, ed)
	})
}

func TestRecipeStepInstrumentsService_ExistenceHandler(T *testing.T) {
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

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepInstrument.ID
		}

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("RecipeStepInstrumentExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return(true, nil)
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager)
	})

	T.Run("with no such recipe step instrument in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepInstrument.ID
		}

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("RecipeStepInstrumentExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return(false, sql.ErrNoRows)
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager)
	})

	T.Run("with error fetching recipe step instrument from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepInstrument.ID
		}

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("RecipeStepInstrumentExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return(false, errors.New("blah"))
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager)
	})
}

func TestRecipeStepInstrumentsService_ReadHandler(T *testing.T) {
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

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepInstrument.ID
		}

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("GetRecipeStepInstrument", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return(exampleRecipeStepInstrument, nil)
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepInstrument")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager, ed)
	})

	T.Run("with no such recipe step instrument in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepInstrument.ID
		}

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("GetRecipeStepInstrument", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return((*models.RecipeStepInstrument)(nil), sql.ErrNoRows)
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager)
	})

	T.Run("with error fetching recipe step instrument from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepInstrument.ID
		}

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("GetRecipeStepInstrument", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return((*models.RecipeStepInstrument)(nil), errors.New("blah"))
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepInstrument.ID
		}

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("GetRecipeStepInstrument", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return(exampleRecipeStepInstrument, nil)
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepInstrument")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager, ed)
	})
}

func TestRecipeStepInstrumentsService_UpdateHandler(T *testing.T) {
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

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepInstrumentUpdateInputFromRecipeStepInstrument(exampleRecipeStepInstrument)

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepInstrument.ID
		}

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("GetRecipeStepInstrument", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return(exampleRecipeStepInstrument, nil)
		recipeStepInstrumentDataManager.On("UpdateRecipeStepInstrument", mock.Anything, mock.AnythingOfType("*models.RecipeStepInstrument")).Return(nil)
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepInstrument")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, r, recipeStepInstrumentDataManager, ed)
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

	T.Run("with no rows fetching recipe step instrument", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepInstrumentUpdateInputFromRecipeStepInstrument(exampleRecipeStepInstrument)

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepInstrument.ID
		}

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("GetRecipeStepInstrument", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return((*models.RecipeStepInstrument)(nil), sql.ErrNoRows)
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager)
	})

	T.Run("with error fetching recipe step instrument", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepInstrumentUpdateInputFromRecipeStepInstrument(exampleRecipeStepInstrument)

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepInstrument.ID
		}

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("GetRecipeStepInstrument", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return((*models.RecipeStepInstrument)(nil), errors.New("blah"))
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager)
	})

	T.Run("with error updating recipe step instrument", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepInstrumentUpdateInputFromRecipeStepInstrument(exampleRecipeStepInstrument)

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepInstrument.ID
		}

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("GetRecipeStepInstrument", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return(exampleRecipeStepInstrument, nil)
		recipeStepInstrumentDataManager.On("UpdateRecipeStepInstrument", mock.Anything, mock.AnythingOfType("*models.RecipeStepInstrument")).Return(errors.New("blah"))
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepInstrumentUpdateInputFromRecipeStepInstrument(exampleRecipeStepInstrument)

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepInstrument.ID
		}

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("GetRecipeStepInstrument", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return(exampleRecipeStepInstrument, nil)
		recipeStepInstrumentDataManager.On("UpdateRecipeStepInstrument", mock.Anything, mock.AnythingOfType("*models.RecipeStepInstrument")).Return(nil)
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepInstrument")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, r, recipeStepInstrumentDataManager, ed)
	})
}

func TestRecipeStepInstrumentsService_ArchiveHandler(T *testing.T) {
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

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepInstrument.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("ArchiveRecipeStepInstrument", mock.Anything, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return(nil)
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.recipeStepInstrumentCounter = mc

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepInstrumentDataManager, mc, r)
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

	T.Run("with no recipe step instrument in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepInstrument.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("ArchiveRecipeStepInstrument", mock.Anything, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return(sql.ErrNoRows)
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepInstrumentDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepInstrument.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepInstrumentDataManager := &mockmodels.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On("ArchiveRecipeStepInstrument", mock.Anything, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID).Return(errors.New("blah"))
		s.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepInstrumentDataManager)
	})
}
