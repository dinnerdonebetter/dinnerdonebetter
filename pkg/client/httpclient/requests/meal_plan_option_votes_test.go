package requests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func TestBuilder_BuildGetMealPlanOptionVoteRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/meal_plan_options/%s/meal_plan_option_votes/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleMealPlanID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)

		actual, err := helper.builder.BuildGetMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildGetMealPlanOptionVoteRequest(helper.ctx, "", exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildGetMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, "", exampleMealPlanOptionVote.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid meal plan option vote ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanOptionID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildGetMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetMealPlanOptionVotesRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/meal_plan_options/%s/meal_plan_option_votes"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPathFormat, exampleMealPlanID, exampleMealPlanOptionID)

		actual, err := helper.builder.BuildGetMealPlanOptionVotesRequest(helper.ctx, exampleMealPlanID, exampleMealPlanOptionID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanOptionID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetMealPlanOptionVotesRequest(helper.ctx, "", exampleMealPlanOptionID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetMealPlanOptionVotesRequest(helper.ctx, exampleMealPlanID, "", filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetMealPlanOptionVotesRequest(helper.ctx, exampleMealPlanID, exampleMealPlanOptionID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateMealPlanOptionVoteRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/meal_plans/%s/vote"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, exampleMealPlanID)

		actual, err := helper.builder.BuildCreateMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()

		actual, err := helper.builder.BuildCreateMealPlanOptionVoteRequest(helper.ctx, "", exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildCreateMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildCreateMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, &types.MealPlanOptionVoteCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()

		actual, err := helper.builder.BuildCreateMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateMealPlanOptionVoteRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/meal_plan_options/%s/meal_plan_option_votes/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleMealPlanID, exampleMealPlanOptionVote.BelongsToMealPlanOption, exampleMealPlanOptionVote.ID)

		actual, err := helper.builder.BuildUpdateMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanOptionVote)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildUpdateMealPlanOptionVoteRequest(helper.ctx, "", exampleMealPlanOptionVote)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildUpdateMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildUpdateMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanOptionVote)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveMealPlanOptionVoteRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plans/%s/meal_plan_options/%s/meal_plan_option_votes/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleMealPlanID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)

		actual, err := helper.builder.BuildArchiveMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildArchiveMealPlanOptionVoteRequest(helper.ctx, "", exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildArchiveMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, "", exampleMealPlanOptionVote.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid meal plan option vote ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildArchiveMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanOptionID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildArchiveMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
