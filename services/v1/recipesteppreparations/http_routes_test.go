package recipesteppreparations

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

func TestRecipeStepPreparationsService_ListHandler(T *testing.T) {
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

		exampleRecipeStepPreparationList := fakemodels.BuildFakeRecipeStepPreparationList()

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("GetRecipeStepPreparations", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRecipeStepPreparationList, nil)
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepPreparationList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeStepPreparationDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("GetRecipeStepPreparations", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RecipeStepPreparationList)(nil), sql.ErrNoRows)
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepPreparationList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeStepPreparationDataManager, ed)
	})

	T.Run("with error fetching recipe step preparations from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("GetRecipeStepPreparations", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RecipeStepPreparationList)(nil), errors.New("blah"))
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepPreparationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparationList := fakemodels.BuildFakeRecipeStepPreparationList()

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("GetRecipeStepPreparations", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRecipeStepPreparationList, nil)
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepPreparationList")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, recipeStepPreparationDataManager, ed)
	})
}

func TestRecipeStepPreparationsService_CreateHandler(T *testing.T) {
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

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("CreateRecipeStepPreparation", mock.Anything, mock.AnythingOfType("*models.RecipeStepPreparationCreationInput")).Return(exampleRecipeStepPreparation, nil)
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepPreparationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepPreparation")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepPreparationDataManager, mc, r, ed)
	})

	T.Run("with nonexistent recipe", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)

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
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)

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

	T.Run("with nonexistent recipe step", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(false, nil)
		s.recipeStepDataManager = recipeStepDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager)
	})

	T.Run("with error checking recipe step existence", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, errors.New("blah"))
		s.recipeStepDataManager = recipeStepDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

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

	T.Run("with error creating recipe step preparation", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("CreateRecipeStepPreparation", mock.Anything, mock.AnythingOfType("*models.RecipeStepPreparationCreationInput")).Return(exampleRecipeStepPreparation, errors.New("blah"))
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepPreparationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("CreateRecipeStepPreparation", mock.Anything, mock.AnythingOfType("*models.RecipeStepPreparationCreationInput")).Return(exampleRecipeStepPreparation, nil)
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepPreparationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepPreparation")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepPreparationDataManager, mc, r, ed)
	})
}

func TestRecipeStepPreparationsService_ExistenceHandler(T *testing.T) {
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

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepPreparation.ID
		}

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("RecipeStepPreparationExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return(true, nil)
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepPreparationDataManager)
	})

	T.Run("with no such recipe step preparation in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepPreparation.ID
		}

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("RecipeStepPreparationExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return(false, sql.ErrNoRows)
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepPreparationDataManager)
	})

	T.Run("with error fetching recipe step preparation from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepPreparation.ID
		}

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("RecipeStepPreparationExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return(false, errors.New("blah"))
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepPreparationDataManager)
	})
}

func TestRecipeStepPreparationsService_ReadHandler(T *testing.T) {
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

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepPreparation.ID
		}

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("GetRecipeStepPreparation", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return(exampleRecipeStepPreparation, nil)
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepPreparation")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeStepPreparationDataManager, ed)
	})

	T.Run("with no such recipe step preparation in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepPreparation.ID
		}

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("GetRecipeStepPreparation", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return((*models.RecipeStepPreparation)(nil), sql.ErrNoRows)
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepPreparationDataManager)
	})

	T.Run("with error fetching recipe step preparation from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepPreparation.ID
		}

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("GetRecipeStepPreparation", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return((*models.RecipeStepPreparation)(nil), errors.New("blah"))
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepPreparationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepPreparation.ID
		}

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("GetRecipeStepPreparation", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return(exampleRecipeStepPreparation, nil)
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepPreparation")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, recipeStepPreparationDataManager, ed)
	})
}

func TestRecipeStepPreparationsService_UpdateHandler(T *testing.T) {
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

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepPreparationUpdateInputFromRecipeStepPreparation(exampleRecipeStepPreparation)

		s.recipeStepPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepPreparation.ID
		}

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("GetRecipeStepPreparation", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return(exampleRecipeStepPreparation, nil)
		recipeStepPreparationDataManager.On("UpdateRecipeStepPreparation", mock.Anything, mock.AnythingOfType("*models.RecipeStepPreparation")).Return(nil)
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepPreparation")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, r, recipeStepPreparationDataManager, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

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

	T.Run("with no rows fetching recipe step preparation", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepPreparationUpdateInputFromRecipeStepPreparation(exampleRecipeStepPreparation)

		s.recipeStepPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepPreparation.ID
		}

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("GetRecipeStepPreparation", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return((*models.RecipeStepPreparation)(nil), sql.ErrNoRows)
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepPreparationDataManager)
	})

	T.Run("with error fetching recipe step preparation", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepPreparationUpdateInputFromRecipeStepPreparation(exampleRecipeStepPreparation)

		s.recipeStepPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepPreparation.ID
		}

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("GetRecipeStepPreparation", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return((*models.RecipeStepPreparation)(nil), errors.New("blah"))
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepPreparationDataManager)
	})

	T.Run("with error updating recipe step preparation", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepPreparationUpdateInputFromRecipeStepPreparation(exampleRecipeStepPreparation)

		s.recipeStepPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepPreparation.ID
		}

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("GetRecipeStepPreparation", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return(exampleRecipeStepPreparation, nil)
		recipeStepPreparationDataManager.On("UpdateRecipeStepPreparation", mock.Anything, mock.AnythingOfType("*models.RecipeStepPreparation")).Return(errors.New("blah"))
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeStepPreparationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepPreparationUpdateInputFromRecipeStepPreparation(exampleRecipeStepPreparation)

		s.recipeStepPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepPreparation.ID
		}

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("GetRecipeStepPreparation", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return(exampleRecipeStepPreparation, nil)
		recipeStepPreparationDataManager.On("UpdateRecipeStepPreparation", mock.Anything, mock.AnythingOfType("*models.RecipeStepPreparation")).Return(nil)
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepPreparation")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, r, recipeStepPreparationDataManager, ed)
	})
}

func TestRecipeStepPreparationsService_ArchiveHandler(T *testing.T) {
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

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepPreparation.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("ArchiveRecipeStepPreparation", mock.Anything, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return(nil)
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.recipeStepPreparationCounter = mc

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepPreparationDataManager, mc, r)
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
		s.recipeStepIDFetcher = recipeStepIDFetcher

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
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

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
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager)
	})

	T.Run("with no recipe step preparation in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepPreparation.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("ArchiveRecipeStepPreparation", mock.Anything, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return(sql.ErrNoRows)
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepPreparationDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepPreparation.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepPreparationDataManager := &mockmodels.RecipeStepPreparationDataManager{}
		recipeStepPreparationDataManager.On("ArchiveRecipeStepPreparation", mock.Anything, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID).Return(errors.New("blah"))
		s.recipeStepPreparationDataManager = recipeStepPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepPreparationDataManager)
	})
}
