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

func createValidIngredientStateForTest(t *testing.T, ctx context.Context, exampleValidIngredientState *types.ValidIngredientState, dbc *Querier) *types.ValidIngredientState {
	t.Helper()

	// create
	if exampleValidIngredientState == nil {
		exampleValidIngredientState = fakes.BuildFakeValidIngredientState()
	}
	dbInput := converters.ConvertValidIngredientStateToValidIngredientStateDatabaseCreationInput(exampleValidIngredientState)

	created, err := dbc.CreateValidIngredientState(ctx, dbInput)
	exampleValidIngredientState.CreatedAt = created.CreatedAt
	assert.NoError(t, err)
	assert.Equal(t, exampleValidIngredientState, created)

	validIngredientState, err := dbc.GetValidIngredientState(ctx, created.ID)
	exampleValidIngredientState.CreatedAt = validIngredientState.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, validIngredientState, exampleValidIngredientState)

	return created
}

func TestQuerier_Integration_ValidIngredientStates(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	exampleValidIngredientState := fakes.BuildFakeValidIngredientState()
	createdValidIngredientStates := []*types.ValidIngredientState{}

	// create
	createdValidIngredientStates = append(createdValidIngredientStates, createValidIngredientStateForTest(t, ctx, exampleValidIngredientState, dbc))

	// update
	updatedValidIngredientState := fakes.BuildFakeValidIngredientState()
	updatedValidIngredientState.ID = createdValidIngredientStates[0].ID
	assert.NoError(t, dbc.UpdateValidIngredientState(ctx, updatedValidIngredientState))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidIngredientState()
		input.Name = fmt.Sprintf("%s %d", updatedValidIngredientState.Name, i)
		createdValidIngredientStates = append(createdValidIngredientStates, createValidIngredientStateForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validIngredientStates, err := dbc.GetValidIngredientStates(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validIngredientStates.Data)
	assert.Equal(t, len(createdValidIngredientStates), len(validIngredientStates.Data))

	// fetch as list of IDs
	validIngredientStateIDs := []string{}
	for _, validIngredientState := range createdValidIngredientStates {
		validIngredientStateIDs = append(validIngredientStateIDs, validIngredientState.ID)
	}

	byIDs, err := dbc.GetValidIngredientStatesWithIDs(ctx, validIngredientStateIDs)
	assert.NoError(t, err)
	assert.Equal(t, validIngredientStates.Data, byIDs)

	// fetch via name search
	byName, err := dbc.SearchForValidIngredientStates(ctx, updatedValidIngredientState.Name)
	assert.NoError(t, err)
	assert.Equal(t, validIngredientStates.Data, byName)

	// delete
	for _, validIngredientState := range createdValidIngredientStates {
		assert.NoError(t, dbc.ArchiveValidIngredientState(ctx, validIngredientState.ID))

		var exists bool
		exists, err = dbc.ValidIngredientStateExists(ctx, validIngredientState.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ValidIngredientState
		y, err = dbc.GetValidIngredientState(ctx, validIngredientState.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
