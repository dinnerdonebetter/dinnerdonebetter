package requiredpreparationinstruments

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

func TestRequiredPreparationInstrumentsService_ListHandler(T *testing.T) {
	T.Parallel()

	exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
	validPreparationIDFetcher := func(_ *http.Request) uint64 {
		return exampleValidPreparation.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrumentList := fakemodels.BuildFakeRequiredPreparationInstrumentList()

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstruments", mock.Anything, exampleValidPreparation.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRequiredPreparationInstrumentList, nil)
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RequiredPreparationInstrumentList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, requiredPreparationInstrumentDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstruments", mock.Anything, exampleValidPreparation.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RequiredPreparationInstrumentList)(nil), sql.ErrNoRows)
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RequiredPreparationInstrumentList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, requiredPreparationInstrumentDataManager, ed)
	})

	T.Run("with error fetching required preparation instruments from database", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstruments", mock.Anything, exampleValidPreparation.ID, mock.AnythingOfType("*models.QueryFilter")).Return((*models.RequiredPreparationInstrumentList)(nil), errors.New("blah"))
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, requiredPreparationInstrumentDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrumentList := fakemodels.BuildFakeRequiredPreparationInstrumentList()

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstruments", mock.Anything, exampleValidPreparation.ID, mock.AnythingOfType("*models.QueryFilter")).Return(exampleRequiredPreparationInstrumentList, nil)
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RequiredPreparationInstrumentList")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, requiredPreparationInstrumentDataManager, ed)
	})
}

func TestRequiredPreparationInstrumentsService_CreateHandler(T *testing.T) {
	T.Parallel()

	exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
	validPreparationIDFetcher := func(_ *http.Request) uint64 {
		return exampleValidPreparation.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		exampleInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("ValidPreparationExists", mock.Anything, exampleValidPreparation.ID).Return(true, nil)
		s.validPreparationDataManager = validPreparationDataManager

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("CreateRequiredPreparationInstrument", mock.Anything, mock.AnythingOfType("*models.RequiredPreparationInstrumentCreationInput")).Return(exampleRequiredPreparationInstrument, nil)
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.requiredPreparationInstrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RequiredPreparationInstrument")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, requiredPreparationInstrumentDataManager, mc, r, ed)
	})

	T.Run("with nonexistent valid preparation", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		exampleInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("ValidPreparationExists", mock.Anything, exampleValidPreparation.ID).Return(false, nil)
		s.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with error checking valid preparation existence", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		exampleInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("ValidPreparationExists", mock.Anything, exampleValidPreparation.ID).Return(true, errors.New("blah"))
		s.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

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

	T.Run("with error creating required preparation instrument", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		exampleInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("ValidPreparationExists", mock.Anything, exampleValidPreparation.ID).Return(true, nil)
		s.validPreparationDataManager = validPreparationDataManager

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("CreateRequiredPreparationInstrument", mock.Anything, mock.AnythingOfType("*models.RequiredPreparationInstrumentCreationInput")).Return(exampleRequiredPreparationInstrument, errors.New("blah"))
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, requiredPreparationInstrumentDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		exampleInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("ValidPreparationExists", mock.Anything, exampleValidPreparation.ID).Return(true, nil)
		s.validPreparationDataManager = validPreparationDataManager

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("CreateRequiredPreparationInstrument", mock.Anything, mock.AnythingOfType("*models.RequiredPreparationInstrumentCreationInput")).Return(exampleRequiredPreparationInstrument, nil)
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.requiredPreparationInstrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RequiredPreparationInstrument")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, requiredPreparationInstrumentDataManager, mc, r, ed)
	})
}

func TestRequiredPreparationInstrumentsService_ExistenceHandler(T *testing.T) {
	T.Parallel()

	exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
	validPreparationIDFetcher := func(_ *http.Request) uint64 {
		return exampleValidPreparation.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRequiredPreparationInstrument.ID
		}

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("RequiredPreparationInstrumentExists", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return(true, nil)
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, requiredPreparationInstrumentDataManager)
	})

	T.Run("with no such required preparation instrument in database", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRequiredPreparationInstrument.ID
		}

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("RequiredPreparationInstrumentExists", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return(false, sql.ErrNoRows)
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, requiredPreparationInstrumentDataManager)
	})

	T.Run("with error fetching required preparation instrument from database", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRequiredPreparationInstrument.ID
		}

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("RequiredPreparationInstrumentExists", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return(false, errors.New("blah"))
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, requiredPreparationInstrumentDataManager)
	})
}

func TestRequiredPreparationInstrumentsService_ReadHandler(T *testing.T) {
	T.Parallel()

	exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
	validPreparationIDFetcher := func(_ *http.Request) uint64 {
		return exampleValidPreparation.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRequiredPreparationInstrument.ID
		}

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstrument", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return(exampleRequiredPreparationInstrument, nil)
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RequiredPreparationInstrument")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, requiredPreparationInstrumentDataManager, ed)
	})

	T.Run("with no such required preparation instrument in database", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRequiredPreparationInstrument.ID
		}

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstrument", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return((*models.RequiredPreparationInstrument)(nil), sql.ErrNoRows)
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, requiredPreparationInstrumentDataManager)
	})

	T.Run("with error fetching required preparation instrument from database", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRequiredPreparationInstrument.ID
		}

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstrument", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return((*models.RequiredPreparationInstrument)(nil), errors.New("blah"))
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, requiredPreparationInstrumentDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRequiredPreparationInstrument.ID
		}

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstrument", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return(exampleRequiredPreparationInstrument, nil)
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RequiredPreparationInstrument")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, requiredPreparationInstrumentDataManager, ed)
	})
}

func TestRequiredPreparationInstrumentsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
	validPreparationIDFetcher := func(_ *http.Request) uint64 {
		return exampleValidPreparation.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		exampleInput := fakemodels.BuildFakeRequiredPreparationInstrumentUpdateInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRequiredPreparationInstrument.ID
		}

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstrument", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return(exampleRequiredPreparationInstrument, nil)
		requiredPreparationInstrumentDataManager.On("UpdateRequiredPreparationInstrument", mock.Anything, mock.AnythingOfType("*models.RequiredPreparationInstrument")).Return(nil)
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RequiredPreparationInstrument")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, r, requiredPreparationInstrumentDataManager, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

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

	T.Run("with no rows fetching required preparation instrument", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		exampleInput := fakemodels.BuildFakeRequiredPreparationInstrumentUpdateInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRequiredPreparationInstrument.ID
		}

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstrument", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return((*models.RequiredPreparationInstrument)(nil), sql.ErrNoRows)
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, requiredPreparationInstrumentDataManager)
	})

	T.Run("with error fetching required preparation instrument", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		exampleInput := fakemodels.BuildFakeRequiredPreparationInstrumentUpdateInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRequiredPreparationInstrument.ID
		}

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstrument", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return((*models.RequiredPreparationInstrument)(nil), errors.New("blah"))
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, requiredPreparationInstrumentDataManager)
	})

	T.Run("with error updating required preparation instrument", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		exampleInput := fakemodels.BuildFakeRequiredPreparationInstrumentUpdateInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRequiredPreparationInstrument.ID
		}

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstrument", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return(exampleRequiredPreparationInstrument, nil)
		requiredPreparationInstrumentDataManager.On("UpdateRequiredPreparationInstrument", mock.Anything, mock.AnythingOfType("*models.RequiredPreparationInstrument")).Return(errors.New("blah"))
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, requiredPreparationInstrumentDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		exampleInput := fakemodels.BuildFakeRequiredPreparationInstrumentUpdateInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRequiredPreparationInstrument.ID
		}

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("GetRequiredPreparationInstrument", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return(exampleRequiredPreparationInstrument, nil)
		requiredPreparationInstrumentDataManager.On("UpdateRequiredPreparationInstrument", mock.Anything, mock.AnythingOfType("*models.RequiredPreparationInstrument")).Return(nil)
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.RequiredPreparationInstrument")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, r, requiredPreparationInstrumentDataManager, ed)
	})
}

func TestRequiredPreparationInstrumentsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
	validPreparationIDFetcher := func(_ *http.Request) uint64 {
		return exampleValidPreparation.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRequiredPreparationInstrument.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("ValidPreparationExists", mock.Anything, exampleValidPreparation.ID).Return(true, nil)
		s.validPreparationDataManager = validPreparationDataManager

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("ArchiveRequiredPreparationInstrument", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return(nil)
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.requiredPreparationInstrumentCounter = mc

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, requiredPreparationInstrumentDataManager, mc, r)
	})

	T.Run("with nonexistent valid preparation", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("ValidPreparationExists", mock.Anything, exampleValidPreparation.ID).Return(false, nil)
		s.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with error checking valid preparation existence", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("ValidPreparationExists", mock.Anything, exampleValidPreparation.ID).Return(true, errors.New("blah"))
		s.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with no required preparation instrument in database", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRequiredPreparationInstrument.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("ValidPreparationExists", mock.Anything, exampleValidPreparation.ID).Return(true, nil)
		s.validPreparationDataManager = validPreparationDataManager

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("ArchiveRequiredPreparationInstrument", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return(sql.ErrNoRows)
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, requiredPreparationInstrumentDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		s.validPreparationIDFetcher = validPreparationIDFetcher

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleRequiredPreparationInstrument.BelongsToValidPreparation = exampleValidPreparation.ID
		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return exampleRequiredPreparationInstrument.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("ValidPreparationExists", mock.Anything, exampleValidPreparation.ID).Return(true, nil)
		s.validPreparationDataManager = validPreparationDataManager

		requiredPreparationInstrumentDataManager := &mockmodels.RequiredPreparationInstrumentDataManager{}
		requiredPreparationInstrumentDataManager.On("ArchiveRequiredPreparationInstrument", mock.Anything, exampleValidPreparation.ID, exampleRequiredPreparationInstrument.ID).Return(errors.New("blah"))
		s.requiredPreparationInstrumentDataManager = requiredPreparationInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, requiredPreparationInstrumentDataManager)
	})
}
