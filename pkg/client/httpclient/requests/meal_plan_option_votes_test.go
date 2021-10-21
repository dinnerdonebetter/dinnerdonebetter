package requests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
)

func TestBuilder_BuildGetMealPlanOptionVoteRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plan_option_votes/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleMealPlanOptionVote.ID)

		actual, err := helper.builder.BuildGetMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanOptionVote.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan option vote ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetMealPlanOptionVoteRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildGetMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanOptionVote.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetMealPlanOptionVotesRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plan_option_votes"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetMealPlanOptionVotesRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetMealPlanOptionVotesRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateMealPlanOptionVoteRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/meal_plan_option_votes"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildCreateMealPlanOptionVoteRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateMealPlanOptionVoteRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateMealPlanOptionVoteRequest(helper.ctx, &types.MealPlanOptionVoteCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInput()

		actual, err := helper.builder.BuildCreateMealPlanOptionVoteRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateMealPlanOptionVoteRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plan_option_votes/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleMealPlanOptionVote.ID)

		actual, err := helper.builder.BuildUpdateMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanOptionVote)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateMealPlanOptionVoteRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildUpdateMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanOptionVote)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveMealPlanOptionVoteRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meal_plan_option_votes/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleMealPlanOptionVote.ID)

		actual, err := helper.builder.BuildArchiveMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanOptionVote.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid meal plan option vote ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveMealPlanOptionVoteRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		actual, err := helper.builder.BuildArchiveMealPlanOptionVoteRequest(helper.ctx, exampleMealPlanOptionVote.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
