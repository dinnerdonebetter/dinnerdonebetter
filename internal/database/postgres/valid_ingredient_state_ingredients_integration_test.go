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

func createValidIngredientStateIngredientForTest(t *testing.T, ctx context.Context, exampleValidIngredientStateIngredient *types.ValidIngredientStateIngredient, dbc *Querier) *types.ValidIngredientStateIngredient {
	t.Helper()

	// create
	if exampleValidIngredientStateIngredient == nil {
		exampleValidIngredient := createValidIngredientForTest(t, ctx, nil, dbc)
		exampleValidIngredientState := createValidIngredientStateForTest(t, ctx, nil, dbc)
		exampleValidIngredientStateIngredient = fakes.BuildFakeValidIngredientStateIngredient()
		exampleValidIngredientStateIngredient.Ingredient = *exampleValidIngredient
		exampleValidIngredientStateIngredient.IngredientState = *exampleValidIngredientState
	}

	dbInput := converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientDatabaseCreationInput(exampleValidIngredientStateIngredient)

	created, err := dbc.CreateValidIngredientStateIngredient(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)
	exampleValidIngredientStateIngredient.CreatedAt = created.CreatedAt
	exampleValidIngredientStateIngredient.IngredientState = types.ValidIngredientState{ID: exampleValidIngredientStateIngredient.IngredientState.ID}
	exampleValidIngredientStateIngredient.Ingredient = types.ValidIngredient{ID: exampleValidIngredientStateIngredient.Ingredient.ID}
	assert.Equal(t, exampleValidIngredientStateIngredient, created)

	validIngredientStateIngredient, err := dbc.GetValidIngredientStateIngredient(ctx, created.ID)
	exampleValidIngredientStateIngredient.CreatedAt = validIngredientStateIngredient.CreatedAt
	exampleValidIngredientStateIngredient.IngredientState = validIngredientStateIngredient.IngredientState
	exampleValidIngredientStateIngredient.Ingredient = validIngredientStateIngredient.Ingredient

	assert.NoError(t, err)
	assert.Equal(t, validIngredientStateIngredient, exampleValidIngredientStateIngredient)

	return created
}

func TestQuerier_Integration_ValidIngredientStateIngredients(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	exampleValidIngredient := createValidIngredientForTest(t, ctx, nil, dbc)
	exampleValidIngredientState := createValidIngredientStateForTest(t, ctx, nil, dbc)
	exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()
	exampleValidIngredientStateIngredient.IngredientState = *exampleValidIngredientState
	exampleValidIngredientStateIngredient.Ingredient = *exampleValidIngredient
	createdValidIngredientStateIngredients := []*types.ValidIngredientStateIngredient{}

	// create
	createdValidIngredientStateIngredients = append(createdValidIngredientStateIngredients, createValidIngredientStateIngredientForTest(t, ctx, exampleValidIngredientStateIngredient, dbc))

	// update
	updatedValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()
	updatedValidIngredientStateIngredient.ID = createdValidIngredientStateIngredients[0].ID
	updatedValidIngredientStateIngredient.IngredientState = createdValidIngredientStateIngredients[0].IngredientState
	updatedValidIngredientStateIngredient.Ingredient = createdValidIngredientStateIngredients[0].Ingredient
	assert.NoError(t, dbc.UpdateValidIngredientStateIngredient(ctx, updatedValidIngredientStateIngredient))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidIngredientStateIngredient()
		input.IngredientState = createdValidIngredientStateIngredients[0].IngredientState
		input.Ingredient = createdValidIngredientStateIngredients[0].Ingredient
		createdValidIngredientStateIngredients = append(createdValidIngredientStateIngredients, createValidIngredientStateIngredientForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validIngredientStateIngredients, err := dbc.GetValidIngredientStateIngredients(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validIngredientStateIngredients.Data)
	assert.Equal(t, len(createdValidIngredientStateIngredients), len(validIngredientStateIngredients.Data))

	forIngredientState, err := dbc.GetValidIngredientStateIngredientsForIngredientState(ctx, createdValidIngredientStateIngredients[0].IngredientState.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, forIngredientState.Data)

	forIngredient, err := dbc.GetValidIngredientStateIngredientsForIngredient(ctx, createdValidIngredientStateIngredients[0].Ingredient.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, forIngredient.Data)

	// delete
	for _, validIngredientStateIngredient := range createdValidIngredientStateIngredients {
		assert.NoError(t, dbc.ArchiveValidIngredientStateIngredient(ctx, validIngredientStateIngredient.ID))

		var exists bool
		exists, err = dbc.ValidIngredientStateIngredientExists(ctx, validIngredientStateIngredient.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ValidIngredientStateIngredient
		y, err = dbc.GetValidIngredientStateIngredient(ctx, validIngredientStateIngredient.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
