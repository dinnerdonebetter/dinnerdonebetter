package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"
)

func TestValidIngredients(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validIngredientsTestSuite))
}

type validIngredientsBaseSuite struct {
	suite.Suite

	ctx                    context.Context
	exampleValidIngredient *types.ValidIngredient
}

var _ suite.SetupTestSuite = (*validIngredientsBaseSuite)(nil)

func (s *validIngredientsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidIngredient = fakes.BuildFakeValidIngredient()
}

type validIngredientsTestSuite struct {
	suite.Suite

	validIngredientsBaseSuite
}

func (s *validIngredientsTestSuite) TestClient_GetValidIngredient() {
	const expectedPathFormat = "/api/v1/valid_ingredients/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidIngredient.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredient)
		actual, err := c.GetValidIngredient(s.ctx, s.exampleValidIngredient.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredient, actual)
	})

	s.Run("with invalid valid ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredient(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredient(s.ctx, s.exampleValidIngredient.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidIngredient.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredient(s.ctx, s.exampleValidIngredient.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientsTestSuite) TestClient_GetRandomValidIngredient() {
	const expectedPath = "/api/v1/valid_ingredients/random"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredient)
		actual, err := c.GetRandomValidIngredient(s.ctx)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredient, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRandomValidIngredient(s.ctx)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRandomValidIngredient(s.ctx)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientsTestSuite) TestClient_GetValidIngredients() {
	const expectedPath = "/api/v1/valid_ingredients"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleValidIngredientList := fakes.BuildFakeValidIngredientList()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidIngredientList)
		actual, err := c.GetValidIngredients(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredients(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredients(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientsTestSuite) TestClient_SearchValidIngredients() {
	const expectedPath = "/api/v1/valid_ingredients/search"

	exampleQuery := "whatever"

	s.Run("standard", func() {
		t := s.T()

		exampleValidIngredientList := fakes.BuildFakeValidIngredientList()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&q=whatever", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidIngredientList.ValidIngredients)
		actual, err := c.SearchValidIngredients(s.ctx, exampleQuery, 0)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientList.ValidIngredients, actual)
	})

	s.Run("with empty query", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.SearchValidIngredients(s.ctx, "", 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.SearchValidIngredients(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with bad response from server", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&q=whatever", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchValidIngredients(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientsTestSuite) TestClient_CreateValidIngredient() {
	const expectedPath = "/api/v1/valid_ingredients"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidIngredientCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredient)

		actual, err := c.CreateValidIngredient(s.ctx, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleValidIngredient, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidIngredient(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidIngredientCreationRequestInput{}

		actual, err := c.CreateValidIngredient(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(s.exampleValidIngredient)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidIngredient(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(s.exampleValidIngredient)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidIngredient(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientsTestSuite) TestClient_UpdateValidIngredient() {
	const expectedPathFormat = "/api/v1/valid_ingredients/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidIngredient.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredient)

		err := c.UpdateValidIngredient(s.ctx, s.exampleValidIngredient)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidIngredient(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidIngredient(s.ctx, s.exampleValidIngredient)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidIngredient(s.ctx, s.exampleValidIngredient)
		assert.Error(t, err)
	})
}

func (s *validIngredientsTestSuite) TestClient_ArchiveValidIngredient() {
	const expectedPathFormat = "/api/v1/valid_ingredients/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidIngredient.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveValidIngredient(s.ctx, s.exampleValidIngredient.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidIngredient(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidIngredient(s.ctx, s.exampleValidIngredient.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidIngredient(s.ctx, s.exampleValidIngredient.ID)
		assert.Error(t, err)
	})
}
