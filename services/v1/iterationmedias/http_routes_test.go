package iterationmedias

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

func TestIterationMediasService_ListHandler(T *testing.T) {
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

	exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
	exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
	recipeIterationIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipeIteration.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMediaList := fakemodels.BuildFakeIterationMediaList()

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("GetIterationMedias", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleIterationMediaList, nil)
		s.iterationMediaDataManager = iterationMediaDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IterationMediaList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, iterationMediaDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("GetIterationMedias", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.IterationMediaList)(nil), sql.ErrNoRows)
		s.iterationMediaDataManager = iterationMediaDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IterationMediaList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, iterationMediaDataManager, ed)
	})

	T.Run("with error fetching iteration medias from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("GetIterationMedias", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.IterationMediaList)(nil), errors.New("blah"))
		s.iterationMediaDataManager = iterationMediaDataManager

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

		mock.AssertExpectationsForObjects(t, iterationMediaDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMediaList := fakemodels.BuildFakeIterationMediaList()

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("GetIterationMedias", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleIterationMediaList, nil)
		s.iterationMediaDataManager = iterationMediaDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IterationMediaList")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, iterationMediaDataManager, ed)
	})
}

func TestIterationMediasService_CreateHandler(T *testing.T) {
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

	exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
	exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
	recipeIterationIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipeIteration.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("RecipeIterationExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(true, nil)
		s.recipeIterationDataManager = recipeIterationDataManager

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("CreateIterationMedia", mock.Anything, mock.AnythingOfType("*models.IterationMediaCreationInput")).Return(exampleIterationMedia, nil)
		s.iterationMediaDataManager = iterationMediaDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.iterationMediaCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IterationMedia")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationDataManager, iterationMediaDataManager, mc, r, ed)
	})

	T.Run("with nonexistent recipe", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)

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
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)

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

	T.Run("with nonexistent recipe iteration", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("RecipeIterationExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(false, nil)
		s.recipeIterationDataManager = recipeIterationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationDataManager)
	})

	T.Run("with error checking recipe iteration existence", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("RecipeIterationExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(true, errors.New("blah"))
		s.recipeIterationDataManager = recipeIterationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationDataManager)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

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

	T.Run("with error creating iteration media", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("RecipeIterationExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(true, nil)
		s.recipeIterationDataManager = recipeIterationDataManager

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("CreateIterationMedia", mock.Anything, mock.AnythingOfType("*models.IterationMediaCreationInput")).Return(exampleIterationMedia, errors.New("blah"))
		s.iterationMediaDataManager = iterationMediaDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationDataManager, iterationMediaDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("RecipeIterationExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(true, nil)
		s.recipeIterationDataManager = recipeIterationDataManager

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("CreateIterationMedia", mock.Anything, mock.AnythingOfType("*models.IterationMediaCreationInput")).Return(exampleIterationMedia, nil)
		s.iterationMediaDataManager = iterationMediaDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.iterationMediaCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IterationMedia")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationDataManager, iterationMediaDataManager, mc, r, ed)
	})
}

func TestIterationMediasService_ExistenceHandler(T *testing.T) {
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

	exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
	exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
	recipeIterationIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipeIteration.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		s.iterationMediaIDFetcher = func(req *http.Request) uint64 {
			return exampleIterationMedia.ID
		}

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("IterationMediaExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return(true, nil)
		s.iterationMediaDataManager = iterationMediaDataManager

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

		mock.AssertExpectationsForObjects(t, iterationMediaDataManager)
	})

	T.Run("with no such iteration media in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		s.iterationMediaIDFetcher = func(req *http.Request) uint64 {
			return exampleIterationMedia.ID
		}

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("IterationMediaExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return(false, sql.ErrNoRows)
		s.iterationMediaDataManager = iterationMediaDataManager

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

		mock.AssertExpectationsForObjects(t, iterationMediaDataManager)
	})

	T.Run("with error fetching iteration media from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		s.iterationMediaIDFetcher = func(req *http.Request) uint64 {
			return exampleIterationMedia.ID
		}

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("IterationMediaExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return(false, errors.New("blah"))
		s.iterationMediaDataManager = iterationMediaDataManager

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

		mock.AssertExpectationsForObjects(t, iterationMediaDataManager)
	})
}

func TestIterationMediasService_ReadHandler(T *testing.T) {
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

	exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
	exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
	recipeIterationIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipeIteration.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		s.iterationMediaIDFetcher = func(req *http.Request) uint64 {
			return exampleIterationMedia.ID
		}

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("GetIterationMedia", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return(exampleIterationMedia, nil)
		s.iterationMediaDataManager = iterationMediaDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IterationMedia")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, iterationMediaDataManager, ed)
	})

	T.Run("with no such iteration media in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		s.iterationMediaIDFetcher = func(req *http.Request) uint64 {
			return exampleIterationMedia.ID
		}

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("GetIterationMedia", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return((*models.IterationMedia)(nil), sql.ErrNoRows)
		s.iterationMediaDataManager = iterationMediaDataManager

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

		mock.AssertExpectationsForObjects(t, iterationMediaDataManager)
	})

	T.Run("with error fetching iteration media from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		s.iterationMediaIDFetcher = func(req *http.Request) uint64 {
			return exampleIterationMedia.ID
		}

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("GetIterationMedia", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return((*models.IterationMedia)(nil), errors.New("blah"))
		s.iterationMediaDataManager = iterationMediaDataManager

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

		mock.AssertExpectationsForObjects(t, iterationMediaDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		s.iterationMediaIDFetcher = func(req *http.Request) uint64 {
			return exampleIterationMedia.ID
		}

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("GetIterationMedia", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return(exampleIterationMedia, nil)
		s.iterationMediaDataManager = iterationMediaDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IterationMedia")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, iterationMediaDataManager, ed)
	})
}

func TestIterationMediasService_UpdateHandler(T *testing.T) {
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

	exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
	exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
	recipeIterationIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipeIteration.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaUpdateInputFromIterationMedia(exampleIterationMedia)

		s.iterationMediaIDFetcher = func(req *http.Request) uint64 {
			return exampleIterationMedia.ID
		}

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("GetIterationMedia", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return(exampleIterationMedia, nil)
		iterationMediaDataManager.On("UpdateIterationMedia", mock.Anything, mock.AnythingOfType("*models.IterationMedia")).Return(nil)
		s.iterationMediaDataManager = iterationMediaDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IterationMedia")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, r, iterationMediaDataManager, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

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

	T.Run("with no rows fetching iteration media", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaUpdateInputFromIterationMedia(exampleIterationMedia)

		s.iterationMediaIDFetcher = func(req *http.Request) uint64 {
			return exampleIterationMedia.ID
		}

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("GetIterationMedia", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return((*models.IterationMedia)(nil), sql.ErrNoRows)
		s.iterationMediaDataManager = iterationMediaDataManager

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

		mock.AssertExpectationsForObjects(t, iterationMediaDataManager)
	})

	T.Run("with error fetching iteration media", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaUpdateInputFromIterationMedia(exampleIterationMedia)

		s.iterationMediaIDFetcher = func(req *http.Request) uint64 {
			return exampleIterationMedia.ID
		}

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("GetIterationMedia", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return((*models.IterationMedia)(nil), errors.New("blah"))
		s.iterationMediaDataManager = iterationMediaDataManager

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

		mock.AssertExpectationsForObjects(t, iterationMediaDataManager)
	})

	T.Run("with error updating iteration media", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaUpdateInputFromIterationMedia(exampleIterationMedia)

		s.iterationMediaIDFetcher = func(req *http.Request) uint64 {
			return exampleIterationMedia.ID
		}

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("GetIterationMedia", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return(exampleIterationMedia, nil)
		iterationMediaDataManager.On("UpdateIterationMedia", mock.Anything, mock.AnythingOfType("*models.IterationMedia")).Return(errors.New("blah"))
		s.iterationMediaDataManager = iterationMediaDataManager

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

		mock.AssertExpectationsForObjects(t, iterationMediaDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaUpdateInputFromIterationMedia(exampleIterationMedia)

		s.iterationMediaIDFetcher = func(req *http.Request) uint64 {
			return exampleIterationMedia.ID
		}

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("GetIterationMedia", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return(exampleIterationMedia, nil)
		iterationMediaDataManager.On("UpdateIterationMedia", mock.Anything, mock.AnythingOfType("*models.IterationMedia")).Return(nil)
		s.iterationMediaDataManager = iterationMediaDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IterationMedia")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, r, iterationMediaDataManager, ed)
	})
}

func TestIterationMediasService_ArchiveHandler(T *testing.T) {
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

	exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
	exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
	recipeIterationIDFetcher := func(_ *http.Request) uint64 {
		return exampleRecipeIteration.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		s.iterationMediaIDFetcher = func(req *http.Request) uint64 {
			return exampleIterationMedia.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("RecipeIterationExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(true, nil)
		s.recipeIterationDataManager = recipeIterationDataManager

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("ArchiveIterationMedia", mock.Anything, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return(nil)
		s.iterationMediaDataManager = iterationMediaDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.iterationMediaCounter = mc

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationDataManager, iterationMediaDataManager, mc, r)
	})

	T.Run("with nonexistent recipe", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

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
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

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

	T.Run("with nonexistent recipe iteration", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("RecipeIterationExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(false, nil)
		s.recipeIterationDataManager = recipeIterationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationDataManager)
	})

	T.Run("with error checking recipe iteration existence", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("RecipeIterationExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(true, errors.New("blah"))
		s.recipeIterationDataManager = recipeIterationDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationDataManager)
	})

	T.Run("with no iteration media in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		s.iterationMediaIDFetcher = func(req *http.Request) uint64 {
			return exampleIterationMedia.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("RecipeIterationExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(true, nil)
		s.recipeIterationDataManager = recipeIterationDataManager

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("ArchiveIterationMedia", mock.Anything, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return(sql.ErrNoRows)
		s.iterationMediaDataManager = iterationMediaDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationDataManager, iterationMediaDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher
		s.recipeIterationIDFetcher = recipeIterationIDFetcher

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		s.iterationMediaIDFetcher = func(req *http.Request) uint64 {
			return exampleIterationMedia.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeIterationDataManager := &mockmodels.RecipeIterationDataManager{}
		recipeIterationDataManager.On("RecipeIterationExists", mock.Anything, exampleRecipe.ID, exampleRecipeIteration.ID).Return(true, nil)
		s.recipeIterationDataManager = recipeIterationDataManager

		iterationMediaDataManager := &mockmodels.IterationMediaDataManager{}
		iterationMediaDataManager.On("ArchiveIterationMedia", mock.Anything, exampleRecipeIteration.ID, exampleIterationMedia.ID).Return(errors.New("blah"))
		s.iterationMediaDataManager = iterationMediaDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeIterationDataManager, iterationMediaDataManager)
	})
}
