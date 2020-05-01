package recipestepingredients

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

func TestRecipeStepIngredientsService_ListHandler(T *testing.T) {
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

		exampleRecipeStepIngredientList := fakemodels.BuildFakeRecipeStepIngredientList()

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("GetRecipeStepIngredients", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRecipeStepIngredientList, nil)
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepIngredientList")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepIngredientDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("GetRecipeStepIngredients", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RecipeStepIngredientList)(nil), sql.ErrNoRows)
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepIngredientList")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepIngredientDataManager, ed)
	})

	T.Run("with error fetching recipe step ingredients from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("GetRecipeStepIngredients", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RecipeStepIngredientList)(nil), errors.New("blah"))
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepIngredientDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepIngredientList := fakemodels.BuildFakeRecipeStepIngredientList()

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("GetRecipeStepIngredients", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRecipeStepIngredientList, nil)
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepIngredientList")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepIngredientDataManager, ed)
	})
}

func TestRecipeStepIngredientsService_CreateHandler(T *testing.T) {
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

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("CreateRecipeStepIngredient", mock.Anything, mock.AnythingOfType("*models.RecipeStepIngredientCreationInput")).Return(exampleRecipeStepIngredient, nil)
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepIngredientCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepIngredient")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepIngredientDataManager, mc, r, ed)
	})

	T.Run("with nonexistent recipe", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

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

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

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

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

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

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

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
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error creating recipe step ingredient", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("CreateRecipeStepIngredient", mock.Anything, mock.AnythingOfType("*models.RecipeStepIngredientCreationInput")).Return(exampleRecipeStepIngredient, errors.New("blah"))
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepIngredientDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("CreateRecipeStepIngredient", mock.Anything, mock.AnythingOfType("*models.RecipeStepIngredientCreationInput")).Return(exampleRecipeStepIngredient, nil)
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepIngredientCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepIngredient")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepIngredientDataManager, mc, r, ed)
	})
}

func TestRecipeStepIngredientsService_ExistenceHandler(T *testing.T) {
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

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("RecipeStepIngredientExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return(true, nil)
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepIngredientDataManager)
	})

	T.Run("with no such recipe step ingredient in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("RecipeStepIngredientExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return(false, sql.ErrNoRows)
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepIngredientDataManager)
	})

	T.Run("with error fetching recipe step ingredient from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("RecipeStepIngredientExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return(false, errors.New("blah"))
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepIngredientDataManager)
	})
}

func TestRecipeStepIngredientsService_ReadHandler(T *testing.T) {
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

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("GetRecipeStepIngredient", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return(exampleRecipeStepIngredient, nil)
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepIngredient")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepIngredientDataManager, ed)
	})

	T.Run("with no such recipe step ingredient in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("GetRecipeStepIngredient", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return((*models.RecipeStepIngredient)(nil), sql.ErrNoRows)
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepIngredientDataManager)
	})

	T.Run("with error fetching recipe step ingredient from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("GetRecipeStepIngredient", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return((*models.RecipeStepIngredient)(nil), errors.New("blah"))
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepIngredientDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("GetRecipeStepIngredient", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return(exampleRecipeStepIngredient, nil)
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepIngredient")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepIngredientDataManager, ed)
	})
}

func TestRecipeStepIngredientsService_UpdateHandler(T *testing.T) {
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

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientUpdateInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("GetRecipeStepIngredient", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return(exampleRecipeStepIngredient, nil)
		recipeStepIngredientDataManager.On("UpdateRecipeStepIngredient", mock.Anything, mock.AnythingOfType("*models.RecipeStepIngredient")).Return(nil)
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepIngredient")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, r, recipeStepIngredientDataManager, ed)
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

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with no rows fetching recipe step ingredient", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientUpdateInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("GetRecipeStepIngredient", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return((*models.RecipeStepIngredient)(nil), sql.ErrNoRows)
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepIngredientDataManager)
	})

	T.Run("with error fetching recipe step ingredient", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientUpdateInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("GetRecipeStepIngredient", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return((*models.RecipeStepIngredient)(nil), errors.New("blah"))
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepIngredientDataManager)
	})

	T.Run("with error updating recipe step ingredient", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientUpdateInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("GetRecipeStepIngredient", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return(exampleRecipeStepIngredient, nil)
		recipeStepIngredientDataManager.On("UpdateRecipeStepIngredient", mock.Anything, mock.AnythingOfType("*models.RecipeStepIngredient")).Return(errors.New("blah"))
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepIngredientDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientUpdateInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("GetRecipeStepIngredient", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return(exampleRecipeStepIngredient, nil)
		recipeStepIngredientDataManager.On("UpdateRecipeStepIngredient", mock.Anything, mock.AnythingOfType("*models.RecipeStepIngredient")).Return(nil)
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeStepIngredient")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, r, recipeStepIngredientDataManager, ed)
	})
}

func TestRecipeStepIngredientsService_ArchiveHandler(T *testing.T) {
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

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("ArchiveRecipeStepIngredient", mock.Anything, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return(nil)
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.recipeStepIngredientCounter = mc

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, http.StatusNoContent, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepIngredientDataManager, mc, r)
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
			"http://todo.verygoodsoftwarenotvirus.ru",
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
			"http://todo.verygoodsoftwarenotvirus.ru",
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
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager)
	})

	T.Run("with no recipe step ingredient in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("ArchiveRecipeStepIngredient", mock.Anything, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return(sql.ErrNoRows)
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepIngredientDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeStepIDFetcher = recipeStepIDFetcher

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeStepIngredient.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeStepDataManager := &mockmodels.RecipeStepDataManager{}
		recipeStepDataManager.On("RecipeStepExists", mock.Anything, exampleRecipe.ID, exampleRecipeStep.ID).Return(true, nil)
		s.recipeStepDataManager = recipeStepDataManager

		recipeStepIngredientDataManager := &mockmodels.RecipeStepIngredientDataManager{}
		recipeStepIngredientDataManager.On("ArchiveRecipeStepIngredient", mock.Anything, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID).Return(errors.New("blah"))
		s.recipeStepIngredientDataManager = recipeStepIngredientDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeStepDataManager, recipeStepIngredientDataManager)
	})
}
