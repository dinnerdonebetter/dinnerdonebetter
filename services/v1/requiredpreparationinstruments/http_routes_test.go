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
	mockmodels "gitlab.com/prixfixe/prixfixe/models/v1/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mocknewsman "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
)

func TestRequiredPreparationInstrumentsService_List(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrumentList{
			RequiredPreparationInstruments: []models.RequiredPreparationInstrument{
				{
					ID: 123,
				},
			},
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("GetRequiredPreparationInstruments", mock.Anything, mock.Anything, requestingUser.ID).Return(expected, nil)
		s.requiredPreparationInstrumentDatabase = id

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.Anything).Return(nil)
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

		assert.Equal(t, res.Code, http.StatusOK)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("GetRequiredPreparationInstruments", mock.Anything, mock.Anything, requestingUser.ID).Return((*models.RequiredPreparationInstrumentList)(nil), sql.ErrNoRows)
		s.requiredPreparationInstrumentDatabase = id

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.Anything).Return(nil)
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

		assert.Equal(t, res.Code, http.StatusOK)
	})

	T.Run("with error fetching required preparation instruments from database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("GetRequiredPreparationInstruments", mock.Anything, mock.Anything, requestingUser.ID).Return((*models.RequiredPreparationInstrumentList)(nil), errors.New("blah"))
		s.requiredPreparationInstrumentDatabase = id

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.Anything).Return(nil)
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

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrumentList{
			RequiredPreparationInstruments: []models.RequiredPreparationInstrument{},
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("GetRequiredPreparationInstruments", mock.Anything, mock.Anything, requestingUser.ID).Return(expected, nil)
		s.requiredPreparationInstrumentDatabase = id

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.Anything).Return(errors.New("blah"))
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

		assert.Equal(t, res.Code, http.StatusOK)
	})
}

func TestRequiredPreparationInstrumentsService_Create(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.requiredPreparationInstrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("CreateRequiredPreparationInstrument", mock.Anything, mock.Anything).Return(expected, nil)
		s.requiredPreparationInstrumentDatabase = id

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.Anything).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.RequiredPreparationInstrumentCreationInput{
			InstrumentID:  expected.InstrumentID,
			PreparationID: expected.PreparationID,
			Notes:         expected.Notes,
		}
		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusCreated)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.Anything).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.CreateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusBadRequest)
	})

	T.Run("with error creating required preparation instrument", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("CreateRequiredPreparationInstrument", mock.Anything, mock.Anything).Return((*models.RequiredPreparationInstrument)(nil), errors.New("blah"))
		s.requiredPreparationInstrumentDatabase = id

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.Anything).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.RequiredPreparationInstrumentCreationInput{
			InstrumentID:  expected.InstrumentID,
			PreparationID: expected.PreparationID,
			Notes:         expected.Notes,
		}
		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.requiredPreparationInstrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("CreateRequiredPreparationInstrument", mock.Anything, mock.Anything).Return(expected, nil)
		s.requiredPreparationInstrumentDatabase = id

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.Anything).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.RequiredPreparationInstrumentCreationInput{
			InstrumentID:  expected.InstrumentID,
			PreparationID: expected.PreparationID,
			Notes:         expected.Notes,
		}
		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusCreated)
	})
}

func TestRequiredPreparationInstrumentsService_Read(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("GetRequiredPreparationInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		s.requiredPreparationInstrumentDatabase = id

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.Anything).Return(nil)
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

		assert.Equal(t, res.Code, http.StatusOK)
	})

	T.Run("with no such required preparation instrument in database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("GetRequiredPreparationInstrument", mock.Anything, expected.ID, requestingUser.ID).Return((*models.RequiredPreparationInstrument)(nil), sql.ErrNoRows)
		s.requiredPreparationInstrumentDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusNotFound)
	})

	T.Run("with error fetching required preparation instrument from database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("GetRequiredPreparationInstrument", mock.Anything, expected.ID, requestingUser.ID).Return((*models.RequiredPreparationInstrument)(nil), errors.New("blah"))
		s.requiredPreparationInstrumentDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("GetRequiredPreparationInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		s.requiredPreparationInstrumentDatabase = id

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.Anything).Return(errors.New("blah"))
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

		assert.Equal(t, res.Code, http.StatusOK)
	})
}

func TestRequiredPreparationInstrumentsService_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.requiredPreparationInstrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("GetRequiredPreparationInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateRequiredPreparationInstrument", mock.Anything, mock.Anything).Return(nil)
		s.requiredPreparationInstrumentDatabase = id

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.Anything).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.RequiredPreparationInstrumentUpdateInput{
			InstrumentID:  expected.InstrumentID,
			PreparationID: expected.PreparationID,
			Notes:         expected.Notes,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusOK)
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

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusBadRequest)
	})

	T.Run("with no rows fetching required preparation instrument", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("GetRequiredPreparationInstrument", mock.Anything, expected.ID, requestingUser.ID).Return((*models.RequiredPreparationInstrument)(nil), sql.ErrNoRows)
		s.requiredPreparationInstrumentDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.RequiredPreparationInstrumentUpdateInput{
			InstrumentID:  expected.InstrumentID,
			PreparationID: expected.PreparationID,
			Notes:         expected.Notes,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusNotFound)
	})

	T.Run("with error fetching required preparation instrument", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("GetRequiredPreparationInstrument", mock.Anything, expected.ID, requestingUser.ID).Return((*models.RequiredPreparationInstrument)(nil), errors.New("blah"))
		s.requiredPreparationInstrumentDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.RequiredPreparationInstrumentUpdateInput{
			InstrumentID:  expected.InstrumentID,
			PreparationID: expected.PreparationID,
			Notes:         expected.Notes,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error updating required preparation instrument", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.requiredPreparationInstrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("GetRequiredPreparationInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateRequiredPreparationInstrument", mock.Anything, mock.Anything).Return(errors.New("blah"))
		s.requiredPreparationInstrumentDatabase = id

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.Anything).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.RequiredPreparationInstrumentUpdateInput{
			InstrumentID:  expected.InstrumentID,
			PreparationID: expected.PreparationID,
			Notes:         expected.Notes,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.requiredPreparationInstrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("GetRequiredPreparationInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateRequiredPreparationInstrument", mock.Anything, mock.Anything).Return(nil)
		s.requiredPreparationInstrumentDatabase = id

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.Anything).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.RequiredPreparationInstrumentUpdateInput{
			InstrumentID:  expected.InstrumentID,
			PreparationID: expected.PreparationID,
			Notes:         expected.Notes,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusOK)
	})
}

func TestRequiredPreparationInstrumentsService_Archive(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement").Return()
		s.requiredPreparationInstrumentCounter = mc

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("ArchiveRequiredPreparationInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(nil)
		s.requiredPreparationInstrumentDatabase = id

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.Anything).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusNoContent)
	})

	T.Run("with no required preparation instrument in database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("ArchiveRequiredPreparationInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(sql.ErrNoRows)
		s.requiredPreparationInstrumentDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusNotFound)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RequiredPreparationInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.requiredPreparationInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RequiredPreparationInstrumentDataManager{}
		id.On("ArchiveRequiredPreparationInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(errors.New("blah"))
		s.requiredPreparationInstrumentDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})
}
