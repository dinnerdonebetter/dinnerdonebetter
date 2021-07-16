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

func TestRecipeSteps(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(recipeStepsTestSuite))
}

type recipeStepsBaseSuite struct {
	suite.Suite
	ctx               context.Context
	exampleRecipeStep *types.RecipeStep
	exampleRecipeID   uint64
}

var _ suite.SetupTestSuite = (*recipeStepsBaseSuite)(nil)

func (s *recipeStepsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleRecipeID = fakes.BuildFakeID()
	s.exampleRecipeStep = fakes.BuildFakeRecipeStep()
	s.exampleRecipeStep.BelongsToRecipe = s.exampleRecipeID
}

type recipeStepsTestSuite struct {
	suite.Suite

	recipeStepsBaseSuite
}

func (s *recipeStepsTestSuite) TestClient_RecipeStepExists() {
	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStep.ID)

		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)
		actual, err := c.RecipeStepExists(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)

		assert.NoError(t, err)
		assert.True(t, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.RecipeStepExists(s.ctx, 0, s.exampleRecipeStep.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.RecipeStepExists(s.ctx, s.exampleRecipeID, 0)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.RecipeStepExists(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		actual, err := c.RecipeStepExists(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func (s *recipeStepsTestSuite) TestClient_GetRecipeStep() {
	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStep.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStep)
		actual, err := c.GetRecipeStep(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeStep, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStep(s.ctx, 0, s.exampleRecipeStep.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStep(s.ctx, s.exampleRecipeID, 0)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStep(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStep.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStep(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepsTestSuite) TestClient_GetRecipeSteps() {
	const expectedPath = "/api/v1/recipes/%d/recipe_steps"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleRecipeStepList := fakes.BuildFakeRecipeStepList()

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath, s.exampleRecipeID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleRecipeStepList)
		actual, err := c.GetRecipeSteps(s.ctx, s.exampleRecipeID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepList, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeSteps(s.ctx, 0, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeSteps(s.ctx, s.exampleRecipeID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath, s.exampleRecipeID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeSteps(s.ctx, s.exampleRecipeID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepsTestSuite) TestClient_CreateRecipeStep() {
	const expectedPath = "/api/v1/recipes/%d/recipe_steps"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepCreationInput()
		exampleInput.BelongsToRecipe = s.exampleRecipeID

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleRecipeID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStep)

		actual, err := c.CreateRecipeStep(s.ctx, exampleInput)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleRecipeStep, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStep(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.RecipeStepCreationInput{}

		actual, err := c.CreateRecipeStep(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(s.exampleRecipeStep)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateRecipeStep(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(s.exampleRecipeStep)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateRecipeStep(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepsTestSuite) TestClient_UpdateRecipeStep() {
	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStep.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStep)

		err := c.UpdateRecipeStep(s.ctx, s.exampleRecipeStep)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeStep(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateRecipeStep(s.ctx, s.exampleRecipeStep)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateRecipeStep(s.ctx, s.exampleRecipeStep)
		assert.Error(t, err)
	})
}

func (s *recipeStepsTestSuite) TestClient_ArchiveRecipeStep() {
	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStep.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveRecipeStep(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStep(s.ctx, 0, s.exampleRecipeStep.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStep(s.ctx, s.exampleRecipeID, 0)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveRecipeStep(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveRecipeStep(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)
		assert.Error(t, err)
	})
}

func (s *recipeStepsTestSuite) TestClient_GetAuditLogForRecipeStep() {
	const (
		expectedPath   = "/api/v1/recipes/%d/recipe_steps/%d/audit"
		expectedMethod = http.MethodGet
	)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, expectedMethod, "", expectedPath, s.exampleRecipeID, s.exampleRecipeStep.ID)
		exampleAuditLogEntryList := fakes.BuildFakeAuditLogEntryList().Entries

		c, _ := buildTestClientWithJSONResponse(t, spec, exampleAuditLogEntryList)

		actual, err := c.GetAuditLogForRecipeStep(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleAuditLogEntryList, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetAuditLogForRecipeStep(s.ctx, 0, s.exampleRecipeStep.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetAuditLogForRecipeStep(s.ctx, s.exampleRecipeID, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetAuditLogForRecipeStep(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetAuditLogForRecipeStep(s.ctx, s.exampleRecipeID, s.exampleRecipeStep.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
