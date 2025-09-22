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

func createValidPreparationVesselForTest(t *testing.T) (*types.ValidPreparation, *types.ValidVessel, *types.ValidPreparationVessel) {
	t.Helper()
	ctx := t.Context()

	createdValidPreparation := createValidPreparationForTest(t)
	createdValidVessel := createValidVesselForTest(t)

	exampleValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()
	exampleValidPreparationVessel.Vessel = *createdValidVessel
	exampleValidPreparationVessel.Preparation = *createdValidPreparation

	exampleValidPreparationVesselInput := mealplanningconverters.ConvertCreateValidPreparationVesselRequestToGRPCValidPreparationVesselCreationRequestInput(converters.ConvertValidPreparationVesselToValidPreparationVesselCreationRequestInput(exampleValidPreparationVessel))
	createdValidPreparationVessel, err := adminClient.CreateValidPreparationVessel(ctx, &mealplanningsvc.CreateValidPreparationVesselRequest{Input: exampleValidPreparationVesselInput})
	require.NoError(t, err)
	require.NotNil(t, createdValidPreparationVessel)

	validPrepVesselRes, err := adminClient.GetValidPreparationVessel(ctx, &mealplanningsvc.GetValidPreparationVesselRequest{
		ValidPreparationVesselID: createdValidPreparationVessel.Result.ID,
	})
	require.NoError(t, err)
	require.NotNil(t, validPrepVesselRes.Result)

	return createdValidPreparation, createdValidVessel, mealplanningconverters.ConvertGRPCValidPreparationVesselToValidPreparationVessel(validPrepVesselRes.Result)
}

func TestValidPreparationVessels_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		createValidPreparationVesselForTest(t)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		exampleValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()
		exampleValidPreparationVesselInput := mealplanningconverters.ConvertCreateValidPreparationVesselRequestToGRPCValidPreparationVesselCreationRequestInput(converters.ConvertValidPreparationVesselToValidPreparationVesselCreationRequestInput(exampleValidPreparationVessel))
		exampleValidPreparationVesselInput.ValidVesselID = ""
		exampleValidPreparationVesselInput.ValidPreparationID = ""

		createdValidPreparationVessel, err := adminClient.CreateValidPreparationVessel(ctx, &mealplanningsvc.CreateValidPreparationVesselRequest{Input: exampleValidPreparationVesselInput})
		require.Error(t, err)
		require.Nil(t, createdValidPreparationVessel)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.CreateValidPreparationVessel(ctx, &mealplanningsvc.CreateValidPreparationVesselRequest{})
		assert.Error(t, err)
	})
}

func TestValidPreparationVessels_Listing(T *testing.T) {
	T.Parallel()

	createdValidPreparationVessels := []*types.ValidPreparationVessel{}
	validPreparation, validVessel, created := createValidPreparationVesselForTest(T)
	createdValidPreparationVessels = append(createdValidPreparationVessels, created)
	for range exampleQuantity - 1 {
		exampleValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()
		exampleValidPreparationVesselInput := mealplanningconverters.ConvertCreateValidPreparationVesselRequestToGRPCValidPreparationVesselCreationRequestInput(converters.ConvertValidPreparationVesselToValidPreparationVesselCreationRequestInput(exampleValidPreparationVessel))
		exampleValidPreparationVesselInput.ValidVesselID = validVessel.ID
		exampleValidPreparationVesselInput.ValidPreparationID = validPreparation.ID

		createdValidPreparationVessel, err := adminClient.CreateValidPreparationVessel(T.Context(), &mealplanningsvc.CreateValidPreparationVesselRequest{Input: exampleValidPreparationVesselInput})
		require.NoError(T, err)
		require.NotNil(T, createdValidPreparationVessel)

		createdValidPreparationVessels = append(createdValidPreparationVessels, mealplanningconverters.ConvertGRPCValidPreparationVesselToValidPreparationVessel(createdValidPreparationVessel.Result))
	}

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidPreparationVessels(ctx, &mealplanningsvc.GetValidPreparationVesselsRequest{})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidPreparationVessels))
	})

	T.Run("by vessel", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidPreparationVesselsByVessel(ctx, &mealplanningsvc.GetValidPreparationVesselsByVesselRequest{ValidVesselID: validVessel.ID})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidPreparationVessels))
	})

	T.Run("by preparation", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidPreparationVesselsByPreparation(ctx, &mealplanningsvc.GetValidPreparationVesselsByPreparationRequest{ValidPreparationID: validPreparation.ID})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidPreparationVessels))
	})
}
