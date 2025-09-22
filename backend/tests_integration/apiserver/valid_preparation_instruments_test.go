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

func createValidPreparationInstrumentForTest(t *testing.T) (*types.ValidPreparation, *types.ValidInstrument, *types.ValidPreparationInstrument) {
	t.Helper()
	ctx := t.Context()

	createdValidPreparation := createValidPreparationForTest(t)
	createdValidInstrument := createValidInstrumentForTest(t)

	exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
	exampleValidPreparationInstrument.Instrument = *createdValidInstrument
	exampleValidPreparationInstrument.Preparation = *createdValidPreparation

	exampleValidPreparationInstrumentInput := mealplanningconverters.ConvertCreateValidPreparationInstrumentRequestToGRPCValidPreparationInstrumentCreationRequestInput(converters.ConvertValidPreparationInstrumentToValidPreparationInstrumentCreationRequestInput(exampleValidPreparationInstrument))
	createdValidPreparationInstrument, err := adminClient.CreateValidPreparationInstrument(ctx, &mealplanningsvc.CreateValidPreparationInstrumentRequest{Input: exampleValidPreparationInstrumentInput})
	require.NoError(t, err)
	require.NotNil(t, createdValidPreparationInstrument)

	validPrepInstrumentRes, err := adminClient.GetValidPreparationInstrument(ctx, &mealplanningsvc.GetValidPreparationInstrumentRequest{
		ValidPreparationInstrumentID: createdValidPreparationInstrument.Result.ID,
	})
	require.NoError(t, err)
	require.NotNil(t, validPrepInstrumentRes.Result)

	return createdValidPreparation, createdValidInstrument, mealplanningconverters.ConvertGRPCValidPreparationInstrumentToValidPreparationInstrument(validPrepInstrumentRes.Result)
}

func TestValidPreparationInstruments_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		createValidPreparationInstrumentForTest(t)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		exampleValidPreparationInstrumentInput := mealplanningconverters.ConvertCreateValidPreparationInstrumentRequestToGRPCValidPreparationInstrumentCreationRequestInput(converters.ConvertValidPreparationInstrumentToValidPreparationInstrumentCreationRequestInput(exampleValidPreparationInstrument))
		exampleValidPreparationInstrumentInput.ValidInstrumentID = ""
		exampleValidPreparationInstrumentInput.ValidPreparationID = ""

		createdValidPreparationInstrument, err := adminClient.CreateValidPreparationInstrument(ctx, &mealplanningsvc.CreateValidPreparationInstrumentRequest{Input: exampleValidPreparationInstrumentInput})
		require.Error(t, err)
		require.Nil(t, createdValidPreparationInstrument)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.CreateValidPreparationInstrument(ctx, &mealplanningsvc.CreateValidPreparationInstrumentRequest{})
		assert.Error(t, err)
	})
}

func TestValidPreparationInstruments_Listing(T *testing.T) {
	T.Parallel()

	createdValidPreparationInstruments := []*types.ValidPreparationInstrument{}
	validPreparation, validInstrument, created := createValidPreparationInstrumentForTest(T)
	createdValidPreparationInstruments = append(createdValidPreparationInstruments, created)
	for range exampleQuantity - 1 {
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		exampleValidPreparationInstrumentInput := mealplanningconverters.ConvertCreateValidPreparationInstrumentRequestToGRPCValidPreparationInstrumentCreationRequestInput(converters.ConvertValidPreparationInstrumentToValidPreparationInstrumentCreationRequestInput(exampleValidPreparationInstrument))
		exampleValidPreparationInstrumentInput.ValidInstrumentID = validInstrument.ID
		exampleValidPreparationInstrumentInput.ValidPreparationID = validPreparation.ID

		createdValidPreparationInstrument, err := adminClient.CreateValidPreparationInstrument(T.Context(), &mealplanningsvc.CreateValidPreparationInstrumentRequest{Input: exampleValidPreparationInstrumentInput})
		require.NoError(T, err)
		require.NotNil(T, createdValidPreparationInstrument)

		createdValidPreparationInstruments = append(createdValidPreparationInstruments, mealplanningconverters.ConvertGRPCValidPreparationInstrumentToValidPreparationInstrument(createdValidPreparationInstrument.Result))
	}

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidPreparationInstruments(ctx, &mealplanningsvc.GetValidPreparationInstrumentsRequest{})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidPreparationInstruments))
	})

	T.Run("by Instrument", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidPreparationInstrumentsByInstrument(ctx, &mealplanningsvc.GetValidPreparationInstrumentsByInstrumentRequest{ValidInstrumentID: validInstrument.ID})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidPreparationInstruments))
	})

	T.Run("by preparation", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		results, err := adminClient.GetValidPreparationInstrumentsByPreparation(ctx, &mealplanningsvc.GetValidPreparationInstrumentsByPreparationRequest{ValidPreparationID: validPreparation.ID})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.True(t, len(results.Results) >= len(createdValidPreparationInstruments))
	})
}
