package recipeiterationsteps

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

func TestRecipeIterationStepsService_ListHandler(T *testing.T) {
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

		exampleRecipeIterationStepList := fakemodels.BuildFakeRecipeIterationStepList()

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("GetRecipeIterationSteps", mock.Anything, exampleRecipe.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRecipeIterationStepList, nil)
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIterationStepList")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeIterationStepDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("GetRecipeIterationSteps", mock.Anything, exampleRecipe.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RecipeIterationStepList)(nil), sql.ErrNoRows)
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIterationStepList")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeIterationStepDataManager, ed)
	})

	T.Run("with error fetching recipe iteration steps from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("GetRecipeIterationSteps", mock.Anything, exampleRecipe.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RecipeIterationStepList)(nil), errors.New("blah"))
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeIterationStepDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationStepList := fakemodels.BuildFakeRecipeIterationStepList()

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("GetRecipeIterationSteps", mock.Anything, exampleRecipe.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRecipeIterationStepList, nil)
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIterationStepList")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeIterationStepDataManager, ed)
	})
}

func TestRecipeIterationStepsService_CreateHandler(T *testing.T) {
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

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("CreateRecipeIterationStep", mock.Anything, mock.AnythingOfType("*models.RecipeIterationStepCreationInput")).Return(exampleRecipeIterationStep, nil)
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeIterationStepCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIterationStep")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationStepDataManager, mc, r, ed)
	})

	T.Run("with nonexistent recipe", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(false, nil)
		s.recipeDataManager = recipeDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("with error checking recipe existence", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, errors.New("blah"))
		s.recipeDataManager = recipeDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error creating recipe iteration step", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("CreateRecipeIterationStep", mock.Anything, mock.AnythingOfType("*models.RecipeIterationStepCreationInput")).Return(exampleRecipeIterationStep, errors.New("blah"))
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationStepDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("CreateRecipeIterationStep", mock.Anything, mock.AnythingOfType("*models.RecipeIterationStepCreationInput")).Return(exampleRecipeIterationStep, nil)
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeIterationStepCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIterationStep")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationStepDataManager, mc, r, ed)
	})
}

func TestRecipeIterationStepsService_ExistenceHandler(T *testing.T) {
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

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationStepIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIterationStep.ID
		}

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("RecipeIterationStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return(true, nil)
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeIterationStepDataManager)
	})

	T.Run("with no such recipe iteration step in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationStepIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIterationStep.ID
		}

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("RecipeIterationStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return(false, sql.ErrNoRows)
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeIterationStepDataManager)
	})

	T.Run("with error fetching recipe iteration step from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationStepIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIterationStep.ID
		}

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("RecipeIterationStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return(false, errors.New("blah"))
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeIterationStepDataManager)
	})
}

func TestRecipeIterationStepsService_ReadHandler(T *testing.T) {
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

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationStepIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIterationStep.ID
		}

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("GetRecipeIterationStep", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return(exampleRecipeIterationStep, nil)
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIterationStep")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeIterationStepDataManager, ed)
	})

	T.Run("with no such recipe iteration step in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationStepIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIterationStep.ID
		}

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("GetRecipeIterationStep", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return((*models.RecipeIterationStep)(nil), sql.ErrNoRows)
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeIterationStepDataManager)
	})

	T.Run("with error fetching recipe iteration step from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationStepIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIterationStep.ID
		}

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("GetRecipeIterationStep", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return((*models.RecipeIterationStep)(nil), errors.New("blah"))
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeIterationStepDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationStepIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIterationStep.ID
		}

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("GetRecipeIterationStep", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return(exampleRecipeIterationStep, nil)
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIterationStep")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeIterationStepDataManager, ed)
	})
}

func TestRecipeIterationStepsService_UpdateHandler(T *testing.T) {
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

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationStepUpdateInputFromRecipeIterationStep(exampleRecipeIterationStep)

		s.recipeIterationStepIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIterationStep.ID
		}

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("GetRecipeIterationStep", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return(exampleRecipeIterationStep, nil)
		recipeIterationStepDataManager.On("UpdateRecipeIterationStep", mock.Anything, mock.AnythingOfType("*models.RecipeIterationStep")).Return(nil)
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIterationStep")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, r, recipeIterationStepDataManager, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with no rows fetching recipe iteration step", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationStepUpdateInputFromRecipeIterationStep(exampleRecipeIterationStep)

		s.recipeIterationStepIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIterationStep.ID
		}

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("GetRecipeIterationStep", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return((*models.RecipeIterationStep)(nil), sql.ErrNoRows)
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeIterationStepDataManager)
	})

	T.Run("with error fetching recipe iteration step", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationStepUpdateInputFromRecipeIterationStep(exampleRecipeIterationStep)

		s.recipeIterationStepIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIterationStep.ID
		}

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("GetRecipeIterationStep", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return((*models.RecipeIterationStep)(nil), errors.New("blah"))
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeIterationStepDataManager)
	})

	T.Run("with error updating recipe iteration step", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationStepUpdateInputFromRecipeIterationStep(exampleRecipeIterationStep)

		s.recipeIterationStepIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIterationStep.ID
		}

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("GetRecipeIterationStep", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return(exampleRecipeIterationStep, nil)
		recipeIterationStepDataManager.On("UpdateRecipeIterationStep", mock.Anything, mock.AnythingOfType("*models.RecipeIterationStep")).Return(errors.New("blah"))
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeIterationStepDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationStepUpdateInputFromRecipeIterationStep(exampleRecipeIterationStep)

		s.recipeIterationStepIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIterationStep.ID
		}

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("GetRecipeIterationStep", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return(exampleRecipeIterationStep, nil)
		recipeIterationStepDataManager.On("UpdateRecipeIterationStep", mock.Anything, mock.AnythingOfType("*models.RecipeIterationStep")).Return(nil)
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeIterationStep")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, r, recipeIterationStepDataManager, ed)
	})
}

func TestRecipeIterationStepsService_ArchiveHandler(T *testing.T) {
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

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationStepIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIterationStep.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("ArchiveRecipeIterationStep", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return(nil)
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.recipeIterationStepCounter = mc

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, http.StatusNoContent, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationStepDataManager, mc, r)
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
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

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
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("with no recipe iteration step in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationStepIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIterationStep.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("ArchiveRecipeIterationStep", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return(sql.ErrNoRows)
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationStepDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		s.recipeIterationStepIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeIterationStep.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationStepDataManager := &mockmodels.RecipeIterationStepDataManager{}
		recipeIterationStepDataManager.On("ArchiveRecipeIterationStep", mock.Anything, exampleRecipe.ID, exampleRecipeIterationStep.ID).Return(errors.New("blah"))
		s.recipeIterationStepDataManager = recipeIterationStepDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationStepDataManager)
	})
}
