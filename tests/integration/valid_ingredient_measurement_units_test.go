package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func checkValidIngredientMeasurementUnitEquality(t *testing.T, expected, actual *types.ValidIngredientMeasurementUnit) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for valid ingredient measurement unit %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.ValidMeasurementUnitID, actual.ValidMeasurementUnitID, "expected ValidMeasurementUnitID for valid ingredient measurement unit %s to be %v, but it was %v", expected.ID, expected.ValidMeasurementUnitID, actual.ValidMeasurementUnitID)
	assert.Equal(t, expected.ValidIngredientID, actual.ValidIngredientID, "expected ValidIngredientID for valid ingredient measurement unit %s to be %v, but it was %v", expected.ID, expected.ValidIngredientID, actual.ValidIngredientID)
	assert.NotZero(t, actual.CreatedOn)
}

// convertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitUpdateInput creates an ValidIngredientMeasurementUnitUpdateRequestInput struct from a valid ingredient measurement unit.
func convertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitUpdateInput(x *types.ValidIngredientMeasurementUnit) *types.ValidIngredientMeasurementUnitUpdateRequestInput {
	return &types.ValidIngredientMeasurementUnitUpdateRequestInput{
		Notes:                  &x.Notes,
		ValidMeasurementUnitID: &x.ValidMeasurementUnitID,
		ValidIngredientID:      &x.ValidIngredientID,
	}
}

func (s *TestSuite) TestValidIngredientMeasurementUnits_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite valid measurement unit")
			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnitInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInputFromValidMeasurementUnit(exampleValidMeasurementUnit)
			createdValidMeasurementUnit, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", createdValidMeasurementUnit.ID)

			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			createdValidMeasurementUnit, err = testClients.main.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			t.Log("creating prerequisite valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)

			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.main.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)
			t.Logf("valid ingredient %q created", createdValidIngredient.ID)

			t.Log("creating valid ingredient measurement unit")
			exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()
			exampleValidIngredientMeasurementUnit.ValidIngredientID = createdValidIngredient.ID
			exampleValidIngredientMeasurementUnit.ValidMeasurementUnitID = createdValidMeasurementUnit.ID
			exampleValidIngredientMeasurementUnitInput := fakes.BuildFakeValidIngredientMeasurementUnitCreationRequestInputFromValidIngredientMeasurementUnit(exampleValidIngredientMeasurementUnit)
			createdValidIngredientMeasurementUnit, err := testClients.admin.CreateValidIngredientMeasurementUnit(ctx, exampleValidIngredientMeasurementUnitInput)
			require.NoError(t, err)
			t.Logf("valid ingredient measurement unit %q created", createdValidIngredientMeasurementUnit.ID)

			checkValidIngredientMeasurementUnitEquality(t, exampleValidIngredientMeasurementUnit, createdValidIngredientMeasurementUnit)

			createdValidIngredientMeasurementUnit, err = testClients.main.GetValidIngredientMeasurementUnit(ctx, createdValidIngredientMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, createdValidIngredientMeasurementUnit, err)

			checkValidIngredientMeasurementUnitEquality(t, exampleValidIngredientMeasurementUnit, createdValidIngredientMeasurementUnit)

			t.Log("changing valid ingredient measurement unit")
			newValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()
			newValidIngredientMeasurementUnit.ValidIngredientID = createdValidIngredient.ID
			newValidIngredientMeasurementUnit.ValidMeasurementUnitID = createdValidMeasurementUnit.ID
			createdValidIngredientMeasurementUnit.Update(convertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitUpdateInput(newValidIngredientMeasurementUnit))
			assert.NoError(t, testClients.admin.UpdateValidIngredientMeasurementUnit(ctx, createdValidIngredientMeasurementUnit))

			t.Log("fetching changed valid ingredient measurement unit")
			actual, err := testClients.main.GetValidIngredientMeasurementUnit(ctx, createdValidIngredientMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient measurement unit equality
			checkValidIngredientMeasurementUnitEquality(t, newValidIngredientMeasurementUnit, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up valid ingredient measurement unit")
			assert.NoError(t, testClients.admin.ArchiveValidIngredientMeasurementUnit(ctx, createdValidIngredientMeasurementUnit.ID))
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
				exampleValidMeasurementUnitInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInputFromValidMeasurementUnit(exampleValidMeasurementUnit)
				createdValidMeasurementUnit, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
				require.NoError(t, err)
				t.Logf("valid measurement unit %q created", createdValidMeasurementUnit.ID)

				checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

				createdValidMeasurementUnit, err = testClients.main.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
				requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
				checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

				t.Log("creating prerequisite valid ingredient")
				exampleValidIngredient := fakes.BuildFakeValidIngredient()
				exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
				createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
				require.NoError(t, err)

				checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

				createdValidIngredient, err = testClients.main.GetValidIngredient(ctx, createdValidIngredient.ID)
				requireNotNilAndNoProblems(t, createdValidIngredient, err)
				checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)
				t.Logf("valid ingredient %q created", createdValidIngredient.ID)

				exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()
				exampleValidIngredientMeasurementUnit.ValidIngredientID = createdValidIngredient.ID
				exampleValidIngredientMeasurementUnit.ValidMeasurementUnitID = createdValidMeasurementUnit.ID
				exampleValidIngredientMeasurementUnitInput := fakes.BuildFakeValidIngredientMeasurementUnitCreationRequestInputFromValidIngredientMeasurementUnit(exampleValidIngredientMeasurementUnit)
				createdValidIngredientMeasurementUnit, createdValidIngredientMeasurementUnitErr := testClients.admin.CreateValidIngredientMeasurementUnit(ctx, exampleValidIngredientMeasurementUnitInput)
				require.NoError(t, createdValidIngredientMeasurementUnitErr)

				checkValidIngredientMeasurementUnitEquality(t, exampleValidIngredientMeasurementUnit, createdValidIngredientMeasurementUnit)

				expected = append(expected, createdValidIngredientMeasurementUnit)
			}

			// assert valid ingredient measurement unit list equality
			actual, err := testClients.main.GetValidIngredientMeasurementUnits(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidIngredientMeasurementUnits),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidIngredientMeasurementUnits),
			)

			t.Log("cleaning up")
			for _, createdValidIngredientMeasurementUnit := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidIngredientMeasurementUnit(ctx, createdValidIngredientMeasurementUnit.ID))
			}
		}
	})
}
