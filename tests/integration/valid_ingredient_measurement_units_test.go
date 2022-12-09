package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"
)

func checkValidIngredientMeasurementUnitEquality(t *testing.T, expected, actual *types.ValidIngredientMeasurementUnit) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected StatusExplanation for valid ingredient measurement unit %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.MeasurementUnit.ID, actual.MeasurementUnit.ID, "expected MeasurementUnit for valid ingredient measurement unit %s to be %v, but it was %v", expected.ID, expected.MeasurementUnit.ID, actual.MeasurementUnit.ID)
	assert.Equal(t, expected.Ingredient.ID, actual.Ingredient.ID, "expected Ingredient for valid ingredient measurement unit %s to be %v, but it was %v", expected.ID, expected.Ingredient.ID, actual.Ingredient.ID)
	assert.Equal(t, expected.MinimumAllowableQuantity, actual.MinimumAllowableQuantity, "expected MinimumAllowableQuantity for valid ingredient measurement unit %s to be %v, but it was %v", expected.ID, expected.MinimumAllowableQuantity, actual.MinimumAllowableQuantity)
	assert.Equal(t, expected.MaximumAllowableQuantity, actual.MaximumAllowableQuantity, "expected MaximumAllowableQuantity for valid ingredient measurement unit %s to be %v, but it was %v", expected.ID, expected.MaximumAllowableQuantity, actual.MaximumAllowableQuantity)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestValidIngredientMeasurementUnits_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite valid measurement unit")
			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnitInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleValidMeasurementUnit)
			createdValidMeasurementUnit, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", createdValidMeasurementUnit.ID)

			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			createdValidMeasurementUnit, err = testClients.user.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

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

			t.Log("creating valid ingredient measurement unit")
			exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()
			exampleValidIngredientMeasurementUnit.Ingredient = *createdValidIngredient
			exampleValidIngredientMeasurementUnit.MeasurementUnit = *createdValidMeasurementUnit
			exampleValidIngredientMeasurementUnitInput := converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(exampleValidIngredientMeasurementUnit)
			createdValidIngredientMeasurementUnit, err := testClients.admin.CreateValidIngredientMeasurementUnit(ctx, exampleValidIngredientMeasurementUnitInput)
			require.NoError(t, err)
			t.Logf("valid ingredient measurement unit %q created", createdValidIngredientMeasurementUnit.ID)

			checkValidIngredientMeasurementUnitEquality(t, exampleValidIngredientMeasurementUnit, createdValidIngredientMeasurementUnit)

			createdValidIngredientMeasurementUnit, err = testClients.user.GetValidIngredientMeasurementUnit(ctx, createdValidIngredientMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, createdValidIngredientMeasurementUnit, err)

			checkValidIngredientMeasurementUnitEquality(t, exampleValidIngredientMeasurementUnit, createdValidIngredientMeasurementUnit)

			t.Log("changing valid ingredient measurement unit")
			newValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()
			newValidIngredientMeasurementUnit.Ingredient = *createdValidIngredient
			newValidIngredientMeasurementUnit.MeasurementUnit = *createdValidMeasurementUnit
			createdValidIngredientMeasurementUnit.Update(converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitUpdateRequestInput(newValidIngredientMeasurementUnit))
			assert.NoError(t, testClients.admin.UpdateValidIngredientMeasurementUnit(ctx, createdValidIngredientMeasurementUnit))

			t.Log("fetching changed valid ingredient measurement unit")
			actual, err := testClients.user.GetValidIngredientMeasurementUnit(ctx, createdValidIngredientMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient measurement unit equality
			checkValidIngredientMeasurementUnitEquality(t, newValidIngredientMeasurementUnit, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			searchedMeasurementUnits, err := testClients.user.SearchValidMeasurementUnitsByIngredientID(ctx, createdValidIngredient.ID, types.DefaultQueryFilter())
			requireNotNilAndNoProblems(t, searchedMeasurementUnits, err)
			assert.GreaterOrEqual(t, len(searchedMeasurementUnits.Data), 1)

			t.Log("cleaning up valid ingredient measurement unit")
			assert.NoError(t, testClients.admin.ArchiveValidIngredientMeasurementUnit(ctx, createdValidIngredientMeasurementUnit.ID))

			t.Log("cleaning up valid measurement unit")
			assert.NoError(t, testClients.admin.ArchiveValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID))

			t.Log("cleaning up valid ingredient")
			assert.NoError(t, testClients.admin.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredientMeasurementUnits_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid ingredient measurement units")
			var expected []*types.ValidIngredientMeasurementUnit
			for i := 0; i < 5; i++ {
				t.Log("creating prerequisite valid measurement unit")
				exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
				exampleValidMeasurementUnitInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleValidMeasurementUnit)
				createdValidMeasurementUnit, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
				require.NoError(t, err)
				t.Logf("valid measurement unit %q created", createdValidMeasurementUnit.ID)

				checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

				createdValidMeasurementUnit, err = testClients.user.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
				requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
				checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

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

				exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()
				exampleValidIngredientMeasurementUnit.Ingredient = *createdValidIngredient
				exampleValidIngredientMeasurementUnit.MeasurementUnit = *createdValidMeasurementUnit
				exampleValidIngredientMeasurementUnitInput := converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(exampleValidIngredientMeasurementUnit)
				createdValidIngredientMeasurementUnit, createdValidIngredientMeasurementUnitErr := testClients.admin.CreateValidIngredientMeasurementUnit(ctx, exampleValidIngredientMeasurementUnitInput)
				require.NoError(t, createdValidIngredientMeasurementUnitErr)

				checkValidIngredientMeasurementUnitEquality(t, exampleValidIngredientMeasurementUnit, createdValidIngredientMeasurementUnit)

				expected = append(expected, createdValidIngredientMeasurementUnit)
			}

			// assert valid ingredient measurement unit list equality
			actual, err := testClients.user.GetValidIngredientMeasurementUnits(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			t.Log("cleaning up")
			for _, createdValidIngredientMeasurementUnit := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidIngredientMeasurementUnit(ctx, createdValidIngredientMeasurementUnit.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidIngredientMeasurementUnits_Listing_ByValues() {
	s.runForEachClient("should be findable via either member of the bridge type", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite valid measurement unit")
			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnitInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleValidMeasurementUnit)
			createdValidMeasurementUnit, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", createdValidMeasurementUnit.ID)

			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			createdValidMeasurementUnit, err = testClients.user.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

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

			t.Log("creating valid ingredient measurement unit")
			exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()
			exampleValidIngredientMeasurementUnit.Ingredient = *createdValidIngredient
			exampleValidIngredientMeasurementUnit.MeasurementUnit = *createdValidMeasurementUnit
			exampleValidIngredientMeasurementUnitInput := converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(exampleValidIngredientMeasurementUnit)
			createdValidIngredientMeasurementUnit, err := testClients.admin.CreateValidIngredientMeasurementUnit(ctx, exampleValidIngredientMeasurementUnitInput)
			require.NoError(t, err)
			t.Logf("valid ingredient measurement unit %q created", createdValidIngredientMeasurementUnit.ID)

			checkValidIngredientMeasurementUnitEquality(t, exampleValidIngredientMeasurementUnit, createdValidIngredientMeasurementUnit)

			validIngredientMeasurementUnitsForValidIngredient, err := testClients.user.GetValidIngredientMeasurementUnitsForIngredient(ctx, createdValidIngredient.ID, nil)
			requireNotNilAndNoProblems(t, validIngredientMeasurementUnitsForValidIngredient, err)

			require.Len(t, validIngredientMeasurementUnitsForValidIngredient.Data, 1)
			assert.Equal(t, validIngredientMeasurementUnitsForValidIngredient.Data[0].ID, createdValidIngredientMeasurementUnit.ID)

			validIngredientMeasurementUnitsForValidMeasurementUnit, err := testClients.user.GetValidIngredientMeasurementUnitsForMeasurementUnit(ctx, createdValidMeasurementUnit.ID, nil)
			requireNotNilAndNoProblems(t, validIngredientMeasurementUnitsForValidMeasurementUnit, err)

			require.Len(t, validIngredientMeasurementUnitsForValidMeasurementUnit.Data, 1)
			assert.Equal(t, validIngredientMeasurementUnitsForValidMeasurementUnit.Data[0].ID, createdValidIngredientMeasurementUnit.ID)
		}
	})
}
