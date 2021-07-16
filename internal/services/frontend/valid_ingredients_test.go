package frontend

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	database "gitlab.com/prixfixe/prixfixe/internal/database"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
	testutils "gitlab.com/prixfixe/prixfixe/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_fetchValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"GetValidIngredient",
			testutils.ContextMatcher,
			exampleValidIngredient.ID,
		).Return(exampleValidIngredient, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		actual, err := s.service.fetchValidIngredient(s.ctx, req)
		assert.Equal(t, exampleValidIngredient, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with fake mode", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.useFakeData = true

		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		actual, err := s.service.fetchValidIngredient(s.ctx, req)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with error fetching valid ingredient", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"GetValidIngredient",
			testutils.ContextMatcher,
			exampleValidIngredient.ID,
		).Return((*types.ValidIngredient)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		actual, err := s.service.fetchValidIngredient(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachValidIngredientCreationInputToRequest(input *types.ValidIngredientCreationInput) *http.Request {
	form := url.Values{
		validIngredientCreationInputNameFormKey:              {anyToString(input.Name)},
		validIngredientCreationInputVariantFormKey:           {anyToString(input.Variant)},
		validIngredientCreationInputDescriptionFormKey:       {anyToString(input.Description)},
		validIngredientCreationInputWarningFormKey:           {anyToString(input.Warning)},
		validIngredientCreationInputContainsEggFormKey:       {anyToString(input.ContainsEgg)},
		validIngredientCreationInputContainsDairyFormKey:     {anyToString(input.ContainsDairy)},
		validIngredientCreationInputContainsPeanutFormKey:    {anyToString(input.ContainsPeanut)},
		validIngredientCreationInputContainsTreeNutFormKey:   {anyToString(input.ContainsTreeNut)},
		validIngredientCreationInputContainsSoyFormKey:       {anyToString(input.ContainsSoy)},
		validIngredientCreationInputContainsWheatFormKey:     {anyToString(input.ContainsWheat)},
		validIngredientCreationInputContainsShellfishFormKey: {anyToString(input.ContainsShellfish)},
		validIngredientCreationInputContainsSesameFormKey:    {anyToString(input.ContainsSesame)},
		validIngredientCreationInputContainsFishFormKey:      {anyToString(input.ContainsFish)},
		validIngredientCreationInputContainsGlutenFormKey:    {anyToString(input.ContainsGluten)},
		validIngredientCreationInputAnimalFleshFormKey:       {anyToString(input.AnimalFlesh)},
		validIngredientCreationInputAnimalDerivedFormKey:     {anyToString(input.AnimalDerived)},
		validIngredientCreationInputVolumetricFormKey:        {anyToString(input.Volumetric)},
		validIngredientCreationInputIconPathFormKey:          {anyToString(input.IconPath)},
	}

	return httptest.NewRequest(http.MethodPost, "/valid_ingredients", strings.NewReader(form.Encode()))
}

func TestService_buildValidIngredientCreatorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		s.service.buildValidIngredientCreatorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		s.service.buildValidIngredientCreatorView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		s.service.buildValidIngredientCreatorView(false)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		s.service.buildValidIngredientCreatorView(true)(res, req)
	})

	T.Run("without base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		s.service.buildValidIngredientCreatorView(false)(res, req)
	})
}

func TestService_parseFormEncodedValidIngredientCreationInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		expected := fakes.BuildFakeValidIngredientCreationInput()
		req := attachValidIngredientCreationInputToRequest(expected)

		actual := s.service.parseFormEncodedValidIngredientCreationInput(s.ctx, req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with error extracting form from request", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedValidIngredientCreationInput(s.ctx, req)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ValidIngredientCreationInput{}
		req := attachValidIngredientCreationInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedValidIngredientCreationInput(s.ctx, req)
		assert.Nil(t, actual)
	})
}

func TestService_handleValidIngredientCreationRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		exampleInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)

		res := httptest.NewRecorder()
		req := attachValidIngredientCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"CreateValidIngredient",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleValidIngredient, nil)
		s.service.dataStore = mockDB

		s.service.handleValidIngredientCreationRequest(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.NotEmpty(t, res.Header().Get(htmxRedirectionHeader))

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		exampleInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachValidIngredientCreationInputToRequest(exampleInput)

		s.service.handleValidIngredientCreationRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		exampleInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)

		res := httptest.NewRecorder()
		req := attachValidIngredientCreationInputToRequest(&types.ValidIngredientCreationInput{})

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"CreateValidIngredient",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleValidIngredient, nil)
		s.service.dataStore = mockDB

		s.service.handleValidIngredientCreationRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error creating valid ingredient in database", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		exampleInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)

		res := httptest.NewRecorder()
		req := attachValidIngredientCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"CreateValidIngredient",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return((*types.ValidIngredient)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		s.service.handleValidIngredientCreationRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildValidIngredientEditorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"GetValidIngredient",
			testutils.ContextMatcher,
			exampleValidIngredient.ID,
		).Return(exampleValidIngredient, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		s.service.buildValidIngredientEditorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"GetValidIngredient",
			testutils.ContextMatcher,
			exampleValidIngredient.ID,
		).Return(exampleValidIngredient, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		s.service.buildValidIngredientEditorView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		s.service.buildValidIngredientEditorView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching valid ingredient", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"GetValidIngredient",
			testutils.ContextMatcher,
			exampleValidIngredient.ID,
		).Return((*types.ValidIngredient)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		s.service.buildValidIngredientEditorView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_fetchValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientList := fakes.BuildFakeValidIngredientList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"GetValidIngredients",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidIngredientList, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		actual, err := s.service.fetchValidIngredients(s.ctx, req)
		assert.Equal(t, exampleValidIngredientList, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with fake mode", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.useFakeData = true

		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		actual, err := s.service.fetchValidIngredients(s.ctx, req)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"GetValidIngredients",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidIngredientList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		actual, err := s.service.fetchValidIngredients(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildValidIngredientsTableView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientList := fakes.BuildFakeValidIngredientList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"GetValidIngredients",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidIngredientList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		s.service.buildValidIngredientsTableView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientList := fakes.BuildFakeValidIngredientList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"GetValidIngredients",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidIngredientList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		s.service.buildValidIngredientsTableView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		s.service.buildValidIngredientsTableView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"GetValidIngredients",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidIngredientList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredients", nil)

		s.service.buildValidIngredientsTableView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachValidIngredientUpdateInputToRequest(input *types.ValidIngredientUpdateInput) *http.Request {
	form := url.Values{
		validIngredientUpdateInputNameFormKey:              {anyToString(input.Name)},
		validIngredientUpdateInputVariantFormKey:           {anyToString(input.Variant)},
		validIngredientUpdateInputDescriptionFormKey:       {anyToString(input.Description)},
		validIngredientUpdateInputWarningFormKey:           {anyToString(input.Warning)},
		validIngredientUpdateInputContainsEggFormKey:       {anyToString(input.ContainsEgg)},
		validIngredientUpdateInputContainsDairyFormKey:     {anyToString(input.ContainsDairy)},
		validIngredientUpdateInputContainsPeanutFormKey:    {anyToString(input.ContainsPeanut)},
		validIngredientUpdateInputContainsTreeNutFormKey:   {anyToString(input.ContainsTreeNut)},
		validIngredientUpdateInputContainsSoyFormKey:       {anyToString(input.ContainsSoy)},
		validIngredientUpdateInputContainsWheatFormKey:     {anyToString(input.ContainsWheat)},
		validIngredientUpdateInputContainsShellfishFormKey: {anyToString(input.ContainsShellfish)},
		validIngredientUpdateInputContainsSesameFormKey:    {anyToString(input.ContainsSesame)},
		validIngredientUpdateInputContainsFishFormKey:      {anyToString(input.ContainsFish)},
		validIngredientUpdateInputContainsGlutenFormKey:    {anyToString(input.ContainsGluten)},
		validIngredientUpdateInputAnimalFleshFormKey:       {anyToString(input.AnimalFlesh)},
		validIngredientUpdateInputAnimalDerivedFormKey:     {anyToString(input.AnimalDerived)},
		validIngredientUpdateInputVolumetricFormKey:        {anyToString(input.Volumetric)},
		validIngredientUpdateInputIconPathFormKey:          {anyToString(input.IconPath)},
	}

	return httptest.NewRequest(http.MethodPost, "/valid_ingredients", strings.NewReader(form.Encode()))
}

func TestService_parseFormEncodedValidIngredientUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		expected := fakes.BuildFakeValidIngredientUpdateInputFromValidIngredient(exampleValidIngredient)

		req := attachValidIngredientUpdateInputToRequest(expected)

		actual := s.service.parseFormEncodedValidIngredientUpdateInput(s.ctx, req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedValidIngredientUpdateInput(s.ctx, req)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input attached to valid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ValidIngredientUpdateInput{}

		req := attachValidIngredientUpdateInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedValidIngredientUpdateInput(s.ctx, req)
		assert.Nil(t, actual)
	})
}

func TestService_handleValidIngredientUpdateRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		exampleInput := fakes.BuildFakeValidIngredientUpdateInputFromValidIngredient(exampleValidIngredient)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"GetValidIngredient",
			testutils.ContextMatcher,
			exampleValidIngredient.ID,
		).Return(exampleValidIngredient, nil)

		mockDB.ValidIngredientDataManager.On(
			"UpdateValidIngredient",
			testutils.ContextMatcher,
			exampleValidIngredient,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachValidIngredientUpdateInputToRequest(exampleInput)

		s.service.handleValidIngredientUpdateRequest(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		exampleInput := fakes.BuildFakeValidIngredientUpdateInputFromValidIngredient(exampleValidIngredient)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachValidIngredientUpdateInputToRequest(exampleInput)

		s.service.handleValidIngredientUpdateRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ValidIngredientUpdateInput{}

		res := httptest.NewRecorder()
		req := attachValidIngredientUpdateInputToRequest(exampleInput)

		s.service.handleValidIngredientUpdateRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		exampleInput := fakes.BuildFakeValidIngredientUpdateInputFromValidIngredient(exampleValidIngredient)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"GetValidIngredient",
			testutils.ContextMatcher,
			exampleValidIngredient.ID,
		).Return((*types.ValidIngredient)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachValidIngredientUpdateInputToRequest(exampleInput)

		s.service.handleValidIngredientUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error updating data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		exampleInput := fakes.BuildFakeValidIngredientUpdateInputFromValidIngredient(exampleValidIngredient)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"GetValidIngredient",
			testutils.ContextMatcher,
			exampleValidIngredient.ID,
		).Return(exampleValidIngredient, nil)

		mockDB.ValidIngredientDataManager.On(
			"UpdateValidIngredient",
			testutils.ContextMatcher,
			exampleValidIngredient,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachValidIngredientUpdateInputToRequest(exampleInput)

		s.service.handleValidIngredientUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_handleValidIngredientArchiveRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		exampleValidIngredientList := fakes.BuildFakeValidIngredientList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"ArchiveValidIngredient",
			testutils.ContextMatcher,
			exampleValidIngredient.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.ValidIngredientDataManager.On(
			"GetValidIngredients",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidIngredientList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_ingredients", nil)

		s.service.handleValidIngredientArchiveRequest(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_ingredients", nil)

		s.service.handleValidIngredientArchiveRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error archiving valid ingredient", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"ArchiveValidIngredient",
			testutils.ContextMatcher,
			exampleValidIngredient.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_ingredients", nil)

		s.service.handleValidIngredientArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error retrieving new list of valid ingredients", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		s.service.validIngredientIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredient.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientDataManager.On(
			"ArchiveValidIngredient",
			testutils.ContextMatcher,
			exampleValidIngredient.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.ValidIngredientDataManager.On(
			"GetValidIngredients",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidIngredientList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_ingredients", nil)

		s.service.handleValidIngredientArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
