package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkMealRatingEquality(t *testing.T, expected, actual *types.MealRating) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for meal rating %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.MealID, actual.MealID, "expected MealID for meal rating %s to be %v, but it was %v", expected.ID, expected.MealID, actual.MealID)
	assert.Equal(t, expected.Taste, actual.Taste, "expected Taste for meal rating %s to be %v, but it was %v", expected.ID, expected.Taste, actual.Taste)
	assert.Equal(t, expected.Instructions, actual.Instructions, "expected Instructions for meal rating %s to be %v, but it was %v", expected.ID, expected.Instructions, actual.Instructions)
	assert.Equal(t, expected.Overall, actual.Overall, "expected Overall for meal rating %s to be %v, but it was %v", expected.ID, expected.Overall, actual.Overall)
	assert.Equal(t, expected.Cleanup, actual.Cleanup, "expected Cleanup for meal rating %s to be %v, but it was %v", expected.ID, expected.Cleanup, actual.Cleanup)
	assert.Equal(t, expected.Difficulty, actual.Difficulty, "expected Difficulty for meal rating %s to be %v, but it was %v", expected.ID, expected.Difficulty, actual.Difficulty)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestMealRatings_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdMeal := createMealForTest(ctx, t, testClients.admin, testClients.user, nil)

			t.Log("creating meal rating")
			exampleMealRating := fakes.BuildFakeMealRating()
			exampleMealRating.MealID = createdMeal.ID
			exampleMealRatingInput := converters.ConvertMealRatingToMealRatingCreationRequestInput(exampleMealRating)
			createdMealRating, err := testClients.user.CreateMealRating(ctx, createdMeal.ID, exampleMealRatingInput)
			require.NoError(t, err)
			t.Logf("meal rating %q created", createdMealRating.ID)
			checkMealRatingEquality(t, exampleMealRating, createdMealRating)

			createdMealRating, err = testClients.user.GetMealRating(ctx, createdMeal.ID, createdMealRating.ID)
			requireNotNilAndNoProblems(t, createdMealRating, err)
			checkMealRatingEquality(t, exampleMealRating, createdMealRating)

			t.Log("changing meal rating")
			newMealRating := fakes.BuildFakeMealRating()
			newMealRating.MealID = createdMeal.ID
			createdMealRating.Update(converters.ConvertMealRatingToMealRatingUpdateRequestInput(newMealRating))
			assert.NoError(t, testClients.admin.UpdateMealRating(ctx, createdMealRating))

			t.Log("fetching changed meal rating")
			actual, err := testClients.admin.GetMealRating(ctx, createdMeal.ID, createdMealRating.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert meal rating equality
			checkMealRatingEquality(t, newMealRating, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			t.Log("cleaning up meal rating")
			assert.NoError(t, testClients.admin.ArchiveMealRating(ctx, createdMeal.ID, createdMealRating.ID))
		}
	})
}

func (s *TestSuite) TestMealRatings_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdMeal := createMealForTest(ctx, t, testClients.admin, testClients.user, nil)

			t.Log("creating meal rating")
			exampleMealRating := fakes.BuildFakeMealRating()
			exampleMealRating.MealID = createdMeal.ID
			exampleMealRatingInput := converters.ConvertMealRatingToMealRatingCreationRequestInput(exampleMealRating)
			createdMealRating, err := testClients.user.CreateMealRating(ctx, createdMeal.ID, exampleMealRatingInput)
			require.NoError(t, err)
			t.Logf("meal rating %q created", createdMealRating.ID)
			checkMealRatingEquality(t, exampleMealRating, createdMealRating)

			// assert meal rating list equality
			actual, err := testClients.admin.GetMealRatings(ctx, createdMeal.ID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.Equal(t, len(actual.Data), 1, "expected %d to be <= %d", len(actual.Data), 1)

			assert.NoError(t, testClients.admin.ArchiveMealRating(ctx, createdMeal.ID, createdMealRating.ID))
		}
	})
}
