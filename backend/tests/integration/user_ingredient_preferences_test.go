package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkUserIngredientPreferenceEquality(t *testing.T, expected, actual *types.UserIngredientPreference) {
	t.Helper()

	assert.NotZero(t, actual.ID)

	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for userClient ingredient preference %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.Ingredient.ID, actual.Ingredient.ID, "expected IngredientID for userClient ingredient preference %s to be %v, but it was %v", expected.ID, expected.Ingredient.ID, actual.Ingredient.ID)
	assert.Equal(t, expected.Rating, actual.Rating, "expected Rating for userClient ingredient preference %s to be %v, but it was %v", expected.ID, expected.Rating, actual.Rating)
	assert.Equal(t, expected.Allergy, actual.Allergy, "expected Allergy for userClient ingredient preference %s to be %v, but it was %v", expected.ID, expected.Allergy, actual.Allergy)

	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestUserIngredientPreferences_CompleteLifecycle() {
	s.runTest("should CRUD", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredient := createValidIngredientForTest(t, ctx, testClients.adminClient)

			exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()
			exampleUserIngredientPreference.Ingredient = *createdValidIngredient
			exampleUserIngredientPreferenceInput := converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceCreationRequestInput(exampleUserIngredientPreference)
			created, err := testClients.adminClient.CreateUserIngredientPreference(ctx, exampleUserIngredientPreferenceInput)
			require.NoError(t, err)

			require.Len(t, created, 1)
			createdUserIngredientPreference := created[0]

			checkUserIngredientPreferenceEquality(t, exampleUserIngredientPreference, createdUserIngredientPreference)

			createdUserIngredientPreferences, err := testClients.adminClient.GetUserIngredientPreferences(ctx, filtering.DefaultQueryFilter())
			requireNotNilAndNoProblems(t, createdUserIngredientPreference, err)
			checkUserIngredientPreferenceEquality(t, exampleUserIngredientPreference, createdUserIngredientPreference)
			require.NotEmpty(t, createdUserIngredientPreferences.Data, "expected to find at least one userClient ingredient preference")
			createdUserIngredientPreference = createdUserIngredientPreferences.Data[0]

			createdValidIngredient2 := createValidIngredientForTest(t, ctx, testClients.adminClient)
			newUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()
			newUserIngredientPreference.Ingredient = *createdValidIngredient2
			updateInput := converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceUpdateRequestInput(newUserIngredientPreference)
			createdUserIngredientPreference.Update(updateInput)
			assert.NoError(t, testClients.adminClient.UpdateUserIngredientPreference(ctx, createdUserIngredientPreference.ID, updateInput))

			newResults, err := testClients.adminClient.GetUserIngredientPreferences(ctx, filtering.DefaultQueryFilter())
			requireNotNilAndNoProblems(t, newResults, err)
			require.NotEmpty(t, newResults.Data, "expected to find at least one userClient ingredient preference")
			actual := newResults.Data[0]

			// assert userClient ingredient preference equality
			checkUserIngredientPreferenceEquality(t, newUserIngredientPreference, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.adminClient.ArchiveUserIngredientPreference(ctx, createdUserIngredientPreference.ID))
		}
	})
}

func (s *TestSuite) TestUserIngredientPreferences_CreatedFromIngredientGroup() {
	s.runTest("should CRUD", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredientGroup := createValidIngredientGroupForTest(t, ctx, nil, testClients.adminClient)
			logJSON(t, createdValidIngredientGroup)

			exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()
			exampleUserIngredientPreferenceInput := converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceCreationRequestInput(exampleUserIngredientPreference)
			exampleUserIngredientPreferenceInput.ValidIngredientGroupID = createdValidIngredientGroup.ID
			exampleUserIngredientPreferenceInput.ValidIngredientID = ""

			logJSON(t, exampleUserIngredientPreferenceInput)
			logJSON(t, createdValidIngredientGroup)

			created, err := testClients.adminClient.CreateUserIngredientPreference(ctx, exampleUserIngredientPreferenceInput)
			require.NoError(t, err)

			require.Equal(t, len(created), len(createdValidIngredientGroup.Members))

			for _, createdUserIngredientPreference := range created {
				assert.NoError(t, testClients.adminClient.ArchiveUserIngredientPreference(ctx, createdUserIngredientPreference.ID))
			}
		}
	})
}

func (s *TestSuite) TestUserIngredientPreferences_Listing() {
	s.runTest("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.UserIngredientPreference
			for i := 0; i < 5; i++ {
				createdValidIngredient := createValidIngredientForTest(t, ctx, testClients.adminClient)
				exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()
				exampleUserIngredientPreference.Ingredient = *createdValidIngredient
				exampleUserIngredientPreferenceInput := converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceCreationRequestInput(exampleUserIngredientPreference)
				created, createdUserIngredientPreferenceErr := testClients.adminClient.CreateUserIngredientPreference(ctx, exampleUserIngredientPreferenceInput)
				require.NoError(t, createdUserIngredientPreferenceErr)

				require.Len(t, created, 1)
				createdUserIngredientPreference := created[0]

				checkUserIngredientPreferenceEquality(t, exampleUserIngredientPreference, createdUserIngredientPreference)

				expected = append(expected, createdUserIngredientPreference)
			}

			// assert userClient ingredient preference list equality
			actual, err := testClients.adminClient.GetUserIngredientPreferences(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdUserIngredientPreference := range expected {
				assert.NoError(t, testClients.adminClient.ArchiveUserIngredientPreference(ctx, createdUserIngredientPreference.ID))
			}
		}
	})
}
