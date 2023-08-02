package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, validVessels.Data, byIDs)

	// fetch via name search
	byName, err := dbc.SearchForValidVessels(ctx, updatedValidVessel.Name)
	assert.NoError(t, err)
	assert.Equal(t, validVessels.Data, byName)

	// delete
	for _, validVessel := range createdValidVessels {
		assert.NoError(t, dbc.ArchiveValidVessel(ctx, validVessel.ID))

		var y *types.ValidVessel
		y, err = dbc.GetValidVessel(ctx, validVessel.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
