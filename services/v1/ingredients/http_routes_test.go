package ingredients

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

func TestIngredientsService_List(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.IngredientList{
			Ingredients: []models.Ingredient{
				{
					ID: 123,
				},
			},
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("GetIngredients", mock.Anything, mock.Anything, requestingUser.ID).Return(expected, nil)
		s.ingredientDatabase = id

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

		id := &mockmodels.IngredientDataManager{}
		id.On("GetIngredients", mock.Anything, mock.Anything, requestingUser.ID).Return((*models.IngredientList)(nil), sql.ErrNoRows)
		s.ingredientDatabase = id

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

	T.Run("with error fetching ingredients from database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("GetIngredients", mock.Anything, mock.Anything, requestingUser.ID).Return((*models.IngredientList)(nil), errors.New("blah"))
		s.ingredientDatabase = id

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
		expected := &models.IngredientList{
			Ingredients: []models.Ingredient{},
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("GetIngredients", mock.Anything, mock.Anything, requestingUser.ID).Return(expected, nil)
		s.ingredientDatabase = id

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

func TestIngredientsService_Create(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Ingredient{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.ingredientCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("CreateIngredient", mock.Anything, mock.Anything).Return(expected, nil)
		s.ingredientDatabase = id

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

		exampleInput := &models.IngredientCreationInput{
			Name:              expected.Name,
			Variant:           expected.Variant,
			Description:       expected.Description,
			Warning:           expected.Warning,
			ContainsEgg:       expected.ContainsEgg,
			ContainsDairy:     expected.ContainsDairy,
			ContainsPeanut:    expected.ContainsPeanut,
			ContainsTreeNut:   expected.ContainsTreeNut,
			ContainsSoy:       expected.ContainsSoy,
			ContainsWheat:     expected.ContainsWheat,
			ContainsShellfish: expected.ContainsShellfish,
			ContainsSesame:    expected.ContainsSesame,
			ContainsFish:      expected.ContainsFish,
			ContainsGluten:    expected.ContainsGluten,
			AnimalFlesh:       expected.AnimalFlesh,
			AnimalDerived:     expected.AnimalDerived,
			ConsideredStaple:  expected.ConsideredStaple,
			Icon:              expected.Icon,
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

	T.Run("with error creating ingredient", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Ingredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("CreateIngredient", mock.Anything, mock.Anything).Return((*models.Ingredient)(nil), errors.New("blah"))
		s.ingredientDatabase = id

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

		exampleInput := &models.IngredientCreationInput{
			Name:              expected.Name,
			Variant:           expected.Variant,
			Description:       expected.Description,
			Warning:           expected.Warning,
			ContainsEgg:       expected.ContainsEgg,
			ContainsDairy:     expected.ContainsDairy,
			ContainsPeanut:    expected.ContainsPeanut,
			ContainsTreeNut:   expected.ContainsTreeNut,
			ContainsSoy:       expected.ContainsSoy,
			ContainsWheat:     expected.ContainsWheat,
			ContainsShellfish: expected.ContainsShellfish,
			ContainsSesame:    expected.ContainsSesame,
			ContainsFish:      expected.ContainsFish,
			ContainsGluten:    expected.ContainsGluten,
			AnimalFlesh:       expected.AnimalFlesh,
			AnimalDerived:     expected.AnimalDerived,
			ConsideredStaple:  expected.ConsideredStaple,
			Icon:              expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Ingredient{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.ingredientCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("CreateIngredient", mock.Anything, mock.Anything).Return(expected, nil)
		s.ingredientDatabase = id

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

		exampleInput := &models.IngredientCreationInput{
			Name:              expected.Name,
			Variant:           expected.Variant,
			Description:       expected.Description,
			Warning:           expected.Warning,
			ContainsEgg:       expected.ContainsEgg,
			ContainsDairy:     expected.ContainsDairy,
			ContainsPeanut:    expected.ContainsPeanut,
			ContainsTreeNut:   expected.ContainsTreeNut,
			ContainsSoy:       expected.ContainsSoy,
			ContainsWheat:     expected.ContainsWheat,
			ContainsShellfish: expected.ContainsShellfish,
			ContainsSesame:    expected.ContainsSesame,
			ContainsFish:      expected.ContainsFish,
			ContainsGluten:    expected.ContainsGluten,
			AnimalFlesh:       expected.AnimalFlesh,
			AnimalDerived:     expected.AnimalDerived,
			ConsideredStaple:  expected.ConsideredStaple,
			Icon:              expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusCreated)
	})
}

func TestIngredientsService_Read(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Ingredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.ingredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("GetIngredient", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		s.ingredientDatabase = id

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

	T.Run("with no such ingredient in database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Ingredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.ingredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("GetIngredient", mock.Anything, expected.ID, requestingUser.ID).Return((*models.Ingredient)(nil), sql.ErrNoRows)
		s.ingredientDatabase = id

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

	T.Run("with error fetching ingredient from database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Ingredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.ingredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("GetIngredient", mock.Anything, expected.ID, requestingUser.ID).Return((*models.Ingredient)(nil), errors.New("blah"))
		s.ingredientDatabase = id

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
		expected := &models.Ingredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.ingredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("GetIngredient", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		s.ingredientDatabase = id

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

func TestIngredientsService_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Ingredient{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.ingredientCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.ingredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("GetIngredient", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateIngredient", mock.Anything, mock.Anything).Return(nil)
		s.ingredientDatabase = id

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

		exampleInput := &models.IngredientUpdateInput{
			Name:              expected.Name,
			Variant:           expected.Variant,
			Description:       expected.Description,
			Warning:           expected.Warning,
			ContainsEgg:       expected.ContainsEgg,
			ContainsDairy:     expected.ContainsDairy,
			ContainsPeanut:    expected.ContainsPeanut,
			ContainsTreeNut:   expected.ContainsTreeNut,
			ContainsSoy:       expected.ContainsSoy,
			ContainsWheat:     expected.ContainsWheat,
			ContainsShellfish: expected.ContainsShellfish,
			ContainsSesame:    expected.ContainsSesame,
			ContainsFish:      expected.ContainsFish,
			ContainsGluten:    expected.ContainsGluten,
			AnimalFlesh:       expected.AnimalFlesh,
			AnimalDerived:     expected.AnimalDerived,
			ConsideredStaple:  expected.ConsideredStaple,
			Icon:              expected.Icon,
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

	T.Run("with no rows fetching ingredient", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Ingredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.ingredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("GetIngredient", mock.Anything, expected.ID, requestingUser.ID).Return((*models.Ingredient)(nil), sql.ErrNoRows)
		s.ingredientDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.IngredientUpdateInput{
			Name:              expected.Name,
			Variant:           expected.Variant,
			Description:       expected.Description,
			Warning:           expected.Warning,
			ContainsEgg:       expected.ContainsEgg,
			ContainsDairy:     expected.ContainsDairy,
			ContainsPeanut:    expected.ContainsPeanut,
			ContainsTreeNut:   expected.ContainsTreeNut,
			ContainsSoy:       expected.ContainsSoy,
			ContainsWheat:     expected.ContainsWheat,
			ContainsShellfish: expected.ContainsShellfish,
			ContainsSesame:    expected.ContainsSesame,
			ContainsFish:      expected.ContainsFish,
			ContainsGluten:    expected.ContainsGluten,
			AnimalFlesh:       expected.AnimalFlesh,
			AnimalDerived:     expected.AnimalDerived,
			ConsideredStaple:  expected.ConsideredStaple,
			Icon:              expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusNotFound)
	})

	T.Run("with error fetching ingredient", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Ingredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.ingredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("GetIngredient", mock.Anything, expected.ID, requestingUser.ID).Return((*models.Ingredient)(nil), errors.New("blah"))
		s.ingredientDatabase = id

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		exampleInput := &models.IngredientUpdateInput{
			Name:              expected.Name,
			Variant:           expected.Variant,
			Description:       expected.Description,
			Warning:           expected.Warning,
			ContainsEgg:       expected.ContainsEgg,
			ContainsDairy:     expected.ContainsDairy,
			ContainsPeanut:    expected.ContainsPeanut,
			ContainsTreeNut:   expected.ContainsTreeNut,
			ContainsSoy:       expected.ContainsSoy,
			ContainsWheat:     expected.ContainsWheat,
			ContainsShellfish: expected.ContainsShellfish,
			ContainsSesame:    expected.ContainsSesame,
			ContainsFish:      expected.ContainsFish,
			ContainsGluten:    expected.ContainsGluten,
			AnimalFlesh:       expected.AnimalFlesh,
			AnimalDerived:     expected.AnimalDerived,
			ConsideredStaple:  expected.ConsideredStaple,
			Icon:              expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error updating ingredient", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Ingredient{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.ingredientCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.ingredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("GetIngredient", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateIngredient", mock.Anything, mock.Anything).Return(errors.New("blah"))
		s.ingredientDatabase = id

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

		exampleInput := &models.IngredientUpdateInput{
			Name:              expected.Name,
			Variant:           expected.Variant,
			Description:       expected.Description,
			Warning:           expected.Warning,
			ContainsEgg:       expected.ContainsEgg,
			ContainsDairy:     expected.ContainsDairy,
			ContainsPeanut:    expected.ContainsPeanut,
			ContainsTreeNut:   expected.ContainsTreeNut,
			ContainsSoy:       expected.ContainsSoy,
			ContainsWheat:     expected.ContainsWheat,
			ContainsShellfish: expected.ContainsShellfish,
			ContainsSesame:    expected.ContainsSesame,
			ContainsFish:      expected.ContainsFish,
			ContainsGluten:    expected.ContainsGluten,
			AnimalFlesh:       expected.AnimalFlesh,
			AnimalDerived:     expected.AnimalDerived,
			ConsideredStaple:  expected.ConsideredStaple,
			Icon:              expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Ingredient{
			ID: 123,
		}

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.ingredientCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.ingredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("GetIngredient", mock.Anything, expected.ID, requestingUser.ID).Return(expected, nil)
		id.On("UpdateIngredient", mock.Anything, mock.Anything).Return(nil)
		s.ingredientDatabase = id

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

		exampleInput := &models.IngredientUpdateInput{
			Name:              expected.Name,
			Variant:           expected.Variant,
			Description:       expected.Description,
			Warning:           expected.Warning,
			ContainsEgg:       expected.ContainsEgg,
			ContainsDairy:     expected.ContainsDairy,
			ContainsPeanut:    expected.ContainsPeanut,
			ContainsTreeNut:   expected.ContainsTreeNut,
			ContainsSoy:       expected.ContainsSoy,
			ContainsWheat:     expected.ContainsWheat,
			ContainsShellfish: expected.ContainsShellfish,
			ContainsSesame:    expected.ContainsSesame,
			ContainsFish:      expected.ContainsFish,
			ContainsGluten:    expected.ContainsGluten,
			AnimalFlesh:       expected.AnimalFlesh,
			AnimalDerived:     expected.AnimalDerived,
			ConsideredStaple:  expected.ConsideredStaple,
			Icon:              expected.Icon,
		}
		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, res.Code, http.StatusOK)
	})
}

func TestIngredientsService_Archive(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Ingredient{
			ID: 123,
		}

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.Anything).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement").Return()
		s.ingredientCounter = mc

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.ingredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("ArchiveIngredient", mock.Anything, expected.ID, requestingUser.ID).Return(nil)
		s.ingredientDatabase = id

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

	T.Run("with no ingredient in database", func(t *testing.T) {
		s := buildTestService()

		requestingUser := &models.User{ID: 1}
		expected := &models.Ingredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.ingredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("ArchiveIngredient", mock.Anything, expected.ID, requestingUser.ID).Return(sql.ErrNoRows)
		s.ingredientDatabase = id

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
		expected := &models.Ingredient{
			ID: 123,
		}

		s.userIDFetcher = func(req *http.Request) uint64 {
			return requestingUser.ID
		}

		s.ingredientIDFetcher = func(req *http.Request) uint64 {
			return expected.ID
		}

		id := &mockmodels.IngredientDataManager{}
		id.On("ArchiveIngredient", mock.Anything, expected.ID, requestingUser.ID).Return(errors.New("blah"))
		s.ingredientDatabase = id

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
