package validpreparations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockencoding "gitlab.com/prixfixe/prixfixe/internal/v1/encoding/mock"
	mockmetrics "gitlab.com/prixfixe/prixfixe/internal/v1/metrics/mock"
	mocksearch "gitlab.com/prixfixe/prixfixe/internal/v1/search/mock"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"
	mockmodels "gitlab.com/prixfixe/prixfixe/models/v1/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mocknewsman "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
)

func TestValidPreparationsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparationList := fakemodels.BuildFakeValidPreparationList()

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparations", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return(exampleValidPreparationList, nil)
		s.validPreparationDataManager = validPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidPreparationList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparations", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return((*models.ValidPreparationList)(nil), sql.ErrNoRows)
		s.validPreparationDataManager = validPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidPreparationList")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, ed)
	})

	T.Run("with error fetching valid preparations from database", func(t *testing.T) {
		s := buildTestService()

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparations", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return((*models.ValidPreparationList)(nil), errors.New("blah"))
		s.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparationList := fakemodels.BuildFakeValidPreparationList()

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparations", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return(exampleValidPreparationList, nil)
		s.validPreparationDataManager = validPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidPreparationList")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, ed)
	})
}

func TestValidPreparationsService_SearchHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleQuery := "whatever"
		exampleLimit := uint8(123)
		exampleValidPreparationList := fakemodels.BuildFakeValidPreparationList().ValidPreparations
		var exampleValidPreparationIDs []uint64
		for _, x := range exampleValidPreparationList {
			exampleValidPreparationIDs = append(exampleValidPreparationIDs, x.ID)
		}

		si := &mocksearch.IndexManager{}
		si.On("Search", mock.Anything, exampleQuery).Return(exampleValidPreparationIDs, nil)
		s.search = si

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparationsWithIDs", mock.Anything, exampleLimit, exampleValidPreparationIDs).Return(exampleValidPreparationList, nil)
		s.validPreparationDataManager = validPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("[]models.ValidPreparation")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, si, validPreparationDataManager, ed)
	})

	T.Run("with error conducting search", func(t *testing.T) {
		s := buildTestService()

		exampleQuery := "whatever"
		exampleLimit := uint8(123)

		si := &mocksearch.IndexManager{}
		si.On("Search", mock.Anything, exampleQuery).Return([]uint64{}, errors.New("blah"))
		s.search = si

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, si)
	})

	T.Run("with now rows returned", func(t *testing.T) {
		s := buildTestService()

		exampleQuery := "whatever"
		exampleLimit := uint8(123)
		exampleValidPreparationList := fakemodels.BuildFakeValidPreparationList().ValidPreparations
		var exampleValidPreparationIDs []uint64
		for _, x := range exampleValidPreparationList {
			exampleValidPreparationIDs = append(exampleValidPreparationIDs, x.ID)
		}

		si := &mocksearch.IndexManager{}
		si.On("Search", mock.Anything, exampleQuery).Return(exampleValidPreparationIDs, nil)
		s.search = si

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparationsWithIDs", mock.Anything, exampleLimit, exampleValidPreparationIDs).Return([]models.ValidPreparation{}, sql.ErrNoRows)
		s.validPreparationDataManager = validPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("[]models.ValidPreparation")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, si, validPreparationDataManager, ed)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		s := buildTestService()

		exampleQuery := "whatever"
		exampleLimit := uint8(123)
		exampleValidPreparationList := fakemodels.BuildFakeValidPreparationList().ValidPreparations
		var exampleValidPreparationIDs []uint64
		for _, x := range exampleValidPreparationList {
			exampleValidPreparationIDs = append(exampleValidPreparationIDs, x.ID)
		}

		si := &mocksearch.IndexManager{}
		si.On("Search", mock.Anything, exampleQuery).Return(exampleValidPreparationIDs, nil)
		s.search = si

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparationsWithIDs", mock.Anything, exampleLimit, exampleValidPreparationIDs).Return([]models.ValidPreparation{}, errors.New("blah"))
		s.validPreparationDataManager = validPreparationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, si, validPreparationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleQuery := "whatever"
		exampleLimit := uint8(123)
		exampleValidPreparationList := fakemodels.BuildFakeValidPreparationList().ValidPreparations
		var exampleValidPreparationIDs []uint64
		for _, x := range exampleValidPreparationList {
			exampleValidPreparationIDs = append(exampleValidPreparationIDs, x.ID)
		}

		si := &mocksearch.IndexManager{}
		si.On("Search", mock.Anything, exampleQuery).Return(exampleValidPreparationIDs, nil)
		s.search = si

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparationsWithIDs", mock.Anything, exampleLimit, exampleValidPreparationIDs).Return(exampleValidPreparationList, nil)
		s.validPreparationDataManager = validPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("[]models.ValidPreparation")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, si, validPreparationDataManager, ed)
	})
}

func TestValidPreparationsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("CreateValidPreparation", mock.Anything, mock.AnythingOfType("*models.ValidPreparationCreationInput")).Return(exampleValidPreparation, nil)
		s.validPreparationDataManager = validPreparationDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.validPreparationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		si := &mocksearch.IndexManager{}
		si.On("Index", mock.Anything, exampleValidPreparation.ID, exampleValidPreparation).Return(nil)
		s.search = si

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidPreparation")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, mc, r, si, ed)
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

	T.Run("with error creating valid preparation", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("CreateValidPreparation", mock.Anything, mock.AnythingOfType("*models.ValidPreparationCreationInput")).Return((*models.ValidPreparation)(nil), errors.New("blah"))
		s.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("CreateValidPreparation", mock.Anything, mock.AnythingOfType("*models.ValidPreparationCreationInput")).Return(exampleValidPreparation, nil)
		s.validPreparationDataManager = validPreparationDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.validPreparationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		si := &mocksearch.IndexManager{}
		si.On("Index", mock.Anything, exampleValidPreparation.ID, exampleValidPreparation).Return(nil)
		s.search = si

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidPreparation")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, mc, r, si, ed)
	})
}

func TestValidPreparationsService_ExistenceHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		s.validPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("ValidPreparationExists", mock.Anything, exampleValidPreparation.ID).Return(true, nil)
		s.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with no such valid preparation in database", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		s.validPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("ValidPreparationExists", mock.Anything, exampleValidPreparation.ID).Return(false, sql.ErrNoRows)
		s.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with error fetching valid preparation from database", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		s.validPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("ValidPreparationExists", mock.Anything, exampleValidPreparation.ID).Return(false, errors.New("blah"))
		s.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})
}

func TestValidPreparationsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		s.validPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparation", mock.Anything, exampleValidPreparation.ID).Return(exampleValidPreparation, nil)
		s.validPreparationDataManager = validPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidPreparation")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, ed)
	})

	T.Run("with no such valid preparation in database", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		s.validPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparation", mock.Anything, exampleValidPreparation.ID).Return((*models.ValidPreparation)(nil), sql.ErrNoRows)
		s.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with error fetching valid preparation from database", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		s.validPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparation", mock.Anything, exampleValidPreparation.ID).Return((*models.ValidPreparation)(nil), errors.New("blah"))
		s.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		s.validPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparation", mock.Anything, exampleValidPreparation.ID).Return(exampleValidPreparation, nil)
		s.validPreparationDataManager = validPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidPreparation")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, ed)
	})
}

func TestValidPreparationsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleInput := fakemodels.BuildFakeValidPreparationUpdateInputFromValidPreparation(exampleValidPreparation)

		s.validPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparation", mock.Anything, exampleValidPreparation.ID).Return(exampleValidPreparation, nil)
		validPreparationDataManager.On("UpdateValidPreparation", mock.Anything, mock.AnythingOfType("*models.ValidPreparation")).Return(nil)
		s.validPreparationDataManager = validPreparationDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		si := &mocksearch.IndexManager{}
		si.On("Index", mock.Anything, exampleValidPreparation.ID, exampleValidPreparation).Return(nil)
		s.search = si

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidPreparation")).Return(nil)
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

		mock.AssertExpectationsForObjects(t, r, validPreparationDataManager, ed)
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

	T.Run("with no rows fetching valid preparation", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleInput := fakemodels.BuildFakeValidPreparationUpdateInputFromValidPreparation(exampleValidPreparation)

		s.validPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparation", mock.Anything, exampleValidPreparation.ID).Return((*models.ValidPreparation)(nil), sql.ErrNoRows)
		s.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with error fetching valid preparation", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleInput := fakemodels.BuildFakeValidPreparationUpdateInputFromValidPreparation(exampleValidPreparation)

		s.validPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparation", mock.Anything, exampleValidPreparation.ID).Return((*models.ValidPreparation)(nil), errors.New("blah"))
		s.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with error updating valid preparation", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleInput := fakemodels.BuildFakeValidPreparationUpdateInputFromValidPreparation(exampleValidPreparation)

		s.validPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparation", mock.Anything, exampleValidPreparation.ID).Return(exampleValidPreparation, nil)
		validPreparationDataManager.On("UpdateValidPreparation", mock.Anything, mock.AnythingOfType("*models.ValidPreparation")).Return(errors.New("blah"))
		s.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleInput := fakemodels.BuildFakeValidPreparationUpdateInputFromValidPreparation(exampleValidPreparation)

		s.validPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("GetValidPreparation", mock.Anything, exampleValidPreparation.ID).Return(exampleValidPreparation, nil)
		validPreparationDataManager.On("UpdateValidPreparation", mock.Anything, mock.AnythingOfType("*models.ValidPreparation")).Return(nil)
		s.validPreparationDataManager = validPreparationDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		si := &mocksearch.IndexManager{}
		si.On("Index", mock.Anything, exampleValidPreparation.ID, exampleValidPreparation).Return(nil)
		s.search = si

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidPreparation")).Return(errors.New("blah"))
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

		mock.AssertExpectationsForObjects(t, r, validPreparationDataManager, ed)
	})
}

func TestValidPreparationsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		s.validPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("ArchiveValidPreparation", mock.Anything, exampleValidPreparation.ID).Return(nil)
		s.validPreparationDataManager = validPreparationDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		si := &mocksearch.IndexManager{}
		si.On("Delete", mock.Anything, exampleValidPreparation.ID).Return(nil)
		s.search = si

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.validPreparationCounter = mc

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, mc, r)
	})

	T.Run("with no valid preparation in database", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		s.validPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("ArchiveValidPreparation", mock.Anything, exampleValidPreparation.ID).Return(sql.ErrNoRows)
		s.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		s.validPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		validPreparationDataManager := &mockmodels.ValidPreparationDataManager{}
		validPreparationDataManager.On("ArchiveValidPreparation", mock.Anything, exampleValidPreparation.ID).Return(errors.New("blah"))
		s.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})
}
