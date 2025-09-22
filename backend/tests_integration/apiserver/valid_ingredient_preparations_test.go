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

func createValidIngredientPreparationForTest(t *testing.T) (*types.ValidIngredient, *types.ValidPreparation, *types.ValidIngredientPreparation) {
	t.Helper()
	ctx := t.Context()

	createdValidIngredient := createValidIngredientForTest(t)
	createdValidPreparation := createValidPreparationForTest(t)

	exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
	exampleValidIngredientPreparation.Preparation = *createdValidPreparation
	exampleValidIngredientPreparation.Ingredient = *createdValidIngredient

	exampleValidIngredientPreparationInput := mealplanningconverters.ConvertCreateValidIngredientPreparationRequestToGRPCValidIngredientPreparationCreationRequestInput(converters.ConvertValidIngredientPreparationToValidIngredientPreparationCreationRequestInput(exampleValidIngredientPreparation))
	createdValidIngredientPreparation, err := adminClient.CreateValidIngredientPreparation(ctx, &mealplanningsvc.CreateValidIngredientPreparationRequest{Input: exampleValidIngredientPreparationInput})
	require.NoError(t, err)
	require.NotNil(t, createdValidIngredientPreparation)

	validPrepPreparationRes, err := adminClient.GetValidIngredientPreparation(ctx, &mealplanningsvc.GetValidIngredientPreparationRequest{
		ValidIngredientPreparationID: createdValidIngredientPreparation.Result.ID,
	})
	require.NoError(t, err)
	require.NotNil(t, validPrepPreparationRes.Result)

	return createdValidIngredient, createdValidPreparation, mealplanningconverters.ConvertGRPCValidIngredientPreparationToValidIngredientPreparation(validPrepPreparationRes.Result)
}

func TestValidIngredientPreparations_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		createValidIngredientPreparationForTest(t)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparationInput := mealplanningconverters.ConvertCreateValidIngredientPreparationRequestToGRPCValidIngredientPreparationCreationRequestInput(converters.ConvertValidIngredientPreparationToValidIngredientPreparationCreationRequestInput(exampleValidIngredientPreparation))
		exampleValidIngredientPreparationInput.ValidPreparationID = ""
		exampleValidIngredientPreparationInput.ValidIngredientID = ""

		createdValidIngredientPreparation, err := adminClient.CreateValidIngredientPreparation(ctx, &mealplanningsvc.CreateValidIngredientPreparationRequest{Input: exampleValidIngredientPreparationInput})
		require.Error(t, err)
		require.Nil(t, createdValidIngredientPreparation)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.CreateValidIngredientPreparation(ctx, &mealplanningsvc.CreateValidIngredientPreparationRequest{})
		assert.Error(t, err)
	})
}

func TestValidIngredientPreparations_Listing(T *testing.T) {
	T.Parallel()

	createdValidIngredientPreparations := []*types.ValidIngredientPreparation{}
	validIngredient, validPreparation, created := createValidIngredientPreparationForTest(T)
	createdValidIngredientPreparations = append(createdValidIngredientPreparations, created)
	for range exampleQuantity - 1 {
		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparationInput := mealplanningconverters.ConvertCreateValidIngredientPreparationRequestToGRPCValidIngredientPreparationCreationRequestInput(converters.ConvertValidIngredientPreparationToValidIngredientPreparationCreationRequestInput(exampleValidIngredientPreparation))
		exampleValidIngredientPreparationInput.ValidPreparationID = validPreparation.ID
		exampleValidIngredientPreparationInput.ValidIngredientID = validIngredient.ID

		createdValidIngredientPreparation, err := adminClient.CreateValidIngredientPreparation(T.Context(), &mealplanningsvc.CreateValidIngredientPreparationRequest{Input: exampleValidIngredientPreparationInput})
		require.NoError(T, err)
		require.NotNil(T, createdValidIngredientPreparation)

		createdValidIngredientPreparations = append(createdValidIngredientPreparations, mealplanningconverters.ConvertGRPCValidIngredientPreparationToValidIngredientPreparation(createdValidIngredientPreparation.Result))
	}

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidIngredientPreparations(ctx, &mealplanningsvc.GetValidIngredientPreparationsRequest{})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidIngredientPreparations))
	})

	T.Run("by Preparation", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidIngredientPreparationsByPreparation(ctx, &mealplanningsvc.GetValidIngredientPreparationsByPreparationRequest{ValidPreparationID: validPreparation.ID})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidIngredientPreparations))
	})

	T.Run("by preparation", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidIngredientPreparationsByIngredient(ctx, &mealplanningsvc.GetValidIngredientPreparationsByIngredientRequest{ValidIngredientID: validIngredient.ID})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidIngredientPreparations))
	})
}
