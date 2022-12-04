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

func TestRecipeStepConditions(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(recipeStepConditionsTestSuite))
}

type recipeStepConditionsBaseSuite struct {
	suite.Suite
	ctx                        context.Context
	exampleRecipeStepCondition *types.RecipeStepCondition
	exampleRecipeID            string
	exampleRecipeStepID        string
}

var _ suite.SetupTestSuite = (*recipeStepConditionsBaseSuite)(nil)

func (s *recipeStepConditionsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleRecipeID = fakes.BuildFakeID()
	s.exampleRecipeStepID = fakes.BuildFakeID()
	s.exampleRecipeStepCondition = fakes.BuildFakeRecipeStepCondition()
	s.exampleRecipeStepCondition.BelongsToRecipeStep = s.exampleRecipeStepID
}

type recipeStepConditionsTestSuite struct {
	suite.Suite

	recipeStepConditionsBaseSuite
}

func (s *recipeStepConditionsTestSuite) TestClient_GetRecipeStepCondition() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/conditions/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCondition.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepCondition)
		actual, err := c.GetRecipeStepCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCondition.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeStepCondition, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepCondition(s.ctx, "", s.exampleRecipeStepID, s.exampleRecipeStepCondition.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepCondition(s.ctx, s.exampleRecipeID, "", s.exampleRecipeStepCondition.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCondition.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCondition.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCondition.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepConditionsTestSuite) TestClient_GetRecipeStepConditions() {
	const expectedPath = "/api/v1/recipes/%s/steps/%s/conditions"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleRecipeStepConditionList := fakes.BuildFakeRecipeStepConditionList()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleRecipeStepConditionList)
		actual, err := c.GetRecipeStepConditions(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepConditionList, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepConditions(s.ctx, "", s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepConditions(s.ctx, s.exampleRecipeID, "", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepConditions(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepConditions(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepConditionsTestSuite) TestClient_CreateRecipeStepCondition() {
	const expectedPath = "/api/v1/recipes/%s/steps/%s/conditions"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepConditionCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepCondition)

		actual, err := c.CreateRecipeStepCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeStepCondition, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepConditionCreationRequestInput()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepCondition(s.ctx, "", s.exampleRecipeStepID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.RecipeStepConditionCreationRequestInput{}

		actual, err := c.CreateRecipeStepCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertRecipeStepConditionToRecipeStepConditionCreationRequestInput(s.exampleRecipeStepCondition)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateRecipeStepCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertRecipeStepConditionToRecipeStepConditionCreationRequestInput(s.exampleRecipeStepCondition)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateRecipeStepCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepConditionsTestSuite) TestClient_UpdateRecipeStepCondition() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/conditions/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCondition.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepCondition)

		err := c.UpdateRecipeStepCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepCondition)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeStepCondition(s.ctx, "", s.exampleRecipeStepCondition)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeStepCondition(s.ctx, s.exampleRecipeID, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateRecipeStepCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepCondition)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateRecipeStepCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepCondition)
		assert.Error(t, err)
	})
}

func (s *recipeStepConditionsTestSuite) TestClient_ArchiveRecipeStepCondition() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/conditions/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCondition.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveRecipeStepCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCondition.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepCondition(s.ctx, "", s.exampleRecipeStepID, s.exampleRecipeStepCondition.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepCondition(s.ctx, s.exampleRecipeID, "", s.exampleRecipeStepCondition.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveRecipeStepCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCondition.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveRecipeStepCondition(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepCondition.ID)
		assert.Error(t, err)
	})
}
