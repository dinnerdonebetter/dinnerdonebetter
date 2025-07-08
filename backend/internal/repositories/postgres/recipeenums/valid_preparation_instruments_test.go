package recipeenums

import (
	"context"
	"database/sql"
	"testing"

	types "github.com/dinnerdonebetter/backend/internal/domain/recipeenums"
	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums/fakes"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuerier_Integration_ValidPreparationInstruments(t *testing.T) {
	if !pgtesting.RunContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t)

	databaseURI, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, databaseURI)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	exampleValidInstrument := CreateValidInstrumentForTest(t, ctx, nil, dbc)
	exampleValidPreparation := CreateValidPreparationForTest(t, ctx, nil, dbc)
	exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
	exampleValidPreparationInstrument.Preparation = *exampleValidPreparation
	exampleValidPreparationInstrument.Instrument = *exampleValidInstrument
	createdValidPreparationInstruments := []*types.ValidPreparationInstrument{}

	// create
	createdValidPreparationInstruments = append(createdValidPreparationInstruments, CreateValidPreparationInstrumentForTest(t, ctx, exampleValidPreparationInstrument, dbc))

	// update
	updatedValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
	updatedValidPreparationInstrument.ID = createdValidPreparationInstruments[0].ID
	updatedValidPreparationInstrument.Preparation = createdValidPreparationInstruments[0].Preparation
	updatedValidPreparationInstrument.Instrument = createdValidPreparationInstruments[0].Instrument
	assert.NoError(t, dbc.UpdateValidPreparationInstrument(ctx, updatedValidPreparationInstrument))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidPreparationInstrument()
		input.Preparation = createdValidPreparationInstruments[0].Preparation
		input.Instrument = createdValidPreparationInstruments[0].Instrument
		createdValidPreparationInstruments = append(createdValidPreparationInstruments, CreateValidPreparationInstrumentForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validPreparationInstruments, err := dbc.GetValidPreparationInstruments(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validPreparationInstruments.Data)
	assert.Equal(t, len(createdValidPreparationInstruments), len(validPreparationInstruments.Data))

	forPreparation, err := dbc.GetValidPreparationInstrumentsForPreparation(ctx, createdValidPreparationInstruments[0].Preparation.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, forPreparation.Data)

	forInstrument, err := dbc.GetValidPreparationInstrumentsForInstrument(ctx, createdValidPreparationInstruments[0].Instrument.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, forInstrument.Data)

	// delete
	for _, validPreparationInstrument := range createdValidPreparationInstruments {
		assert.NoError(t, dbc.ArchiveValidPreparationInstrument(ctx, validPreparationInstrument.ID))

		var exists bool
		exists, err = dbc.ValidPreparationInstrumentExists(ctx, validPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ValidPreparationInstrument
		y, err = dbc.GetValidPreparationInstrument(ctx, validPreparationInstrument.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_ValidPreparationInstrumentExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c := buildInertClientForTest(t)

		actual, err := c.ValidPreparationInstrumentExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.GetValidPreparationInstrument(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.CreateValidPreparationInstrument(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		assert.Error(t, c.UpdateValidPreparationInstrument(ctx, nil))
	})
}

func TestQuerier_ArchiveValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		assert.Error(t, c.ArchiveValidPreparationInstrument(ctx, ""))
	})
}
