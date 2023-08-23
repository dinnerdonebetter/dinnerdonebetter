package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createValidMeasurementUnitConversionForTest(t *testing.T, ctx context.Context, exampleValidMeasurementUnitConversion *types.ValidMeasurementUnitConversion, dbc *Querier) *types.ValidMeasurementUnitConversion {
	t.Helper()

	// create
	if exampleValidMeasurementUnitConversion == nil {
		exampleValidMeasurementUnitConversion = fakes.BuildFakeValidMeasurementUnitConversion()
	}
	dbInput := converters.ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionDatabaseCreationInput(exampleValidMeasurementUnitConversion)

	created, err := dbc.CreateValidMeasurementUnitConversion(ctx, dbInput)
	require.NoError(t, err)
	require.NotNil(t, created)

	validMeasurementUnitConversion, err := dbc.GetValidMeasurementUnitConversion(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, validMeasurementUnitConversion)

	exampleValidMeasurementUnitConversion.CreatedAt = validMeasurementUnitConversion.CreatedAt
	assert.Equal(t, exampleValidMeasurementUnitConversion.From.ID, validMeasurementUnitConversion.From.ID)
	exampleValidMeasurementUnitConversion.From.CreatedAt = validMeasurementUnitConversion.From.CreatedAt
	assert.Equal(t, exampleValidMeasurementUnitConversion.To.ID, validMeasurementUnitConversion.To.ID)
	exampleValidMeasurementUnitConversion.To.CreatedAt = validMeasurementUnitConversion.To.CreatedAt
	exampleValidMeasurementUnitConversion.OnlyForIngredient = validMeasurementUnitConversion.OnlyForIngredient

	assert.Equal(t, validMeasurementUnitConversion, exampleValidMeasurementUnitConversion)

	return created
}

func TestQuerier_Integration_ValidMeasurementUnitConversions(t *testing.T) {
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

	from := createValidMeasurementUnitForTest(t, ctx, nil, dbc)
	to := createValidMeasurementUnitForTest(t, ctx, nil, dbc)

	exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()
	exampleValidMeasurementUnitConversion.From = *from
	exampleValidMeasurementUnitConversion.To = *to
	exampleValidMeasurementUnitConversion.OnlyForIngredient = nil
	createdValidMeasurementUnitConversions := []*types.ValidMeasurementUnitConversion{}

	// create
	createdValidMeasurementUnitConversions = append(createdValidMeasurementUnitConversions, createValidMeasurementUnitConversionForTest(t, ctx, exampleValidMeasurementUnitConversion, dbc))

	// update
	updatedValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()
	updatedValidMeasurementUnitConversion.ID = createdValidMeasurementUnitConversions[0].ID
	updatedValidMeasurementUnitConversion.From = *from
	updatedValidMeasurementUnitConversion.To = *to
	assert.NoError(t, dbc.UpdateValidMeasurementUnitConversion(ctx, updatedValidMeasurementUnitConversion))

	toUnits, err := dbc.GetValidMeasurementUnitConversionsToUnit(ctx, updatedValidMeasurementUnitConversion.To.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, toUnits)

	fromUnits, err := dbc.GetValidMeasurementUnitConversionsFromUnit(ctx, updatedValidMeasurementUnitConversion.From.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, fromUnits)

	// delete
	for _, validMeasurementUnitConversion := range createdValidMeasurementUnitConversions {
		assert.NoError(t, dbc.ArchiveValidMeasurementUnitConversion(ctx, validMeasurementUnitConversion.ID))

		var exists bool
		exists, err = dbc.ValidMeasurementUnitConversionExists(ctx, validMeasurementUnitConversion.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ValidMeasurementUnitConversion
		y, err = dbc.GetValidMeasurementUnitConversion(ctx, validMeasurementUnitConversion.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
