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

func checkUserIngredientPreferenceEquality(t *testing.T, expected, actual *types.UserIngredientPreference) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for user ingredient preference %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestUserIngredientPreferences_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredient := createValidIngredientForTest(t, ctx, testClients.admin)

			t.Log("creating user ingredient preference")
			exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()
			exampleUserIngredientPreferenceInput := converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceCreationRequestInput(exampleUserIngredientPreference)
			exampleUserIngredientPreferenceInput.IngredientID = createdValidIngredient.ID
			createdUserIngredientPreference, err := testClients.admin.CreateUserIngredientPreference(ctx, exampleUserIngredientPreferenceInput)
			require.NoError(t, err)
			t.Logf("user ingredient preference %q created", createdUserIngredientPreference.ID)
			checkUserIngredientPreferenceEquality(t, exampleUserIngredientPreference, createdUserIngredientPreference)

			createdUserIngredientPreferences, err := testClients.admin.GetUserIngredientPreferences(ctx, types.DefaultQueryFilter())
			requireNotNilAndNoProblems(t, createdUserIngredientPreference, err)
			checkUserIngredientPreferenceEquality(t, exampleUserIngredientPreference, createdUserIngredientPreference)
			require.NotEmpty(t, createdUserIngredientPreferences.Data, "expected to find at least one user ingredient preference")
			createdUserIngredientPreference = createdUserIngredientPreferences.Data[0]

			t.Log("changing user ingredient preference")
			newUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()
			createdValidIngredient2 := createValidIngredientForTest(t, ctx, testClients.admin)
			newUserIngredientPreference.Ingredient = *createdValidIngredient2
			createdUserIngredientPreference.Update(converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceUpdateRequestInput(newUserIngredientPreference))
			assert.NoError(t, testClients.admin.UpdateUserIngredientPreference(ctx, createdUserIngredientPreference))

			t.Log("fetching changed user ingredient preference")
			newResults, err := testClients.admin.GetUserIngredientPreferences(ctx, types.DefaultQueryFilter())
			requireNotNilAndNoProblems(t, newResults, err)
			require.NotEmpty(t, newResults.Data, "expected to find at least one user ingredient preference")
			actual := newResults.Data[0]

			// assert user ingredient preference equality
			checkUserIngredientPreferenceEquality(t, newUserIngredientPreference, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			t.Log("cleaning up user ingredient preference")
			assert.NoError(t, testClients.admin.ArchiveUserIngredientPreference(ctx, createdUserIngredientPreference.ID))
		}
	})
}

func (s *TestSuite) TestUserIngredientPreferences_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating user ingredient preferences")
			var expected []*types.UserIngredientPreference
			for i := 0; i < 5; i++ {
				createdValidIngredient := createValidIngredientForTest(t, ctx, testClients.admin)
				exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()
				exampleUserIngredientPreferenceInput := converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceCreationRequestInput(exampleUserIngredientPreference)
				exampleUserIngredientPreferenceInput.IngredientID = createdValidIngredient.ID
				createdUserIngredientPreference, createdUserIngredientPreferenceErr := testClients.admin.CreateUserIngredientPreference(ctx, exampleUserIngredientPreferenceInput)
				require.NoError(t, createdUserIngredientPreferenceErr)
				t.Logf("user ingredient preference %q created", createdUserIngredientPreference.ID)

				checkUserIngredientPreferenceEquality(t, exampleUserIngredientPreference, createdUserIngredientPreference)

				expected = append(expected, createdUserIngredientPreference)
			}

			// assert user ingredient preference list equality
			actual, err := testClients.admin.GetUserIngredientPreferences(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			t.Log("cleaning up")
			for _, createdUserIngredientPreference := range expected {
				assert.NoError(t, testClients.admin.ArchiveUserIngredientPreference(ctx, createdUserIngredientPreference.ID))
			}
		}
	})
}
