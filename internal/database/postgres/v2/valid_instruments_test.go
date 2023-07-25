package v2

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createValidInstrumentForTest(t *testing.T, ctx context.Context, exampleValidInstrument *types.ValidInstrument, dbc *DatabaseClient) *types.ValidInstrument {
	t.Helper()

	// create
	if exampleValidInstrument == nil {
		exampleValidInstrument = fakes.BuildFakeValidInstrument()
	}
	var x ValidInstrument
	require.NoError(t, copier.Copy(&x, exampleValidInstrument))

	created, err := dbc.CreateValidInstrument(ctx, &x)
	assert.NoError(t, err)
	assert.Equal(t, exampleValidInstrument, created)

	validInstrument, err := dbc.GetValidInstrument(ctx, created.ID)
	exampleValidInstrument.CreatedAt = validInstrument.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, validInstrument, exampleValidInstrument)

	return created
}

func TestDatabaseClient_ValidInstruments(t *testing.T) {
	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	exampleValidInstrument := fakes.BuildFakeValidInstrument()
	createdValidInstruments := []*types.ValidInstrument{}

	// create
	createdValidInstruments = append(createdValidInstruments, createValidInstrumentForTest(t, ctx, exampleValidInstrument, dbc))

	// update
	updatedValidInstrument := fakes.BuildFakeValidInstrument()
	updatedValidInstrument.ID = createdValidInstruments[0].ID
	var x ValidInstrument
	require.NoError(t, copier.Copy(&x, updatedValidInstrument))
	assert.NoError(t, dbc.UpdateValidInstrument(ctx, updatedValidInstrument))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidInstrument()
		input.Name = fmt.Sprintf("%s %d", exampleValidInstrument.Name, i)
		createdValidInstruments = append(createdValidInstruments, createValidInstrumentForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validInstruments, err := dbc.GetValidInstruments(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validInstruments.Data)
	assert.Equal(t, len(createdValidInstruments), len(validInstruments.Data))

	// fetch as list of IDs
	validInstrumentIDs := []string{}
	for _, validInstrument := range createdValidInstruments {
		validInstrumentIDs = append(validInstrumentIDs, validInstrument.ID)
	}

	byIDs, err := dbc.GetValidInstrumentsWithIDs(ctx, validInstrumentIDs)
	assert.NoError(t, err)
	assert.Equal(t, validInstruments.Data, byIDs)

	// fetch via name search
	byName, err := dbc.SearchForValidInstruments(ctx, exampleValidInstrument.Name, nil)
	assert.NoError(t, err)
	assert.Equal(t, validInstruments, byName)

	// delete
	for _, validInstrument := range createdValidInstruments {
		assert.NoError(t, dbc.ArchiveValidInstrument(ctx, validInstrument.ID))

		var y *types.ValidInstrument
		y, err = dbc.GetValidInstrument(ctx, validInstrument.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}