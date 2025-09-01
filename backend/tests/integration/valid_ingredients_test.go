package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createValidIngredientForTest(t *testing.T) *mealplanning.ValidIngredient {
	t.Helper()

	ctx := t.Context()

	creationRequestInput := fakes.BuildFakeValidIngredientCreationRequestInput()
	convertedInput := grpcconverters.ConvertValidIngredientCreationRequestInputToGRPCValidIngredientCreationRequestInput(creationRequestInput)

	created, err := adminClient.CreateValidIngredient(ctx, &mealplanningsvc.CreateValidIngredientRequest{
		Input: convertedInput,
	})
	require.NoError(t, err)
	assert.NotNil(t, created)

	return grpcconverters.ConvertGRPCValidIngredientToValidIngredient(created.Result)
}

func TestValidIngredients_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		createValidIngredientForTest(t)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeValidIngredientCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidIngredientCreationRequestInputToGRPCValidIngredientCreationRequestInput(creationRequestInput)

		c := buildUnauthenticatedGRPCClientForTest(t)
		created, err := c.CreateValidIngredient(ctx, &mealplanningsvc.CreateValidIngredientRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeValidIngredientCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidIngredientCreationRequestInputToGRPCValidIngredientCreationRequestInput(creationRequestInput)
		// this is not allowed
		convertedInput.Name = ""

		created, err := adminClient.CreateValidIngredient(ctx, &mealplanningsvc.CreateValidIngredientRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("non-admin users are forbidden from creating", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)

		creationRequestInput := fakes.BuildFakeValidIngredientCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidIngredientCreationRequestInputToGRPCValidIngredientCreationRequestInput(creationRequestInput)

		created, err := testClient.CreateValidIngredient(ctx, &mealplanningsvc.CreateValidIngredientRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})
}

func TestValidIngredients_Reading(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientForTest(t)

		retrieved, err := testClient.GetValidIngredient(ctx, &mealplanningsvc.GetValidIngredientRequest{ValidIngredientID: created.ID})
		assert.NoError(t, err)

		converted := grpcconverters.ConvertGRPCValidIngredientToValidIngredient(retrieved.Result)

		assertRoughEquality(t, created, converted, "CreatedAt", "LastUpdatedAt", "ArchivedAt")
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientForTest(t)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.GetValidIngredient(ctx, &mealplanningsvc.GetValidIngredientRequest{ValidIngredientID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.GetValidIngredient(ctx, &mealplanningsvc.GetValidIngredientRequest{ValidIngredientID: nonexistentID})
		assert.Error(t, err)
	})
}

func TestValidIngredients_Updating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientForTest(t)

		updateInput := fakes.BuildFakeValidIngredientUpdateRequestInput()
		created.Update(updateInput)

		response, err := adminClient.UpdateValidIngredient(ctx, &mealplanningsvc.UpdateValidIngredientRequest{
			ValidIngredientID: created.ID,
			Input:             grpcconverters.ConvertValidIngredientUpdateRequestInputToGRPCValidIngredientUpdateRequestInput(updateInput),
		})
		assert.NoError(t, err)

		updated := response.Result
		// Ensure UpdatedAt was set
		require.NotNil(t, updated.LastUpdatedAt)

		assertRoughEquality(t, created, updated, "CreatedAt", "LastUpdatedAt", "ArchivedAt")
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientForTest(t)

		updateInput := fakes.BuildFakeValidIngredientUpdateRequestInput()
		created.Update(updateInput)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.UpdateValidIngredient(ctx, &mealplanningsvc.UpdateValidIngredientRequest{
			ValidIngredientID: created.ID,
			Input:             grpcconverters.ConvertValidIngredientUpdateRequestInputToGRPCValidIngredientUpdateRequestInput(updateInput),
		})
		assert.Error(t, err)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()

		/*
			there's no way to provide invalid input to this method, but
			I want to make it explicit that tests should be written the moment that changes
		*/
	})

	T.Run("non-admin users are forbidden from updating", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)

		created := createValidIngredientForTest(t)

		response, err := testClient.UpdateValidIngredient(ctx, &mealplanningsvc.UpdateValidIngredientRequest{
			ValidIngredientID: created.ID,
			Input: &mealplanningsvc.ValidIngredientUpdateRequestInput{
				Name: pointer.To("doesn't matter"),
			},
		})
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestValidIngredients_Archiving(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientForTest(t)

		_, err := adminClient.ArchiveValidIngredient(ctx, &mealplanningsvc.ArchiveValidIngredientRequest{ValidIngredientID: created.ID})
		assert.NoError(t, err)

		x, err := adminClient.GetValidIngredient(ctx, &mealplanningsvc.GetValidIngredientRequest{ValidIngredientID: created.ID})
		assert.Nil(t, x)
		assert.Error(t, err)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientForTest(t)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.ArchiveValidIngredient(ctx, &mealplanningsvc.ArchiveValidIngredientRequest{ValidIngredientID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.ArchiveValidIngredient(ctx, &mealplanningsvc.ArchiveValidIngredientRequest{ValidIngredientID: nonexistentID})
		assert.Error(t, err)
	})

	T.Run("non-admin users are forbidden from archiving", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientForTest(t)
		_, testClient := createUserAndClientForTest(T)

		_, err := testClient.ArchiveValidIngredient(ctx, &mealplanningsvc.ArchiveValidIngredientRequest{ValidIngredientID: created.ID})
		assert.Error(t, err)
	})
}

func TestValidIngredients_GetRandom(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		// in case we haven't already
		createValidIngredientForTest(t)

		response, err := testClient.GetRandomValidIngredient(ctx, &mealplanningsvc.GetRandomValidIngredientRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, response)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		response, err := c.GetRandomValidIngredient(ctx, &mealplanningsvc.GetRandomValidIngredientRequest{})
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestValidIngredients_Listing(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)
	createdValidIngredients := []*mealplanning.ValidIngredient{}
	for range exampleQuantity {
		created := createValidIngredientForTest(T)
		createdValidIngredients = append(createdValidIngredients, created)
	}

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		retrieved, err := testClient.GetValidIngredients(ctx, &mealplanningsvc.GetValidIngredientsRequest{})
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.True(t, len(retrieved.Results) >= len(createdValidIngredients))
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.GetValidIngredients(ctx, &mealplanningsvc.GetValidIngredientsRequest{})
		assert.Error(t, err)
	})
}

func TestValidIngredients_Searching(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)
		created := createValidIngredientForTest(t)

		retrieved, err := testClient.SearchForValidIngredients(ctx, &mealplanningsvc.SearchForValidIngredientsRequest{
			Query: created.Name[:2],
		})
		require.NoError(t, err)
		require.NotNil(t, retrieved)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.SearchForValidIngredients(ctx, &mealplanningsvc.SearchForValidIngredientsRequest{})
		assert.Error(t, err)
	})
}
