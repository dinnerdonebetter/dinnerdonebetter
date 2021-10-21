package httpclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
)

func TestMealPlanOptionVotes(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(mealPlanOptionVotesTestSuite))
}

type mealPlanOptionVotesBaseSuite struct {
	suite.Suite

	ctx                       context.Context
	exampleMealPlanOptionVote *types.MealPlanOptionVote
}

var _ suite.SetupTestSuite = (*mealPlanOptionVotesBaseSuite)(nil)

func (s *mealPlanOptionVotesBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleMealPlanOptionVote = fakes.BuildFakeMealPlanOptionVote()
}

type mealPlanOptionVotesTestSuite struct {
	suite.Suite

	mealPlanOptionVotesBaseSuite
}

func (s *mealPlanOptionVotesTestSuite) TestClient_GetMealPlanOptionVote() {
	const expectedPathFormat = "/api/v1/meal_plan_option_votes/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlanOptionVote.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanOptionVote)
		actual, err := c.GetMealPlanOptionVote(s.ctx, s.exampleMealPlanOptionVote.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealPlanOptionVote, actual)
	})

	s.Run("with invalid meal plan option vote ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanOptionVote(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanOptionVote(s.ctx, s.exampleMealPlanOptionVote.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlanOptionVote.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanOptionVote(s.ctx, s.exampleMealPlanOptionVote.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanOptionVotesTestSuite) TestClient_GetMealPlanOptionVotes() {
	const expectedPath = "/api/v1/meal_plan_option_votes"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleMealPlanOptionVoteList := fakes.BuildFakeMealPlanOptionVoteList()

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleMealPlanOptionVoteList)
		actual, err := c.GetMealPlanOptionVotes(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanOptionVoteList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanOptionVotes(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanOptionVotes(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanOptionVotesTestSuite) TestClient_CreateMealPlanOptionVote() {
	const expectedPath = "/api/v1/meal_plan_option_votes"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()
		exampleInput.BelongsToAccount = ""

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, &types.PreWriteResponse{ID: s.exampleMealPlanOptionVote.ID})

		actual, err := c.CreateMealPlanOptionVote(s.ctx, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleMealPlanOptionVote.ID, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateMealPlanOptionVote(s.ctx, nil)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.MealPlanOptionVoteCreationRequestInput{}

		actual, err := c.CreateMealPlanOptionVote(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(s.exampleMealPlanOptionVote)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateMealPlanOptionVote(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(s.exampleMealPlanOptionVote)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateMealPlanOptionVote(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanOptionVotesTestSuite) TestClient_UpdateMealPlanOptionVote() {
	const expectedPathFormat = "/api/v1/meal_plan_option_votes/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleMealPlanOptionVote.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanOptionVote)

		err := c.UpdateMealPlanOptionVote(s.ctx, s.exampleMealPlanOptionVote)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateMealPlanOptionVote(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateMealPlanOptionVote(s.ctx, s.exampleMealPlanOptionVote)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateMealPlanOptionVote(s.ctx, s.exampleMealPlanOptionVote)
		assert.Error(t, err)
	})
}

func (s *mealPlanOptionVotesTestSuite) TestClient_ArchiveMealPlanOptionVote() {
	const expectedPathFormat = "/api/v1/meal_plan_option_votes/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleMealPlanOptionVote.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveMealPlanOptionVote(s.ctx, s.exampleMealPlanOptionVote.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid meal plan option vote ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveMealPlanOptionVote(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveMealPlanOptionVote(s.ctx, s.exampleMealPlanOptionVote.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveMealPlanOptionVote(s.ctx, s.exampleMealPlanOptionVote.ID)
		assert.Error(t, err)
	})
}
