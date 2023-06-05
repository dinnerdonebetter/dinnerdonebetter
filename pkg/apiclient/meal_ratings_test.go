package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestMealRatings(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(mealRatingsTestSuite))
}

type mealRatingsBaseSuite struct {
	suite.Suite

	ctx               context.Context
	exampleMeal       *types.Meal
	exampleMealRating *types.MealRating
}

var _ suite.SetupTestSuite = (*mealRatingsBaseSuite)(nil)

func (s *mealRatingsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleMeal = fakes.BuildFakeMeal()
	s.exampleMealRating = fakes.BuildFakeMealRating()
}

type mealRatingsTestSuite struct {
	suite.Suite

	mealRatingsBaseSuite
}

func (s *mealRatingsTestSuite) TestClient_GetMealRating() {
	const expectedPathFormat = "/api/v1/meals/%s/ratings/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMeal.ID, s.exampleMealRating.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealRating)
		actual, err := c.GetMealRating(s.ctx, s.exampleMeal.ID, s.exampleMealRating.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealRating, actual)
	})

	s.Run("with invalid valid instrument ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealRating(s.ctx, s.exampleMeal.ID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealRating(s.ctx, s.exampleMeal.ID, s.exampleMealRating.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMeal.ID, s.exampleMealRating.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealRating(s.ctx, s.exampleMeal.ID, s.exampleMealRating.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealRatingsTestSuite) TestClient_GetMealRatings() {
	const expectedPath = "/api/v1/meals/%s/ratings"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleMealRatingList := fakes.BuildFakeMealRatingList()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath, s.exampleMeal.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleMealRatingList)
		actual, err := c.GetMealRatings(s.ctx, s.exampleMeal.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealRatingList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealRatings(s.ctx, s.exampleMeal.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath, s.exampleMeal.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealRatings(s.ctx, s.exampleMeal.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealRatingsTestSuite) TestClient_CreateMealRating() {
	const expectedPath = "/api/v1/meals/%s/ratings"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealRatingCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleMeal.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealRating)

		actual, err := c.CreateMealRating(s.ctx, s.exampleMeal.ID, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealRating, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateMealRating(s.ctx, s.exampleMeal.ID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.MealRatingCreationRequestInput{}

		actual, err := c.CreateMealRating(s.ctx, s.exampleMeal.ID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertMealRatingToMealRatingCreationRequestInput(s.exampleMealRating)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateMealRating(s.ctx, s.exampleMeal.ID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertMealRatingToMealRatingCreationRequestInput(s.exampleMealRating)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateMealRating(s.ctx, s.exampleMeal.ID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealRatingsTestSuite) TestClient_UpdateMealRating() {
	const expectedPathFormat = "/api/v1/meals/%s/ratings/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleMealRating.MealID, s.exampleMealRating.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealRating)

		err := c.UpdateMealRating(s.ctx, s.exampleMealRating)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateMealRating(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateMealRating(s.ctx, s.exampleMealRating)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateMealRating(s.ctx, s.exampleMealRating)
		assert.Error(t, err)
	})
}

func (s *mealRatingsTestSuite) TestClient_ArchiveMealRating() {
	const expectedPathFormat = "/api/v1/meals/%s/ratings/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleMeal.ID, s.exampleMealRating.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveMealRating(s.ctx, s.exampleMeal.ID, s.exampleMealRating.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid instrument ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveMealRating(s.ctx, s.exampleMeal.ID, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveMealRating(s.ctx, s.exampleMeal.ID, s.exampleMealRating.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveMealRating(s.ctx, s.exampleMeal.ID, s.exampleMealRating.ID)
		assert.Error(t, err)
	})
}
