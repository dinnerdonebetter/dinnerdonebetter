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
)

func createValidMeasurementUnitForTest(t *testing.T, ctx context.Context, exampleValidMeasurementUnit *types.ValidMeasurementUnit, dbc *Querier) *types.ValidMeasurementUnit {
	t.Helper()

	// create
	if exampleValidMeasurementUnit == nil {
		exampleValidMeasurementUnit = fakes.BuildFakeValidMeasurementUnit()
	}
	dbInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitDatabaseCreationInput(exampleValidMeasurementUnit)

	created, err := dbc.CreateValidMeasurementUnit(ctx, dbInput)
	exampleValidMeasurementUnit.CreatedAt = created.CreatedAt
	assert.NoError(t, err)
	assert.Equal(t, exampleValidMeasurementUnit, created)

	validMeasurementUnit, err := dbc.GetValidMeasurementUnit(ctx, created.ID)
	exampleValidMeasurementUnit.CreatedAt = validMeasurementUnit.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, validMeasurementUnit, exampleValidMeasurementUnit)

	return created
}

func TestQuerier_Integration_ValidMeasurementUnits(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
	createdValidMeasurementUnits := []*types.ValidMeasurementUnit{}

	// create
	createdValidMeasurementUnits = append(createdValidMeasurementUnits, createValidMeasurementUnitForTest(t, ctx, exampleValidMeasurementUnit, dbc))

	// update
	updatedValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
	updatedValidMeasurementUnit.ID = createdValidMeasurementUnits[0].ID
	assert.NoError(t, dbc.UpdateValidMeasurementUnit(ctx, updatedValidMeasurementUnit))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidMeasurementUnit()
		input.Name = fmt.Sprintf("%s %d", updatedValidMeasurementUnit.Name, i)
		createdValidMeasurementUnits = append(createdValidMeasurementUnits, createValidMeasurementUnitForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validMeasurementUnits, err := dbc.GetValidMeasurementUnits(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validMeasurementUnits.Data)
	assert.GreaterOrEqual(t, len(validMeasurementUnits.Data), len(createdValidMeasurementUnits))

	// fetch as list of IDs
	validMeasurementUnitIDs := []string{}
	for _, validMeasurementUnit := range createdValidMeasurementUnits {
		validMeasurementUnitIDs = append(validMeasurementUnitIDs, validMeasurementUnit.ID)
	}

	byIDs, err := dbc.GetValidMeasurementUnitsWithIDs(ctx, validMeasurementUnitIDs)
	assert.NoError(t, err)
	assert.Subset(t, validMeasurementUnits.Data, byIDs)

	// fetch via name search
	byName, err := dbc.SearchForValidMeasurementUnits(ctx, updatedValidMeasurementUnit.Name)
	assert.NoError(t, err)
	assert.Subset(t, validMeasurementUnits.Data, byName)

	// delete
	for _, validMeasurementUnit := range createdValidMeasurementUnits {
		assert.NoError(t, dbc.ArchiveValidMeasurementUnit(ctx, validMeasurementUnit.ID))

		var y *types.ValidMeasurementUnit
		y, err = dbc.GetValidMeasurementUnit(ctx, validMeasurementUnit.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}