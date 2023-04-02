package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestValidIngredientStateIngredients(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validIngredientStateIngredientsTestSuite))
}

type validIngredientStateIngredientsBaseSuite struct {
	suite.Suite

	ctx                                   context.Context
	exampleValidIngredientStateIngredient *types.ValidIngredientStateIngredient
}

var _ suite.SetupTestSuite = (*validIngredientStateIngredientsBaseSuite)(nil)

func (s *validIngredientStateIngredientsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidIngredientStateIngredient = fakes.BuildFakeValidIngredientStateIngredient()
}

type validIngredientStateIngredientsTestSuite struct {
	suite.Suite

	validIngredientStateIngredientsBaseSuite
}

func (s *validIngredientStateIngredientsTestSuite) TestClient_GetValidIngredientStateIngredient() {
	const expectedPathFormat = "/api/v1/valid_ingredient_state_ingredients/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidIngredientStateIngredient.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientStateIngredient)
		actual, err := c.GetValidIngredientStateIngredient(s.ctx, s.exampleValidIngredientStateIngredient.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientStateIngredient, actual)
	})

	s.Run("with invalid valid ingredient preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientStateIngredient(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientStateIngredient(s.ctx, s.exampleValidIngredientStateIngredient.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidIngredientStateIngredient.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientStateIngredient(s.ctx, s.exampleValidIngredientStateIngredient.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientStateIngredientsTestSuite) TestClient_GetValidIngredientStateIngredients() {
	const expectedPath = "/api/v1/valid_ingredient_state_ingredients"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleValidIngredientStateIngredientList := fakes.BuildFakeValidIngredientStateIngredientList()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidIngredientStateIngredientList)
		actual, err := c.GetValidIngredientStateIngredients(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientStateIngredientList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientStateIngredients(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientStateIngredients(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientMeasurementUnitsTestSuite) TestClient_GetValidIngredientStateIngredientsForIngredient() {
	const expectedPath = "/api/v1/valid_ingredient_state_ingredients/by_ingredient/%s"

	exampleValidIngredient := fakes.BuildFakeValidIngredient()

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleValidIngredientMeasurementUnitList := fakes.BuildFakeValidIngredientStateIngredientList()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath, exampleValidIngredient.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidIngredientMeasurementUnitList)
		actual, err := c.GetValidIngredientStateIngredientsForIngredient(s.ctx, exampleValidIngredient.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientMeasurementUnitList, actual)
	})

	s.Run("with invalid ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientStateIngredientsForIngredient(s.ctx, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientStateIngredientsForIngredient(s.ctx, exampleValidIngredient.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath, exampleValidIngredient.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientStateIngredientsForIngredient(s.ctx, exampleValidIngredient.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientMeasurementUnitsTestSuite) TestClient_GetValidIngredientStateIngredientsForPreparation() {
	const expectedPath = "/api/v1/valid_ingredient_state_ingredients/by_ingredient_state/%s"

	exampleValidPreparation := fakes.BuildFakeValidPreparation()

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleValidIngredientMeasurementUnitList := fakes.BuildFakeValidIngredientStateIngredientList()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath, exampleValidPreparation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidIngredientMeasurementUnitList)
		actual, err := c.GetValidIngredientStateIngredientsForIngredientState(s.ctx, exampleValidPreparation.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientMeasurementUnitList, actual)
	})

	s.Run("with invalid ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientStateIngredientsForIngredientState(s.ctx, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientStateIngredientsForIngredientState(s.ctx, exampleValidPreparation.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath, exampleValidPreparation.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientStateIngredientsForIngredientState(s.ctx, exampleValidPreparation.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientStateIngredientsTestSuite) TestClient_CreateValidIngredientStateIngredient() {
	const expectedPath = "/api/v1/valid_ingredient_state_ingredients"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidIngredientStateIngredientCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientStateIngredient)

		actual, err := c.CreateValidIngredientStateIngredient(s.ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientStateIngredient, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidIngredientStateIngredient(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidIngredientStateIngredientCreationRequestInput{}

		actual, err := c.CreateValidIngredientStateIngredient(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientCreationRequestInput(s.exampleValidIngredientStateIngredient)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidIngredientStateIngredient(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientCreationRequestInput(s.exampleValidIngredientStateIngredient)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidIngredientStateIngredient(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientStateIngredientsTestSuite) TestClient_UpdateValidIngredientStateIngredient() {
	const expectedPathFormat = "/api/v1/valid_ingredient_state_ingredients/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidIngredientStateIngredient.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientStateIngredient)

		err := c.UpdateValidIngredientStateIngredient(s.ctx, s.exampleValidIngredientStateIngredient)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidIngredientStateIngredient(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidIngredientStateIngredient(s.ctx, s.exampleValidIngredientStateIngredient)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidIngredientStateIngredient(s.ctx, s.exampleValidIngredientStateIngredient)
		assert.Error(t, err)
	})
}

func (s *validIngredientStateIngredientsTestSuite) TestClient_ArchiveValidIngredientStateIngredient() {
	const expectedPathFormat = "/api/v1/valid_ingredient_state_ingredients/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidIngredientStateIngredient.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveValidIngredientStateIngredient(s.ctx, s.exampleValidIngredientStateIngredient.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid ingredient preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidIngredientStateIngredient(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidIngredientStateIngredient(s.ctx, s.exampleValidIngredientStateIngredient.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidIngredientStateIngredient(s.ctx, s.exampleValidIngredientStateIngredient.ID)
		assert.Error(t, err)
	})
}
