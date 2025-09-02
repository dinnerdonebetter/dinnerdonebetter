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

func createValidIngredientGroupForTest(t *testing.T) *mealplanning.ValidIngredientGroup {
	t.Helper()

	ctx := t.Context()

	creationRequestInput := fakes.BuildFakeValidIngredientGroupCreationRequestInput()
	convertedInput := grpcconverters.ConvertValidIngredientGroupCreationRequestInputToGRPCValidIngredientGroupCreationRequestInput(creationRequestInput)

	created, err := adminClient.CreateValidIngredientGroup(ctx, &mealplanningsvc.CreateValidIngredientGroupRequest{
		Input: convertedInput,
	})
	require.NoError(t, err)
	assert.NotNil(t, created)

	return grpcconverters.ConvertGRPCValidIngredientGroupToValidIngredientGroup(created.Result)
}

func TestValidIngredientGroups_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		createValidIngredientGroupForTest(t)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeValidIngredientGroupCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidIngredientGroupCreationRequestInputToGRPCValidIngredientGroupCreationRequestInput(creationRequestInput)

		c := buildUnauthenticatedGRPCClientForTest(t)
		created, err := c.CreateValidIngredientGroup(ctx, &mealplanningsvc.CreateValidIngredientGroupRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeValidIngredientGroupCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidIngredientGroupCreationRequestInputToGRPCValidIngredientGroupCreationRequestInput(creationRequestInput)
		// this is not allowed
		convertedInput.Name = ""

		created, err := adminClient.CreateValidIngredientGroup(ctx, &mealplanningsvc.CreateValidIngredientGroupRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("non-admin users are forbidden from creating", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)

		creationRequestInput := fakes.BuildFakeValidIngredientGroupCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidIngredientGroupCreationRequestInputToGRPCValidIngredientGroupCreationRequestInput(creationRequestInput)

		created, err := testClient.CreateValidIngredientGroup(ctx, &mealplanningsvc.CreateValidIngredientGroupRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})
}

func TestValidIngredientGroups_Reading(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientGroupForTest(t)

		retrieved, err := testClient.GetValidIngredientGroup(ctx, &mealplanningsvc.GetValidIngredientGroupRequest{ValidIngredientGroupID: created.ID})
		assert.NoError(t, err)

		converted := grpcconverters.ConvertGRPCValidIngredientGroupToValidIngredientGroup(retrieved.Result)

		assertRoughEquality(t, created, converted, defaultIgnoredFields()...)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientGroupForTest(t)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.GetValidIngredientGroup(ctx, &mealplanningsvc.GetValidIngredientGroupRequest{ValidIngredientGroupID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.GetValidIngredientGroup(ctx, &mealplanningsvc.GetValidIngredientGroupRequest{ValidIngredientGroupID: nonexistentID})
		assert.Error(t, err)
	})
}

func TestValidIngredientGroups_Updating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientGroupForTest(t)

		updateInput := fakes.BuildFakeValidIngredientGroupUpdateRequestInput()
		created.Update(updateInput)

		response, err := adminClient.UpdateValidIngredientGroup(ctx, &mealplanningsvc.UpdateValidIngredientGroupRequest{
			ValidIngredientGroupID: created.ID,
			Input:                  grpcconverters.ConvertValidIngredientGroupUpdateRequestInputToGRPCValidIngredientGroupUpdateRequestInput(updateInput),
		})
		assert.NoError(t, err)
		updated := grpcconverters.ConvertGRPCValidIngredientGroupToValidIngredientGroup(response.Result)

		// Ensure UpdatedAt was set
		require.NotNil(t, updated.LastUpdatedAt)

		assertRoughEquality(t, created, updated, defaultIgnoredFields()...)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientGroupForTest(t)

		updateInput := fakes.BuildFakeValidIngredientGroupUpdateRequestInput()
		created.Update(updateInput)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.UpdateValidIngredientGroup(ctx, &mealplanningsvc.UpdateValidIngredientGroupRequest{
			ValidIngredientGroupID: created.ID,
			Input:                  grpcconverters.ConvertValidIngredientGroupUpdateRequestInputToGRPCValidIngredientGroupUpdateRequestInput(updateInput),
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

		created := createValidIngredientGroupForTest(t)

		response, err := testClient.UpdateValidIngredientGroup(ctx, &mealplanningsvc.UpdateValidIngredientGroupRequest{
			ValidIngredientGroupID: created.ID,
			Input: &mealplanningsvc.ValidIngredientGroupUpdateRequestInput{
				Name: pointer.To("doesn't matter"),
			},
		})
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestValidIngredientGroups_Archiving(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientGroupForTest(t)

		_, err := adminClient.ArchiveValidIngredientGroup(ctx, &mealplanningsvc.ArchiveValidIngredientGroupRequest{ValidIngredientGroupID: created.ID})
		assert.NoError(t, err)

		x, err := adminClient.GetValidIngredientGroup(ctx, &mealplanningsvc.GetValidIngredientGroupRequest{ValidIngredientGroupID: created.ID})
		assert.Nil(t, x)
		assert.Error(t, err)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientGroupForTest(t)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.ArchiveValidIngredientGroup(ctx, &mealplanningsvc.ArchiveValidIngredientGroupRequest{ValidIngredientGroupID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.ArchiveValidIngredientGroup(ctx, &mealplanningsvc.ArchiveValidIngredientGroupRequest{ValidIngredientGroupID: nonexistentID})
		assert.Error(t, err)
	})

	T.Run("non-admin users are forbidden from archiving", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidIngredientGroupForTest(t)
		_, testClient := createUserAndClientForTest(T)

		_, err := testClient.ArchiveValidIngredientGroup(ctx, &mealplanningsvc.ArchiveValidIngredientGroupRequest{ValidIngredientGroupID: created.ID})
		assert.Error(t, err)
	})
}

/* // TODO: we have this functionality for other types, why not this one?

func TestValidIngredientGroups_GetRandom(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		// in case we haven't already
		createValidIngredientGroupForTest(t)

		response, err := testClient.GetRandomValidIngredientGroup(ctx, &mealplanningsvc.GetRandomValidIngredientGroupRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, response)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		response, err := c.GetRandomValidIngredientGroup(ctx, &mealplanningsvc.GetRandomValidIngredientGroupRequest{})
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}
*/

func TestValidIngredientGroups_Listing(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)
	createdValidIngredientGroups := []*mealplanning.ValidIngredientGroup{}
	for range exampleQuantity {
		created := createValidIngredientGroupForTest(T)
		createdValidIngredientGroups = append(createdValidIngredientGroups, created)
	}

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		retrieved, err := testClient.GetValidIngredientGroups(ctx, &mealplanningsvc.GetValidIngredientGroupsRequest{})
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.True(t, len(retrieved.Results) >= len(createdValidIngredientGroups))
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.GetValidIngredientGroups(ctx, &mealplanningsvc.GetValidIngredientGroupsRequest{})
		assert.Error(t, err)
	})
}

func TestValidIngredientGroups_Searching(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)
		created := createValidIngredientGroupForTest(t)

		retrieved, err := testClient.SearchForValidIngredientGroups(ctx, &mealplanningsvc.SearchForValidIngredientGroupsRequest{
			Query: created.Name[:2],
		})
		require.NoError(t, err)
		require.NotNil(t, retrieved)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.SearchForValidIngredientGroups(ctx, &mealplanningsvc.SearchForValidIngredientGroupsRequest{})
		assert.Error(t, err)
	})
}
