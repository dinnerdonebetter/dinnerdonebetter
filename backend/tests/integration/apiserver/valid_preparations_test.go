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

func createValidPreparationForTest(t *testing.T) *mealplanning.ValidPreparation {
	t.Helper()

	ctx := t.Context()

	creationRequestInput := fakes.BuildFakeValidPreparationCreationRequestInput()
	convertedInput := grpcconverters.ConvertValidPreparationCreationRequestInputToGRPCValidPreparationCreationRequestInput(creationRequestInput)

	created, err := adminClient.CreateValidPreparation(ctx, &mealplanningsvc.CreateValidPreparationRequest{
		Input: convertedInput,
	})
	require.NoError(t, err)
	assert.NotNil(t, created)

	return grpcconverters.ConvertGRPCValidPreparationToValidPreparation(created.Result)
}

func TestValidPreparations_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		createValidPreparationForTest(t)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeValidPreparationCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidPreparationCreationRequestInputToGRPCValidPreparationCreationRequestInput(creationRequestInput)

		c := buildUnauthenticatedGRPCClientForTest(t)
		created, err := c.CreateValidPreparation(ctx, &mealplanningsvc.CreateValidPreparationRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeValidPreparationCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidPreparationCreationRequestInputToGRPCValidPreparationCreationRequestInput(creationRequestInput)
		// this is not allowed
		convertedInput.Name = ""

		created, err := adminClient.CreateValidPreparation(ctx, &mealplanningsvc.CreateValidPreparationRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("non-admin users are forbidden from creating", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)

		creationRequestInput := fakes.BuildFakeValidPreparationCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidPreparationCreationRequestInputToGRPCValidPreparationCreationRequestInput(creationRequestInput)

		created, err := testClient.CreateValidPreparation(ctx, &mealplanningsvc.CreateValidPreparationRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})
}

func TestValidPreparations_Reading(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidPreparationForTest(t)

		retrieved, err := testClient.GetValidPreparation(ctx, &mealplanningsvc.GetValidPreparationRequest{ValidPreparationID: created.ID})
		assert.NoError(t, err)

		converted := grpcconverters.ConvertGRPCValidPreparationToValidPreparation(retrieved.Result)

		assertRoughEquality(t, created, converted, defaultIgnoredFields()...)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidPreparationForTest(t)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.GetValidPreparation(ctx, &mealplanningsvc.GetValidPreparationRequest{ValidPreparationID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.GetValidPreparation(ctx, &mealplanningsvc.GetValidPreparationRequest{ValidPreparationID: nonexistentID})
		assert.Error(t, err)
	})
}

func TestValidPreparations_Updating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidPreparationForTest(t)

		updateInput := fakes.BuildFakeValidPreparationUpdateRequestInput()
		created.Update(updateInput)

		response, err := adminClient.UpdateValidPreparation(ctx, &mealplanningsvc.UpdateValidPreparationRequest{
			ValidPreparationID: created.ID,
			Input:              grpcconverters.ConvertValidPreparationUpdateRequestInputToGRPCValidPreparationUpdateRequestInput(updateInput),
		})
		assert.NoError(t, err)

		updated := grpcconverters.ConvertGRPCValidPreparationToValidPreparation(response.Result)
		// Ensure UpdatedAt was set
		require.NotNil(t, updated.LastUpdatedAt)

		assertRoughEquality(t, created, updated, defaultIgnoredFields()...)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidPreparationForTest(t)

		updateInput := fakes.BuildFakeValidPreparationUpdateRequestInput()
		created.Update(updateInput)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.UpdateValidPreparation(ctx, &mealplanningsvc.UpdateValidPreparationRequest{
			ValidPreparationID: created.ID,
			Input:              grpcconverters.ConvertValidPreparationUpdateRequestInputToGRPCValidPreparationUpdateRequestInput(updateInput),
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

		created := createValidPreparationForTest(t)

		response, err := testClient.UpdateValidPreparation(ctx, &mealplanningsvc.UpdateValidPreparationRequest{
			ValidPreparationID: created.ID,
			Input: &mealplanningsvc.ValidPreparationUpdateRequestInput{
				Name: pointer.To("doesn't matter"),
			},
		})
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestValidPreparations_Archiving(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidPreparationForTest(t)

		_, err := adminClient.ArchiveValidPreparation(ctx, &mealplanningsvc.ArchiveValidPreparationRequest{ValidPreparationID: created.ID})
		assert.NoError(t, err)

		x, err := adminClient.GetValidPreparation(ctx, &mealplanningsvc.GetValidPreparationRequest{ValidPreparationID: created.ID})
		assert.Nil(t, x)
		assert.Error(t, err)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidPreparationForTest(t)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.ArchiveValidPreparation(ctx, &mealplanningsvc.ArchiveValidPreparationRequest{ValidPreparationID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.ArchiveValidPreparation(ctx, &mealplanningsvc.ArchiveValidPreparationRequest{ValidPreparationID: nonexistentID})
		assert.Error(t, err)
	})

	T.Run("non-admin users are forbidden from archiving", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidPreparationForTest(t)
		_, testClient := createUserAndClientForTest(T)

		_, err := testClient.ArchiveValidPreparation(ctx, &mealplanningsvc.ArchiveValidPreparationRequest{ValidPreparationID: created.ID})
		assert.Error(t, err)
	})
}

func TestValidPreparations_GetRandom(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		// in case we haven't already
		createValidPreparationForTest(t)

		response, err := testClient.GetRandomValidPreparation(ctx, &mealplanningsvc.GetRandomValidPreparationRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, response)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		response, err := c.GetRandomValidPreparation(ctx, &mealplanningsvc.GetRandomValidPreparationRequest{})
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestValidPreparations_Listing(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)
	createdValidPreparations := []*mealplanning.ValidPreparation{}
	for range exampleQuantity {
		created := createValidPreparationForTest(T)
		createdValidPreparations = append(createdValidPreparations, created)
	}

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		retrieved, err := testClient.GetValidPreparations(ctx, &mealplanningsvc.GetValidPreparationsRequest{})
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.True(t, len(retrieved.Results) >= len(createdValidPreparations))
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.GetValidPreparations(ctx, &mealplanningsvc.GetValidPreparationsRequest{})
		assert.Error(t, err)
	})
}

func TestValidPreparations_Searching(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)
		created := createValidPreparationForTest(t)

		retrieved, err := testClient.SearchForValidPreparations(ctx, &mealplanningsvc.SearchForValidPreparationsRequest{
			Query: created.Name[:2],
		})
		require.NoError(t, err)
		require.NotNil(t, retrieved)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.SearchForValidPreparations(ctx, &mealplanningsvc.SearchForValidPreparationsRequest{})
		assert.Error(t, err)
	})
}
