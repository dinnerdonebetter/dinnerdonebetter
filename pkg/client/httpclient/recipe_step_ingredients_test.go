package httpclient

import (
	"context"
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestRecipeStepIngredients(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(recipeStepIngredientsTestSuite))
}

type recipeStepIngredientsBaseSuite struct {
	suite.Suite
	ctx                         context.Context
	exampleRecipeStepIngredient *types.RecipeStepIngredient
	exampleRecipeID             uint64
	exampleRecipeStepID         uint64
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

func (s *recipeStepIngredientsTestSuite) TestClient_RecipeStepIngredientExists() {
	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)

		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)
		actual, err := c.RecipeStepIngredientExists(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)

		assert.NoError(t, err)
		assert.True(t, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.RecipeStepIngredientExists(s.ctx, 0, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.RecipeStepIngredientExists(s.ctx, s.exampleRecipeID, 0, s.exampleRecipeStepIngredient.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with invalid recipe step ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.RecipeStepIngredientExists(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, 0)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.RecipeStepIngredientExists(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		actual, err := c.RecipeStepIngredientExists(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func (s *recipeStepIngredientsTestSuite) TestClient_GetRecipeStepIngredient() {
	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients/%d"

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
		actual, err := c.GetRecipeStepIngredient(s.ctx, 0, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepIngredient(s.ctx, s.exampleRecipeID, 0, s.exampleRecipeStepIngredient.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, 0)

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
	const expectedPath = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients"

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
		actual, err := c.GetRecipeStepIngredients(s.ctx, 0, s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepIngredients(s.ctx, s.exampleRecipeID, 0, filter)

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
	const expectedPath = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInput()
		exampleInput.BelongsToRecipeStep = s.exampleRecipeStepID

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepIngredient)

		actual, err := c.CreateRecipeStepIngredient(s.ctx, s.exampleRecipeID, exampleInput)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleRecipeStepIngredient, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInput()
		exampleInput.BelongsToRecipeStep = s.exampleRecipeStepID

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepIngredient(s.ctx, 0, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepIngredient(s.ctx, s.exampleRecipeID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.RecipeStepIngredientCreationInput{}

		actual, err := c.CreateRecipeStepIngredient(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(s.exampleRecipeStepIngredient)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateRecipeStepIngredient(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(s.exampleRecipeStepIngredient)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateRecipeStepIngredient(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepIngredientsTestSuite) TestClient_UpdateRecipeStepIngredient() {
	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients/%d"

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

		err := c.UpdateRecipeStepIngredient(s.ctx, 0, s.exampleRecipeStepIngredient)
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
	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients/%d"

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

		err := c.ArchiveRecipeStepIngredient(s.ctx, 0, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepIngredient(s.ctx, s.exampleRecipeID, 0, s.exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, 0)
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

func (s *recipeStepIngredientsTestSuite) TestClient_GetAuditLogForRecipeStepIngredient() {
	const (
		expectedPath   = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_ingredients/%d/audit"
		expectedMethod = http.MethodGet
	)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, expectedMethod, "", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)
		exampleAuditLogEntryList := fakes.BuildFakeAuditLogEntryList().Entries

		c, _ := buildTestClientWithJSONResponse(t, spec, exampleAuditLogEntryList)

		actual, err := c.GetAuditLogForRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleAuditLogEntryList, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetAuditLogForRecipeStepIngredient(s.ctx, 0, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetAuditLogForRecipeStepIngredient(s.ctx, s.exampleRecipeID, 0, s.exampleRecipeStepIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetAuditLogForRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetAuditLogForRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetAuditLogForRecipeStepIngredient(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
