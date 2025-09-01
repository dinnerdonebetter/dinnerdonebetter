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

func createValidMeasurementUnitForTest(t *testing.T) *mealplanning.ValidMeasurementUnit {
	t.Helper()

	ctx := t.Context()

	creationRequestInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInput()
	convertedInput := grpcconverters.ConvertValidMeasurementUnitCreationRequestInputToGRPCValidMeasurementUnitCreationRequestInput(creationRequestInput)

	created, err := adminClient.CreateValidMeasurementUnit(ctx, &mealplanningsvc.CreateValidMeasurementUnitRequest{
		Input: convertedInput,
	})
	require.NoError(t, err)
	assert.NotNil(t, created)

	return grpcconverters.ConvertGRPCValidMeasurementUnitToValidMeasurementUnit(created.Result)
}

func TestValidMeasurementUnits_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		createValidMeasurementUnitForTest(t)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidMeasurementUnitCreationRequestInputToGRPCValidMeasurementUnitCreationRequestInput(creationRequestInput)

		c := buildUnauthenticatedGRPCClientForTest(t)
		created, err := c.CreateValidMeasurementUnit(ctx, &mealplanningsvc.CreateValidMeasurementUnitRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidMeasurementUnitCreationRequestInputToGRPCValidMeasurementUnitCreationRequestInput(creationRequestInput)
		// this is not allowed
		convertedInput.Name = ""

		created, err := adminClient.CreateValidMeasurementUnit(ctx, &mealplanningsvc.CreateValidMeasurementUnitRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("non-admin users are forbidden from creating", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)

		creationRequestInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidMeasurementUnitCreationRequestInputToGRPCValidMeasurementUnitCreationRequestInput(creationRequestInput)

		created, err := testClient.CreateValidMeasurementUnit(ctx, &mealplanningsvc.CreateValidMeasurementUnitRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})
}

func TestValidMeasurementUnits_Reading(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidMeasurementUnitForTest(t)

		retrieved, err := testClient.GetValidMeasurementUnit(ctx, &mealplanningsvc.GetValidMeasurementUnitRequest{ValidMeasurementUnitID: created.ID})
		assert.NoError(t, err)

		converted := grpcconverters.ConvertGRPCValidMeasurementUnitToValidMeasurementUnit(retrieved.Result)

		assertRoughEquality(t, created, converted, "CreatedAt", "LastUpdatedAt", "ArchivedAt")
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidMeasurementUnitForTest(t)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.GetValidMeasurementUnit(ctx, &mealplanningsvc.GetValidMeasurementUnitRequest{ValidMeasurementUnitID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.GetValidMeasurementUnit(ctx, &mealplanningsvc.GetValidMeasurementUnitRequest{ValidMeasurementUnitID: nonexistentID})
		assert.Error(t, err)
	})
}

func TestValidMeasurementUnits_Updating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidMeasurementUnitForTest(t)

		updateInput := fakes.BuildFakeValidMeasurementUnitUpdateRequestInput()
		created.Update(updateInput)

		response, err := adminClient.UpdateValidMeasurementUnit(ctx, &mealplanningsvc.UpdateValidMeasurementUnitRequest{
			ValidMeasurementUnitID: created.ID,
			Input:                  grpcconverters.ConvertValidMeasurementUnitUpdateRequestInputToGRPCValidMeasurementUnitUpdateRequestInput(updateInput),
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

		created := createValidMeasurementUnitForTest(t)

		updateInput := fakes.BuildFakeValidMeasurementUnitUpdateRequestInput()
		created.Update(updateInput)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.UpdateValidMeasurementUnit(ctx, &mealplanningsvc.UpdateValidMeasurementUnitRequest{
			ValidMeasurementUnitID: created.ID,
			Input:                  grpcconverters.ConvertValidMeasurementUnitUpdateRequestInputToGRPCValidMeasurementUnitUpdateRequestInput(updateInput),
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

		created := createValidMeasurementUnitForTest(t)

		response, err := testClient.UpdateValidMeasurementUnit(ctx, &mealplanningsvc.UpdateValidMeasurementUnitRequest{
			ValidMeasurementUnitID: created.ID,
			Input: &mealplanningsvc.ValidMeasurementUnitUpdateRequestInput{
				Name: pointer.To("doesn't matter"),
			},
		})
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestValidMeasurementUnits_Archiving(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidMeasurementUnitForTest(t)

		_, err := adminClient.ArchiveValidMeasurementUnit(ctx, &mealplanningsvc.ArchiveValidMeasurementUnitRequest{ValidMeasurementUnitID: created.ID})
		assert.NoError(t, err)

		x, err := adminClient.GetValidMeasurementUnit(ctx, &mealplanningsvc.GetValidMeasurementUnitRequest{ValidMeasurementUnitID: created.ID})
		assert.Nil(t, x)
		assert.Error(t, err)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidMeasurementUnitForTest(t)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.ArchiveValidMeasurementUnit(ctx, &mealplanningsvc.ArchiveValidMeasurementUnitRequest{ValidMeasurementUnitID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.ArchiveValidMeasurementUnit(ctx, &mealplanningsvc.ArchiveValidMeasurementUnitRequest{ValidMeasurementUnitID: nonexistentID})
		assert.Error(t, err)
	})

	T.Run("non-admin users are forbidden from archiving", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidMeasurementUnitForTest(t)
		_, testClient := createUserAndClientForTest(T)

		_, err := testClient.ArchiveValidMeasurementUnit(ctx, &mealplanningsvc.ArchiveValidMeasurementUnitRequest{ValidMeasurementUnitID: created.ID})
		assert.Error(t, err)
	})
}

/* TODO: I have this functionality, don't I?
func TestValidMeasurementUnits_GetRandom(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		// in case we haven't already
		createValidMeasurementUnitForTest(t)

		response, err := testClient.GetRandomValidMeasurementUnit(ctx, &mealplanningsvc.GetRandomValidMeasurementUnitRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, response)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		response, err := c.GetRandomValidMeasurementUnit(ctx, &mealplanningsvc.GetRandomValidMeasurementUnitRequest{})
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}
*/

func TestValidMeasurementUnits_Listing(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)
	createdValidMeasurementUnits := []*mealplanning.ValidMeasurementUnit{}
	for range exampleQuantity {
		created := createValidMeasurementUnitForTest(T)
		createdValidMeasurementUnits = append(createdValidMeasurementUnits, created)
	}

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		retrieved, err := testClient.GetValidMeasurementUnits(ctx, &mealplanningsvc.GetValidMeasurementUnitsRequest{})
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.True(t, len(retrieved.Results) >= len(createdValidMeasurementUnits))
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.GetValidMeasurementUnits(ctx, &mealplanningsvc.GetValidMeasurementUnitsRequest{})
		assert.Error(t, err)
	})
}

func TestValidMeasurementUnits_Searching(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)
		created := createValidMeasurementUnitForTest(t)

		retrieved, err := testClient.SearchForValidMeasurementUnits(ctx, &mealplanningsvc.SearchForValidMeasurementUnitsRequest{
			Query: created.Name[:2],
		})
		require.NoError(t, err)
		require.NotNil(t, retrieved)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.SearchForValidMeasurementUnits(ctx, &mealplanningsvc.SearchForValidMeasurementUnitsRequest{})
		assert.Error(t, err)
	})
}
