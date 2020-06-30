package validingredienttags

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

func TestValidIngredientTagsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTagList := fakemodels.BuildFakeValidIngredientTagList()

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("GetValidIngredientTags", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return(exampleValidIngredientTagList, nil)
		s.validIngredientTagDataManager = validIngredientTagDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientTagList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("GetValidIngredientTags", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return((*models.ValidIngredientTagList)(nil), sql.ErrNoRows)
		s.validIngredientTagDataManager = validIngredientTagDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientTagList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager, ed)
	})

	T.Run("with error fetching valid ingredient tags from database", func(t *testing.T) {
		s := buildTestService()

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("GetValidIngredientTags", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return((*models.ValidIngredientTagList)(nil), errors.New("blah"))
		s.validIngredientTagDataManager = validIngredientTagDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTagList := fakemodels.BuildFakeValidIngredientTagList()

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("GetValidIngredientTags", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return(exampleValidIngredientTagList, nil)
		s.validIngredientTagDataManager = validIngredientTagDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientTagList")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager, ed)
	})
}

func TestValidIngredientTagsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		exampleInput := fakemodels.BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(exampleValidIngredientTag)

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("CreateValidIngredientTag", mock.Anything, mock.AnythingOfType("*models.ValidIngredientTagCreationInput")).Return(exampleValidIngredientTag, nil)
		s.validIngredientTagDataManager = validIngredientTagDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.validIngredientTagCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientTag")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager, mc, r, ed)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

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

	T.Run("with error creating valid ingredient tag", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		exampleInput := fakemodels.BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(exampleValidIngredientTag)

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("CreateValidIngredientTag", mock.Anything, mock.AnythingOfType("*models.ValidIngredientTagCreationInput")).Return(exampleValidIngredientTag, errors.New("blah"))
		s.validIngredientTagDataManager = validIngredientTagDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		exampleInput := fakemodels.BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(exampleValidIngredientTag)

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("CreateValidIngredientTag", mock.Anything, mock.AnythingOfType("*models.ValidIngredientTagCreationInput")).Return(exampleValidIngredientTag, nil)
		s.validIngredientTagDataManager = validIngredientTagDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.validIngredientTagCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientTag")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager, mc, r, ed)
	})
}

func TestValidIngredientTagsService_ExistenceHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		s.validIngredientTagIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientTag.ID
		}

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("ValidIngredientTagExists", mock.Anything, exampleValidIngredientTag.ID).Return(true, nil)
		s.validIngredientTagDataManager = validIngredientTagDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager)
	})

	T.Run("with no such valid ingredient tag in database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		s.validIngredientTagIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientTag.ID
		}

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("ValidIngredientTagExists", mock.Anything, exampleValidIngredientTag.ID).Return(false, sql.ErrNoRows)
		s.validIngredientTagDataManager = validIngredientTagDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager)
	})

	T.Run("with error fetching valid ingredient tag from database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		s.validIngredientTagIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientTag.ID
		}

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("ValidIngredientTagExists", mock.Anything, exampleValidIngredientTag.ID).Return(false, errors.New("blah"))
		s.validIngredientTagDataManager = validIngredientTagDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager)
	})
}

func TestValidIngredientTagsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		s.validIngredientTagIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientTag.ID
		}

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("GetValidIngredientTag", mock.Anything, exampleValidIngredientTag.ID).Return(exampleValidIngredientTag, nil)
		s.validIngredientTagDataManager = validIngredientTagDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientTag")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager, ed)
	})

	T.Run("with no such valid ingredient tag in database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		s.validIngredientTagIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientTag.ID
		}

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("GetValidIngredientTag", mock.Anything, exampleValidIngredientTag.ID).Return((*models.ValidIngredientTag)(nil), sql.ErrNoRows)
		s.validIngredientTagDataManager = validIngredientTagDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager)
	})

	T.Run("with error fetching valid ingredient tag from database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		s.validIngredientTagIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientTag.ID
		}

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("GetValidIngredientTag", mock.Anything, exampleValidIngredientTag.ID).Return((*models.ValidIngredientTag)(nil), errors.New("blah"))
		s.validIngredientTagDataManager = validIngredientTagDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		s.validIngredientTagIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientTag.ID
		}

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("GetValidIngredientTag", mock.Anything, exampleValidIngredientTag.ID).Return(exampleValidIngredientTag, nil)
		s.validIngredientTagDataManager = validIngredientTagDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientTag")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager, ed)
	})
}

func TestValidIngredientTagsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		exampleInput := fakemodels.BuildFakeValidIngredientTagUpdateInputFromValidIngredientTag(exampleValidIngredientTag)

		s.validIngredientTagIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientTag.ID
		}

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("GetValidIngredientTag", mock.Anything, exampleValidIngredientTag.ID).Return(exampleValidIngredientTag, nil)
		validIngredientTagDataManager.On("UpdateValidIngredientTag", mock.Anything, mock.AnythingOfType("*models.ValidIngredientTag")).Return(nil)
		s.validIngredientTagDataManager = validIngredientTagDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientTag")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, r, validIngredientTagDataManager, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

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

	T.Run("with no rows fetching valid ingredient tag", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		exampleInput := fakemodels.BuildFakeValidIngredientTagUpdateInputFromValidIngredientTag(exampleValidIngredientTag)

		s.validIngredientTagIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientTag.ID
		}

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("GetValidIngredientTag", mock.Anything, exampleValidIngredientTag.ID).Return((*models.ValidIngredientTag)(nil), sql.ErrNoRows)
		s.validIngredientTagDataManager = validIngredientTagDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager)
	})

	T.Run("with error fetching valid ingredient tag", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		exampleInput := fakemodels.BuildFakeValidIngredientTagUpdateInputFromValidIngredientTag(exampleValidIngredientTag)

		s.validIngredientTagIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientTag.ID
		}

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("GetValidIngredientTag", mock.Anything, exampleValidIngredientTag.ID).Return((*models.ValidIngredientTag)(nil), errors.New("blah"))
		s.validIngredientTagDataManager = validIngredientTagDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager)
	})

	T.Run("with error updating valid ingredient tag", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		exampleInput := fakemodels.BuildFakeValidIngredientTagUpdateInputFromValidIngredientTag(exampleValidIngredientTag)

		s.validIngredientTagIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientTag.ID
		}

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("GetValidIngredientTag", mock.Anything, exampleValidIngredientTag.ID).Return(exampleValidIngredientTag, nil)
		validIngredientTagDataManager.On("UpdateValidIngredientTag", mock.Anything, mock.AnythingOfType("*models.ValidIngredientTag")).Return(errors.New("blah"))
		s.validIngredientTagDataManager = validIngredientTagDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		exampleInput := fakemodels.BuildFakeValidIngredientTagUpdateInputFromValidIngredientTag(exampleValidIngredientTag)

		s.validIngredientTagIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientTag.ID
		}

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("GetValidIngredientTag", mock.Anything, exampleValidIngredientTag.ID).Return(exampleValidIngredientTag, nil)
		validIngredientTagDataManager.On("UpdateValidIngredientTag", mock.Anything, mock.AnythingOfType("*models.ValidIngredientTag")).Return(nil)
		s.validIngredientTagDataManager = validIngredientTagDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientTag")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, r, validIngredientTagDataManager, ed)
	})
}

func TestValidIngredientTagsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		s.validIngredientTagIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientTag.ID
		}

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("ArchiveValidIngredientTag", mock.Anything, exampleValidIngredientTag.ID).Return(nil)
		s.validIngredientTagDataManager = validIngredientTagDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.validIngredientTagCounter = mc

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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager, mc, r)
	})

	T.Run("with no valid ingredient tag in database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		s.validIngredientTagIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientTag.ID
		}

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("ArchiveValidIngredientTag", mock.Anything, exampleValidIngredientTag.ID).Return(sql.ErrNoRows)
		s.validIngredientTagDataManager = validIngredientTagDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		s.validIngredientTagIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientTag.ID
		}

		validIngredientTagDataManager := &mockmodels.ValidIngredientTagDataManager{}
		validIngredientTagDataManager.On("ArchiveValidIngredientTag", mock.Anything, exampleValidIngredientTag.ID).Return(errors.New("blah"))
		s.validIngredientTagDataManager = validIngredientTagDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientTagDataManager)
	})
}
