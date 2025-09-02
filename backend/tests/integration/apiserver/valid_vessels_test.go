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

func createValidVesselForTest(t *testing.T) *mealplanning.ValidVessel {
	t.Helper()

	ctx := t.Context()

	creationRequestInput := fakes.BuildFakeValidVesselCreationRequestInput()
	convertedInput := grpcconverters.ConvertValidVesselCreationRequestInputToGRPCValidVesselCreationRequestInput(creationRequestInput)

	measurementUnit := createValidMeasurementUnitForTest(t)
	convertedInput.CapacityUnitID = &measurementUnit.ID

	created, err := adminClient.CreateValidVessel(ctx, &mealplanningsvc.CreateValidVesselRequest{
		Input: convertedInput,
	})
	require.NoError(t, err)
	assert.NotNil(t, created)

	return grpcconverters.ConvertGRPCValidVesselToValidVessel(created.Result)
}

func TestValidVessels_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		createValidVesselForTest(t)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeValidVesselCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidVesselCreationRequestInputToGRPCValidVesselCreationRequestInput(creationRequestInput)

		c := buildUnauthenticatedGRPCClientForTest(t)
		created, err := c.CreateValidVessel(ctx, &mealplanningsvc.CreateValidVesselRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeValidVesselCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidVesselCreationRequestInputToGRPCValidVesselCreationRequestInput(creationRequestInput)
		// this is not allowed
		convertedInput.Name = ""

		created, err := adminClient.CreateValidVessel(ctx, &mealplanningsvc.CreateValidVesselRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("non-admin users are forbidden from creating", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)

		creationRequestInput := fakes.BuildFakeValidVesselCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidVesselCreationRequestInputToGRPCValidVesselCreationRequestInput(creationRequestInput)

		created, err := testClient.CreateValidVessel(ctx, &mealplanningsvc.CreateValidVesselRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})
}

func TestValidVessels_Reading(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidVesselForTest(t)

		retrieved, err := testClient.GetValidVessel(ctx, &mealplanningsvc.GetValidVesselRequest{ValidVesselID: created.ID})
		assert.NoError(t, err)

		converted := grpcconverters.ConvertGRPCValidVesselToValidVessel(retrieved.Result)

		assertRoughEquality(t, created, converted, append(defaultIgnoredFields(), "CapacityUnit")...)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidVesselForTest(t)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.GetValidVessel(ctx, &mealplanningsvc.GetValidVesselRequest{ValidVesselID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.GetValidVessel(ctx, &mealplanningsvc.GetValidVesselRequest{ValidVesselID: nonexistentID})
		assert.Error(t, err)
	})
}

func TestValidVessels_Updating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidVesselForTest(t)

		updateInput := fakes.BuildFakeValidVesselUpdateRequestInput()
		updateInput.CapacityUnitID = nil
		created.Update(updateInput)

		response, err := adminClient.UpdateValidVessel(ctx, &mealplanningsvc.UpdateValidVesselRequest{
			ValidVesselID: created.ID,
			Input:         grpcconverters.ConvertValidVesselUpdateRequestInputToGRPCValidVesselUpdateRequestInput(updateInput),
		})
		assert.NoError(t, err)

		updated := grpcconverters.ConvertGRPCValidVesselToValidVessel(response.Result)
		// Ensure UpdatedAt was set
		require.NotNil(t, updated.LastUpdatedAt)

		assertRoughEquality(t, created, updated, append(defaultIgnoredFields(), "CapacityUnit")...)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidVesselForTest(t)

		updateInput := fakes.BuildFakeValidVesselUpdateRequestInput()
		created.Update(updateInput)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.UpdateValidVessel(ctx, &mealplanningsvc.UpdateValidVesselRequest{
			ValidVesselID: created.ID,
			Input:         grpcconverters.ConvertValidVesselUpdateRequestInputToGRPCValidVesselUpdateRequestInput(updateInput),
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

		created := createValidVesselForTest(t)

		response, err := testClient.UpdateValidVessel(ctx, &mealplanningsvc.UpdateValidVesselRequest{
			ValidVesselID: created.ID,
			Input: &mealplanningsvc.ValidVesselUpdateRequestInput{
				Name: pointer.To("doesn't matter"),
			},
		})
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestValidVessels_Archiving(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidVesselForTest(t)

		_, err := adminClient.ArchiveValidVessel(ctx, &mealplanningsvc.ArchiveValidVesselRequest{ValidVesselID: created.ID})
		assert.NoError(t, err)

		x, err := adminClient.GetValidVessel(ctx, &mealplanningsvc.GetValidVesselRequest{ValidVesselID: created.ID})
		assert.Nil(t, x)
		assert.Error(t, err)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidVesselForTest(t)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.ArchiveValidVessel(ctx, &mealplanningsvc.ArchiveValidVesselRequest{ValidVesselID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.ArchiveValidVessel(ctx, &mealplanningsvc.ArchiveValidVesselRequest{ValidVesselID: nonexistentID})
		assert.Error(t, err)
	})

	T.Run("non-admin users are forbidden from archiving", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidVesselForTest(t)
		_, testClient := createUserAndClientForTest(T)

		_, err := testClient.ArchiveValidVessel(ctx, &mealplanningsvc.ArchiveValidVesselRequest{ValidVesselID: created.ID})
		assert.Error(t, err)
	})
}

func TestValidVessels_GetRandom(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		// in case we haven't already
		createValidVesselForTest(t)

		response, err := testClient.GetRandomValidVessel(ctx, &mealplanningsvc.GetRandomValidVesselRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, response)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		response, err := c.GetRandomValidVessel(ctx, &mealplanningsvc.GetRandomValidVesselRequest{})
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestValidVessels_Listing(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)
	createdValidVessels := []*mealplanning.ValidVessel{}
	for range exampleQuantity {
		created := createValidVesselForTest(T)
		createdValidVessels = append(createdValidVessels, created)
	}

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		retrieved, err := testClient.GetValidVessels(ctx, &mealplanningsvc.GetValidVesselsRequest{})
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.True(t, len(retrieved.Results) >= len(createdValidVessels))
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.GetValidVessels(ctx, &mealplanningsvc.GetValidVesselsRequest{})
		assert.Error(t, err)
	})
}

func TestValidVessels_Searching(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)
		created := createValidVesselForTest(t)

		retrieved, err := testClient.SearchForValidVessels(ctx, &mealplanningsvc.SearchForValidVesselsRequest{
			Query: created.Name[:2],
		})
		require.NoError(t, err)
		require.NotNil(t, retrieved)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.SearchForValidVessels(ctx, &mealplanningsvc.SearchForValidVesselsRequest{})
		assert.Error(t, err)
	})
}
