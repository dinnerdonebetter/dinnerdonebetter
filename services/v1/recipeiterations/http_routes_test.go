package recipeiterations

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

func TestRecipeIterationsService_List(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeIterationList{
			RecipeIterations: []models.RecipeIteration{
				{
					ID: 123,
				},
			},
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("GetRecipeIterations", mock.Anything, mock.Anything, requestingUser.ID).Return(expected, nil)
		s.recipeIterationDatabase = id

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

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("GetRecipeIterations", mock.Anything, mock.Anything, requestingUser.ID).Return((*models.RecipeIterationList)(nil), sql.ErrNoRows)
		s.recipeIterationDatabase = id

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

	T.Run("with error fetching recipe iterations from database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("GetRecipeIterations", mock.Anything, mock.Anything, requestingUser.ID).Return((*models.RecipeIterationList)(nil), errors.New("blah"))
		s.recipeIterationDatabase = id

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
		expected := &models.RecipeIterationList{
			RecipeIterations: []models.RecipeIteration{},
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("GetRecipeIterations", mock.Anything, mock.Anything, requestingUser.ID).Return(expected, nil)
		s.recipeIterationDatabase = id

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

func TestRecipeIterationsService_Create(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeIteration{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeIterationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("CreateRecipeIteration", mock.Anything, mock.Anything).Return(expected, nil)
		s.recipeIterationDatabase = id

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

		exampleInput := &models.RecipeIterationCreationInput{
			RecipeID:            expected.RecipeID,
			EndDifficultyRating: expected.EndDifficultyRating,
			EndComplexityRating: expected.EndComplexityRating,
			EndTasteRating:      expected.EndTasteRating,
			EndOverallRating:    expected.EndOverallRating,
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

	T.Run("with error creating recipe iteration", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeIteration{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("CreateRecipeIteration", mock.Anything, mock.Anything).Return((*models.RecipeIteration)(nil), errors.New("blah"))
		s.recipeIterationDatabase = id

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

		exampleInput := &models.RecipeIterationCreationInput{
			RecipeID:            expected.RecipeID,
			EndDifficultyRating: expected.EndDifficultyRating,
			EndComplexityRating: expected.EndComplexityRating,
			EndTasteRating:      expected.EndTasteRating,
			EndOverallRating:    expected.EndOverallRating,
		}
		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeIteration{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeIterationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("CreateRecipeIteration", mock.Anything, mock.Anything).Return(expected, nil)
		s.recipeIterationDatabase = id

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

		exampleInput := &models.RecipeIterationCreationInput{
			RecipeID:            expected.RecipeID,
			EndDifficultyRating: expected.EndDifficultyRating,
			EndComplexityRating: expected.EndComplexityRating,
			EndTasteRating:      expected.EndTasteRating,
			EndOverallRating:    expected.EndOverallRating,
		}
		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusCreated)
	})
}

func TestRecipeIterationsService_Read(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeIteration{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("GetRecipeIteration", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		s.recipeIterationDatabase = id

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

	T.Run("with no such recipe iteration in database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeIteration{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("GetRecipeIteration", mock.Anything, expected.ID, requestingUser.ID).Return((*models.RecipeIteration)(nil), sql.ErrNoRows)
		s.recipeIterationDatabase = id

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

	T.Run("with error fetching recipe iteration from database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeIteration{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("GetRecipeIteration", mock.Anything, expected.ID, requestingUser.ID).Return((*models.RecipeIteration)(nil), errors.New("blah"))
		s.recipeIterationDatabase = id

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
		expected := &models.RecipeIteration{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("GetRecipeIteration", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		s.recipeIterationDatabase = id

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

func TestRecipeIterationsService_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeIteration{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeIterationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("GetRecipeIteration", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateRecipeIteration", mock.Anything, mock.Anything).Return(nil)
		s.recipeIterationDatabase = id

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

		exampleInput := &models.RecipeIterationUpdateInput{
			RecipeID:            expected.RecipeID,
			EndDifficultyRating: expected.EndDifficultyRating,
			EndComplexityRating: expected.EndComplexityRating,
			EndTasteRating:      expected.EndTasteRating,
			EndOverallRating:    expected.EndOverallRating,
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

	T.Run("with no rows fetching recipe iteration", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeIteration{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("GetRecipeIteration", mock.Anything, expected.ID, requestingUser.ID).Return((*models.RecipeIteration)(nil), sql.ErrNoRows)
		s.recipeIterationDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.RecipeIterationUpdateInput{
			RecipeID:            expected.RecipeID,
			EndDifficultyRating: expected.EndDifficultyRating,
			EndComplexityRating: expected.EndComplexityRating,
			EndTasteRating:      expected.EndTasteRating,
			EndOverallRating:    expected.EndOverallRating,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusNotFound)
	})

	T.Run("with error fetching recipe iteration", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeIteration{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("GetRecipeIteration", mock.Anything, expected.ID, requestingUser.ID).Return((*models.RecipeIteration)(nil), errors.New("blah"))
		s.recipeIterationDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.RecipeIterationUpdateInput{
			RecipeID:            expected.RecipeID,
			EndDifficultyRating: expected.EndDifficultyRating,
			EndComplexityRating: expected.EndComplexityRating,
			EndTasteRating:      expected.EndTasteRating,
			EndOverallRating:    expected.EndOverallRating,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error updating recipe iteration", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeIteration{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeIterationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("GetRecipeIteration", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateRecipeIteration", mock.Anything, mock.Anything).Return(errors.New("blah"))
		s.recipeIterationDatabase = id

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

		exampleInput := &models.RecipeIterationUpdateInput{
			RecipeID:            expected.RecipeID,
			EndDifficultyRating: expected.EndDifficultyRating,
			EndComplexityRating: expected.EndComplexityRating,
			EndTasteRating:      expected.EndTasteRating,
			EndOverallRating:    expected.EndOverallRating,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeIteration{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeIterationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("GetRecipeIteration", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateRecipeIteration", mock.Anything, mock.Anything).Return(nil)
		s.recipeIterationDatabase = id

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

		exampleInput := &models.RecipeIterationUpdateInput{
			RecipeID:            expected.RecipeID,
			EndDifficultyRating: expected.EndDifficultyRating,
			EndComplexityRating: expected.EndComplexityRating,
			EndTasteRating:      expected.EndTasteRating,
			EndOverallRating:    expected.EndOverallRating,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusOK)
	})
}

func TestRecipeIterationsService_Archive(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeIteration{
			ID: 123,
		}

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement").Return()
		s.recipeIterationCounter = mc

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("ArchiveRecipeIteration", mock.Anything, expected.ID, requestingUser.ID).Return(nil)
		s.recipeIterationDatabase = id

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

	T.Run("with no recipe iteration in database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeIteration{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("ArchiveRecipeIteration", mock.Anything, expected.ID, requestingUser.ID).Return(sql.ErrNoRows)
		s.recipeIterationDatabase = id

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
		expected := &models.RecipeIteration{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeIterationIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeIterationDataManager{}
		id.On("ArchiveRecipeIteration", mock.Anything, expected.ID, requestingUser.ID).Return(errors.New("blah"))
		s.recipeIterationDatabase = id

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
