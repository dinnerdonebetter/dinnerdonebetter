package preparations

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

func TestPreparationsService_List(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.PreparationList{
			Preparations: []models.Preparation{
				{
					ID: 123,
				},
			},
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("GetPreparations", mock.Anything, mock.Anything).Return(expected, nil)
		s.preparationDatabase = id

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

		id := &mockmodels.PreparationDataManager{}
		id.On("GetPreparations", mock.Anything, mock.Anything).Return((*models.PreparationList)(nil), sql.ErrNoRows)
		s.preparationDatabase = id

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

	T.Run("with error fetching preparations from database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("GetPreparations", mock.Anything, mock.Anything).Return((*models.PreparationList)(nil), errors.New("blah"))
		s.preparationDatabase = id

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
		expected := &models.PreparationList{
			Preparations: []models.Preparation{},
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("GetPreparations", mock.Anything, mock.Anything).Return(expected, nil)
		s.preparationDatabase = id

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

func TestPreparationsService_Create(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Preparation{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.preparationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("CreatePreparation", mock.Anything, mock.Anything).Return(expected, nil)
		s.preparationDatabase = id

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

		exampleInput := &models.PreparationCreationInput{
			Name:           expected.Name,
			Variant:        expected.Variant,
			Description:    expected.Description,
			AllergyWarning: expected.AllergyWarning,
			Icon:           expected.Icon,
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

	T.Run("with error creating preparation", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Preparation{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("CreatePreparation", mock.Anything, mock.Anything).Return((*models.Preparation)(nil), errors.New("blah"))
		s.preparationDatabase = id

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

		exampleInput := &models.PreparationCreationInput{
			Name:           expected.Name,
			Variant:        expected.Variant,
			Description:    expected.Description,
			AllergyWarning: expected.AllergyWarning,
			Icon:           expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Preparation{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.preparationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("CreatePreparation", mock.Anything, mock.Anything).Return(expected, nil)
		s.preparationDatabase = id

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

		exampleInput := &models.PreparationCreationInput{
			Name:           expected.Name,
			Variant:        expected.Variant,
			Description:    expected.Description,
			AllergyWarning: expected.AllergyWarning,
			Icon:           expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusCreated)
	})
}

func TestPreparationsService_Read(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Preparation{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.preparationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("GetPreparation", mock.Anything, expected.ID).Return(expected, nil)
		s.preparationDatabase = id

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

	T.Run("with no such preparation in database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Preparation{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.preparationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("GetPreparation", mock.Anything, expected.ID).Return((*models.Preparation)(nil), sql.ErrNoRows)
		s.preparationDatabase = id

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

	T.Run("with error fetching preparation from database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Preparation{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.preparationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("GetPreparation", mock.Anything, expected.ID).Return((*models.Preparation)(nil), errors.New("blah"))
		s.preparationDatabase = id

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
		expected := &models.Preparation{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.preparationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("GetPreparation", mock.Anything, expected.ID).Return(expected, nil)
		s.preparationDatabase = id

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

func TestPreparationsService_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Preparation{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.preparationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.preparationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("GetPreparation", mock.Anything, expected.ID).Return(expected, nil)
		id.On("UpdatePreparation", mock.Anything, mock.Anything).Return(nil)
		s.preparationDatabase = id

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

		exampleInput := &models.PreparationUpdateInput{
			Name:           expected.Name,
			Variant:        expected.Variant,
			Description:    expected.Description,
			AllergyWarning: expected.AllergyWarning,
			Icon:           expected.Icon,
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

	T.Run("with no rows fetching preparation", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Preparation{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.preparationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("GetPreparation", mock.Anything, expected.ID).Return((*models.Preparation)(nil), sql.ErrNoRows)
		s.preparationDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.PreparationUpdateInput{
			Name:           expected.Name,
			Variant:        expected.Variant,
			Description:    expected.Description,
			AllergyWarning: expected.AllergyWarning,
			Icon:           expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusNotFound)
	})

	T.Run("with error fetching preparation", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Preparation{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.preparationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("GetPreparation", mock.Anything, expected.ID).Return((*models.Preparation)(nil), errors.New("blah"))
		s.preparationDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.PreparationUpdateInput{
			Name:           expected.Name,
			Variant:        expected.Variant,
			Description:    expected.Description,
			AllergyWarning: expected.AllergyWarning,
			Icon:           expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error updating preparation", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Preparation{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.preparationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.preparationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("GetPreparation", mock.Anything, expected.ID).Return(expected, nil)
		id.On("UpdatePreparation", mock.Anything, mock.Anything).Return(errors.New("blah"))
		s.preparationDatabase = id

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

		exampleInput := &models.PreparationUpdateInput{
			Name:           expected.Name,
			Variant:        expected.Variant,
			Description:    expected.Description,
			AllergyWarning: expected.AllergyWarning,
			Icon:           expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Preparation{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.preparationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.preparationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("GetPreparation", mock.Anything, expected.ID).Return(expected, nil)
		id.On("UpdatePreparation", mock.Anything, mock.Anything).Return(nil)
		s.preparationDatabase = id

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

		exampleInput := &models.PreparationUpdateInput{
			Name:           expected.Name,
			Variant:        expected.Variant,
			Description:    expected.Description,
			AllergyWarning: expected.AllergyWarning,
			Icon:           expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusOK)
	})
}

func TestPreparationsService_Archive(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Preparation{
			ID: 123,
		}

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement").Return()
		s.preparationCounter = mc

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.preparationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("ArchivePreparation", mock.Anything, expected.ID).Return(nil)
		s.preparationDatabase = id

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

	T.Run("with no preparation in database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Preparation{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.preparationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("ArchivePreparation", mock.Anything, expected.ID).Return(sql.ErrNoRows)
		s.preparationDatabase = id

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
		expected := &models.Preparation{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.preparationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.PreparationDataManager{}
		id.On("ArchivePreparation", mock.Anything, expected.ID).Return(errors.New("blah"))
		s.preparationDatabase = id

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
