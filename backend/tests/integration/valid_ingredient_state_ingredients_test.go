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

func checkValidIngredientStateIngredientEquality(t *testing.T, expected, actual *types.ValidIngredientStateIngredient) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected StatusExplanation for valid ingredient state ingredient %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.IngredientState.ID, actual.IngredientState.ID, "expected IngredientState for valid ingredient state ingredient %s to be %v, but it was %v", expected.ID, expected.IngredientState.ID, actual.IngredientState.ID)
	assert.Equal(t, expected.Ingredient.ID, actual.Ingredient.ID, "expected Ingredient for valid ingredient state ingredient %s to be %v, but it was %v", expected.ID, expected.Ingredient.ID, actual.Ingredient.ID)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestValidIngredientStateIngredients_CompleteLifecycle() {
	s.runTest("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredientState := createValidIngredientStateForTest(t, ctx, testClients.adminClient)

			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, err := testClients.adminClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)

			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.userClient.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()
			exampleValidIngredientStateIngredient.Ingredient = *createdValidIngredient
			exampleValidIngredientStateIngredient.IngredientState = *createdValidIngredientState
			exampleValidIngredientStateIngredientInput := converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientCreationRequestInput(exampleValidIngredientStateIngredient)
			createdValidIngredientStateIngredient, err := testClients.adminClient.CreateValidIngredientStateIngredient(ctx, exampleValidIngredientStateIngredientInput)
			require.NoError(t, err)

			checkValidIngredientStateIngredientEquality(t, exampleValidIngredientStateIngredient, createdValidIngredientStateIngredient)

			createdValidIngredientStateIngredient, err = testClients.userClient.GetValidIngredientStateIngredient(ctx, createdValidIngredientStateIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredientStateIngredient, err)

			checkValidIngredientStateIngredientEquality(t, exampleValidIngredientStateIngredient, createdValidIngredientStateIngredient)

			newValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()
			newValidIngredientStateIngredient.Ingredient = *createdValidIngredient
			newValidIngredientStateIngredient.IngredientState = *createdValidIngredientState
			createdValidIngredientStateIngredient.Update(converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientUpdateRequestInput(newValidIngredientStateIngredient))
			require.NoError(t, testClients.adminClient.UpdateValidIngredientStateIngredient(ctx, createdValidIngredientStateIngredient))

			actual, err := testClients.userClient.GetValidIngredientStateIngredient(ctx, createdValidIngredientStateIngredient.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient state ingredient equality
			checkValidIngredientStateIngredientEquality(t, newValidIngredientStateIngredient, actual)
			require.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.adminClient.ArchiveValidIngredientStateIngredient(ctx, createdValidIngredientStateIngredient.ID))

			assert.NoError(t, testClients.adminClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))

			assert.NoError(t, testClients.adminClient.ArchiveValidIngredientState(ctx, createdValidIngredientState.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredientStateIngredients_Listing() {
	s.runTest("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidIngredientStateIngredient
			for i := 0; i < 5; i++ {

				createdValidIngredientState := createValidIngredientStateForTest(t, ctx, testClients.adminClient)

				exampleValidIngredient := fakes.BuildFakeValidIngredient()
				exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
				createdValidIngredient, err := testClients.adminClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
				require.NoError(t, err)

				checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

				createdValidIngredient, err = testClients.userClient.GetValidIngredient(ctx, createdValidIngredient.ID)
				requireNotNilAndNoProblems(t, createdValidIngredient, err)
				checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

				exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()
				exampleValidIngredientStateIngredient.Ingredient = *createdValidIngredient
				exampleValidIngredientStateIngredient.IngredientState = *createdValidIngredientState
				exampleValidIngredientStateIngredientInput := converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientCreationRequestInput(exampleValidIngredientStateIngredient)
				createdValidIngredientStateIngredient, createdValidIngredientStateIngredientErr := testClients.adminClient.CreateValidIngredientStateIngredient(ctx, exampleValidIngredientStateIngredientInput)
				require.NoError(t, createdValidIngredientStateIngredientErr)

				exampleValidIngredientStateIngredient.Ingredient = types.ValidIngredient{ID: createdValidIngredient.ID}
				exampleValidIngredientStateIngredient.IngredientState = types.ValidIngredientState{ID: createdValidIngredientState.ID}

				checkValidIngredientStateIngredientEquality(t, exampleValidIngredientStateIngredient, createdValidIngredientStateIngredient)

				expected = append(expected, createdValidIngredientStateIngredient)
			}

			// assert valid ingredient state ingredient list equality
			actual, err := testClients.userClient.GetValidIngredientStateIngredients(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdValidIngredientStateIngredient := range expected {
				assert.NoError(t, testClients.adminClient.ArchiveValidIngredientStateIngredient(ctx, createdValidIngredientStateIngredient.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidIngredientStateIngredients_Listing_ByValues() {
	s.runTest("should be findable via either member of the bridge type", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredientState := createValidIngredientStateForTest(t, ctx, testClients.adminClient)

			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, err := testClients.adminClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)

			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.userClient.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()
			exampleValidIngredientStateIngredient.Ingredient = *createdValidIngredient
			exampleValidIngredientStateIngredient.IngredientState = *createdValidIngredientState
			exampleValidIngredientStateIngredientInput := converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientCreationRequestInput(exampleValidIngredientStateIngredient)
			createdValidIngredientStateIngredient, err := testClients.adminClient.CreateValidIngredientStateIngredient(ctx, exampleValidIngredientStateIngredientInput)
			require.NoError(t, err)

			checkValidIngredientStateIngredientEquality(t, exampleValidIngredientStateIngredient, createdValidIngredientStateIngredient)

			validIngredientMeasurementUnitsForValidIngredient, err := testClients.userClient.GetValidIngredientStateIngredientsForIngredient(ctx, createdValidIngredient.ID, nil)
			requireNotNilAndNoProblems(t, validIngredientMeasurementUnitsForValidIngredient, err)

			require.Len(t, validIngredientMeasurementUnitsForValidIngredient.Data, 1)
			assert.Equal(t, validIngredientMeasurementUnitsForValidIngredient.Data[0].ID, createdValidIngredientStateIngredient.ID)

			validIngredientMeasurementUnitsForValidMeasurementUnit, err := testClients.userClient.GetValidIngredientStateIngredientsForIngredientState(ctx, createdValidIngredientState.ID, nil)
			requireNotNilAndNoProblems(t, validIngredientMeasurementUnitsForValidMeasurementUnit, err)

			require.Len(t, validIngredientMeasurementUnitsForValidMeasurementUnit.Data, 1)
			assert.Equal(t, validIngredientMeasurementUnitsForValidMeasurementUnit.Data[0].ID, createdValidIngredientStateIngredient.ID)

			assert.NoError(t, testClients.adminClient.ArchiveValidIngredientStateIngredient(ctx, createdValidIngredientStateIngredient.ID))

			assert.NoError(t, testClients.adminClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))

			assert.NoError(t, testClients.adminClient.ArchiveValidIngredientState(ctx, createdValidIngredientState.ID))
		}
	})
}
