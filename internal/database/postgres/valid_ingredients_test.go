package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestQuerier_ValidIngredientExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidIngredientExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredient(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_SearchForValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)
		filter := types.DefaultQueryFilter()

		actual, err := c.SearchForValidIngredients(ctx, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_SearchForValidIngredientsForPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidIngredientsForPreparation(ctx, exampleValidPreparation.ID, exampleValidIngredient.Name, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidIngredient(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidIngredient(ctx, nil))
	})
}

func TestQuerier_ArchiveValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidIngredient(ctx, ""))
	})
}

func TestQuerier_MarkValidIngredientAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkValidIngredientAsIndexed(ctx, ""))
	})
}
