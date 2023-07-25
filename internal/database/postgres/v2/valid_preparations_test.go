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

func createValidPreparationForTest(t *testing.T, ctx context.Context, exampleValidPreparation *types.ValidPreparation, dbc *DatabaseClient) *types.ValidPreparation {
	t.Helper()

	// create
	if exampleValidPreparation == nil {
		exampleValidPreparation = fakes.BuildFakeValidPreparation()
	}
	var x ValidPreparation
	require.NoError(t, copier.Copy(&x, exampleValidPreparation))

	created, err := dbc.CreateValidPreparation(ctx, &x)
	assert.NoError(t, err)
	assert.Equal(t, exampleValidPreparation, created)

	validPreparation, err := dbc.GetValidPreparation(ctx, created.ID)
	exampleValidPreparation.CreatedAt = validPreparation.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, validPreparation, exampleValidPreparation)

	return created
}

func TestDatabaseClient_ValidPreparations(t *testing.T) {
	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	exampleValidPreparation := fakes.BuildFakeValidPreparation()
	createdValidPreparations := []*types.ValidPreparation{}

	// create
	createdValidPreparations = append(createdValidPreparations, createValidPreparationForTest(t, ctx, exampleValidPreparation, dbc))

	// update
	updatedValidPreparation := fakes.BuildFakeValidPreparation()
	updatedValidPreparation.ID = createdValidPreparations[0].ID
	var x ValidPreparation
	require.NoError(t, copier.Copy(&x, updatedValidPreparation))
	assert.NoError(t, dbc.UpdateValidPreparation(ctx, updatedValidPreparation))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidPreparation()
		input.Name = fmt.Sprintf("%s %d", exampleValidPreparation.Name, i)
		createdValidPreparations = append(createdValidPreparations, createValidPreparationForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validPreparations, err := dbc.GetValidPreparations(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validPreparations.Data)
	assert.Equal(t, len(createdValidPreparations), len(validPreparations.Data))

	// fetch as list of IDs
	validPreparationIDs := []string{}
	for _, validPreparation := range createdValidPreparations {
		validPreparationIDs = append(validPreparationIDs, validPreparation.ID)
	}

	byIDs, err := dbc.GetValidPreparationsWithIDs(ctx, validPreparationIDs)
	assert.NoError(t, err)
	assert.Equal(t, validPreparations.Data, byIDs)

	// fetch via name search
	byName, err := dbc.SearchForValidPreparations(ctx, exampleValidPreparation.Name, nil)
	assert.NoError(t, err)
	assert.Equal(t, validPreparations, byName)

	// delete
	for _, validPreparation := range createdValidPreparations {
		assert.NoError(t, dbc.ArchiveValidPreparation(ctx, validPreparation.ID))

		var y *types.ValidPreparation
		y, err = dbc.GetValidPreparation(ctx, validPreparation.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
