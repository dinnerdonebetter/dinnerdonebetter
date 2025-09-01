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

func createValidInstrumentForTest(t *testing.T) *mealplanning.ValidInstrument {
	t.Helper()

	ctx := t.Context()

	creationRequestInput := fakes.BuildFakeValidInstrumentCreationRequestInput()
	convertedInput := grpcconverters.ConvertValidInstrumentCreationRequestInputToGRPCValidInstrumentCreationRequestInput(creationRequestInput)

	created, err := adminClient.CreateValidInstrument(ctx, &mealplanningsvc.CreateValidInstrumentRequest{
		Input: convertedInput,
	})
	require.NoError(t, err)
	assert.NotNil(t, created)

	return grpcconverters.ConvertGRPCValidInstrumentToValidInstrument(created.Result)
}

func TestValidInstruments_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		createValidInstrumentForTest(t)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeValidInstrumentCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidInstrumentCreationRequestInputToGRPCValidInstrumentCreationRequestInput(creationRequestInput)

		c := buildUnauthenticatedGRPCClientForTest(t)
		created, err := c.CreateValidInstrument(ctx, &mealplanningsvc.CreateValidInstrumentRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeValidInstrumentCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidInstrumentCreationRequestInputToGRPCValidInstrumentCreationRequestInput(creationRequestInput)
		// this is not allowed
		convertedInput.Name = ""

		created, err := adminClient.CreateValidInstrument(ctx, &mealplanningsvc.CreateValidInstrumentRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("non-admin users are forbidden from creating", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)

		creationRequestInput := fakes.BuildFakeValidInstrumentCreationRequestInput()
		convertedInput := grpcconverters.ConvertValidInstrumentCreationRequestInputToGRPCValidInstrumentCreationRequestInput(creationRequestInput)

		created, err := testClient.CreateValidInstrument(ctx, &mealplanningsvc.CreateValidInstrumentRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})
}

func TestValidInstruments_Reading(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidInstrumentForTest(t)

		retrieved, err := testClient.GetValidInstrument(ctx, &mealplanningsvc.GetValidInstrumentRequest{ValidInstrumentID: created.ID})
		assert.NoError(t, err)

		converted := grpcconverters.ConvertGRPCValidInstrumentToValidInstrument(retrieved.Result)

		assertRoughEquality(t, created, converted, "CreatedAt", "LastUpdatedAt", "ArchivedAt")
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidInstrumentForTest(t)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.GetValidInstrument(ctx, &mealplanningsvc.GetValidInstrumentRequest{ValidInstrumentID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.GetValidInstrument(ctx, &mealplanningsvc.GetValidInstrumentRequest{ValidInstrumentID: nonexistentID})
		assert.Error(t, err)
	})
}

func TestValidInstruments_Updating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidInstrumentForTest(t)

		updateInput := fakes.BuildFakeValidInstrumentUpdateRequestInput()
		created.Update(updateInput)

		response, err := adminClient.UpdateValidInstrument(ctx, &mealplanningsvc.UpdateValidInstrumentRequest{
			ValidInstrumentID: created.ID,
			Input:             grpcconverters.ConvertValidInstrumentUpdateRequestInputToGRPCValidInstrumentUpdateRequestInput(updateInput),
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

		created := createValidInstrumentForTest(t)

		updateInput := fakes.BuildFakeValidInstrumentUpdateRequestInput()
		created.Update(updateInput)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.UpdateValidInstrument(ctx, &mealplanningsvc.UpdateValidInstrumentRequest{
			ValidInstrumentID: created.ID,
			Input:             grpcconverters.ConvertValidInstrumentUpdateRequestInputToGRPCValidInstrumentUpdateRequestInput(updateInput),
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

		created := createValidInstrumentForTest(t)

		response, err := testClient.UpdateValidInstrument(ctx, &mealplanningsvc.UpdateValidInstrumentRequest{
			ValidInstrumentID: created.ID,
			Input: &mealplanningsvc.ValidInstrumentUpdateRequestInput{
				Name: pointer.To("doesn't matter"),
			},
		})
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestValidInstruments_Archiving(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidInstrumentForTest(t)

		_, err := adminClient.ArchiveValidInstrument(ctx, &mealplanningsvc.ArchiveValidInstrumentRequest{ValidInstrumentID: created.ID})
		assert.NoError(t, err)

		x, err := adminClient.GetValidInstrument(ctx, &mealplanningsvc.GetValidInstrumentRequest{ValidInstrumentID: created.ID})
		assert.Nil(t, x)
		assert.Error(t, err)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidInstrumentForTest(t)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.ArchiveValidInstrument(ctx, &mealplanningsvc.ArchiveValidInstrumentRequest{ValidInstrumentID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.ArchiveValidInstrument(ctx, &mealplanningsvc.ArchiveValidInstrumentRequest{ValidInstrumentID: nonexistentID})
		assert.Error(t, err)
	})

	T.Run("non-admin users are forbidden from archiving", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createValidInstrumentForTest(t)
		_, testClient := createUserAndClientForTest(T)

		_, err := testClient.ArchiveValidInstrument(ctx, &mealplanningsvc.ArchiveValidInstrumentRequest{ValidInstrumentID: created.ID})
		assert.Error(t, err)
	})
}

func TestValidInstruments_GetRandom(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		// in case we haven't already
		createValidInstrumentForTest(t)

		response, err := testClient.GetRandomValidInstrument(ctx, &mealplanningsvc.GetRandomValidInstrumentRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, response)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		response, err := c.GetRandomValidInstrument(ctx, &mealplanningsvc.GetRandomValidInstrumentRequest{})
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestValidInstruments_Listing(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)
	createdValidInstruments := []*mealplanning.ValidInstrument{}
	for range exampleQuantity {
		created := createValidInstrumentForTest(T)
		createdValidInstruments = append(createdValidInstruments, created)
	}

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		retrieved, err := testClient.GetValidInstruments(ctx, &mealplanningsvc.GetValidInstrumentsRequest{})
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.True(t, len(retrieved.Results) >= len(createdValidInstruments))
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.GetValidInstruments(ctx, &mealplanningsvc.GetValidInstrumentsRequest{})
		assert.Error(t, err)
	})
}

func TestValidInstruments_Searching(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)
		created := createValidInstrumentForTest(t)

		retrieved, err := testClient.SearchForValidInstruments(ctx, &mealplanningsvc.SearchForValidInstrumentsRequest{
			Query: created.Name[:2],
		})
		require.NoError(t, err)
		require.NotNil(t, retrieved)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.SearchForValidInstruments(ctx, &mealplanningsvc.SearchForValidInstrumentsRequest{})
		assert.Error(t, err)
	})
}
