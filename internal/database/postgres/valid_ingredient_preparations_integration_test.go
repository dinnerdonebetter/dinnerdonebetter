package postgres

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func createValidIngredientPreparationForTest(t *testing.T, ctx context.Context, exampleValidIngredientPreparation *types.ValidIngredientPreparation, dbc *Querier) *types.ValidIngredientPreparation {
	t.Helper()

	// create
	if exampleValidIngredientPreparation == nil {
		exampleValidIngredient := createValidIngredientForTest(t, ctx, nil, dbc)
		exampleValidPreparation := createValidPreparationForTest(t, ctx, nil, dbc)
		exampleValidIngredientPreparation = fakes.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.Ingredient = *exampleValidIngredient
		exampleValidIngredientPreparation.Preparation = *exampleValidPreparation

	}
	dbInput := converters.ConvertValidIngredientPreparationToValidIngredientPreparationDatabaseCreationInput(exampleValidIngredientPreparation)

	created, err := dbc.CreateValidIngredientPreparation(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)
	exampleValidIngredientPreparation.CreatedAt = created.CreatedAt
	exampleValidIngredientPreparation.Preparation = types.ValidPreparation{ID: exampleValidIngredientPreparation.Preparation.ID}
	exampleValidIngredientPreparation.Ingredient = types.ValidIngredient{ID: exampleValidIngredientPreparation.Ingredient.ID}
	assert.Equal(t, exampleValidIngredientPreparation, created)

	validIngredientPreparation, err := dbc.GetValidIngredientPreparation(ctx, created.ID)
	exampleValidIngredientPreparation.CreatedAt = validIngredientPreparation.CreatedAt
	exampleValidIngredientPreparation.Preparation = validIngredientPreparation.Preparation
	exampleValidIngredientPreparation.Ingredient = validIngredientPreparation.Ingredient

	assert.NoError(t, err)
	assert.Equal(t, validIngredientPreparation, exampleValidIngredientPreparation)

	return created
}

func TestQuerier_Integration_ValidIngredientPreparations(t *testing.T) {
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
	exampleValidPreparation := createValidPreparationForTest(t, ctx, nil, dbc)
	exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
	exampleValidIngredientPreparation.Preparation = *exampleValidPreparation
	exampleValidIngredientPreparation.Ingredient = *exampleValidIngredient
	createdValidIngredientPreparations := []*types.ValidIngredientPreparation{}

	// create
	createdValidIngredientPreparations = append(createdValidIngredientPreparations, createValidIngredientPreparationForTest(t, ctx, exampleValidIngredientPreparation, dbc))

	// update
	updatedValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
	updatedValidIngredientPreparation.ID = createdValidIngredientPreparations[0].ID
	updatedValidIngredientPreparation.Preparation = createdValidIngredientPreparations[0].Preparation
	updatedValidIngredientPreparation.Ingredient = createdValidIngredientPreparations[0].Ingredient
	assert.NoError(t, dbc.UpdateValidIngredientPreparation(ctx, updatedValidIngredientPreparation))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidIngredientPreparation()
		input.Preparation = createdValidIngredientPreparations[0].Preparation
		input.Ingredient = createdValidIngredientPreparations[0].Ingredient
		createdValidIngredientPreparations = append(createdValidIngredientPreparations, createValidIngredientPreparationForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validIngredientPreparations, err := dbc.GetValidIngredientPreparations(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validIngredientPreparations.Data)
	assert.Equal(t, len(createdValidIngredientPreparations), len(validIngredientPreparations.Data))

	// fetch as list of IDs
	validIngredientPreparationIDs := []string{}
	for _, validIngredientPreparation := range createdValidIngredientPreparations {
		validIngredientPreparationIDs = append(validIngredientPreparationIDs, validIngredientPreparation.ID)
	}

	// delete
	for _, validIngredientPreparation := range createdValidIngredientPreparations {
		assert.NoError(t, dbc.ArchiveValidIngredientPreparation(ctx, validIngredientPreparation.ID))

		var y *types.ValidIngredientPreparation
		y, err = dbc.GetValidIngredientPreparation(ctx, validIngredientPreparation.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
