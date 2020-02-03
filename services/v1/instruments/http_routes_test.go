package instruments

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

func TestInstrumentsService_List(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.InstrumentList{
			Instruments: []models.Instrument{
				{
					ID: 123,
				},
			},
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("GetInstruments", mock.Anything, mock.Anything, requestingUser.ID).Return(expected, nil)
		s.instrumentDatabase = id

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

		id := &mockmodels.InstrumentDataManager{}
		id.On("GetInstruments", mock.Anything, mock.Anything, requestingUser.ID).Return((*models.InstrumentList)(nil), sql.ErrNoRows)
		s.instrumentDatabase = id

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

	T.Run("with error fetching instruments from database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("GetInstruments", mock.Anything, mock.Anything, requestingUser.ID).Return((*models.InstrumentList)(nil), errors.New("blah"))
		s.instrumentDatabase = id

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
		expected := &models.InstrumentList{
			Instruments: []models.Instrument{},
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("GetInstruments", mock.Anything, mock.Anything, requestingUser.ID).Return(expected, nil)
		s.instrumentDatabase = id

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

func TestInstrumentsService_Create(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Instrument{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.instrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("CreateInstrument", mock.Anything, mock.Anything).Return(expected, nil)
		s.instrumentDatabase = id

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

		exampleInput := &models.InstrumentCreationInput{
			Name:        expected.Name,
			Variant:     expected.Variant,
			Description: expected.Description,
			Icon:        expected.Icon,
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

	T.Run("with error creating instrument", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Instrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("CreateInstrument", mock.Anything, mock.Anything).Return((*models.Instrument)(nil), errors.New("blah"))
		s.instrumentDatabase = id

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

		exampleInput := &models.InstrumentCreationInput{
			Name:        expected.Name,
			Variant:     expected.Variant,
			Description: expected.Description,
			Icon:        expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Instrument{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.instrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("CreateInstrument", mock.Anything, mock.Anything).Return(expected, nil)
		s.instrumentDatabase = id

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

		exampleInput := &models.InstrumentCreationInput{
			Name:        expected.Name,
			Variant:     expected.Variant,
			Description: expected.Description,
			Icon:        expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusCreated)
	})
}

func TestInstrumentsService_Read(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Instrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.instrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("GetInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		s.instrumentDatabase = id

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

	T.Run("with no such instrument in database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Instrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.instrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("GetInstrument", mock.Anything, expected.ID, requestingUser.ID).Return((*models.Instrument)(nil), sql.ErrNoRows)
		s.instrumentDatabase = id

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

	T.Run("with error fetching instrument from database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Instrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.instrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("GetInstrument", mock.Anything, expected.ID, requestingUser.ID).Return((*models.Instrument)(nil), errors.New("blah"))
		s.instrumentDatabase = id

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
		expected := &models.Instrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.instrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("GetInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		s.instrumentDatabase = id

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

func TestInstrumentsService_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Instrument{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.instrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.instrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("GetInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateInstrument", mock.Anything, mock.Anything).Return(nil)
		s.instrumentDatabase = id

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

		exampleInput := &models.InstrumentUpdateInput{
			Name:        expected.Name,
			Variant:     expected.Variant,
			Description: expected.Description,
			Icon:        expected.Icon,
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

	T.Run("with no rows fetching instrument", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Instrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.instrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("GetInstrument", mock.Anything, expected.ID, requestingUser.ID).Return((*models.Instrument)(nil), sql.ErrNoRows)
		s.instrumentDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.InstrumentUpdateInput{
			Name:        expected.Name,
			Variant:     expected.Variant,
			Description: expected.Description,
			Icon:        expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusNotFound)
	})

	T.Run("with error fetching instrument", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Instrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.instrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("GetInstrument", mock.Anything, expected.ID, requestingUser.ID).Return((*models.Instrument)(nil), errors.New("blah"))
		s.instrumentDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.InstrumentUpdateInput{
			Name:        expected.Name,
			Variant:     expected.Variant,
			Description: expected.Description,
			Icon:        expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error updating instrument", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Instrument{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.instrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.instrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("GetInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateInstrument", mock.Anything, mock.Anything).Return(errors.New("blah"))
		s.instrumentDatabase = id

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

		exampleInput := &models.InstrumentUpdateInput{
			Name:        expected.Name,
			Variant:     expected.Variant,
			Description: expected.Description,
			Icon:        expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Instrument{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.instrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.instrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("GetInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateInstrument", mock.Anything, mock.Anything).Return(nil)
		s.instrumentDatabase = id

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

		exampleInput := &models.InstrumentUpdateInput{
			Name:        expected.Name,
			Variant:     expected.Variant,
			Description: expected.Description,
			Icon:        expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusOK)
	})
}

func TestInstrumentsService_Archive(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Instrument{
			ID: 123,
		}

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement").Return()
		s.instrumentCounter = mc

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.instrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("ArchiveInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(nil)
		s.instrumentDatabase = id

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

	T.Run("with no instrument in database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Instrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.instrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("ArchiveInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(sql.ErrNoRows)
		s.instrumentDatabase = id

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
		expected := &models.Instrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.instrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.InstrumentDataManager{}
		id.On("ArchiveInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(errors.New("blah"))
		s.instrumentDatabase = id

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
