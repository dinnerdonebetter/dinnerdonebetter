package recipeenums

import (
	"context"
	"database/sql"
	"testing"

	types "github.com/dinnerdonebetter/backend/internal/domain/recipeenums"
	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums/fakes"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuerier_Integration_ValidIngredientPreparations(t *testing.T) {
	if !pgtesting.RunContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t)

	databaseURI, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, databaseURI)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	exampleValidIngredient := CreateValidIngredientForTest(t, ctx, nil, dbc)
	exampleValidPreparation := CreateValidPreparationForTest(t, ctx, nil, dbc)
	exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
	exampleValidIngredientPreparation.Preparation = *exampleValidPreparation
	exampleValidIngredientPreparation.Ingredient = *exampleValidIngredient
	createdValidIngredientPreparations := []*types.ValidIngredientPreparation{}

	// create
	createdValidIngredientPreparations = append(createdValidIngredientPreparations, CreateValidIngredientPreparationForTest(t, ctx, exampleValidIngredientPreparation, dbc))

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
		createdValidIngredientPreparations = append(createdValidIngredientPreparations, CreateValidIngredientPreparationForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validIngredientPreparations, err := dbc.GetValidIngredientPreparations(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validIngredientPreparations.Data)
	assert.Equal(t, len(createdValidIngredientPreparations), len(validIngredientPreparations.Data))

	forIngredient, err := dbc.GetValidIngredientPreparationsForIngredient(ctx, createdValidIngredientPreparations[0].Ingredient.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, forIngredient.Data)

	forPreparation, err := dbc.GetValidIngredientPreparationsForPreparation(ctx, createdValidIngredientPreparations[0].Preparation.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, forPreparation.Data)

	// delete
	for _, validIngredientPreparation := range createdValidIngredientPreparations {
		assert.NoError(t, dbc.ArchiveValidIngredientPreparation(ctx, validIngredientPreparation.ID))

		var exists bool
		exists, err = dbc.ValidIngredientPreparationExists(ctx, validIngredientPreparation.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ValidIngredientPreparation
		y, err = dbc.GetValidIngredientPreparation(ctx, validIngredientPreparation.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_ValidIngredientPreparationExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c := buildInertClientForTest(t)

		actual, err := c.ValidIngredientPreparationExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.GetValidIngredientPreparation(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetValidIngredientPreparationsForIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.GetValidIngredientPreparationsForIngredient(ctx, "", nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetValidIngredientPreparationsForPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.GetValidIngredientPreparationsForPreparation(ctx, "", nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.CreateValidIngredientPreparation(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		assert.Error(t, c.UpdateValidIngredientPreparation(ctx, nil))
	})
}

func TestQuerier_ArchiveValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		assert.Error(t, c.ArchiveValidIngredientPreparation(ctx, ""))
	})
}
