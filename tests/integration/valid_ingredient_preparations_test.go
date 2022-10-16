package integration

import (
	"github.com/prixfixeco/api_server/pkg/types/converters"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func checkValidIngredientPreparationEquality(t *testing.T, expected, actual *types.ValidIngredientPreparation) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected StatusExplanation for valid ingredient preparation %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.Preparation.ID, actual.Preparation.ID, "expected Preparation for valid ingredient preparation %s to be %v, but it was %v", expected.ID, expected.Preparation.ID, actual.Preparation.ID)
	assert.Equal(t, expected.Ingredient.ID, actual.Ingredient.ID, "expected Ingredient for valid ingredient preparation %s to be %v, but it was %v", expected.ID, expected.Ingredient.ID, actual.Ingredient.ID)
	assert.NotZero(t, actual.CreatedAt)
}

// convertValidIngredientPreparationToValidIngredientPreparationUpdateInput creates an ValidIngredientPreparationUpdateRequestInput struct from a valid ingredient preparation.
func convertValidIngredientPreparationToValidIngredientPreparationUpdateInput(x *types.ValidIngredientPreparation) *types.ValidIngredientPreparationUpdateRequestInput {
	return &types.ValidIngredientPreparationUpdateRequestInput{
		Notes:              &x.Notes,
		ValidPreparationID: &x.Preparation.ID,
		ValidIngredientID:  &x.Ingredient.ID,
	}
}

func (s *TestSuite) TestValidIngredientPreparations_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparation.ID)

			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.user.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			t.Log("creating prerequisite valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)

			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.user.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)
			t.Logf("valid ingredient %q created", createdValidIngredient.ID)

			t.Log("creating valid ingredient preparation")
			exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparation.Ingredient = *createdValidIngredient
			exampleValidIngredientPreparation.Preparation = *createdValidPreparation
			exampleValidIngredientPreparationInput := converters.ConvertValidIngredientPreparationToValidIngredientPreparationCreationRequestInput(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := testClients.admin.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			require.NoError(t, err)
			t.Logf("valid ingredient preparation %q created", createdValidIngredientPreparation.ID)

			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, createdValidIngredientPreparation)

			createdValidIngredientPreparation, err = testClients.user.GetValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidIngredientPreparation, err)

			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, createdValidIngredientPreparation)

			t.Log("changing valid ingredient preparation")
			newValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			newValidIngredientPreparation.Ingredient = *createdValidIngredient
			newValidIngredientPreparation.Preparation = *createdValidPreparation
			createdValidIngredientPreparation.Update(convertValidIngredientPreparationToValidIngredientPreparationUpdateInput(newValidIngredientPreparation))
			assert.NoError(t, testClients.admin.UpdateValidIngredientPreparation(ctx, createdValidIngredientPreparation))

			t.Log("fetching changed valid ingredient preparation")
			actual, err := testClients.user.GetValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient preparation equality
			checkValidIngredientPreparationEquality(t, newValidIngredientPreparation, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			t.Log("cleaning up valid ingredient preparation")
			assert.NoError(t, testClients.admin.ArchiveValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID))

			t.Log("cleaning up valid ingredient preparation")
			assert.NoError(t, testClients.admin.ArchiveValidIngredient(ctx, createdValidIngredient.ID))

			t.Log("cleaning up valid ingredient preparation")
			assert.NoError(t, testClients.admin.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredientPreparations_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid ingredient preparations")
			var expected []*types.ValidIngredientPreparation
			for i := 0; i < 5; i++ {
				t.Log("creating prerequisite valid preparation")
				exampleValidPreparation := fakes.BuildFakeValidPreparation()
				exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
				createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
				require.NoError(t, err)
				t.Logf("valid preparation %q created", createdValidPreparation.ID)

				checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

				createdValidPreparation, err = testClients.user.GetValidPreparation(ctx, createdValidPreparation.ID)
				requireNotNilAndNoProblems(t, createdValidPreparation, err)
				checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

				t.Log("creating prerequisite valid ingredient")
				exampleValidIngredient := fakes.BuildFakeValidIngredient()
				exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
				createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
				require.NoError(t, err)

				checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

				createdValidIngredient, err = testClients.user.GetValidIngredient(ctx, createdValidIngredient.ID)
				requireNotNilAndNoProblems(t, createdValidIngredient, err)
				checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)
				t.Logf("valid ingredient %q created", createdValidIngredient.ID)

				exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
				exampleValidIngredientPreparation.Ingredient = *createdValidIngredient
				exampleValidIngredientPreparation.Preparation = *createdValidPreparation
				exampleValidIngredientPreparationInput := converters.ConvertValidIngredientPreparationToValidIngredientPreparationCreationRequestInput(exampleValidIngredientPreparation)
				createdValidIngredientPreparation, createdValidIngredientPreparationErr := testClients.admin.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
				require.NoError(t, createdValidIngredientPreparationErr)

				exampleValidIngredientPreparation.Ingredient = types.ValidIngredient{ID: createdValidIngredient.ID}
				exampleValidIngredientPreparation.Preparation = types.ValidPreparation{ID: createdValidPreparation.ID}

				checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, createdValidIngredientPreparation)

				expected = append(expected, createdValidIngredientPreparation)
			}

			// assert valid ingredient preparation list equality
			actual, err := testClients.user.GetValidIngredientPreparations(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidIngredientPreparations),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidIngredientPreparations),
			)

			t.Log("cleaning up")
			for _, createdValidIngredientPreparation := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidIngredientPreparations_Listing_ByValues() {
	s.runForEachClient("should be findable via either member of the bridge type", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparation.ID)

			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.user.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			t.Log("creating prerequisite valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)

			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.user.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)
			t.Logf("valid ingredient %q created", createdValidIngredient.ID)

			t.Log("creating valid ingredient preparation")
			exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparation.Ingredient = *createdValidIngredient
			exampleValidIngredientPreparation.Preparation = *createdValidPreparation
			exampleValidIngredientPreparationInput := converters.ConvertValidIngredientPreparationToValidIngredientPreparationCreationRequestInput(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := testClients.admin.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			require.NoError(t, err)
			t.Logf("valid ingredient preparation %q created", createdValidIngredientPreparation.ID)

			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, createdValidIngredientPreparation)

			validIngredientMeasurementUnitsForValidIngredient, err := testClients.user.GetValidIngredientPreparationsForIngredient(ctx, createdValidIngredient.ID, nil)
			requireNotNilAndNoProblems(t, validIngredientMeasurementUnitsForValidIngredient, err)

			require.Len(t, validIngredientMeasurementUnitsForValidIngredient.ValidIngredientPreparations, 1)
			assert.Equal(t, validIngredientMeasurementUnitsForValidIngredient.ValidIngredientPreparations[0].ID, createdValidIngredientPreparation.ID)

			validIngredientMeasurementUnitsForValidMeasurementUnit, err := testClients.user.GetValidIngredientPreparationsForPreparation(ctx, createdValidPreparation.ID, nil)
			requireNotNilAndNoProblems(t, validIngredientMeasurementUnitsForValidMeasurementUnit, err)

			require.Len(t, validIngredientMeasurementUnitsForValidMeasurementUnit.ValidIngredientPreparations, 1)
			assert.Equal(t, validIngredientMeasurementUnitsForValidMeasurementUnit.ValidIngredientPreparations[0].ID, createdValidIngredientPreparation.ID)

			t.Log("cleaning up valid ingredient preparation")
			assert.NoError(t, testClients.admin.ArchiveValidIngredientPreparation(ctx, createdValidIngredientPreparation.ID))

			t.Log("cleaning up valid ingredient")
			assert.NoError(t, testClients.admin.ArchiveValidIngredient(ctx, createdValidIngredient.ID))

			t.Log("cleaning up valid preparation")
			assert.NoError(t, testClients.admin.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		}
	})
}
