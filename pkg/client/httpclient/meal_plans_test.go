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

func TestMealPlans(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(mealPlansTestSuite))
}

type mealPlansBaseSuite struct {
	suite.Suite

	ctx             context.Context
	exampleMealPlan *types.MealPlan
}

var _ suite.SetupTestSuite = (*mealPlansBaseSuite)(nil)

func (s *mealPlansBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleMealPlan = fakes.BuildFakeMealPlan()
}

type mealPlansTestSuite struct {
	suite.Suite

	mealPlansBaseSuite
}

func (s *mealPlansTestSuite) TestClient_GetMealPlan() {
	const expectedPathFormat = "/api/v1/meal_plans/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlan.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlan)
		actual, err := c.GetMealPlan(s.ctx, s.exampleMealPlan.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealPlan, actual)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlan(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlan(s.ctx, s.exampleMealPlan.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlan.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlan(s.ctx, s.exampleMealPlan.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlansTestSuite) TestClient_GetMealPlans() {
	const expectedPath = "/api/v1/meal_plans"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleMealPlanList := fakes.BuildFakeMealPlanList()

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleMealPlanList)
		actual, err := c.GetMealPlans(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlans(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlans(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlansTestSuite) TestClient_CreateMealPlan() {
	const expectedPath = "/api/v1/meal_plans"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealPlanCreationRequestInput()
		exampleInput.BelongsToHousehold = ""

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlan)

		actual, err := c.CreateMealPlan(s.ctx, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleMealPlan, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateMealPlan(s.ctx, nil)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.MealPlanCreationRequestInput{}

		actual, err := c.CreateMealPlan(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealPlanCreationRequestInputFromMealPlan(s.exampleMealPlan)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateMealPlan(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealPlanCreationRequestInputFromMealPlan(s.exampleMealPlan)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateMealPlan(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlansTestSuite) TestClient_UpdateMealPlan() {
	const expectedPathFormat = "/api/v1/meal_plans/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleMealPlan.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlan)

		err := c.UpdateMealPlan(s.ctx, s.exampleMealPlan)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateMealPlan(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateMealPlan(s.ctx, s.exampleMealPlan)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateMealPlan(s.ctx, s.exampleMealPlan)
		assert.Error(t, err)
	})
}

func (s *mealPlansTestSuite) TestClient_ArchiveMealPlan() {
	const expectedPathFormat = "/api/v1/meal_plans/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleMealPlan.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveMealPlan(s.ctx, s.exampleMealPlan.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveMealPlan(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveMealPlan(s.ctx, s.exampleMealPlan.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveMealPlan(s.ctx, s.exampleMealPlan.ID)
		assert.Error(t, err)
	})
}
