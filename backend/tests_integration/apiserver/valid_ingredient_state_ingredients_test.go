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

func createValidIngredientStateIngredientForTest(t *testing.T) (*types.ValidIngredientState, *types.ValidIngredient, *types.ValidIngredientStateIngredient) {
	t.Helper()
	ctx := t.Context()

	createdValidIngredientState := createValidIngredientStateForTest(t)
	createdValidIngredient := createValidIngredientForTest(t)

	exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()
	exampleValidIngredientStateIngredient.Ingredient = *createdValidIngredient
	exampleValidIngredientStateIngredient.IngredientState = *createdValidIngredientState

	exampleValidIngredientStateIngredientInput := mealplanningconverters.ConvertCreateValidIngredientStateIngredientRequestToGRPCValidIngredientStateIngredientCreationRequestInput(converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientCreationRequestInput(exampleValidIngredientStateIngredient))
	createdValidIngredientStateIngredient, err := adminClient.CreateValidIngredientStateIngredient(ctx, &mealplanningsvc.CreateValidIngredientStateIngredientRequest{Input: exampleValidIngredientStateIngredientInput})
	require.NoError(t, err)
	require.NotNil(t, createdValidIngredientStateIngredient)

	validPrepIngredientRes, err := adminClient.GetValidIngredientStateIngredient(ctx, &mealplanningsvc.GetValidIngredientStateIngredientRequest{
		ValidIngredientStateIngredientID: createdValidIngredientStateIngredient.Result.ID,
	})
	require.NoError(t, err)
	require.NotNil(t, validPrepIngredientRes.Result)

	return createdValidIngredientState, createdValidIngredient, mealplanningconverters.ConvertGRPCValidIngredientStateIngredientToValidIngredientStateIngredient(validPrepIngredientRes.Result)
}

func TestValidIngredientStateIngredients_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		createValidIngredientStateIngredientForTest(t)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()
		exampleValidIngredientStateIngredientInput := mealplanningconverters.ConvertCreateValidIngredientStateIngredientRequestToGRPCValidIngredientStateIngredientCreationRequestInput(converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientCreationRequestInput(exampleValidIngredientStateIngredient))
		exampleValidIngredientStateIngredientInput.ValidIngredientID = ""
		exampleValidIngredientStateIngredientInput.ValidIngredientStateID = ""

		createdValidIngredientStateIngredient, err := adminClient.CreateValidIngredientStateIngredient(ctx, &mealplanningsvc.CreateValidIngredientStateIngredientRequest{Input: exampleValidIngredientStateIngredientInput})
		require.Error(t, err)
		require.Nil(t, createdValidIngredientStateIngredient)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.CreateValidIngredientStateIngredient(ctx, &mealplanningsvc.CreateValidIngredientStateIngredientRequest{})
		assert.Error(t, err)
	})
}

func TestValidIngredientStateIngredients_Listing(T *testing.T) {
	T.Parallel()

	createdValidIngredientStateIngredients := []*types.ValidIngredientStateIngredient{}
	validIngredientState, validIngredient, created := createValidIngredientStateIngredientForTest(T)
	createdValidIngredientStateIngredients = append(createdValidIngredientStateIngredients, created)
	for range exampleQuantity - 1 {
		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()
		exampleValidIngredientStateIngredientInput := mealplanningconverters.ConvertCreateValidIngredientStateIngredientRequestToGRPCValidIngredientStateIngredientCreationRequestInput(converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientCreationRequestInput(exampleValidIngredientStateIngredient))
		exampleValidIngredientStateIngredientInput.ValidIngredientID = validIngredient.ID
		exampleValidIngredientStateIngredientInput.ValidIngredientStateID = validIngredientState.ID

		createdValidIngredientStateIngredient, err := adminClient.CreateValidIngredientStateIngredient(T.Context(), &mealplanningsvc.CreateValidIngredientStateIngredientRequest{Input: exampleValidIngredientStateIngredientInput})
		require.NoError(T, err)
		require.NotNil(T, createdValidIngredientStateIngredient)

		createdValidIngredientStateIngredients = append(createdValidIngredientStateIngredients, mealplanningconverters.ConvertGRPCValidIngredientStateIngredientToValidIngredientStateIngredient(createdValidIngredientStateIngredient.Result))
	}

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidIngredientStateIngredients(ctx, &mealplanningsvc.GetValidIngredientStateIngredientsRequest{})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidIngredientStateIngredients))
	})

	T.Run("by ingredient", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidIngredientStateIngredientsByIngredient(ctx, &mealplanningsvc.GetValidIngredientStateIngredientsByIngredientRequest{ValidIngredientID: validIngredient.ID})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidIngredientStateIngredients))
	})

	T.Run("by ingredient state", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidIngredientStateIngredientsByIngredientState(ctx, &mealplanningsvc.GetValidIngredientStateIngredientsByIngredientStateRequest{ValidIngredientStateID: validIngredientState.ID})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidIngredientStateIngredients))
	})
}
