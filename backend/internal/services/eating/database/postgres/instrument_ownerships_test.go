package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/converters"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createInstrumentOwnershipForTest(t *testing.T, ctx context.Context, exampleInstrumentOwnership *types.InstrumentOwnership, dbc *Querier) *types.InstrumentOwnership {
	t.Helper()

	// create
	if exampleInstrumentOwnership == nil {
		exampleInstrumentOwnership = fakes.BuildFakeInstrumentOwnership()
	}
	dbInput := converters.ConvertInstrumentOwnershipToInstrumentOwnershipDatabaseCreationInput(exampleInstrumentOwnership)

	created, err := dbc.CreateInstrumentOwnership(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleInstrumentOwnership.CreatedAt = created.CreatedAt
	assert.Equal(t, exampleInstrumentOwnership.Instrument.ID, created.Instrument.ID)
	exampleInstrumentOwnership.Instrument = created.Instrument
	assert.Equal(t, exampleInstrumentOwnership, created)

	householdInstrumentOwnership, err := dbc.GetInstrumentOwnership(ctx, created.ID, created.BelongsToHousehold)
	exampleInstrumentOwnership.CreatedAt = householdInstrumentOwnership.CreatedAt
	assert.Equal(t, exampleInstrumentOwnership.Instrument.ID, householdInstrumentOwnership.Instrument.ID)
	exampleInstrumentOwnership.Instrument = householdInstrumentOwnership.Instrument

	assert.NoError(t, err)
	assert.Equal(t, householdInstrumentOwnership, exampleInstrumentOwnership)

	return created
}

func TestQuerier_Integration_InstrumentOwnerships(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := t.Context()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	databaseURI, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, databaseURI)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := createUserForTest(t, ctx, nil, dbc)
	householdID, err := dbc.GetDefaultHouseholdIDForUser(ctx, user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, householdID)

	instrument := createValidInstrumentForTest(t, ctx, nil, dbc)

	exampleInstrumentOwnership := fakes.BuildFakeInstrumentOwnership()
	exampleInstrumentOwnership.BelongsToHousehold = householdID
	exampleInstrumentOwnership.Instrument = *instrument
	createdInstrumentOwnerships := []*types.InstrumentOwnership{}

	// create
	createdInstrumentOwnerships = append(createdInstrumentOwnerships, createInstrumentOwnershipForTest(t, ctx, exampleInstrumentOwnership, dbc))

	// update
	assert.NoError(t, dbc.UpdateInstrumentOwnership(ctx, createdInstrumentOwnerships[0]))

	// create more
	for range exampleQuantity {
		newInstrument := createValidInstrumentForTest(t, ctx, nil, dbc)
		input := fakes.BuildFakeInstrumentOwnership()
		input.BelongsToHousehold = householdID
		input.Instrument = *newInstrument
		createdInstrumentOwnerships = append(createdInstrumentOwnerships, createInstrumentOwnershipForTest(t, ctx, input, dbc))
	}

	// fetch as list
	householdInstrumentOwnerships, err := dbc.GetInstrumentOwnerships(ctx, householdID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, householdInstrumentOwnerships.Data)
	assert.Equal(t, len(createdInstrumentOwnerships), len(householdInstrumentOwnerships.Data))

	// delete
	for _, householdInstrumentOwnership := range createdInstrumentOwnerships {
		assert.NoError(t, dbc.ArchiveInstrumentOwnership(ctx, householdInstrumentOwnership.ID, householdID))

		var exists bool
		exists, err = dbc.InstrumentOwnershipExists(ctx, householdInstrumentOwnership.ID, householdID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.InstrumentOwnership
		y, err = dbc.GetInstrumentOwnership(ctx, householdInstrumentOwnership.ID, householdID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_InstrumentOwnershipExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := t.Context()
		c, _ := buildTestClient(t)

		actual, err := c.InstrumentOwnershipExists(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetHouseholdInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := t.Context()
		c, _ := buildTestClient(t)

		actual, err := c.GetInstrumentOwnership(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildTestClient(t)

		actual, err := c.CreateInstrumentOwnership(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateInstrumentOwnership(ctx, nil))
	})
}

func TestQuerier_ArchiveInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := t.Context()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveInstrumentOwnership(ctx, "", exampleHouseholdID))
	})
}
