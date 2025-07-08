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

func TestQuerier_Integration_ValidPreparationVessels(t *testing.T) {
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

	exampleValidVessel := CreateValidVesselForTest(t, ctx, nil, dbc)
	exampleValidPreparation := CreateValidPreparationForTest(t, ctx, nil, dbc)
	exampleValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()
	exampleValidPreparationVessel.Preparation = *exampleValidPreparation
	exampleValidPreparationVessel.Vessel = *exampleValidVessel
	createdValidPreparationVessels := []*types.ValidPreparationVessel{}

	// create
	createdValidPreparationVessels = append(createdValidPreparationVessels, CreateValidPreparationVesselForTest(t, ctx, exampleValidPreparationVessel, dbc))

	// update
	updatedValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()
	updatedValidPreparationVessel.ID = createdValidPreparationVessels[0].ID
	updatedValidPreparationVessel.Preparation = createdValidPreparationVessels[0].Preparation
	updatedValidPreparationVessel.Vessel = createdValidPreparationVessels[0].Vessel
	assert.NoError(t, dbc.UpdateValidPreparationVessel(ctx, updatedValidPreparationVessel))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidPreparationVessel()
		input.Preparation = createdValidPreparationVessels[0].Preparation
		input.Vessel = createdValidPreparationVessels[0].Vessel
		createdValidPreparationVessels = append(createdValidPreparationVessels, CreateValidPreparationVesselForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validPreparationVessels, err := dbc.GetValidPreparationVessels(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validPreparationVessels.Data)
	assert.Equal(t, len(createdValidPreparationVessels), len(validPreparationVessels.Data))

	forPreparation, err := dbc.GetValidPreparationVesselsForPreparation(ctx, createdValidPreparationVessels[0].Preparation.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, forPreparation.Data)

	forVessel, err := dbc.GetValidPreparationVesselsForVessel(ctx, createdValidPreparationVessels[0].Vessel.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, forVessel.Data)

	// delete
	for _, validPreparationVessel := range createdValidPreparationVessels {
		assert.NoError(t, dbc.ArchiveValidPreparationVessel(ctx, validPreparationVessel.ID))

		var exists bool
		exists, err = dbc.ValidPreparationVesselExists(ctx, validPreparationVessel.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ValidPreparationVessel
		y, err = dbc.GetValidPreparationVessel(ctx, validPreparationVessel.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_ValidPreparationVesselExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c := buildInertClientForTest(t)

		actual, err := c.ValidPreparationVesselExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.GetValidPreparationVessel(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.CreateValidPreparationVessel(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		assert.Error(t, c.UpdateValidPreparationVessel(ctx, nil))
	})
}

func TestQuerier_ArchiveValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		assert.Error(t, c.ArchiveValidPreparationVessel(ctx, ""))
	})
}
