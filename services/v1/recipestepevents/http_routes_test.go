package recipestepevents

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

func TestRecipeStepEventsService_ListHandler(T *testing.T) {
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

		exampleRecipeStepEventList := fakemodels.BuildFakeRecipeStepEventList()

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("GetRecipeStepEvents", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRecipeStepEventList, nil)
		s.recipeStepEventDataManager = recipeStepEventDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepEventList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("GetRecipeStepEvents", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RecipeStepEventList)(nil), sql.ErrNoRows)
		s.recipeStepEventDataManager = recipeStepEventDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepEventList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager, ed)
	})

	T.Run("with error fetching recipe step events from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("GetRecipeStepEvents", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RecipeStepEventList)(nil), errors.New("blah"))
		s.recipeStepEventDataManager = recipeStepEventDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepEventList := fakemodels.BuildFakeRecipeStepEventList()

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("GetRecipeStepEvents", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRecipeStepEventList, nil)
		s.recipeStepEventDataManager = recipeStepEventDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepEventList")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager, ed)
	})
}

func TestRecipeStepEventsService_CreateHandler(T *testing.T) {
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

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("CreateRecipeStepEvent", mock.Anything, mock.AnythingOfType("*models.RecipeStepEventCreationInput")).Return(exampleRecipeStepEvent, nil)
		s.recipeStepEventDataManager = recipeStepEventDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepEventCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepEvent")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager, mc, r, ed)
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

	T.Run("with error creating recipe step event", func(t *testing.T) {
		s := buildTestService()

		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher
		s.userIDFetcher = userIDFetcher

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("CreateRecipeStepEvent", mock.Anything, mock.AnythingOfType("*models.RecipeStepEventCreationInput")).Return((*models.RecipeStepEvent)(nil), errors.New("blah"))
		s.recipeStepEventDataManager = recipeStepEventDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher
		s.userIDFetcher = userIDFetcher

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("CreateRecipeStepEvent", mock.Anything, mock.AnythingOfType("*models.RecipeStepEventCreationInput")).Return(exampleRecipeStepEvent, nil)
		s.recipeStepEventDataManager = recipeStepEventDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepEventCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepEvent")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager, mc, r, ed)
	})
}

func TestRecipeStepEventsService_ExistenceHandler(T *testing.T) {
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

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepEventIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepEvent.ID
		}

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("RecipeStepEventExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return(true, nil)
		s.recipeStepEventDataManager = recipeStepEventDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager)
	})

	T.Run("with no such recipe step event in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepEventIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepEvent.ID
		}

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("RecipeStepEventExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return(false, sql.ErrNoRows)
		s.recipeStepEventDataManager = recipeStepEventDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager)
	})

	T.Run("with error fetching recipe step event from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepEventIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepEvent.ID
		}

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("RecipeStepEventExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return(false, errors.New("blah"))
		s.recipeStepEventDataManager = recipeStepEventDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager)
	})
}

func TestRecipeStepEventsService_ReadHandler(T *testing.T) {
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

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepEventIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepEvent.ID
		}

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("GetRecipeStepEvent", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return(exampleRecipeStepEvent, nil)
		s.recipeStepEventDataManager = recipeStepEventDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepEvent")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager, ed)
	})

	T.Run("with no such recipe step event in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepEventIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepEvent.ID
		}

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("GetRecipeStepEvent", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return((*models.RecipeStepEvent)(nil), sql.ErrNoRows)
		s.recipeStepEventDataManager = recipeStepEventDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager)
	})

	T.Run("with error fetching recipe step event from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepEventIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepEvent.ID
		}

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("GetRecipeStepEvent", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return((*models.RecipeStepEvent)(nil), errors.New("blah"))
		s.recipeStepEventDataManager = recipeStepEventDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepEventIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepEvent.ID
		}

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("GetRecipeStepEvent", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return(exampleRecipeStepEvent, nil)
		s.recipeStepEventDataManager = recipeStepEventDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepEvent")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager, ed)
	})
}

func TestRecipeStepEventsService_UpdateHandler(T *testing.T) {
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

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepEventUpdateInputFromRecipeStepEvent(exampleRecipeStepEvent)

		s.recipeStepEventIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepEvent.ID
		}

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("GetRecipeStepEvent", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return(exampleRecipeStepEvent, nil)
		recipeStepEventDataManager.On("UpdateRecipeStepEvent", mock.Anything, mock.AnythingOfType("*models.RecipeStepEvent")).Return(nil)
		s.recipeStepEventDataManager = recipeStepEventDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepEvent")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, r, recipeStepEventDataManager, ed)
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

	T.Run("with no rows fetching recipe step event", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepEventUpdateInputFromRecipeStepEvent(exampleRecipeStepEvent)

		s.recipeStepEventIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepEvent.ID
		}

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("GetRecipeStepEvent", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return((*models.RecipeStepEvent)(nil), sql.ErrNoRows)
		s.recipeStepEventDataManager = recipeStepEventDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager)
	})

	T.Run("with error fetching recipe step event", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepEventUpdateInputFromRecipeStepEvent(exampleRecipeStepEvent)

		s.recipeStepEventIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepEvent.ID
		}

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("GetRecipeStepEvent", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return((*models.RecipeStepEvent)(nil), errors.New("blah"))
		s.recipeStepEventDataManager = recipeStepEventDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager)
	})

	T.Run("with error updating recipe step event", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepEventUpdateInputFromRecipeStepEvent(exampleRecipeStepEvent)

		s.recipeStepEventIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepEvent.ID
		}

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("GetRecipeStepEvent", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return(exampleRecipeStepEvent, nil)
		recipeStepEventDataManager.On("UpdateRecipeStepEvent", mock.Anything, mock.AnythingOfType("*models.RecipeStepEvent")).Return(errors.New("blah"))
		s.recipeStepEventDataManager = recipeStepEventDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepEventDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepEventUpdateInputFromRecipeStepEvent(exampleRecipeStepEvent)

		s.recipeStepEventIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepEvent.ID
		}

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("GetRecipeStepEvent", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return(exampleRecipeStepEvent, nil)
		recipeStepEventDataManager.On("UpdateRecipeStepEvent", mock.Anything, mock.AnythingOfType("*models.RecipeStepEvent")).Return(nil)
		s.recipeStepEventDataManager = recipeStepEventDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepEvent")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, r, recipeStepEventDataManager, ed)
	})
}

func TestRecipeStepEventsService_ArchiveHandler(T *testing.T) {
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

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepEventIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepEvent.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("ArchiveRecipeStepEvent", mock.Anything, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return(nil)
		s.recipeStepEventDataManager = recipeStepEventDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.recipeStepEventCounter = mc

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepEventDataManager, mc, r)
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

	T.Run("with no recipe step event in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepEventIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepEvent.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("ArchiveRecipeStepEvent", mock.Anything, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return(sql.ErrNoRows)
		s.recipeStepEventDataManager = recipeStepEventDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepEventDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepEventIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepEvent.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepEventDataManager := &mockmodels.RecipeStepEventDataManager{}
		recipeStepEventDataManager.On("ArchiveRecipeStepEvent", mock.Anything, exampleRecipeStep.ID, exampleRecipeStepEvent.ID).Return(errors.New("blah"))
		s.recipeStepEventDataManager = recipeStepEventDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepEventDataManager)
	})
}
