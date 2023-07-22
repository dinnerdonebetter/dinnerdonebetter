package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/require"
)

func checkValidMeasurementConversionEquality(t *testing.T, expected, actual *types.ValidMeasurementUnitConversion) {
	t.Helper()

	require.NotZero(t, actual.ID)

	require.Equal(t, expected.OnlyForIngredient, actual.OnlyForIngredient, "expected OnlyForIngredient for valid measurement conversion %s to be %v, but it was %v", expected.ID, expected.OnlyForIngredient, actual.OnlyForIngredient)
	require.Equal(t, expected.Notes, actual.Notes, "expected Notes for valid measurement conversion %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	require.Equal(t, expected.From.ID, actual.From.ID, "expected From for valid measurement conversion %s to be %v, but it was %v", expected.ID, expected.From.ID, actual.From.ID)
	require.Equal(t, expected.To.ID, actual.To.ID, "expected To for valid measurement conversion %s to be %v, but it was %v", expected.ID, expected.To.ID, actual.To.ID)
	require.Equal(t, expected.Modifier, actual.Modifier, "expected Modifier for valid measurement conversion %s to be %v, but it was %v", expected.ID, expected.Modifier, actual.Modifier)

	require.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestValidMeasurementConversions_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable without ingredient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidMeasurementUnit1 := createValidMeasurementUnitForTest(t, ctx, testClients.admin)
			createdValidMeasurementUnit2 := createValidMeasurementUnitForTest(t, ctx, testClients.admin)

			exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()
			exampleValidMeasurementConversion.From = *createdValidMeasurementUnit1
			exampleValidMeasurementConversion.To = *createdValidMeasurementUnit2
			exampleValidMeasurementConversionInput := converters.ConvertValidMeasurementConversionToValidMeasurementConversionCreationRequestInput(exampleValidMeasurementConversion)

			createdValidMeasurementConversion, err := testClients.admin.CreateValidMeasurementConversion(ctx, exampleValidMeasurementConversionInput)
			require.NoError(t, err)
			checkValidMeasurementConversionEquality(t, exampleValidMeasurementConversion, createdValidMeasurementConversion)

			createdValidMeasurementConversion, err = testClients.admin.GetValidMeasurementUnitConversion(ctx, createdValidMeasurementConversion.ID)
			requireNotNilAndNoProblems(t, createdValidMeasurementConversion, err)

			checkValidMeasurementConversionEquality(t, exampleValidMeasurementConversion, createdValidMeasurementConversion)

			createdValidMeasurementConversion.Modifier = fakes.BuildFakeValidMeasurementConversion().Modifier
			require.NoError(t, testClients.admin.UpdateValidMeasurementConversion(ctx, createdValidMeasurementConversion))

			actual, err := testClients.admin.GetValidMeasurementUnitConversion(ctx, createdValidMeasurementConversion.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid measurement conversion equality
			checkValidMeasurementConversionEquality(t, createdValidMeasurementConversion, actual)
			require.NotNil(t, actual.LastUpdatedAt)

			require.NoError(t, testClients.admin.ArchiveValidMeasurementConversion(ctx, createdValidMeasurementConversion.ID))
		}
	})
}

func (s *TestSuite) TestValidMeasurementConversions_GetFromUnits() {
	s.runForEachClient("should be able to get what a unit converts from", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidMeasurementUnit1 := createValidMeasurementUnitForTest(t, ctx, testClients.admin)
			createdValidMeasurementUnit2 := createValidMeasurementUnitForTest(t, ctx, testClients.admin)

			exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()
			exampleValidMeasurementConversion.From = *createdValidMeasurementUnit1
			exampleValidMeasurementConversion.To = *createdValidMeasurementUnit2
			exampleValidMeasurementConversionInput := converters.ConvertValidMeasurementConversionToValidMeasurementConversionCreationRequestInput(exampleValidMeasurementConversion)

			createdValidMeasurementConversion, err := testClients.admin.CreateValidMeasurementConversion(ctx, exampleValidMeasurementConversionInput)
			require.NoError(t, err)
			checkValidMeasurementConversionEquality(t, exampleValidMeasurementConversion, createdValidMeasurementConversion)

			createdValidMeasurementConversion, err = testClients.admin.GetValidMeasurementUnitConversion(ctx, createdValidMeasurementConversion.ID)
			requireNotNilAndNoProblems(t, createdValidMeasurementConversion, err)
			checkValidMeasurementConversionEquality(t, exampleValidMeasurementConversion, createdValidMeasurementConversion)

			fromUnits, err := testClients.admin.GetValidMeasurementUnitConversionsFromUnit(ctx, createdValidMeasurementUnit1.ID)
			requireNotNilAndNoProblems(t, fromUnits, err)
			require.Equal(t, 1, len(fromUnits))

			require.NoError(t, testClients.admin.ArchiveValidMeasurementConversion(ctx, createdValidMeasurementConversion.ID))
		}
	})
}

func (s *TestSuite) TestValidMeasurementConversions_GetToUnits() {
	s.runForEachClient("should be able to get what a unit converts to", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidMeasurementUnit1 := createValidMeasurementUnitForTest(t, ctx, testClients.admin)
			createdValidMeasurementUnit2 := createValidMeasurementUnitForTest(t, ctx, testClients.admin)

			exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()
			exampleValidMeasurementConversion.From = *createdValidMeasurementUnit1
			exampleValidMeasurementConversion.To = *createdValidMeasurementUnit2
			exampleValidMeasurementConversionInput := converters.ConvertValidMeasurementConversionToValidMeasurementConversionCreationRequestInput(exampleValidMeasurementConversion)

			createdValidMeasurementConversion, err := testClients.admin.CreateValidMeasurementConversion(ctx, exampleValidMeasurementConversionInput)
			require.NoError(t, err)
			checkValidMeasurementConversionEquality(t, exampleValidMeasurementConversion, createdValidMeasurementConversion)

			createdValidMeasurementConversion, err = testClients.admin.GetValidMeasurementUnitConversion(ctx, createdValidMeasurementConversion.ID)
			requireNotNilAndNoProblems(t, createdValidMeasurementConversion, err)
			checkValidMeasurementConversionEquality(t, exampleValidMeasurementConversion, createdValidMeasurementConversion)

			fromUnits, err := testClients.admin.GetValidMeasurementUnitConversionToUnit(ctx, createdValidMeasurementUnit2.ID)
			requireNotNilAndNoProblems(t, fromUnits, err)
			require.Equal(t, 1, len(fromUnits))

			require.NoError(t, testClients.admin.ArchiveValidMeasurementConversion(ctx, createdValidMeasurementConversion.ID))
		}
	})
}
