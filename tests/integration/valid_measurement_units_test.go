package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func checkValidMeasurementUnitEquality(t *testing.T, expected, actual *types.ValidMeasurementUnit) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Volumetric, actual.Volumetric, "expected Volumetric for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.Volumetric, actual.Volumetric)
	assert.Equal(t, expected.IconPath, actual.IconPath, "expected IconPath for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.IconPath, actual.IconPath)
	assert.Equal(t, expected.Universal, actual.Universal, "expected Universal for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.Universal, actual.Universal)
	assert.Equal(t, expected.Metric, actual.Metric, "expected Metric for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.Metric, actual.Metric)
	assert.Equal(t, expected.Imperial, actual.Imperial, "expected Imperial for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.Imperial, actual.Imperial)
	assert.Equal(t, expected.PluralName, actual.PluralName, "expected PluralName for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.PluralName, actual.PluralName)
	assert.NotZero(t, actual.CreatedAt)
}

// convertValidMeasurementUnitToValidMeasurementUnitUpdateInput creates an ValidMeasurementUnitUpdateRequestInput struct from a valid measurement unit.
func convertValidMeasurementUnitToValidMeasurementUnitUpdateInput(x *types.ValidMeasurementUnit) *types.ValidMeasurementUnitUpdateRequestInput {
	return &types.ValidMeasurementUnitUpdateRequestInput{
		Name:        &x.Name,
		Description: &x.Description,
		Volumetric:  &x.Volumetric,
		IconPath:    &x.IconPath,
		Universal:   &x.Universal,
		Metric:      &x.Metric,
		Imperial:    &x.Imperial,
		PluralName:  &x.PluralName,
	}
}

func (s *TestSuite) TestValidMeasurementUnits_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid measurement unit")
			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnitInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInputFromValidMeasurementUnit(exampleValidMeasurementUnit)
			createdValidMeasurementUnit, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", createdValidMeasurementUnit.ID)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			createdValidMeasurementUnit, err = testClients.admin.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			t.Log("changing valid measurement unit")
			newValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			createdValidMeasurementUnit.Update(convertValidMeasurementUnitToValidMeasurementUnitUpdateInput(newValidMeasurementUnit))
			assert.NoError(t, testClients.admin.UpdateValidMeasurementUnit(ctx, createdValidMeasurementUnit))

			t.Log("fetching changed valid measurement unit")
			actual, err := testClients.admin.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid measurement unit equality
			checkValidMeasurementUnitEquality(t, newValidMeasurementUnit, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			t.Log("cleaning up valid measurement unit")
			assert.NoError(t, testClients.admin.ArchiveValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID))
		}
	})
}

func (s *TestSuite) TestValidMeasurementUnits_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid measurement units")
			var expected []*types.ValidMeasurementUnit
			for i := 0; i < 5; i++ {
				exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
				exampleValidMeasurementUnitInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInputFromValidMeasurementUnit(exampleValidMeasurementUnit)
				createdValidMeasurementUnit, createdValidMeasurementUnitErr := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
				require.NoError(t, createdValidMeasurementUnitErr)
				t.Logf("valid measurement unit %q created", createdValidMeasurementUnit.ID)

				checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

				expected = append(expected, createdValidMeasurementUnit)
			}

			// assert valid measurement unit list equality
			actual, err := testClients.admin.GetValidMeasurementUnits(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidMeasurementUnits),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidMeasurementUnits),
			)

			t.Log("cleaning up")
			for _, createdValidMeasurementUnit := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidMeasurementUnits_Searching() {
	s.runForEachClient("should be able to be search for valid measurement units", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid measurement units")
			var expected []*types.ValidMeasurementUnit
			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnit.Name = fmt.Sprintf("example_%s", testClients.authType)
			searchQuery := exampleValidMeasurementUnit.Name
			for i := 0; i < 5; i++ {
				exampleValidMeasurementUnit.Name = fmt.Sprintf("%s %d", searchQuery, i)
				exampleValidMeasurementUnitInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInputFromValidMeasurementUnit(exampleValidMeasurementUnit)
				createdValidMeasurementUnit, createdValidMeasurementUnitErr := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
				require.NoError(t, createdValidMeasurementUnitErr)
				checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)
				t.Logf("valid measurement unit %q created", createdValidMeasurementUnit.ID)

				expected = append(expected, createdValidMeasurementUnit)
			}

			exampleLimit := uint8(20)

			// assert valid measurement unit list equality
			actual, err := testClients.admin.SearchValidMeasurementUnits(ctx, searchQuery, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected %d to be <= %d",
				len(expected),
				len(actual),
			)

			t.Log("cleaning up")
			for _, createdValidMeasurementUnit := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID))
			}
		}
	})
}
