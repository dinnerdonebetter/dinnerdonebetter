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

func TestMealPlanOptions(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(mealPlanOptionsTestSuite))
}

type mealPlanOptionsBaseSuite struct {
	suite.Suite
	ctx                   context.Context
	exampleMealPlanOption *types.MealPlanOption
	exampleMealPlanID     string
}

var _ suite.SetupTestSuite = (*mealPlanOptionsBaseSuite)(nil)

func (s *mealPlanOptionsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleMealPlanID = fakes.BuildFakeID()
	s.exampleMealPlanOption = fakes.BuildFakeMealPlanOption()
	s.exampleMealPlanOption.BelongsToMealPlan = s.exampleMealPlanID
}

type mealPlanOptionsTestSuite struct {
	suite.Suite

	mealPlanOptionsBaseSuite
}

func (s *mealPlanOptionsTestSuite) TestClient_GetMealPlanOption() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/meal_plan_options/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanOption.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanOption)
		actual, err := c.GetMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanOption.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealPlanOption, actual)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOption(s.ctx, "", s.exampleMealPlanOption.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid meal plan option ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOption(s.ctx, s.exampleMealPlanID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanOption.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanOption.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanOption.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanOptionsTestSuite) TestClient_GetMealPlanOptions() {
	const expectedPath = "/api/v1/meal_plans/%s/meal_plan_options"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleMealPlanOptionList := fakes.BuildFakeMealPlanOptionList()

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath, s.exampleMealPlanID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleMealPlanOptionList)
		actual, err := c.GetMealPlanOptions(s.ctx, s.exampleMealPlanID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanOptionList, actual)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOptions(s.ctx, "", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanOptions(s.ctx, s.exampleMealPlanID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath, s.exampleMealPlanID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanOptions(s.ctx, s.exampleMealPlanID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanOptionsTestSuite) TestClient_CreateMealPlanOption() {
	const expectedPath = "/api/v1/meal_plans/%s/meal_plan_options"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealPlanOptionCreationRequestInput()
		exampleInput.BelongsToMealPlan = s.exampleMealPlanID

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleMealPlanID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanOption)

		actual, err := c.CreateMealPlanOption(s.ctx, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleMealPlanOption, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateMealPlanOption(s.ctx, nil)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.MealPlanOptionCreationRequestInput{}

		actual, err := c.CreateMealPlanOption(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealPlanOptionCreationRequestInputFromMealPlanOption(s.exampleMealPlanOption)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateMealPlanOption(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealPlanOptionCreationRequestInputFromMealPlanOption(s.exampleMealPlanOption)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateMealPlanOption(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanOptionsTestSuite) TestClient_UpdateMealPlanOption() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/meal_plan_options/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanOption.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanOption)

		err := c.UpdateMealPlanOption(s.ctx, s.exampleMealPlanOption)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateMealPlanOption(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateMealPlanOption(s.ctx, s.exampleMealPlanOption)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateMealPlanOption(s.ctx, s.exampleMealPlanOption)
		assert.Error(t, err)
	})
}

func (s *mealPlanOptionsTestSuite) TestClient_ArchiveMealPlanOption() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/meal_plan_options/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanOption.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanOption.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveMealPlanOption(s.ctx, "", s.exampleMealPlanOption.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid meal plan option ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveMealPlanOption(s.ctx, s.exampleMealPlanID, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanOption.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveMealPlanOption(s.ctx, s.exampleMealPlanID, s.exampleMealPlanOption.ID)
		assert.Error(t, err)
	})
}
