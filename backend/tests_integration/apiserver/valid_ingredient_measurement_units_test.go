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

func createValidIngredientMeasurementUnitForTest(t *testing.T) (*types.ValidIngredient, *types.ValidMeasurementUnit, *types.ValidIngredientMeasurementUnit) {
	t.Helper()
	ctx := t.Context()

	createdValidIngredient := createValidIngredientForTest(t)
	createdValidMeasurementUnit := createValidMeasurementUnitForTest(t)

	exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()
	exampleValidIngredientMeasurementUnit.MeasurementUnit = *createdValidMeasurementUnit
	exampleValidIngredientMeasurementUnit.Ingredient = *createdValidIngredient

	exampleValidIngredientMeasurementUnitInput := mealplanningconverters.ConvertCreateValidIngredientMeasurementUnitRequestToGRPCValidIngredientMeasurementUnitCreationRequestInput(converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(exampleValidIngredientMeasurementUnit))
	createdValidIngredientMeasurementUnit, err := adminClient.CreateValidIngredientMeasurementUnit(ctx, &mealplanningsvc.CreateValidIngredientMeasurementUnitRequest{Input: exampleValidIngredientMeasurementUnitInput})
	require.NoError(t, err)
	require.NotNil(t, createdValidIngredientMeasurementUnit)

	validPrepMeasurementUnitRes, err := adminClient.GetValidIngredientMeasurementUnit(ctx, &mealplanningsvc.GetValidIngredientMeasurementUnitRequest{
		ValidIngredientMeasurementUnitId: createdValidIngredientMeasurementUnit.Result.Id,
	})
	require.NoError(t, err)
	require.NotNil(t, validPrepMeasurementUnitRes.Result)

	return createdValidIngredient, createdValidMeasurementUnit, mealplanningconverters.ConvertGRPCValidIngredientMeasurementUnitToValidIngredientMeasurementUnit(validPrepMeasurementUnitRes.Result)
}

func TestValidIngredientMeasurementUnits_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		createValidIngredientMeasurementUnitForTest(t)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()
		exampleValidIngredientMeasurementUnitInput := mealplanningconverters.ConvertCreateValidIngredientMeasurementUnitRequestToGRPCValidIngredientMeasurementUnitCreationRequestInput(converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(exampleValidIngredientMeasurementUnit))
		exampleValidIngredientMeasurementUnitInput.ValidMeasurementUnitId = ""
		exampleValidIngredientMeasurementUnitInput.ValidIngredientId = ""

		createdValidIngredientMeasurementUnit, err := adminClient.CreateValidIngredientMeasurementUnit(ctx, &mealplanningsvc.CreateValidIngredientMeasurementUnitRequest{Input: exampleValidIngredientMeasurementUnitInput})
		require.Error(t, err)
		require.Nil(t, createdValidIngredientMeasurementUnit)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.CreateValidIngredientMeasurementUnit(ctx, &mealplanningsvc.CreateValidIngredientMeasurementUnitRequest{})
		assert.Error(t, err)
	})
}

func TestValidIngredientMeasurementUnits_Listing(T *testing.T) {
	T.Parallel()

	createdValidIngredientMeasurementUnits := []*types.ValidIngredientMeasurementUnit{}
	validIngredient, validMeasurementUnit, created := createValidIngredientMeasurementUnitForTest(T)
	createdValidIngredientMeasurementUnits = append(createdValidIngredientMeasurementUnits, created)
	for range exampleQuantity - 1 {
		exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()
		exampleValidIngredientMeasurementUnitInput := mealplanningconverters.ConvertCreateValidIngredientMeasurementUnitRequestToGRPCValidIngredientMeasurementUnitCreationRequestInput(converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitCreationRequestInput(exampleValidIngredientMeasurementUnit))
		exampleValidIngredientMeasurementUnitInput.ValidMeasurementUnitId = validMeasurementUnit.ID
		exampleValidIngredientMeasurementUnitInput.ValidIngredientId = validIngredient.ID

		createdValidIngredientMeasurementUnit, err := adminClient.CreateValidIngredientMeasurementUnit(T.Context(), &mealplanningsvc.CreateValidIngredientMeasurementUnitRequest{Input: exampleValidIngredientMeasurementUnitInput})
		require.NoError(T, err)
		require.NotNil(T, createdValidIngredientMeasurementUnit)

		createdValidIngredientMeasurementUnits = append(createdValidIngredientMeasurementUnits, mealplanningconverters.ConvertGRPCValidIngredientMeasurementUnitToValidIngredientMeasurementUnit(createdValidIngredientMeasurementUnit.Result))
	}

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidIngredientMeasurementUnits(ctx, &mealplanningsvc.GetValidIngredientMeasurementUnitsRequest{})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidIngredientMeasurementUnits))
	})

	T.Run("by MeasurementUnit", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidIngredientMeasurementUnitsByMeasurementUnit(ctx, &mealplanningsvc.GetValidIngredientMeasurementUnitsByMeasurementUnitRequest{ValidMeasurementUnitId: validMeasurementUnit.ID})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidIngredientMeasurementUnits))
	})

	T.Run("by preparation", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidIngredientMeasurementUnitsByIngredient(ctx, &mealplanningsvc.GetValidIngredientMeasurementUnitsByIngredientRequest{ValidIngredientId: validIngredient.ID})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidIngredientMeasurementUnits))
	})
}
