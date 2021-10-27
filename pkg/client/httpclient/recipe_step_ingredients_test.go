package httpclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func TestRecipeStepIngredients(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(recipeStepIngredientsTestSuite))
}

type recipeStepIngredientsBaseSuite struct {
	suite.Suite
	ctx                         context.Context
	exampleRecipeStepIngredient *types.RecipeStepIngredient
	exampleRecipeID             string
	exampleRecipeStepID         string
}

var _ suite.SetupTestSuite = (*recipeStepIngredientsBaseSuite)(nil)

func (s *recipeStepIngredientsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleRecipeID = fakes.BuildFakeID()
	s.exampleRecipeStepID = fakes.BuildFakeID()
	s.exampleRecipeStepIngredient = fakes.BuildFakeRecipeStepIngredient()
	s.exampleRecipeStepIngredient.BelongsToRecipeStep = s.exampleRecipeStepID
}

type recipeStepIngredientsTestSuite struct {
	suite.Suite

	recipeStepIngredientsBaseSuite
}

func (s *recipeStepIngredientsTestSuite) TestClient_GetRecipeStepIngredient() {
	const expectedPathFormat = "/api/v1/recipes/%s/recipe_steps/%s/recipe_step_ingredients/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepIngredient)
		actual, err := c.GetRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeStepIngredient, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepIngredient(s.ctx, "", s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepIngredient(s.ctx, s.exampleRecipeID, "", s.exampleRecipeStepIngredient.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepIngredientsTestSuite) TestClient_GetRecipeStepIngredients() {
	const expectedPath = "/api/v1/recipes/%s/recipe_steps/%s/recipe_step_ingredients"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleRecipeStepIngredientList)
		actual, err := c.GetRecipeStepIngredients(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredientList, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepIngredients(s.ctx, "", s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepIngredients(s.ctx, s.exampleRecipeID, "", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepIngredients(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepIngredients(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepIngredientsTestSuite) TestClient_CreateRecipeStepIngredient() {
	const expectedPath = "/api/v1/recipes/%s/recipe_steps/%s/recipe_step_ingredients"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationRequestInput()
		exampleInput.BelongsToRecipeStep = s.exampleRecipeStepID

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, &types.PreWriteResponse{ID: s.exampleRecipeStepIngredient.ID})

		actual, err := c.CreateRecipeStepIngredient(s.ctx, s.exampleRecipeID, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleRecipeStepIngredient.ID, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationRequestInput()
		exampleInput.BelongsToRecipeStep = s.exampleRecipeStepID

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepIngredient(s.ctx, "", exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepIngredient(s.ctx, s.exampleRecipeID, nil)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.RecipeStepIngredientCreationRequestInput{}

		actual, err := c.CreateRecipeStepIngredient(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationRequestInputFromRecipeStepIngredient(s.exampleRecipeStepIngredient)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateRecipeStepIngredient(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationRequestInputFromRecipeStepIngredient(s.exampleRecipeStepIngredient)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateRecipeStepIngredient(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepIngredientsTestSuite) TestClient_UpdateRecipeStepIngredient() {
	const expectedPathFormat = "/api/v1/recipes/%s/recipe_steps/%s/recipe_step_ingredients/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepIngredient)

		err := c.UpdateRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepIngredient)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeStepIngredient(s.ctx, "", s.exampleRecipeStepIngredient)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeStepIngredient(s.ctx, s.exampleRecipeID, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepIngredient)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepIngredient)
		assert.Error(t, err)
	})
}

func (s *recipeStepIngredientsTestSuite) TestClient_ArchiveRecipeStepIngredient() {
	const expectedPathFormat = "/api/v1/recipes/%s/recipe_steps/%s/recipe_step_ingredients/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepIngredient(s.ctx, "", s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepIngredient(s.ctx, s.exampleRecipeID, "", s.exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
	})
}
