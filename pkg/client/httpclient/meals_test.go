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

func TestMeals(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(mealsTestSuite))
}

type mealsBaseSuite struct {
	suite.Suite

	ctx         context.Context
	exampleMeal *types.Meal
}

var _ suite.SetupTestSuite = (*mealsBaseSuite)(nil)

func (s *mealsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleMeal = fakes.BuildFakeMeal()
}

type mealsTestSuite struct {
	suite.Suite

	mealsBaseSuite
}

func (s *mealsTestSuite) TestClient_GetMeal() {
	const expectedPathFormat = "/api/v1/meals/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMeal.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMeal)
		actual, err := c.GetMeal(s.ctx, s.exampleMeal.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMeal, actual)
	})

	s.Run("with invalid meal ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMeal(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMeal(s.ctx, s.exampleMeal.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMeal.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMeal(s.ctx, s.exampleMeal.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealsTestSuite) TestClient_GetMeals() {
	const expectedPath = "/api/v1/meals"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleMealList := fakes.BuildFakeMealList()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleMealList)
		actual, err := c.GetMeals(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMeals(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMeals(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipesTestSuite) TestClient_SearchForMeals() {
	const expectedPath = "/api/v1/meals/search"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleMealList := fakes.BuildFakeMealList()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&q=example&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleMealList)
		actual, err := c.SearchForMeals(s.ctx, "example", filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.SearchForMeals(s.ctx, "example", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&q=example&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchForMeals(s.ctx, "example", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealsTestSuite) TestClient_CreateMeal() {
	const expectedPath = "/api/v1/meals"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealCreationRequestInput()
		exampleInput.CreatedByUser = ""

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMeal)

		actual, err := c.CreateMeal(s.ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMeal, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateMeal(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.MealCreationRequestInput{}

		actual, err := c.CreateMeal(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealCreationRequestInputFromMeal(s.exampleMeal)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateMeal(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealCreationRequestInputFromMeal(s.exampleMeal)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateMeal(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealsTestSuite) TestClient_ArchiveMeal() {
	const expectedPathFormat = "/api/v1/meals/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleMeal.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveMeal(s.ctx, s.exampleMeal.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid meal ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveMeal(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveMeal(s.ctx, s.exampleMeal.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveMeal(s.ctx, s.exampleMeal.ID)
		assert.Error(t, err)
	})
}
