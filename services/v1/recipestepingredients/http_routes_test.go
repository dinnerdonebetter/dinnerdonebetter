package recipestepingredients

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

func TestRecipeStepIngredientsService_List(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepIngredientList{
			RecipeStepIngredients: []models.RecipeStepIngredient{
				{
					ID: 123,
				},
			},
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("GetRecipeStepIngredients", mock.Anything, mock.Anything, requestingUser.ID).Return(expected, nil)
		s.recipeStepIngredientDatabase = id

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

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("GetRecipeStepIngredients", mock.Anything, mock.Anything, requestingUser.ID).Return((*models.RecipeStepIngredientList)(nil), sql.ErrNoRows)
		s.recipeStepIngredientDatabase = id

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

	T.Run("with error fetching recipe step ingredients from database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("GetRecipeStepIngredients", mock.Anything, mock.Anything, requestingUser.ID).Return((*models.RecipeStepIngredientList)(nil), errors.New("blah"))
		s.recipeStepIngredientDatabase = id

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
		expected := &models.RecipeStepIngredientList{
			RecipeStepIngredients: []models.RecipeStepIngredient{},
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("GetRecipeStepIngredients", mock.Anything, mock.Anything, requestingUser.ID).Return(expected, nil)
		s.recipeStepIngredientDatabase = id

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

func TestRecipeStepIngredientsService_Create(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepIngredientCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("CreateRecipeStepIngredient", mock.Anything, mock.Anything).Return(expected, nil)
		s.recipeStepIngredientDatabase = id

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

		exampleInput := &models.RecipeStepIngredientCreationInput{
			IngredientID:    expected.IngredientID,
			QuantityType:    expected.QuantityType,
			QuantityValue:   expected.QuantityValue,
			QuantityNotes:   expected.QuantityNotes,
			ProductOfRecipe: expected.ProductOfRecipe,
			IngredientNotes: expected.IngredientNotes,
			RecipeStepID:    expected.RecipeStepID,
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

	T.Run("with error creating recipe step ingredient", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("CreateRecipeStepIngredient", mock.Anything, mock.Anything).Return((*models.RecipeStepIngredient)(nil), errors.New("blah"))
		s.recipeStepIngredientDatabase = id

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

		exampleInput := &models.RecipeStepIngredientCreationInput{
			IngredientID:    expected.IngredientID,
			QuantityType:    expected.QuantityType,
			QuantityValue:   expected.QuantityValue,
			QuantityNotes:   expected.QuantityNotes,
			ProductOfRecipe: expected.ProductOfRecipe,
			IngredientNotes: expected.IngredientNotes,
			RecipeStepID:    expected.RecipeStepID,
		}
		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepIngredientCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("CreateRecipeStepIngredient", mock.Anything, mock.Anything).Return(expected, nil)
		s.recipeStepIngredientDatabase = id

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

		exampleInput := &models.RecipeStepIngredientCreationInput{
			IngredientID:    expected.IngredientID,
			QuantityType:    expected.QuantityType,
			QuantityValue:   expected.QuantityValue,
			QuantityNotes:   expected.QuantityNotes,
			ProductOfRecipe: expected.ProductOfRecipe,
			IngredientNotes: expected.IngredientNotes,
			RecipeStepID:    expected.RecipeStepID,
		}
		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusCreated)
	})
}

func TestRecipeStepIngredientsService_Read(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("GetRecipeStepIngredient", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		s.recipeStepIngredientDatabase = id

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

	T.Run("with no such recipe step ingredient in database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("GetRecipeStepIngredient", mock.Anything, expected.ID, requestingUser.ID).Return((*models.RecipeStepIngredient)(nil), sql.ErrNoRows)
		s.recipeStepIngredientDatabase = id

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

	T.Run("with error fetching recipe step ingredient from database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("GetRecipeStepIngredient", mock.Anything, expected.ID, requestingUser.ID).Return((*models.RecipeStepIngredient)(nil), errors.New("blah"))
		s.recipeStepIngredientDatabase = id

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
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("GetRecipeStepIngredient", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		s.recipeStepIngredientDatabase = id

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

func TestRecipeStepIngredientsService_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepIngredientCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("GetRecipeStepIngredient", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateRecipeStepIngredient", mock.Anything, mock.Anything).Return(nil)
		s.recipeStepIngredientDatabase = id

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

		exampleInput := &models.RecipeStepIngredientUpdateInput{
			IngredientID:    expected.IngredientID,
			QuantityType:    expected.QuantityType,
			QuantityValue:   expected.QuantityValue,
			QuantityNotes:   expected.QuantityNotes,
			ProductOfRecipe: expected.ProductOfRecipe,
			IngredientNotes: expected.IngredientNotes,
			RecipeStepID:    expected.RecipeStepID,
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

	T.Run("with no rows fetching recipe step ingredient", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("GetRecipeStepIngredient", mock.Anything, expected.ID, requestingUser.ID).Return((*models.RecipeStepIngredient)(nil), sql.ErrNoRows)
		s.recipeStepIngredientDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.RecipeStepIngredientUpdateInput{
			IngredientID:    expected.IngredientID,
			QuantityType:    expected.QuantityType,
			QuantityValue:   expected.QuantityValue,
			QuantityNotes:   expected.QuantityNotes,
			ProductOfRecipe: expected.ProductOfRecipe,
			IngredientNotes: expected.IngredientNotes,
			RecipeStepID:    expected.RecipeStepID,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusNotFound)
	})

	T.Run("with error fetching recipe step ingredient", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("GetRecipeStepIngredient", mock.Anything, expected.ID, requestingUser.ID).Return((*models.RecipeStepIngredient)(nil), errors.New("blah"))
		s.recipeStepIngredientDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.RecipeStepIngredientUpdateInput{
			IngredientID:    expected.IngredientID,
			QuantityType:    expected.QuantityType,
			QuantityValue:   expected.QuantityValue,
			QuantityNotes:   expected.QuantityNotes,
			ProductOfRecipe: expected.ProductOfRecipe,
			IngredientNotes: expected.IngredientNotes,
			RecipeStepID:    expected.RecipeStepID,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error updating recipe step ingredient", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepIngredientCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("GetRecipeStepIngredient", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateRecipeStepIngredient", mock.Anything, mock.Anything).Return(errors.New("blah"))
		s.recipeStepIngredientDatabase = id

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

		exampleInput := &models.RecipeStepIngredientUpdateInput{
			IngredientID:    expected.IngredientID,
			QuantityType:    expected.QuantityType,
			QuantityValue:   expected.QuantityValue,
			QuantityNotes:   expected.QuantityNotes,
			ProductOfRecipe: expected.ProductOfRecipe,
			IngredientNotes: expected.IngredientNotes,
			RecipeStepID:    expected.RecipeStepID,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.recipeStepIngredientCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("GetRecipeStepIngredient", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateRecipeStepIngredient", mock.Anything, mock.Anything).Return(nil)
		s.recipeStepIngredientDatabase = id

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

		exampleInput := &models.RecipeStepIngredientUpdateInput{
			IngredientID:    expected.IngredientID,
			QuantityType:    expected.QuantityType,
			QuantityValue:   expected.QuantityValue,
			QuantityNotes:   expected.QuantityNotes,
			ProductOfRecipe: expected.ProductOfRecipe,
			IngredientNotes: expected.IngredientNotes,
			RecipeStepID:    expected.RecipeStepID,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusOK)
	})
}

func TestRecipeStepIngredientsService_Archive(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement").Return()
		s.recipeStepIngredientCounter = mc

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("ArchiveRecipeStepIngredient", mock.Anything, expected.ID, requestingUser.ID).Return(nil)
		s.recipeStepIngredientDatabase = id

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

	T.Run("with no recipe step ingredient in database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("ArchiveRecipeStepIngredient", mock.Anything, expected.ID, requestingUser.ID).Return(sql.ErrNoRows)
		s.recipeStepIngredientDatabase = id

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
		expected := &models.RecipeStepIngredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.recipeStepIngredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.RecipeStepIngredientDataManager{}
		id.On("ArchiveRecipeStepIngredient", mock.Anything, expected.ID, requestingUser.ID).Return(errors.New("blah"))
		s.recipeStepIngredientDatabase = id

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
