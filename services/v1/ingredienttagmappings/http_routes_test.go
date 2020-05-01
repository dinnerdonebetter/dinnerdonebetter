package ingredienttagmappings

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

func TestIngredientTagMappingsService_ListHandler(T *testing.T) {
	T.Parallel()

	exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
	validIngredientIDFetcher := func(_ *http.Request) uint64 {
		return exampleValidIngredient.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMappingList := fakemodels.BuildFakeIngredientTagMappingList()

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("GetIngredientTagMappings", mock.Anything, exampleValidIngredient.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleIngredientTagMappingList, nil)
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IngredientTagMappingList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, ingredientTagMappingDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("GetIngredientTagMappings", mock.Anything, exampleValidIngredient.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.IngredientTagMappingList)(nil), sql.ErrNoRows)
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IngredientTagMappingList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, ingredientTagMappingDataManager, ed)
	})

	T.Run("with error fetching ingredient tag mappings from database", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("GetIngredientTagMappings", mock.Anything, exampleValidIngredient.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.IngredientTagMappingList)(nil), errors.New("blah"))
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

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

		mock.AssertExpectationsForObjects(t, ingredientTagMappingDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMappingList := fakemodels.BuildFakeIngredientTagMappingList()

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("GetIngredientTagMappings", mock.Anything, exampleValidIngredient.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleIngredientTagMappingList, nil)
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IngredientTagMappingList")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, ingredientTagMappingDataManager, ed)
	})
}

func TestIngredientTagMappingsService_CreateHandler(T *testing.T) {
	T.Parallel()

	exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
	validIngredientIDFetcher := func(_ *http.Request) uint64 {
		return exampleValidIngredient.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("ValidIngredientExists", mock.Anything, exampleValidIngredient.ID).Return(true, nil)
		s.validIngredientDataManager = validIngredientDataManager

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("CreateIngredientTagMapping", mock.Anything, mock.AnythingOfType("*models.IngredientTagMappingCreationInput")).Return(exampleIngredientTagMapping, nil)
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.ingredientTagMappingCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IngredientTagMapping")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager, ingredientTagMappingDataManager, mc, r, ed)
	})

	T.Run("with nonexistent valid ingredient", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("ValidIngredientExists", mock.Anything, exampleValidIngredient.ID).Return(false, nil)
		s.validIngredientDataManager = validIngredientDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager)
	})

	T.Run("with error checking valid ingredient existence", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("ValidIngredientExists", mock.Anything, exampleValidIngredient.ID).Return(true, errors.New("blah"))
		s.validIngredientDataManager = validIngredientDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

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

	T.Run("with error creating ingredient tag mapping", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("ValidIngredientExists", mock.Anything, exampleValidIngredient.ID).Return(true, nil)
		s.validIngredientDataManager = validIngredientDataManager

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("CreateIngredientTagMapping", mock.Anything, mock.AnythingOfType("*models.IngredientTagMappingCreationInput")).Return(exampleIngredientTagMapping, errors.New("blah"))
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager, ingredientTagMappingDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("ValidIngredientExists", mock.Anything, exampleValidIngredient.ID).Return(true, nil)
		s.validIngredientDataManager = validIngredientDataManager

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("CreateIngredientTagMapping", mock.Anything, mock.AnythingOfType("*models.IngredientTagMappingCreationInput")).Return(exampleIngredientTagMapping, nil)
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.ingredientTagMappingCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IngredientTagMapping")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager, ingredientTagMappingDataManager, mc, r, ed)
	})
}

func TestIngredientTagMappingsService_ExistenceHandler(T *testing.T) {
	T.Parallel()

	exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
	validIngredientIDFetcher := func(_ *http.Request) uint64 {
		return exampleValidIngredient.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		s.ingredientTagMappingIDFetcher = func(req *http.Request) uint64 {
			return exampleIngredientTagMapping.ID
		}

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("IngredientTagMappingExists", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return(true, nil)
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

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

		mock.AssertExpectationsForObjects(t, ingredientTagMappingDataManager)
	})

	T.Run("with no such ingredient tag mapping in database", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		s.ingredientTagMappingIDFetcher = func(req *http.Request) uint64 {
			return exampleIngredientTagMapping.ID
		}

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("IngredientTagMappingExists", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return(false, sql.ErrNoRows)
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

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

		mock.AssertExpectationsForObjects(t, ingredientTagMappingDataManager)
	})

	T.Run("with error fetching ingredient tag mapping from database", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		s.ingredientTagMappingIDFetcher = func(req *http.Request) uint64 {
			return exampleIngredientTagMapping.ID
		}

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("IngredientTagMappingExists", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return(false, errors.New("blah"))
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

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

		mock.AssertExpectationsForObjects(t, ingredientTagMappingDataManager)
	})
}

func TestIngredientTagMappingsService_ReadHandler(T *testing.T) {
	T.Parallel()

	exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
	validIngredientIDFetcher := func(_ *http.Request) uint64 {
		return exampleValidIngredient.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		s.ingredientTagMappingIDFetcher = func(req *http.Request) uint64 {
			return exampleIngredientTagMapping.ID
		}

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("GetIngredientTagMapping", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return(exampleIngredientTagMapping, nil)
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IngredientTagMapping")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, ingredientTagMappingDataManager, ed)
	})

	T.Run("with no such ingredient tag mapping in database", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		s.ingredientTagMappingIDFetcher = func(req *http.Request) uint64 {
			return exampleIngredientTagMapping.ID
		}

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("GetIngredientTagMapping", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return((*models.IngredientTagMapping)(nil), sql.ErrNoRows)
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

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

		mock.AssertExpectationsForObjects(t, ingredientTagMappingDataManager)
	})

	T.Run("with error fetching ingredient tag mapping from database", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		s.ingredientTagMappingIDFetcher = func(req *http.Request) uint64 {
			return exampleIngredientTagMapping.ID
		}

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("GetIngredientTagMapping", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return((*models.IngredientTagMapping)(nil), errors.New("blah"))
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

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

		mock.AssertExpectationsForObjects(t, ingredientTagMappingDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		s.ingredientTagMappingIDFetcher = func(req *http.Request) uint64 {
			return exampleIngredientTagMapping.ID
		}

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("GetIngredientTagMapping", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return(exampleIngredientTagMapping, nil)
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IngredientTagMapping")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, ingredientTagMappingDataManager, ed)
	})
}

func TestIngredientTagMappingsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
	validIngredientIDFetcher := func(_ *http.Request) uint64 {
		return exampleValidIngredient.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeIngredientTagMappingUpdateInputFromIngredientTagMapping(exampleIngredientTagMapping)

		s.ingredientTagMappingIDFetcher = func(req *http.Request) uint64 {
			return exampleIngredientTagMapping.ID
		}

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("GetIngredientTagMapping", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return(exampleIngredientTagMapping, nil)
		ingredientTagMappingDataManager.On("UpdateIngredientTagMapping", mock.Anything, mock.AnythingOfType("*models.IngredientTagMapping")).Return(nil)
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IngredientTagMapping")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, r, ingredientTagMappingDataManager, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

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

	T.Run("with no rows fetching ingredient tag mapping", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeIngredientTagMappingUpdateInputFromIngredientTagMapping(exampleIngredientTagMapping)

		s.ingredientTagMappingIDFetcher = func(req *http.Request) uint64 {
			return exampleIngredientTagMapping.ID
		}

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("GetIngredientTagMapping", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return((*models.IngredientTagMapping)(nil), sql.ErrNoRows)
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

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

		mock.AssertExpectationsForObjects(t, ingredientTagMappingDataManager)
	})

	T.Run("with error fetching ingredient tag mapping", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeIngredientTagMappingUpdateInputFromIngredientTagMapping(exampleIngredientTagMapping)

		s.ingredientTagMappingIDFetcher = func(req *http.Request) uint64 {
			return exampleIngredientTagMapping.ID
		}

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("GetIngredientTagMapping", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return((*models.IngredientTagMapping)(nil), errors.New("blah"))
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

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

		mock.AssertExpectationsForObjects(t, ingredientTagMappingDataManager)
	})

	T.Run("with error updating ingredient tag mapping", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeIngredientTagMappingUpdateInputFromIngredientTagMapping(exampleIngredientTagMapping)

		s.ingredientTagMappingIDFetcher = func(req *http.Request) uint64 {
			return exampleIngredientTagMapping.ID
		}

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("GetIngredientTagMapping", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return(exampleIngredientTagMapping, nil)
		ingredientTagMappingDataManager.On("UpdateIngredientTagMapping", mock.Anything, mock.AnythingOfType("*models.IngredientTagMapping")).Return(errors.New("blah"))
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

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

		mock.AssertExpectationsForObjects(t, ingredientTagMappingDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeIngredientTagMappingUpdateInputFromIngredientTagMapping(exampleIngredientTagMapping)

		s.ingredientTagMappingIDFetcher = func(req *http.Request) uint64 {
			return exampleIngredientTagMapping.ID
		}

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("GetIngredientTagMapping", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return(exampleIngredientTagMapping, nil)
		ingredientTagMappingDataManager.On("UpdateIngredientTagMapping", mock.Anything, mock.AnythingOfType("*models.IngredientTagMapping")).Return(nil)
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.IngredientTagMapping")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, r, ingredientTagMappingDataManager, ed)
	})
}

func TestIngredientTagMappingsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
	validIngredientIDFetcher := func(_ *http.Request) uint64 {
		return exampleValidIngredient.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		s.ingredientTagMappingIDFetcher = func(req *http.Request) uint64 {
			return exampleIngredientTagMapping.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("ValidIngredientExists", mock.Anything, exampleValidIngredient.ID).Return(true, nil)
		s.validIngredientDataManager = validIngredientDataManager

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("ArchiveIngredientTagMapping", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return(nil)
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.ingredientTagMappingCounter = mc

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager, ingredientTagMappingDataManager, mc, r)
	})

	T.Run("with nonexistent valid ingredient", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("ValidIngredientExists", mock.Anything, exampleValidIngredient.ID).Return(false, nil)
		s.validIngredientDataManager = validIngredientDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager)
	})

	T.Run("with error checking valid ingredient existence", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("ValidIngredientExists", mock.Anything, exampleValidIngredient.ID).Return(true, errors.New("blah"))
		s.validIngredientDataManager = validIngredientDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager)
	})

	T.Run("with no ingredient tag mapping in database", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		s.ingredientTagMappingIDFetcher = func(req *http.Request) uint64 {
			return exampleIngredientTagMapping.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("ValidIngredientExists", mock.Anything, exampleValidIngredient.ID).Return(true, nil)
		s.validIngredientDataManager = validIngredientDataManager

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("ArchiveIngredientTagMapping", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return(sql.ErrNoRows)
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager, ingredientTagMappingDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		s.validIngredientIDFetcher = validIngredientIDFetcher

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		s.ingredientTagMappingIDFetcher = func(req *http.Request) uint64 {
			return exampleIngredientTagMapping.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("ValidIngredientExists", mock.Anything, exampleValidIngredient.ID).Return(true, nil)
		s.validIngredientDataManager = validIngredientDataManager

		ingredientTagMappingDataManager := &mockmodels.IngredientTagMappingDataManager{}
		ingredientTagMappingDataManager.On("ArchiveIngredientTagMapping", mock.Anything, exampleValidIngredient.ID, exampleIngredientTagMapping.ID).Return(errors.New("blah"))
		s.ingredientTagMappingDataManager = ingredientTagMappingDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager, ingredientTagMappingDataManager)
	})
}
