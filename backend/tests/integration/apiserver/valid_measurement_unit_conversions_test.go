package integration

import (
	"testing"

	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	mealplanningconverters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createValidMeasurementUnitConversionForTest(t *testing.T) (*types.ValidMeasurementUnit, *types.ValidMeasurementUnit, *types.ValidMeasurementUnitConversion) {
	t.Helper()
	ctx := t.Context()

	unit1 := createValidMeasurementUnitForTest(t)
	unit2 := createValidMeasurementUnitForTest(t)

	exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()
	exampleValidMeasurementUnitConversion.To = *unit1
	exampleValidMeasurementUnitConversion.From = *unit2

	exampleValidMeasurementUnitConversionInput := mealplanningconverters.ConvertCreateValidMeasurementUnitConversionRequestToGRPCValidMeasurementUnitConversionCreationRequestInput(converters.ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionCreationRequestInput(exampleValidMeasurementUnitConversion))
	createdValidMeasurementUnitConversion, err := adminClient.CreateValidMeasurementUnitConversion(ctx, &mealplanningsvc.CreateValidMeasurementUnitConversionRequest{Input: exampleValidMeasurementUnitConversionInput})
	require.NoError(t, err)
	require.NotNil(t, createdValidMeasurementUnitConversion)

	validPrepPreparationRes, err := adminClient.GetValidMeasurementUnitConversion(ctx, &mealplanningsvc.GetValidMeasurementUnitConversionRequest{
		ValidMeasurementUnitConversionID: createdValidMeasurementUnitConversion.Result.ID,
	})
	require.NoError(t, err)
	require.NotNil(t, validPrepPreparationRes.Result)

	return unit1, unit2, mealplanningconverters.ConvertGRPCValidMeasurementUnitConversionToValidMeasurementUnitConversion(validPrepPreparationRes.Result)
}

func TestValidMeasurementUnitConversions_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		createValidMeasurementUnitConversionForTest(t)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()
		exampleValidMeasurementUnitConversionInput := mealplanningconverters.ConvertCreateValidMeasurementUnitConversionRequestToGRPCValidMeasurementUnitConversionCreationRequestInput(converters.ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionCreationRequestInput(exampleValidMeasurementUnitConversion))
		exampleValidMeasurementUnitConversionInput.To = ""
		exampleValidMeasurementUnitConversionInput.From = ""

		createdValidMeasurementUnitConversion, err := adminClient.CreateValidMeasurementUnitConversion(ctx, &mealplanningsvc.CreateValidMeasurementUnitConversionRequest{Input: exampleValidMeasurementUnitConversionInput})
		require.Error(t, err)
		require.Nil(t, createdValidMeasurementUnitConversion)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.CreateValidMeasurementUnitConversion(ctx, &mealplanningsvc.CreateValidMeasurementUnitConversionRequest{})
		assert.Error(t, err)
	})
}

func TestValidMeasurementUnitConversions_Listing(T *testing.T) {
	T.Parallel()

	createdValidMeasurementUnitConversions := []*types.ValidMeasurementUnitConversion{}
	toUnit, fromUnit, created := createValidMeasurementUnitConversionForTest(T)
	createdValidMeasurementUnitConversions = append(createdValidMeasurementUnitConversions, created)
	for range exampleQuantity - 1 {
		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()
		exampleValidMeasurementUnitConversionInput := mealplanningconverters.ConvertCreateValidMeasurementUnitConversionRequestToGRPCValidMeasurementUnitConversionCreationRequestInput(converters.ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionCreationRequestInput(exampleValidMeasurementUnitConversion))
		exampleValidMeasurementUnitConversionInput.To = toUnit.ID
		exampleValidMeasurementUnitConversionInput.From = fromUnit.ID

		createdValidMeasurementUnitConversion, err := adminClient.CreateValidMeasurementUnitConversion(T.Context(), &mealplanningsvc.CreateValidMeasurementUnitConversionRequest{Input: exampleValidMeasurementUnitConversionInput})
		require.NoError(T, err)
		require.NotNil(T, createdValidMeasurementUnitConversion)

		createdValidMeasurementUnitConversions = append(createdValidMeasurementUnitConversions, mealplanningconverters.ConvertGRPCValidMeasurementUnitConversionToValidMeasurementUnitConversion(createdValidMeasurementUnitConversion.Result))
	}

	T.Run("to unit", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidMeasurementUnitConversionsToUnit(ctx, &mealplanningsvc.GetValidMeasurementUnitConversionsToUnitRequest{
			ValidMeasurementUnitID: toUnit.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidMeasurementUnitConversions))
	})

	T.Run("from unit", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidMeasurementUnitConversionsFromUnit(ctx, &mealplanningsvc.GetValidMeasurementUnitConversionsFromUnitRequest{
			ValidMeasurementUnitID: fromUnit.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidMeasurementUnitConversions))
	})
}
