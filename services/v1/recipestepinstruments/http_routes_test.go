package recipestepinstruments

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

func TestRecipeStepInstrumentsService_List(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepInstrumentList{
			RecipeStepInstruments: []models.RecipeStepInstrument{
				{
					ID: 123,
				},
			},
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("GetRecipeStepInstruments", mock.Anything, mock.Anything, requestingUser.ID).Return(expected, nil)
		s.recipeStepInstrumentDatabase = id

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

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("GetRecipeStepInstruments", mock.Anything, mock.Anything, requestingUser.ID).Return((*models.RecipeStepInstrumentList)(nil), sql.ErrNoRows)
		s.recipeStepInstrumentDatabase = id

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

	T.Run("with error fetching recipe step instruments from database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("GetRecipeStepInstruments", mock.Anything, mock.Anything, requestingUser.ID).Return((*models.RecipeStepInstrumentList)(nil), errors.New("blah"))
		s.recipeStepInstrumentDatabase = id

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
		expected := &models.RecipeStepInstrumentList{
			RecipeStepInstruments: []models.RecipeStepInstrument{},
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("GetRecipeStepInstruments", mock.Anything, mock.Anything, requestingUser.ID).Return(expected, nil)
		s.recipeStepInstrumentDatabase = id

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

func TestRecipeStepInstrumentsService_Create(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepInstrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("CreateRecipeStepInstrument", mock.Anything, mock.Anything).Return(expected, nil)
		s.recipeStepInstrumentDatabase = id

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

		exampleInput := &models.RecipeStepInstrumentCreationInput{
			InstrumentID: expected.InstrumentID,
			RecipeStepID: expected.RecipeStepID,
			Notes:        expected.Notes,
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

	T.Run("with error creating recipe step instrument", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("CreateRecipeStepInstrument", mock.Anything, mock.Anything).Return((*models.RecipeStepInstrument)(nil), errors.New("blah"))
		s.recipeStepInstrumentDatabase = id

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

		exampleInput := &models.RecipeStepInstrumentCreationInput{
			InstrumentID: expected.InstrumentID,
			RecipeStepID: expected.RecipeStepID,
			Notes:        expected.Notes,
		}
		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepInstrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("CreateRecipeStepInstrument", mock.Anything, mock.Anything).Return(expected, nil)
		s.recipeStepInstrumentDatabase = id

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

		exampleInput := &models.RecipeStepInstrumentCreationInput{
			InstrumentID: expected.InstrumentID,
			RecipeStepID: expected.RecipeStepID,
			Notes:        expected.Notes,
		}
		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusCreated)
	})
}

func TestRecipeStepInstrumentsService_Read(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("GetRecipeStepInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		s.recipeStepInstrumentDatabase = id

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

	T.Run("with no such recipe step instrument in database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("GetRecipeStepInstrument", mock.Anything, expected.ID, requestingUser.ID).Return((*models.RecipeStepInstrument)(nil), sql.ErrNoRows)
		s.recipeStepInstrumentDatabase = id

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

	T.Run("with error fetching recipe step instrument from database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("GetRecipeStepInstrument", mock.Anything, expected.ID, requestingUser.ID).Return((*models.RecipeStepInstrument)(nil), errors.New("blah"))
		s.recipeStepInstrumentDatabase = id

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
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("GetRecipeStepInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		s.recipeStepInstrumentDatabase = id

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

func TestRecipeStepInstrumentsService_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepInstrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("GetRecipeStepInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateRecipeStepInstrument", mock.Anything, mock.Anything).Return(nil)
		s.recipeStepInstrumentDatabase = id

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

		exampleInput := &models.RecipeStepInstrumentUpdateInput{
			InstrumentID: expected.InstrumentID,
			RecipeStepID: expected.RecipeStepID,
			Notes:        expected.Notes,
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

	T.Run("with no rows fetching recipe step instrument", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("GetRecipeStepInstrument", mock.Anything, expected.ID, requestingUser.ID).Return((*models.RecipeStepInstrument)(nil), sql.ErrNoRows)
		s.recipeStepInstrumentDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.RecipeStepInstrumentUpdateInput{
			InstrumentID: expected.InstrumentID,
			RecipeStepID: expected.RecipeStepID,
			Notes:        expected.Notes,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusNotFound)
	})

	T.Run("with error fetching recipe step instrument", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("GetRecipeStepInstrument", mock.Anything, expected.ID, requestingUser.ID).Return((*models.RecipeStepInstrument)(nil), errors.New("blah"))
		s.recipeStepInstrumentDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.RecipeStepInstrumentUpdateInput{
			InstrumentID: expected.InstrumentID,
			RecipeStepID: expected.RecipeStepID,
			Notes:        expected.Notes,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error updating recipe step instrument", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepInstrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("GetRecipeStepInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateRecipeStepInstrument", mock.Anything, mock.Anything).Return(errors.New("blah"))
		s.recipeStepInstrumentDatabase = id

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

		exampleInput := &models.RecipeStepInstrumentUpdateInput{
			InstrumentID: expected.InstrumentID,
			RecipeStepID: expected.RecipeStepID,
			Notes:        expected.Notes,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepInstrumentCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("GetRecipeStepInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateRecipeStepInstrument", mock.Anything, mock.Anything).Return(nil)
		s.recipeStepInstrumentDatabase = id

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

		exampleInput := &models.RecipeStepInstrumentUpdateInput{
			InstrumentID: expected.InstrumentID,
			RecipeStepID: expected.RecipeStepID,
			Notes:        expected.Notes,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusOK)
	})
}

func TestRecipeStepInstrumentsService_Archive(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement").Return()
		s.recipeStepInstrumentCounter = mc

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("ArchiveRecipeStepInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(nil)
		s.recipeStepInstrumentDatabase = id

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

	T.Run("with no recipe step instrument in database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("ArchiveRecipeStepInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(sql.ErrNoRows)
		s.recipeStepInstrumentDatabase = id

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
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepInstrumentIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepInstrumentDataManager{}
		id.On("ArchiveRecipeStepInstrument", mock.Anything, expected.ID, requestingUser.ID).Return(errors.New("blah"))
		s.recipeStepInstrumentDatabase = id

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
