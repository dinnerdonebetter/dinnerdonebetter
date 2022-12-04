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

func TestRecipeStepCompletionConditions(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(recipeStepCompletionConditionsTestSuite))
}

type recipeStepCompletionConditionsBaseSuite struct {
	suite.Suite
	ctx                                  context.Context
	exampleRecipeStepCompletionCondition *types.RecipeStepCompletionCondition
	exampleRecipeID                      string
	exampleRecipeStepID                  string
}

var _ suite.SetupTestSuite = (*recipeStepCompletionConditionsBaseSuite)(nil)

func (s *recipeStepCompletionConditionsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleRecipeID = fakes.BuildFakeID()
	s.exampleRecipeStepID = fakes.BuildFakeID()
	s.exampleRecipeStepCompletionCondition = fakes.BuildFakeRecipeStepCompletionCondition()
	s.exampleRecipeStepCompletionCondition.BelongsToRecipeStep = s.exampleRecipeStepID
}

type recipeStepCompletionConditionsTestSuite struct {
	suite.Suite

	recipeStepCompletionConditionsBaseSuite
}

func (s *recipeStepCompletionConditionsTestSuite) TestClient_GetRecipeStepCompletionCondition() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/completion_conditions/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCompletionCondition.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepCompletionCondition)
		actual, err := c.GetRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCompletionCondition.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeStepCompletionCondition, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepCompletionCondition(s.ctx, "", s.exampleRecipeStepID, s.exampleRecipeStepCompletionCondition.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, "", s.exampleRecipeStepCompletionCondition.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCompletionCondition.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCompletionCondition.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCompletionCondition.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepCompletionConditionsTestSuite) TestClient_GetRecipeStepCompletionConditions() {
	const expectedPath = "/api/v1/recipes/%s/steps/%s/conditions"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleRecipeStepCompletionConditionList := fakes.BuildFakeRecipeStepCompletionConditionList()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleRecipeStepCompletionConditionList)
		actual, err := c.GetRecipeStepCompletionConditions(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepCompletionConditionList, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepCompletionConditions(s.ctx, "", s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepCompletionConditions(s.ctx, s.exampleRecipeID, "", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepCompletionConditions(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepCompletionConditions(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepCompletionConditionsTestSuite) TestClient_CreateRecipeStepCompletionCondition() {
	const expectedPath = "/api/v1/recipes/%s/steps/%s/conditions"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepCompletionConditionCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepCompletionCondition)

		actual, err := c.CreateRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeStepCompletionCondition, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepCompletionConditionCreationRequestInput()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepCompletionCondition(s.ctx, "", s.exampleRecipeStepID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.RecipeStepCompletionConditionCreationRequestInput{}

		actual, err := c.CreateRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionCreationRequestInput(s.exampleRecipeStepCompletionCondition)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionCreationRequestInput(s.exampleRecipeStepCompletionCondition)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepCompletionConditionsTestSuite) TestClient_UpdateRecipeStepCompletionCondition() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/completion_conditions/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCompletionCondition.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepCompletionCondition)

		err := c.UpdateRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepCompletionCondition)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeStepCompletionCondition(s.ctx, "", s.exampleRecipeStepCompletionCondition)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepCompletionCondition)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepCompletionCondition)
		assert.Error(t, err)
	})
}

func (s *recipeStepCompletionConditionsTestSuite) TestClient_ArchiveRecipeStepCompletionCondition() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/completion_conditions/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCompletionCondition.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCompletionCondition.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepCompletionCondition(s.ctx, "", s.exampleRecipeStepID, s.exampleRecipeStepCompletionCondition.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, "", s.exampleRecipeStepCompletionCondition.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCompletionCondition.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveRecipeStepCompletionCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCompletionCondition.ID)
		assert.Error(t, err)
	})
}
