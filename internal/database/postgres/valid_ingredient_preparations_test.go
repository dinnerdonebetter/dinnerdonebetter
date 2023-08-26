package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestQuerier_ValidIngredientPreparationExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

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
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientPreparation(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetValidIngredientPreparationsForIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientPreparationsForIngredient(ctx, exampleValidIngredient.ID, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetValidIngredientPreparationsForPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientPreparationsForPreparation(ctx, exampleValidPreparation.ID, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

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
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidIngredientPreparation(ctx, nil))
	})
}

func TestQuerier_ArchiveValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidIngredientPreparation(ctx, ""))
	})
}
