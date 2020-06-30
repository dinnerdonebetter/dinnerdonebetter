package recipetags

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

func TestRecipeTagsService_ListHandler(T *testing.T) {
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

		exampleRecipeTagList := fakemodels.BuildFakeRecipeTagList()

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("GetRecipeTags", mock.Anything, exampleRecipe.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRecipeTagList, nil)
		s.recipeTagDataManager = recipeTagDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeTagList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeTagDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("GetRecipeTags", mock.Anything, exampleRecipe.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RecipeTagList)(nil), sql.ErrNoRows)
		s.recipeTagDataManager = recipeTagDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeTagList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeTagDataManager, ed)
	})

	T.Run("with error fetching recipe tags from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("GetRecipeTags", mock.Anything, exampleRecipe.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RecipeTagList)(nil), errors.New("blah"))
		s.recipeTagDataManager = recipeTagDataManager

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

		mock.AssertExpectationsForObjects(t, recipeTagDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeTagList := fakemodels.BuildFakeRecipeTagList()

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("GetRecipeTags", mock.Anything, exampleRecipe.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRecipeTagList, nil)
		s.recipeTagDataManager = recipeTagDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeTagList")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, recipeTagDataManager, ed)
	})
}

func TestRecipeTagsService_CreateHandler(T *testing.T) {
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

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("CreateRecipeTag", mock.Anything, mock.AnythingOfType("*models.RecipeTagCreationInput")).Return(exampleRecipeTag, nil)
		s.recipeTagDataManager = recipeTagDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeTagCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeTag")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeTagDataManager, mc, r, ed)
	})

	T.Run("with nonexistent recipe", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)

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

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)

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

	T.Run("with error creating recipe tag", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("CreateRecipeTag", mock.Anything, mock.AnythingOfType("*models.RecipeTagCreationInput")).Return(exampleRecipeTag, errors.New("blah"))
		s.recipeTagDataManager = recipeTagDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeTagDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("CreateRecipeTag", mock.Anything, mock.AnythingOfType("*models.RecipeTagCreationInput")).Return(exampleRecipeTag, nil)
		s.recipeTagDataManager = recipeTagDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeTagCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeTag")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeTagDataManager, mc, r, ed)
	})
}

func TestRecipeTagsService_ExistenceHandler(T *testing.T) {
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

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		s.recipeTagIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeTag.ID
		}

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("RecipeTagExists", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return(true, nil)
		s.recipeTagDataManager = recipeTagDataManager

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

		mock.AssertExpectationsForObjects(t, recipeTagDataManager)
	})

	T.Run("with no such recipe tag in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		s.recipeTagIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeTag.ID
		}

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("RecipeTagExists", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return(false, sql.ErrNoRows)
		s.recipeTagDataManager = recipeTagDataManager

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

		mock.AssertExpectationsForObjects(t, recipeTagDataManager)
	})

	T.Run("with error fetching recipe tag from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		s.recipeTagIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeTag.ID
		}

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("RecipeTagExists", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return(false, errors.New("blah"))
		s.recipeTagDataManager = recipeTagDataManager

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

		mock.AssertExpectationsForObjects(t, recipeTagDataManager)
	})
}

func TestRecipeTagsService_ReadHandler(T *testing.T) {
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

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		s.recipeTagIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeTag.ID
		}

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("GetRecipeTag", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return(exampleRecipeTag, nil)
		s.recipeTagDataManager = recipeTagDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeTag")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, recipeTagDataManager, ed)
	})

	T.Run("with no such recipe tag in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		s.recipeTagIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeTag.ID
		}

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("GetRecipeTag", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return((*models.RecipeTag)(nil), sql.ErrNoRows)
		s.recipeTagDataManager = recipeTagDataManager

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

		mock.AssertExpectationsForObjects(t, recipeTagDataManager)
	})

	T.Run("with error fetching recipe tag from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		s.recipeTagIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeTag.ID
		}

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("GetRecipeTag", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return((*models.RecipeTag)(nil), errors.New("blah"))
		s.recipeTagDataManager = recipeTagDataManager

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

		mock.AssertExpectationsForObjects(t, recipeTagDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		s.recipeTagIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeTag.ID
		}

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("GetRecipeTag", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return(exampleRecipeTag, nil)
		s.recipeTagDataManager = recipeTagDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeTag")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, recipeTagDataManager, ed)
	})
}

func TestRecipeTagsService_UpdateHandler(T *testing.T) {
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

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeTagUpdateInputFromRecipeTag(exampleRecipeTag)

		s.recipeTagIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeTag.ID
		}

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("GetRecipeTag", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return(exampleRecipeTag, nil)
		recipeTagDataManager.On("UpdateRecipeTag", mock.Anything, mock.AnythingOfType("*models.RecipeTag")).Return(nil)
		s.recipeTagDataManager = recipeTagDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeTag")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, r, recipeTagDataManager, ed)
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

	T.Run("with no rows fetching recipe tag", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeTagUpdateInputFromRecipeTag(exampleRecipeTag)

		s.recipeTagIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeTag.ID
		}

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("GetRecipeTag", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return((*models.RecipeTag)(nil), sql.ErrNoRows)
		s.recipeTagDataManager = recipeTagDataManager

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

		mock.AssertExpectationsForObjects(t, recipeTagDataManager)
	})

	T.Run("with error fetching recipe tag", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeTagUpdateInputFromRecipeTag(exampleRecipeTag)

		s.recipeTagIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeTag.ID
		}

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("GetRecipeTag", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return((*models.RecipeTag)(nil), errors.New("blah"))
		s.recipeTagDataManager = recipeTagDataManager

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

		mock.AssertExpectationsForObjects(t, recipeTagDataManager)
	})

	T.Run("with error updating recipe tag", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeTagUpdateInputFromRecipeTag(exampleRecipeTag)

		s.recipeTagIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeTag.ID
		}

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("GetRecipeTag", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return(exampleRecipeTag, nil)
		recipeTagDataManager.On("UpdateRecipeTag", mock.Anything, mock.AnythingOfType("*models.RecipeTag")).Return(errors.New("blah"))
		s.recipeTagDataManager = recipeTagDataManager

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

		mock.AssertExpectationsForObjects(t, recipeTagDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeTagUpdateInputFromRecipeTag(exampleRecipeTag)

		s.recipeTagIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeTag.ID
		}

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("GetRecipeTag", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return(exampleRecipeTag, nil)
		recipeTagDataManager.On("UpdateRecipeTag", mock.Anything, mock.AnythingOfType("*models.RecipeTag")).Return(nil)
		s.recipeTagDataManager = recipeTagDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RecipeTag")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, r, recipeTagDataManager, ed)
	})
}

func TestRecipeTagsService_ArchiveHandler(T *testing.T) {
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

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		s.recipeTagIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeTag.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("ArchiveRecipeTag", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return(nil)
		s.recipeTagDataManager = recipeTagDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.recipeTagCounter = mc

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeTagDataManager, mc, r)
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

	T.Run("with no recipe tag in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		s.recipeTagIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeTag.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("ArchiveRecipeTag", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return(sql.ErrNoRows)
		s.recipeTagDataManager = recipeTagDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeTagDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher
		s.recipeIDFetcher = recipeIDFetcher

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		s.recipeTagIDFetcher = func(req *http.Request) uint64 {
			return exampleRecipeTag.ID
		}

		recipeDataManager := &mockmodels.RecipeDataManager{}
		recipeDataManager.On("RecipeExists", mock.Anything, exampleRecipe.ID).Return(true, nil)
		s.recipeDataManager = recipeDataManager

		recipeTagDataManager := &mockmodels.RecipeTagDataManager{}
		recipeTagDataManager.On("ArchiveRecipeTag", mock.Anything, exampleRecipe.ID, exampleRecipeTag.ID).Return(errors.New("blah"))
		s.recipeTagDataManager = recipeTagDataManager

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

		mock.AssertExpectationsForObjects(t, recipeDataManager, recipeTagDataManager)
	})
}
