package validingredients

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

func TestValidIngredientsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientList := fakemodels.BuildFakeValidIngredientList()

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("GetValidIngredients", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return(exampleValidIngredientList, nil)
		s.validIngredientDataManager = validIngredientDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("GetValidIngredients", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return((*models.ValidIngredientList)(nil), sql.ErrNoRows)
		s.validIngredientDataManager = validIngredientDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager, ed)
	})

	T.Run("with error fetching valid ingredients from database", func(t *testing.T) {
		s := buildTestService()

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("GetValidIngredients", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return((*models.ValidIngredientList)(nil), errors.New("blah"))
		s.validIngredientDataManager = validIngredientDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientList := fakemodels.BuildFakeValidIngredientList()

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("GetValidIngredients", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return(exampleValidIngredientList, nil)
		s.validIngredientDataManager = validIngredientDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientList")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager, ed)
	})
}

func TestValidIngredientsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("CreateValidIngredient", mock.Anything, mock.AnythingOfType("*models.ValidIngredientCreationInput")).Return(exampleValidIngredient, nil)
		s.validIngredientDataManager = validIngredientDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.validIngredientCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredient")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager, mc, r, ed)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

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

	T.Run("with error creating valid ingredient", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("CreateValidIngredient", mock.Anything, mock.AnythingOfType("*models.ValidIngredientCreationInput")).Return((*models.ValidIngredient)(nil), errors.New("blah"))
		s.validIngredientDataManager = validIngredientDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("CreateValidIngredient", mock.Anything, mock.AnythingOfType("*models.ValidIngredientCreationInput")).Return(exampleValidIngredient, nil)
		s.validIngredientDataManager = validIngredientDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.validIngredientCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredient")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager, mc, r, ed)
	})
}

func TestValidIngredientsService_ExistenceHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		s.validIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("ValidIngredientExists", mock.Anything, exampleValidIngredient.ID).Return(true, nil)
		s.validIngredientDataManager = validIngredientDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager)
	})

	T.Run("with no such valid ingredient in database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		s.validIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("ValidIngredientExists", mock.Anything, exampleValidIngredient.ID).Return(false, sql.ErrNoRows)
		s.validIngredientDataManager = validIngredientDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager)
	})

	T.Run("with error fetching valid ingredient from database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		s.validIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("ValidIngredientExists", mock.Anything, exampleValidIngredient.ID).Return(false, errors.New("blah"))
		s.validIngredientDataManager = validIngredientDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager)
	})
}

func TestValidIngredientsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		s.validIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("GetValidIngredient", mock.Anything, exampleValidIngredient.ID).Return(exampleValidIngredient, nil)
		s.validIngredientDataManager = validIngredientDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredient")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager, ed)
	})

	T.Run("with no such valid ingredient in database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		s.validIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("GetValidIngredient", mock.Anything, exampleValidIngredient.ID).Return((*models.ValidIngredient)(nil), sql.ErrNoRows)
		s.validIngredientDataManager = validIngredientDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager)
	})

	T.Run("with error fetching valid ingredient from database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		s.validIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("GetValidIngredient", mock.Anything, exampleValidIngredient.ID).Return((*models.ValidIngredient)(nil), errors.New("blah"))
		s.validIngredientDataManager = validIngredientDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		s.validIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("GetValidIngredient", mock.Anything, exampleValidIngredient.ID).Return(exampleValidIngredient, nil)
		s.validIngredientDataManager = validIngredientDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredient")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager, ed)
	})
}

func TestValidIngredientsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleInput := fakemodels.BuildFakeValidIngredientUpdateInputFromValidIngredient(exampleValidIngredient)

		s.validIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("GetValidIngredient", mock.Anything, exampleValidIngredient.ID).Return(exampleValidIngredient, nil)
		validIngredientDataManager.On("UpdateValidIngredient", mock.Anything, mock.AnythingOfType("*models.ValidIngredient")).Return(nil)
		s.validIngredientDataManager = validIngredientDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredient")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, r, validIngredientDataManager, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

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

	T.Run("with no rows fetching valid ingredient", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleInput := fakemodels.BuildFakeValidIngredientUpdateInputFromValidIngredient(exampleValidIngredient)

		s.validIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("GetValidIngredient", mock.Anything, exampleValidIngredient.ID).Return((*models.ValidIngredient)(nil), sql.ErrNoRows)
		s.validIngredientDataManager = validIngredientDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager)
	})

	T.Run("with error fetching valid ingredient", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleInput := fakemodels.BuildFakeValidIngredientUpdateInputFromValidIngredient(exampleValidIngredient)

		s.validIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("GetValidIngredient", mock.Anything, exampleValidIngredient.ID).Return((*models.ValidIngredient)(nil), errors.New("blah"))
		s.validIngredientDataManager = validIngredientDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager)
	})

	T.Run("with error updating valid ingredient", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleInput := fakemodels.BuildFakeValidIngredientUpdateInputFromValidIngredient(exampleValidIngredient)

		s.validIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("GetValidIngredient", mock.Anything, exampleValidIngredient.ID).Return(exampleValidIngredient, nil)
		validIngredientDataManager.On("UpdateValidIngredient", mock.Anything, mock.AnythingOfType("*models.ValidIngredient")).Return(errors.New("blah"))
		s.validIngredientDataManager = validIngredientDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleInput := fakemodels.BuildFakeValidIngredientUpdateInputFromValidIngredient(exampleValidIngredient)

		s.validIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("GetValidIngredient", mock.Anything, exampleValidIngredient.ID).Return(exampleValidIngredient, nil)
		validIngredientDataManager.On("UpdateValidIngredient", mock.Anything, mock.AnythingOfType("*models.ValidIngredient")).Return(nil)
		s.validIngredientDataManager = validIngredientDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredient")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, r, validIngredientDataManager, ed)
	})
}

func TestValidIngredientsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		s.validIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("ArchiveValidIngredient", mock.Anything, exampleValidIngredient.ID).Return(nil)
		s.validIngredientDataManager = validIngredientDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.validIngredientCounter = mc

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager, mc, r)
	})

	T.Run("with no valid ingredient in database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		s.validIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("ArchiveValidIngredient", mock.Anything, exampleValidIngredient.ID).Return(sql.ErrNoRows)
		s.validIngredientDataManager = validIngredientDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		s.validIngredientIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		validIngredientDataManager := &mockmodels.ValidIngredientDataManager{}
		validIngredientDataManager.On("ArchiveValidIngredient", mock.Anything, exampleValidIngredient.ID).Return(errors.New("blah"))
		s.validIngredientDataManager = validIngredientDataManager

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

		mock.AssertExpectationsForObjects(t, validIngredientDataManager)
	})
}
