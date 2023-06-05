package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetMealRatingRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meals/%s/ratings/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealID := fakes.BuildFakeID()
		exampleMealRating := fakes.BuildFakeMealRating()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleMealID, exampleMealRating.ID)

		actual, err := helper.builder.BuildGetMealRatingRequest(helper.ctx, exampleMealID, exampleMealRating.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		exampleMealID := fakes.BuildFakeID()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetMealRatingRequest(helper.ctx, exampleMealID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealID := fakes.BuildFakeID()
		exampleMealRating := fakes.BuildFakeMealRating()

		actual, err := helper.builder.BuildGetMealRatingRequest(helper.ctx, exampleMealID, exampleMealRating.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetMealRatingsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meals/%s/ratings"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealID := fakes.BuildFakeID()
		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPathFormat, exampleMealID)

		actual, err := helper.builder.BuildGetMealRatingsRequest(helper.ctx, exampleMealID, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealID := fakes.BuildFakeID()
		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetMealRatingsRequest(helper.ctx, exampleMealID, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateMealRatingRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/meals/%s/ratings"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealRatingCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, exampleMealID)

		actual, err := helper.builder.BuildCreateMealRatingRequest(helper.ctx, exampleMealID, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		exampleMealID := fakes.BuildFakeID()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateMealRatingRequest(helper.ctx, exampleMealID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		exampleMealID := fakes.BuildFakeID()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateMealRatingRequest(helper.ctx, exampleMealID, &types.MealRatingCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeMealRatingCreationRequestInput()

		actual, err := helper.builder.BuildCreateMealRatingRequest(helper.ctx, exampleMealID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateMealRatingRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meals/%s/ratings/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealRating := fakes.BuildFakeMealRating()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleMealRating.MealID, exampleMealRating.ID)

		actual, err := helper.builder.BuildUpdateMealRatingRequest(helper.ctx, exampleMealRating)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateMealRatingRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealRating := fakes.BuildFakeMealRating()

		actual, err := helper.builder.BuildUpdateMealRatingRequest(helper.ctx, exampleMealRating)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveMealRatingRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/meals/%s/ratings/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleMealID := fakes.BuildFakeID()
		exampleMealRating := fakes.BuildFakeMealRating()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleMealID, exampleMealRating.ID)

		actual, err := helper.builder.BuildArchiveMealRatingRequest(helper.ctx, exampleMealID, exampleMealRating.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		exampleMealID := fakes.BuildFakeID()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveMealRatingRequest(helper.ctx, exampleMealID, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleMealID := fakes.BuildFakeID()
		exampleMealRating := fakes.BuildFakeMealRating()

		actual, err := helper.builder.BuildArchiveMealRatingRequest(helper.ctx, exampleMealID, exampleMealRating.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
