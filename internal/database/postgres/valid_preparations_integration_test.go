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

func createValidPreparationForTest(t *testing.T, ctx context.Context, exampleValidPreparation *types.ValidPreparation, dbc *Querier) *types.ValidPreparation {
	t.Helper()

	// create
	if exampleValidPreparation == nil {
		exampleValidPreparation = fakes.BuildFakeValidPreparation()
	}
	dbInput := converters.ConvertValidPreparationToValidPreparationDatabaseCreationInput(exampleValidPreparation)

	created, err := dbc.CreateValidPreparation(ctx, dbInput)
	exampleValidPreparation.CreatedAt = created.CreatedAt
	assert.NoError(t, err)
	assert.Equal(t, exampleValidPreparation, created)

	validPreparation, err := dbc.GetValidPreparation(ctx, created.ID)
	exampleValidPreparation.CreatedAt = validPreparation.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, validPreparation, exampleValidPreparation)

	return created
}

func TestQuerier_Integration_ValidPreparations(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

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
	assert.NoError(t, dbc.UpdateValidPreparation(ctx, updatedValidPreparation))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidPreparation()
		input.Name = fmt.Sprintf("%s %d", updatedValidPreparation.Name, i)
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
	byName, err := dbc.SearchForValidPreparations(ctx, updatedValidPreparation.Name)
	assert.NoError(t, err)
	assert.Equal(t, validPreparations.Data, byName)

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
