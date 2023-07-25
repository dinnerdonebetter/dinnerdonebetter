package v2

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createValidIngredientForTest(t *testing.T, ctx context.Context, dbc *DatabaseClient) *types.ValidIngredient {
	t.Helper()

	// create
	exampleValidIngredient := fakes.BuildFakeValidIngredient()
	var x ValidIngredient
	require.NoError(t, copier.Copy(&x, exampleValidIngredient))

	created, err := dbc.CreateValidIngredient(ctx, &x)
	assert.NoError(t, err)
	assert.Equal(t, exampleValidIngredient, created)

	validIngredient, err := dbc.GetValidIngredient(ctx, created.ID)
	exampleValidIngredient.CreatedAt = validIngredient.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, validIngredient, exampleValidIngredient)

	return created
}

func TestDatabaseClient_ValidIngredients(t *testing.T) {
	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	createdValidIngredients := []*types.ValidIngredient{}

	// create
	createdValidIngredients = append(createdValidIngredients, createValidIngredientForTest(t, ctx, dbc))

	// update
	updatedValidIngredient := fakes.BuildFakeValidIngredient()
	updatedValidIngredient.ID = createdValidIngredients[0].ID
	var x ValidIngredient
	require.NoError(t, copier.Copy(&x, updatedValidIngredient))
	assert.NoError(t, dbc.UpdateValidIngredient(ctx, updatedValidIngredient))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		createdValidIngredients = append(createdValidIngredients, createValidIngredientForTest(t, ctx, dbc))
	}

	// fetch as list
	validIngredients, err := dbc.GetValidIngredients(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validIngredients.Data)
	assert.Equal(t, len(createdValidIngredients), len(validIngredients.Data))

	// fetch as list of IDs
	validIngredientIDs := []string{}
	for _, validIngredient := range createdValidIngredients {
		validIngredientIDs = append(validIngredientIDs, validIngredient.ID)
	}

	byIDs, err := dbc.GetValidIngredientsWithIDs(ctx, validIngredientIDs)
	assert.NoError(t, err)
	assert.Equal(t, validIngredients.Data, byIDs)

	// delete
	for _, validIngredient := range createdValidIngredients {
		assert.NoError(t, dbc.ArchiveValidIngredient(ctx, validIngredient.ID))

		var y *types.ValidIngredient
		y, err = dbc.GetValidIngredient(ctx, validIngredient.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
