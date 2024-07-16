package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createValidVesselForTest(t *testing.T, ctx context.Context, exampleValidVessel *types.ValidVessel, dbc *Querier) *types.ValidVessel {
	t.Helper()

	// create
	if exampleValidVessel == nil {
		createdValidMeasurementUnit := createValidMeasurementUnitForTest(t, ctx, nil, dbc)
		exampleValidVessel = fakes.BuildFakeValidVessel()
		exampleValidVessel.CapacityUnit = createdValidMeasurementUnit
	}
	dbInput := converters.ConvertValidVesselToValidVesselDatabaseCreationInput(exampleValidVessel)

	created, err := dbc.CreateValidVessel(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleValidVessel.CreatedAt = created.CreatedAt
	exampleValidVessel.CapacityUnit = &types.ValidMeasurementUnit{ID: exampleValidVessel.CapacityUnit.ID}
	assert.Equal(t, exampleValidVessel, created)

	validVessel, err := dbc.GetValidVessel(ctx, created.ID)
	exampleValidVessel.CreatedAt = validVessel.CreatedAt
	exampleValidVessel.CapacityUnit = validVessel.CapacityUnit

	assert.NoError(t, err)
	assert.Equal(t, validVessel, exampleValidVessel)

	return created
}

func TestQuerier_Integration_ValidVessels(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	databaseURI, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, databaseURI)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	exampleValidVessel := fakes.BuildFakeValidVessel()
	createdValidVessels := []*types.ValidVessel{}

	createdValidMeasurementUnit := createValidMeasurementUnitForTest(t, ctx, nil, dbc)
	exampleValidVessel.CapacityUnit = createdValidMeasurementUnit

	// create
	createdValidVessels = append(createdValidVessels, createValidVesselForTest(t, ctx, exampleValidVessel, dbc))

	// update
	updatedValidVessel := fakes.BuildFakeValidVessel()
	updatedValidVessel.ID = createdValidVessels[0].ID
	updatedValidVessel.CapacityUnit = createdValidMeasurementUnit
	assert.NoError(t, dbc.UpdateValidVessel(ctx, updatedValidVessel))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidVessel()
		input.CapacityUnit = createdValidMeasurementUnit
		input.Name = fmt.Sprintf("%s %d", updatedValidVessel.Name, i)
		createdValidVessels = append(createdValidVessels, createValidVesselForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validVessels, err := dbc.GetValidVessels(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validVessels.Data)
	assert.Equal(t, len(createdValidVessels), len(validVessels.Data))

	// fetch as list of IDs
	validVesselIDs := []string{}
	for _, validVessel := range createdValidVessels {
		validVesselIDs = append(validVesselIDs, validVessel.ID)
	}

	byIDs, err := dbc.GetValidVesselsWithIDs(ctx, validVesselIDs)
	assert.NoError(t, err)

	for i, v := range byIDs {
		validVessels.Data[i].CreatedAt = v.CreatedAt
		validVessels.Data[i].LastUpdatedAt = v.LastUpdatedAt
		validVessels.Data[i].CapacityUnit.CreatedAt = v.CapacityUnit.CreatedAt
		validVessels.Data[i].CapacityUnit.LastUpdatedAt = v.CapacityUnit.LastUpdatedAt
	}

	assert.Equal(t, validVessels.Data, byIDs)

	// fetch via name search
	byName, err := dbc.SearchForValidVessels(ctx, updatedValidVessel.Name)

	for i, v := range byName {
		validVessels.Data[i].CreatedAt = v.CreatedAt
		validVessels.Data[i].LastUpdatedAt = v.LastUpdatedAt
		validVessels.Data[i].CapacityUnit.CreatedAt = v.CapacityUnit.CreatedAt
		validVessels.Data[i].CapacityUnit.LastUpdatedAt = v.CapacityUnit.LastUpdatedAt
	}

	assert.NoError(t, err)
	assert.Equal(t, validVessels.Data, byName)

	random, err := dbc.GetRandomValidVessel(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, random)

	results, err := dbc.GetValidVesselIDsThatNeedSearchIndexing(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, results)

	// delete
	for _, validVessel := range createdValidVessels {
		assert.NoError(t, dbc.MarkValidVesselAsIndexed(ctx, validVessel.ID))
		assert.NoError(t, dbc.ArchiveValidVessel(ctx, validVessel.ID))

		var exists bool
		exists, err = dbc.ValidVesselExists(ctx, validVessel.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ValidVessel
		y, err = dbc.GetValidVessel(ctx, validVessel.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_ValidVesselExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid vessel ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidVesselExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid vessel ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidVessel(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_SearchForValidVessels(T *testing.T) {
	T.Parallel()

	T.Run("with invalid query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidVessels(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetValidVesselsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("with empty IDs list", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.GetValidVesselsWithIDs(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidVesselThatNeedSearchIndexing(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, db := buildTestClient(t)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidVessel(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidVessel(ctx, nil))
	})
}

func TestQuerier_ArchiveValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid vessel ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidVessel(ctx, ""))
	})
}

func TestQuerier_MarkValidVesselAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkValidVesselAsIndexed(ctx, ""))
	})
}
